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
	// endpoints
	return oauth2.Endpoint{
		AuthURL:  "https://" + df.tokenURI() + "/oauth/auth",
		TokenURL: "https://" + df.tokenURI() + "/oauth/token",
	}
}

func (df demoFlag) tokenURI() string {
	if df {
		return "account-d.docusign.com"
	}
	return "account.docusign.com"
}

func (df demoFlag) getUserInfoForToken(ctx context.Context, f ctxclient.Func, tk *oauth2.Token) (*UserInfo, error) {
	// needed to use token credential due to different host and path parameters for op
	var u *UserInfo
	err := (&Op{
		Credential: &tokenCredential{tk, f},
		Method:     "GET",
		Path:       "https://" + df.tokenURI() + "/oauth/userinfo",
	}).Do(ctx, &u)
	return u, err
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
	HTTPClientFunc ctxclient.Func `json:"-"`
}

// codeGrantConfig creates an oauth2 config for refreshing
// and generating a token.
func (c *OAuth2Config) codeGrantConfig(scopes ...string) *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:    c.RedirURL,
		ClientID:       c.IntegratorKey,
		ClientSecret:   c.Secret,
		Scopes:         scopes,
		Endpoint:       demoFlag(c.IsDemo).endpoint(),
		HTTPClientFunc: c.HTTPClientFunc,
	}
}

func addUnique(scopes []string, scope string) []string {
	for _, val := range scopes {
		if val == scope {
			return scopes
		}
	}
	return append(scopes, scope)
}

// AuthURL returns a URL to DocuSign's OAuth 2.0 consent page with
// all appropriate query parmeters for starting 3-legged OAuth2Flow.
//
// If scopes are empty, {"signature"} is assumed.
//
// State is a token to protect the user from CSRF attacks. You must
// always provide a non-zero string and validate that it matches the
// the state query parameter on your redirect callback.
func (c *OAuth2Config) AuthURL(state string, scopes ...string) string {
	if len(scopes) == 0 {
		scopes = []string{"signature"}
	}
	if c.ExtendedLifetime {
		scopes = addUnique(scopes, "extended")
	}
	cfg := c.codeGrantConfig(scopes...)
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
	cfg := c.codeGrantConfig() // scopes are not passed in this step
	// oauth2 exchange
	tk, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	u, err := demoFlag(c.IsDemo).getUserInfoForToken(ctx, cfg.HTTPClientFunc, tk)
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
	cfg := c.codeGrantConfig() // scopes are not passed in this step
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
	// Options can specify how long the token will be valid. DocuSign
	// limits this to 1 hour.  1 hour is assumed if left empty.  Offsets
	// for expiring token may also be used.  Do not set FormValues or Custom Claims.
	Options *jwt.ConfigOptions `json:"expires,omitempty"`
	// if not nil, CacheFunc is called after a new token is created passing
	// the newly created Token and UserInfo.
	CacheFunc func(context.Context, oauth2.Token, UserInfo) `json:"-"`
	// HTTPClientFunc determines client used for oauth2 token calls.  If
	// nil, ctxclient.DefaultClient will be used.
	HTTPClientFunc ctxclient.Func `json:"-"`
}

// UserConsentURL creates a url allowing a user to consent to impersonation
// https://developers.docusign.com/esign-rest-api/guides/authentication/obtaining-consent#individual-consent
func (c *JWTConfig) UserConsentURL(redirectURL string, scopes ...string) string {
	scopeValue := "signature impersonation"
	if len(scopes) > 0 {
		scopeValue = strings.Join(addUnique(scopes, "impersonation"), " ")
	}
	// docusign insists upon %20 not + in scope definition
	return demoFlag(c.IsDemo).endpoint().AuthURL + "?" + replacePlus(url.Values{
		"response_type": {"code"},
		"scope":         {scopeValue},
		"client_id":     {c.IntegratorKey},
		"redirect_uri":  {redirectURL},
	}.Encode())
}

