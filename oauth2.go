// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esign

// oauth2.go contains definitions for DocuSign's oauth2
// authorization scheme.  See the legacy package for
// the previous authorization schemes.

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/jfcote87/ctxclient"
	"github.com/jfcote87/oauth2"
	"github.com/jfcote87/oauth2/jws"
	"github.com/jfcote87/oauth2/jwt"
)

type demoFlag bool

func (df demoFlag) endpoint() oauth2.Endpoint {
	if df {
		// Demo endpoints
		return oauth2.Endpoint{
			AuthURL:  "https://account-d.docusign.com/oauth/auth",
			TokenURL: "https://account-d.docusign.com/oauth/token",
		}
	}
	// Production endpoints
	return oauth2.Endpoint{
		AuthURL:  "https://account.docusign.com/oauth/auth",
		TokenURL: "https://account.docusign.com/oauth/token",
	}
}

func (df demoFlag) tokenURI() string {
	if df {
		return "account-d.docusign.com"
	}
	return "account.docusign.com"
}

func (df demoFlag) userInfoPath() string {
	if df {
		return "https://account-d.docusign.com/oauth/userinfo"
	}
	return "https://account.docusign.com/oauth/userinfo"
}

// OAuth2Config allows for 3-legged oauth via a code grant mechanism
// see https://developers.docusign.com/esign-rest-api/guides/authentication/oauth2-code-grant
type OAuth2Config struct {
	// see "Create integrator key and configure settings" at
	// https://developers.docusign.com/esign-rest-api/guides
	IntegratorKey string `json:"integrator_key,omitempty"`
	// Secret generated when setting up integration in DocuSign. Leave blank for
	// implicit grant.
	Secret string `json:"secret,omitempty"`
	// The redirect URI must exactly match one of those pre-registered for the
	// integrator key. This determines where to redirect the user after they
	// successfully authenticate.
	RedirURL string `json:"redir_url,omitempty"`
	// DocuSign users may have more than one account.  If AccountID is
	// not set then the user's default account will be used.
	AccountID string `json:"account_id,omitempty"`
	// if not nil, CacheFunc is called after a new token is created passing
	// the newly created Token and UserInfo.
	CacheFunc func(context.Context, oauth2.Token, UserInfo) `json:"cache_func,omitempty"`
	// Prompt indicates whether the authentication server will prompt
	// the user for re-authentication, even if they have an active login session.
	Prompt bool `json:"prompt,omitempty"`
	// List of the end-user’s preferred languages, represented as a
	// space-separated list of RFC5646 language tag values ordered by preference.
	// Note: can no longer find in docusign documentation.
	UIlocales []string `json:"u_ilocales,omitempty"`
	// Set to true to obtain an extended lifetime token (i.e. contains refresh token)
	ExtendedLifetime bool `json:"extended_lifetime,omitempty"`
	// Use developer sandbox
	IsDemo bool `json:"is_demo,omitempty"`
	// determines client used for oauth2 token calls.  If
	// nil, ctxclient.Default will be used.
	HTTPClientFunc ctxclient.Func
}

// codeGrantConfig creates an oauth2 config for refreshing
// and generating a token.
func (c *OAuth2Config) codeGrantConfig() *oauth2.Config {
	scopes := []string{"signature"}
	if c.ExtendedLifetime {
		scopes = []string{"signature", "extended"}
	}
	return &oauth2.Config{
		RedirectURL:    c.RedirURL,
		ClientID:       c.IntegratorKey,
		ClientSecret:   c.Secret,
		Scopes:         scopes,
		Endpoint:       demoFlag(c.IsDemo).endpoint(),
		HTTPClientFunc: c.HTTPClientFunc,
	}
}

// AuthURL returns a URL to DocuSign's OAuth 2.0 consent page with
// all appropriate query parmeters for starting 3-legged OAuth2Flow.
//
// State is a token to protect the user from CSRF attacks. You must
// always provide a non-zero string and validate that it matches the
// the state query parameter on your redirect callback.
func (c *OAuth2Config) AuthURL(state string) string {
	cfg := c.codeGrantConfig() // client not needed for this action
	opts := make([]oauth2.AuthCodeOption, 0)
	if c.Prompt {
		opts = append(opts, oauth2.SetAuthURLParam("prompt", "login"))
	}
	if len(c.UIlocales) > 0 {
		opts = append(opts, oauth2.SetAuthURLParam("ui_locales", strings.Join(c.UIlocales, " ")))
	}
	// https://developers.docusign.com/esign-rest-api/guides/authentication/oauth2-code-grant#step-1-request-the-authorization-code
	// DocuSign insists on Path escape for url (i.e. %20 not + for spaces)
	return replacePlus(cfg.AuthCodeURL(state, opts...))
}

