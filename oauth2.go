// Copyright 2017 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esign

// oauth2.go contains definitions for Docusign's oauth2
// authorization scheme.  See the legacy package for
// the previous authorization schemes.

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"github.com/jfcote87/ctxclient"
	"github.com/jfcote87/oauth2"
	"github.com/jfcote87/oauth2/jwt"
)

type demoFlag bool

func (df demoFlag) endpoint() oauth2.Endpoint {
	if df {
		// Demo endpoints
		return oauth2.Endpoint{
			AuthURL:  "https://account.docusign-d.com/oauth/auth",
			TokenURL: "https://account.docusign-d.com/oauth/token",
		}
	}
	// Production endpoints
	return oauth2.Endpoint{
		AuthURL:  "https://account.docusign.com/oauth/auth",
		TokenURL: "https://account.docusign.com/oauth/token",
	}
}

func (df demoFlag) tokenURL() string {
	if df {
		return "https://account.docusign-d.com/oauth/token"
	}
	return "https://account.docusign.com/oauth/token"
}

func (df demoFlag) userInfoPath() string {
	if df {
		return "https://account.docusign-d.com"
	}
	return "https://account.docusign.com"
}

// Oauth2Config allows for 3-legged oauth via a code grant mechanism
// see https://docs.docusign.com/esign/guide/authentication/oa2_auth_code.html
type Oauth2Config struct {
	// see https://docs.docusign.com/esign/guide/authentication/integrator_key.html
	IntegratorKey string
	// Secret generated when setting up integration in DocuSign. Leave blank for
	// implicit grant.
	Secret string
	// The redirect URI must exactly match one of those pre-registered for the
	// integrator key. This determines where to redirect the user after they
	// successfully authenticate.
	RedirURL string
	// DocuSign users may have more than one account.  If AccountID is
	// not set then the user's default account will be used.
	AccountID string
	// CacheFunc is called after a new token is created.  The function
	// will receive the new token and the associated UserInfo
	CacheFunc func(context.Context, oauth2.Token, *UserInfo)
	// Prompt indicates whether the authentication server will prompt
	// the user for re-authentication, even if they have an active login session.
	Prompt bool
	// List of the end-user’s preferred languages, represented as a
	// space-separated list of RFC5646 language tag values ordered by preference.
	UIlocales []string
	// Set to true to obtain an extended lifetime token
	ExtendedLifetime bool
	// Use developer sandbox
	IsDemo bool
	// determines client used for oauth2 token calls.  If
	// nil, ctxclient.Default will be used.
	ctxclient.Func
}

// codeGrantConfig creates an oauth2 config for refreshing
// and generating a token.
func (c *Oauth2Config) codeGrantConfig() *oauth2.Config {
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
		HTTPClientFunc: c.Func,
	}
}

// AuthURL returns a URL to OAuth 2.0 provider's consent page
// that asks for permissions for the required scopes explicitly.
//
// State is a token to protect the user from CSRF attacks. You must
// always provide a non-zero string and validate that it matches the
// the state query parameter on your redirect callback.
func (c *Oauth2Config) AuthURL(state string) string {
	cfg := c.codeGrantConfig()
	opts := make([]oauth2.AuthCodeOption, 0)
	if c.Prompt {
		opts = append(opts, oauth2.SetAuthURLParam("prompt", "login"))
	}
	if len(c.UIlocales) > 0 {
		opts = append(opts, oauth2.SetAuthURLParam("ui_locales", strings.Join(c.UIlocales, " ")))
	}
	// https://docs.docusign.com/esign/guide/authentication/auth_server.html
	// Docusign insists on Path escape for url (i.e. %20 not + for spaces)
	return strings.Replace(cfg.AuthCodeURL(state, opts...), "+", "%20", -1)
}

// Exchange converts an authorization code into a token.
//
// It is used after a resource provider redirects the user back
// to the Redirect URI (the URL obtained from AuthCodeURL).
//
// The code will be in the *http.Request.FormValue("code"). Before
// calling Exchange, be sure to validate FormValue("state").
func (c *Oauth2Config) Exchange(ctx context.Context, code string) (*Oauth2Credential, error) {
	cfg := c.codeGrantConfig()
	tk, err := cfg.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return c.Credential(tk, nil)
}

func (c *Oauth2Config) refresher() func(context.Context, *oauth2.Token) (*oauth2.Token, error) {
	cfg := c.codeGrantConfig()
	return func(ctx context.Context, tk *oauth2.Token) (*oauth2.Token, error) {
		if tk == nil || tk.RefreshToken == "" {
			return nil, errors.New("codeGrantRefresher: empty refresh token")
		}
		return cfg.RefreshToken(ctx, tk.RefreshToken)
	}
}

// Credential returns an *Oauth2Credential using the passed oauth2.Token
// as the starting authorization token.
func (c *Oauth2Config) Credential(tk *oauth2.Token, u *UserInfo) (*Oauth2Credential, error) {
	if tk == nil {
		return nil, errors.New("token may not be nil")
	}
	return &Oauth2Credential{
		accountID:   c.AccountID,
		cachedToken: tk,
		refresher:   c.refresher(),
		cacheFunc:   c.CacheFunc,
		isDemo:      demoFlag(c.IsDemo),
		userInfo:    u,
		Func:        c.Func,
	}, nil
}