// ExternalAdminConsentURL creates a url for beginning external admin consent workflow. See
// https://developers.docusign.com/esign-rest-api/guides/authentication/obtaining-consent#admin-consent-for-external-applications
// for details.
//
// redirectURL is the URI to which DocuSign will redirect the browser after authorization has been granted by the extneral organization's
// admin.  The redirect URI must exactly match one of those pre-registered for the Integrator Key in your DocuSign account.
//
// authType may be either  code (Authorization Code grant) or token (implicit grant).
//
// state holds an optional value that is returned with the authorization code.
//
// prompt determines whether the user is prompted for re-authentication, even with an active login session.
//
// scopes permissions being requested for the application from each user in the organization.  Valid values are
//   signature — allows your application to create and send envelopes, and obtain links for starting signing sessions.
//   extended — issues your application a refresh token that can be used any number of times (Authorization Code flow only).
//   impersonation — allows your application to access a user’s account and act on their behalf via JWT authentication.
func (c *JWTConfig) ExternalAdminConsentURL(redirectURL, authType, state string, prompt bool, scopes ...string) (string, error) {
	if authType != "code" && authType != "token" {
		return "", fmt.Errorf("invalid authType %s, must be code or token", authType)
	}
	if len(scopes) == 0 {
		return "", fmt.Errorf("at least one scope must be specified")
	}
	v := url.Values{
		"scope":               {"openid"},
		"client_id":           {c.IntegratorKey},
		"response_type":       {authType},
		"redirect_uri":        {redirectURL},
		"admin_consent_scope": {strings.Join(scopes, " ")},
	}
	if state > "" {
		v.Set("state", state)
	}
	if prompt {
		v.Set("prompt", "login")
	}
	query := replacePlus(v.Encode())
	return "https://" + demoFlag(c.IsDemo).tokenURI() + "/oauth/auth?" + query, nil
}

// AdminConsentResponse is the response sent to the redirect url of and external admin
// consent
// https://developers.docusign.com/esign-rest-api/guides/authentication/obtaining-consent#admin-consent-for-external-applications
type AdminConsentResponse struct {
	Issuer    string   `json:"iss"`       // domain of integration key
	Audience  string   `json:"aud"`       // the integrator key (also known as client ID) of the application
	ExpiresAt int64    `json:"exp"`       // the datetime when the ID token will expire, in Unix epoch format
	IssuedAt  int64    `json:"iat"`       // the datetime when the ID token was issued, in Unix epoch format
	Subject   string   `json:"sub"`       //  user ID of the admin granting consent for the organization users
	SiteID    int64    `json:"siteid"`    // identifies the docusign server used.
	AuthTime  string   `json:"auth_time"` // The unix epoch datetime when the ID token was created.
	AMR       []string `json:"amr"`       // how the user authenticated
	COID      []string `json:"coid"`      //  list of organization IDs for the organizations whose admin has granted consent
}

func (c *JWTConfig) jwtRefresher(apiUserName string, signer jws.Signer, scopes ...string) func(ctx context.Context, tk *oauth2.Token) (*oauth2.Token, error) {
	if len(scopes) == 0 {
		scopes = []string{"signature", "impersonation"}
	} else {
		scopes = addUnique(scopes, "impersonation")
	}
	cfg := &jwt.Config{
		Issuer:         c.IntegratorKey,
		Signer:         signer,
		Subject:        apiUserName,
		Options:        c.Options,
		Scopes:         scopes,
		Audience:       demoFlag(c.IsDemo).tokenURI(),
		TokenURL:       demoFlag(c.IsDemo).endpoint().TokenURL,
		HTTPClientFunc: c.HTTPClientFunc,
	}
	return func(ctx context.Context, tk *oauth2.Token) (*oauth2.Token, error) {
		return cfg.Token(ctx)
	}
}