// Exchange converts an authorization code into a token.
//
// It is used after a resource provider redirects the user back
// to the Redirect URI (the URL obtained from AuthCodeURL).
//
// The code will be in the *http.Request.FormValue("code"). Before
// calling Exchange, be sure to validate FormValue("state").
func (c *OAuth2Config) Exchange(ctx context.Context, code string) (*OAuth2Credential, error) {
	cfg := c.codeGrantConfig()
	// oauth2 exchange
	tk, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	u, err := getUserInfoForToken(ctx, demoFlag(c.IsDemo).userInfoPath(), tk, cfg.HTTPClientFunc)
	if err != nil {
		return nil, err
	}
	if c.CacheFunc != nil {
		c.CacheFunc(ctx, *tk, *u)
	}
	// create credential
	return c.Credential(tk, u)

}

func (c *OAuth2Config) refresher() func(context.Context, *oauth2.Token) (*oauth2.Token, error) {
	cfg := c.codeGrantConfig()
	return func(ctx context.Context, tk *oauth2.Token) (*oauth2.Token, error) {
		if tk == nil || tk.RefreshToken == "" {
			return nil, errors.New("codeGrantRefresher: empty refresh token")
		}
		return cfg.RefreshToken(ctx, tk.RefreshToken)
	}
}

// Credential returns an *OAuth2Credential using the passed oauth2.Token
// as the starting authorization token.
func (c *OAuth2Config) Credential(tk *oauth2.Token, u *UserInfo) (*OAuth2Credential, error) {
	if c == nil {
		return nil, errors.New("nil configuration")
	}
	if tk == nil {
		return nil, errors.New("token may not be nil")
	}
	tokenIsValid := tk.Valid()
	if !tokenIsValid && tk.RefreshToken == "" {
		return nil, errors.New("empty refresh token")
	}
	var accountID = c.AccountID
	var baseURI *url.URL
	var err error
	if tokenIsValid && u != nil {
		if accountID, baseURI, err = u.getAccountID(accountID); err != nil {
			return nil, err
		}
	}
	return &OAuth2Credential{
		accountID:   c.AccountID,
		baseURI:     baseURI,
		cachedToken: tk,
		refresher:   c.refresher(),
		cacheFunc:   c.CacheFunc,
		isDemo:      demoFlag(c.IsDemo),
		userInfo:    u,
		Func:        c.HTTPClientFunc,
	}, nil
}

// JWTConfig is used to create an OAuth2Credential based upon DocuSign's
// Service Integration Authentication.
//
// See https://developers.docusign.com/esign-rest-api/guides/authentication/oauth2-jsonwebtoken
type JWTConfig struct {
	// see https://developers.docusign.com/esign-rest-api/guides/authentication/oauth2-jsonwebtoken#prerequisites
	IntegratorKey string `json:"integrator_key,omitempty"`
	// Use developer sandbox
	IsDemo bool `json:"is_demo,omitempty"`
	// PEM encoding of an RSA Private Key.
	// see https://developers.docusign.com/esign-rest-api/guides/authentication/oauth2-jsonwebtoken#prerequisites
	// for how to create RSA keys to the application.
	PrivateKey string `json:"private_key,omitempty"`
	KeyPairID  string `json:"key_pair_id,omitempty"`
	// DocuSign users may have more than one account.  If AccountID is
	// not set then the user's default account will be used.
	AccountID string `json:"account_id,omitempty"`
	// (optional)Expires specifies how long the token will be valid. DocuSign
	// limits this to 1 hour.  1 hour is assumed if left empty.
	Expiration *jwt.ExpirationSetting `json:"expires,omitempty"`
	// if not nil, CacheFunc is called after a new token is created passing
	// the newly created Token and UserInfo.
	CacheFunc func(context.Context, oauth2.Token, UserInfo) `json:"cache_func,omitempty"`
	// HTTPClientFunc determines client used for oauth2 token calls.  If
	// nil, ctxclient.DefaultClient will be used.
	HTTPClientFunc ctxclient.Func
}