// JWTConfig is used to create an Oauth2Credential based upon DocuSign's
// Service Integration Authentication.
//
// See https://docs.docusign.com/esign/guide/authentication/oa2_jwt.html
type JWTConfig struct {
	// see https://docs.docusign.com/esign/guide/authentication/integrator_key.html
	IntegratorKey string
	// PEM encoding of an RSA Private Key.
	// see https://docs.docusign.com/esign/guide/authentication/integrator_key.html#rsakeys
	// for how to create a DocuSign private key.
	PrivateKey string
	KeyPairID  string
	// APIUsername may be found on the Edit User page of the docusign admin site.
	// Do not use the email address.
	APIUsername string
	// Expires optionally specifies how long the token is valid for. Docusign
	// limits this to 1 hour regardless if the duration is greater than 1 hour.
	Expires time.Duration
	// DocuSign users may have more than one account.  If AccountID is
	// not set then the user's default account will be used.
	AccountID string
	// CacheFunc is called after a new token is created.  The function
	// will receive the new token and the associate UserInfo
	CacheFunc func(context.Context, oauth2.Token, *UserInfo)
	// Use developer sandbox
	IsDemo bool
	// UserInfo may be set if available
	UserInfo *UserInfo
	// Func determines client used for oauth2 token calls.  If
	// nil, the default client for docusign calls will be used.
	ctxclient.Func
}

func (c *JWTConfig) jwtRefresher() func(ctx context.Context, tk *oauth2.Token) (*oauth2.Token, error) {
	cfg := &jwt.Config{
		Subject:        c.APIUsername,
		PrivateKey:     []byte(c.PrivateKey),
		PrivateKeyID:   c.KeyPairID,
		Email:          c.IntegratorKey,
		Expires:        c.Expires,
		Scopes:         []string{"signature", "impersonation"},
		TokenURL:       demoFlag(c.IsDemo).tokenURL(),
		HTTPClientFunc: c.Func,
	}
	return func(ctx context.Context, tk *oauth2.Token) (*oauth2.Token, error) {
		return cfg.Token(ctx)
	}
}

// Credential returns an *Oauth2Credential using the passed oauth2.Token
// as the starting authorization token.
func (c *JWTConfig) Credential(tk *oauth2.Token) (*Oauth2Credential, error) {
	if tk == nil {
		return nil, errors.New("token may not be nil")
	}
	return &Oauth2Credential{
		accountID:   c.AccountID,
		cachedToken: tk,
		refresher:   c.jwtRefresher(),
		cacheFunc:   c.CacheFunc,
		isDemo:      demoFlag(c.IsDemo),
		userInfo:    c.UserInfo,
		Func:        c.Func,
	}, nil
}

// Oauth2Credential authorizes call requests via DocuSign's oauth2
type Oauth2Credential struct {
	accountID   string
	baseURI     *url.URL // baseURI for calls not token
	cachedToken *oauth2.Token
	refresher   func(context.Context, *oauth2.Token) (*oauth2.Token, error)
	cacheFunc   func(context.Context, oauth2.Token, *UserInfo)
	userInfo    *UserInfo
	isDemo      demoFlag
	mu          sync.Mutex
	ctxclient.Func
}

// Authorize set the authorization header and completes request's url
// with the users's baseURI and account id.
func (cred *Oauth2Credential) Authorize(ctx context.Context, req *http.Request) error {
	t, err := cred.Token(ctx)
	if err != nil {
		return err
	}
	t.SetAuthHeader(req)
	// finalize url
	ResolveDSURL(req.URL, cred.baseURI.Host, cred.accountID)
	return nil
}

// UserInfo returns user data returned from the /oauth/userinfo ednpoint.
// See https://docs.docusign.com/esign/guide/authentication/userinfo.html
func (cred *Oauth2Credential) UserInfo(ctx context.Context) (*UserInfo, error) {
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
// blank, (see https://docs.docusign.com/esign/guide/authentication/userinfo.html).
// If a new token was created, the cache function is passed the new token, the baseURI
// and the accountID.
func (cred *Oauth2Credential) Token(ctx context.Context) (*oauth2.Token, error) {
	if ctx == nil {
		return nil, errors.New("context may not be nil")
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
	// check for userInfo and set AccountID and BaseURI to resolve call urls
	if cred.userInfo == nil {
		var u *UserInfo

		// needed to use new credential due to
		if err := (&Call{
			Credential: &tokenCredential{cred.cachedToken, cred.Func},
			Method:     "GET",
			Path:       cred.isDemo.userInfoPath(),
		}).Do(ctx, &u); err != nil {
			return nil, err
		}
		if cred.accountID, cred.baseURI, err = u.getAccountID(cred.accountID); err != nil {
			return nil, err
		}
		cred.userInfo = u

	} else if cred.baseURI == nil || cred.accountID == "" { // values may be blank if loading userinfo from cache
		if cred.accountID, cred.baseURI, err = cred.userInfo.getAccountID(cred.accountID); err != nil {
			return nil, err
		}
	}
	if isNewToken && cred.cacheFunc != nil {
		cred.cacheFunc(ctx, *cred.cachedToken, cred.userInfo)
	}
	return cred.cachedToken, nil
}

// tokenCredential provides authorization for userInfo calls.
type tokenCredential struct {
	*oauth2.Token
	ctxclient.Func
}

func (t *tokenCredential) Authorize(ctx context.Context, req *http.Request) error {
	t.Token.SetAuthHeader(req)
	return nil
}

// UserInfo provides all account info for a specific user.  Data from
// the /oauth/userinfo call is unmarshaled into this struct.
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
		return "", nil, errors.New("UserInfo is nil")
	}
	for _, a := range u.Accounts {
		if (id == "" && a.IsDefault) || id == a.AccountID {
			ux, err := url.Parse(a.BaseURI)
			return a.AccountID, ux, err
		}
	}

	return "", nil, fmt.Errorf("no account %s for %s", id, u.Email)
}