// Credential returns an *OAuth2Credential.  The passed token will be refreshed
// as needed.  If no scopes listed, signature is assumed.
func (c *JWTConfig) Credential(apiUserName string, token *oauth2.Token, u *UserInfo, scopes ...string) (*OAuth2Credential, error) {
	signer, err := jws.RS256FromPEM([]byte(c.PrivateKey), c.KeyPairID)
	if err != nil {
		return nil, err
	}
	return &OAuth2Credential{
		accountID:   c.AccountID,
		cachedToken: token,
		refresher:   c.jwtRefresher(apiUserName, signer, scopes...),
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

// AuthDo set the authorization header and completes request's url
// with the users's baseURI and account id before sending the request
func (cred *OAuth2Credential) AuthDo(ctx context.Context, req *http.Request, v *APIVersion) (*http.Response, error) {
	t, err := cred.Token(ctx)
	if err != nil {
		if req.Body != nil {
			req.Body.Close()
		}
		return nil, err
	}
	r2 := *req
	h := make(http.Header)
	for k, v := range req.Header {
		h[k] = v
	}
	r2.Header = h

	t.SetAuthHeader(&r2)
	// finalize url
	r2.URL = v.ResolveDSURL(req.URL, cred.baseURI.Host, cred.accountID)
	res, err := cred.Func.Do(ctx, &r2)
	return res, toResponseError(err)
}

// WithAccountID creates a copy the current credential with a new accountID.  An empty
// accountID indicates the user's default account. If the accountID is invalid for the user
// an error will occur when authorizing an operation.  Check for valid account using
// *OAuth2Credential.UserInfo(ctx).
func (cred *OAuth2Credential) WithAccountID(accountID string) *OAuth2Credential {
	if cred == nil {
		return nil
	}
	cred.mu.Lock()
	c := OAuth2Credential{
		accountID:   accountID,
		baseURI:     nil,
		cachedToken: cred.cachedToken,
		refresher:   cred.refresher,
		cacheFunc:   cred.cacheFunc,
		userInfo:    cred.userInfo,
		isDemo:      cred.isDemo,
		Func:        cred.Func,
	}
	cred.mu.Unlock()
	return &c
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
	var updateCache bool
	var err error
	// lock credential during validation and possible update
	cred.mu.Lock()
	defer cred.mu.Unlock()

	if !cred.cachedToken.Valid() {
		if cred.refresher == nil {
			return nil, errors.New("no refresher function for invalid/expired token")
		}
		if cred.cachedToken, err = cred.refresher(ctx, cred.cachedToken); err != nil {
			return nil, err
		}
		updateCache = (cred.cacheFunc != nil)
	}
	// check for userInfo and set AccountID and BaseURI to resolve op urls
	if cred.userInfo == nil {
		cred.userInfo, err = cred.isDemo.getUserInfoForToken(ctx, cred.Func, cred.cachedToken)
		if err != nil {
			return nil, err
		}
		updateCache = (cred.cacheFunc != nil)
	}
	if cred.baseURI == nil || cred.accountID == "" { // values may be blank if loading userinfo from cache
		if cred.accountID, cred.baseURI, err = cred.userInfo.getAccountID(cred.accountID); err != nil {
			return nil, err
		}
	}
	if updateCache {
		cred.cacheFunc(ctx, *cred.cachedToken, *cred.userInfo)
	}
	return cred.cachedToken, nil
}

// SetClientFunc safely replaces the ctxclient.Func for the credential
func (cred *OAuth2Credential) SetClientFunc(f ctxclient.Func) *OAuth2Credential {
	cred.mu.Lock()
	cred.Func = f
	cred.mu.Unlock()
	return cred
}

// SetCacheFunc safely replaces the caching function for the credential
func (cred *OAuth2Credential) SetCacheFunc(f func(context.Context, oauth2.Token, UserInfo)) *OAuth2Credential {
	cred.mu.Lock()
	cred.cacheFunc = f
	cred.mu.Unlock()
	return cred
}

// TokenCredential create a static credential without refresh capabilities.  When
// the token expires, ops will receive a 401 error,
func TokenCredential(accessToken string, isDemo bool) *OAuth2Credential {
	return &OAuth2Credential{
		cachedToken: &oauth2.Token{
			AccessToken: accessToken,
		},
		isDemo: demoFlag(isDemo),
	}
}

// tokenCredential provides authorization for userInfo ops.
type tokenCredential struct {
	*oauth2.Token
	ctxclient.Func
}

func (t *tokenCredential) AuthDo(ctx context.Context, req *http.Request, v *APIVersion) (*http.Response, error) {
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

var expReplacePlusInScope = regexp.MustCompile(`[\?&_]scope=([^\+&]*\+)+`)

func replacePlus(s string) string {
	return expReplacePlusInScope.ReplaceAllStringFunc(s, func(rstr string) string {
		return strings.Replace(rstr, "+", "%20", -1)
	})
}