// UserConsentURL creates a url allowing a user to consent to impersonation
// https://developers.docusign.com/esign-rest-api/guides/authentication/oauth2-jsonwebtoken#step-1-request-the-authorization-code
func (c *JWTConfig) UserConsentURL(redirectURL string) string {
	q := make(url.Values)
	q.Set("response_type", "code")
	q.Set("scope", "signature impersonation")
	q.Set("client_id", c.IntegratorKey)
	q.Set("redirect_uri", redirectURL)
	// docusign insists upon %20 not + in scope definition
	return demoFlag(c.IsDemo).endpoint().AuthURL + "?" + replacePlus(q.Encode())
}

func (c *JWTConfig) jwtRefresher(apiUserName string, signer jws.Signer) func(ctx context.Context, tk *oauth2.Token) (*oauth2.Token, error) {
	cfg := &jwt.Config{
		Issuer:         c.IntegratorKey,
		Signer:         signer,
		Subject:        apiUserName,
		Expiration:     c.Expiration,
		Scopes:         []string{"signature", "impersonation"},
		Audience:       demoFlag(c.IsDemo).tokenURI(),
		TokenURL:       demoFlag(c.IsDemo).endpoint().TokenURL,
		HTTPClientFunc: c.HTTPClientFunc,
	}
	return func(ctx context.Context, tk *oauth2.Token) (*oauth2.Token, error) {
		return cfg.Token(ctx)
	}
}

// Credential returns an *OAuth2Credential.  The passed token will be refreshed
// as needed.
func (c *JWTConfig) Credential(apiUserName string, token *oauth2.Token, u *UserInfo) (*OAuth2Credential, error) {
	signer, err := jws.RS256FromPEM([]byte(c.PrivateKey), c.KeyPairID)
	if err != nil {
		return nil, err
	}
	return &OAuth2Credential{
		accountID:   c.AccountID,
		cachedToken: token,
		refresher:   c.jwtRefresher(apiUserName, signer),
		cacheFunc:   c.CacheFunc,
		isDemo:      demoFlag(c.IsDemo),
		userInfo:    u,
		Func:        c.HTTPClientFunc,
	}, nil
}

// OAuth2Credential authorizes op requests via DocuSign's oauth2 protocol.
type OAuth2Credential struct {
	accountID   string
	baseURI     *url.URL // baseURI for ops not token
	cachedToken *oauth2.Token
	refresher   func(context.Context, *oauth2.Token) (*oauth2.Token, error)
	cacheFunc   func(context.Context, oauth2.Token, UserInfo)
	userInfo    *UserInfo
	isDemo      demoFlag
	mu          sync.Mutex
	ctxclient.Func
}

// Authorize set the authorization header and completes request's url
// with the users's baseURI and account id.
func (cred *OAuth2Credential) Authorize(ctx context.Context, req *http.Request) error {
	t, err := cred.Token(ctx)
	if err != nil {
		return err
	}
	t.SetAuthHeader(req)
	// finalize url
	ResolveDSURL(req.URL, cred.baseURI.Host, cred.accountID)
	return nil
}

// AuthDo set the authorization header and completes request's url
// with the users's baseURI and account id before sending the request
func (cred *OAuth2Credential) AuthDo(ctx context.Context, req *http.Request) (*http.Response, error) {

	t, err := cred.Token(ctx)
	if err != nil {
		if req.Body != nil {
			req.Body.Close()
		}
		return nil, err
	}
	t.SetAuthHeader(req)
	// finalize url
	ResolveDSURL(req.URL, cred.baseURI.Host, cred.accountID)
	res, err := cred.Func.Do(ctx, req)
	return res, toResponseError(err)
}

// WithAccountID creates a copy the current credential with a new accountID.  An empty
// accountID will use the user's default account. If the accountID is invalid for the user
// an error will occur when authorizing and operation.  Check for a valid account using
// *OAuth2Credential.UserInfo(ctx).
func (cred *OAuth2Credential) WithAccountID(accountID string) *OAuth2Credential {
	if cred == nil {
		return nil
	}

	cred.mu.Lock()
	defer cred.mu.Unlock()

	return &OAuth2Credential{
		accountID:   accountID,
		baseURI:     cred.baseURI,
		cachedToken: cred.cachedToken,
		refresher:   cred.refresher,
		cacheFunc:   cred.cacheFunc,
		userInfo:    cred.userInfo,
		isDemo:      cred.isDemo,
		Func:        cred.Func,
	}

}

// UserInfo returns user data returned from the /oauth/userinfo ednpoint.
// See https://developers.docusign.com/esign-rest-api/guides/authentication/user-info-endpoints
func (cred *OAuth2Credential) UserInfo(ctx context.Context) (*UserInfo, error) {
	cred.mu.Lock()
	if cred.userInfo == nil {
		cred.mu.Unlock() // release lock b/c Token locks
		// Token assignes userInfo if nil...
		if _, err := cred.Token(ctx); err != nil {
			return nil, err
		}
		cred.mu.Lock()
	}
	u := *cred.userInfo
	cred.mu.Unlock()
	return &u, nil
}

// Token checks where the cachedToken is valid.  If not it attempts to obtain
// a new token via the refresher.  Next accountID and baseURI are updated if
// blank, (see https://developers.docusign.com/esign-rest-api/guides/authentication/user-info-endpoints).
func (cred *OAuth2Credential) Token(ctx context.Context) (*oauth2.Token, error) {
	if ctx == nil {
		return nil, errors.New("context may not be nil")
	}
	if cred == nil {
		return nil, errors.New("nil credential")
	}
	var isNewToken bool
	var err error
	// lock credential during validation and possible update
	cred.mu.Lock()
	defer cred.mu.Unlock()

	if !cred.cachedToken.Valid() {
		if cred.cachedToken, err = cred.refresher(ctx, cred.cachedToken); err != nil {
			return nil, err
		}
		isNewToken = true
	}
	// check for userInfo and set AccountID and BaseURI to resolve op urls
	if cred.userInfo == nil {
		cred.userInfo, err = getUserInfoForToken(ctx, cred.isDemo.userInfoPath(), cred.cachedToken, cred.Func)
		if err != nil {
			return nil, err
		}
	}
	if cred.baseURI == nil || cred.accountID == "" { // values may be blank if loading userinfo from cache
		if cred.accountID, cred.baseURI, err = cred.userInfo.getAccountID(cred.accountID); err != nil {
			return nil, err
		}
	}
	if isNewToken && cred.cacheFunc != nil {
		cred.cacheFunc(ctx, *cred.cachedToken, *cred.userInfo)
	}
	return cred.cachedToken, nil
}

func getUserInfoForToken(ctx context.Context, path string, tk *oauth2.Token, f ctxclient.Func) (*UserInfo, error) {
	// needed to use token credential due to different host and path parameters for op
	var u UserInfo
	err := (&Op{
		Credential: &tokenCredential{tk, f},
		Method:     "GET",
		Path:       path,
	}).Do(ctx, &u)
	return &u, err

}

// tokenCredential provides authorization for userInfo ops.
type tokenCredential struct {
	*oauth2.Token
	ctxclient.Func
}

func (t *tokenCredential) AuthDo(ctx context.Context, req *http.Request) (*http.Response, error) {
	t.Token.SetAuthHeader(req)
	res, err := t.Func.Do(ctx, req)
	return res, toResponseError(err)
}

func toResponseError(err error) error {
	if nsErr, ok := err.(*ctxclient.NotSuccess); ok {
		return NewResponseError(nsErr.Body, nsErr.StatusCode)
	}
	return err
}

// UserInfo provides all account info for a specific user.  Data from
// the /oauth/userinfo op is unmarshaled into this struct.
type UserInfo struct {
	APIUsername string            `json:"sub"`
	Accounts    []UserInfoAccount `json:"accounts"`
	Name        string            `json:"name"`
	GivenName   string            `json:"given_name"`
	FamilyName  string            `json:"family_name"`
	Email       string            `json:"email"`
}

// UserInfoAccount contains the account information for a UserInfo
type UserInfoAccount struct {
	AccountID   string `json:"account_id"`
	IsDefault   bool   `json:"is_default"`
	AccountName string `json:"account_name"`
	BaseURI     string `json:"base_uri"`
}

// getAccountID returns the user AccountID and BaseURI
// for the given id.  If id is blank, return from the
// default account.
func (u *UserInfo) getAccountID(id string) (string, *url.URL, error) {
	if u == nil {
		return "", nil, errors.New("userInfo is nil")
	}
	for _, a := range u.Accounts {
		if (id == "" && a.IsDefault) || id == a.AccountID {
			ux, err := url.Parse(a.BaseURI)
			return a.AccountID, ux, err
		}
	}

	return "", nil, fmt.Errorf("no account %s for %s", id, u.Email)
}

var expReplacePlusInScope = regexp.MustCompile(`[\?&]scope=([^\+&]*\+)+`)

func replacePlus(s string) string {
	return expReplacePlusInScope.ReplaceAllStringFunc(s, func(rstr string) string {
		return strings.Replace(rstr, "+", "%20", -1)
	})
}
