// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package oauth2 provides support for making
// OAuth2 authorized and authenticated HTTP requests.
// It can additionally grant authorization with Bearer JWT.
package oauth2 // import "github.com/jfcote87/oauth2"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jfcote87/ctxclient"
)

// NoContext is the default background context
var NoContext = context.Background()

// Config describes a typical 3-legged OAuth2 flow, with both the
// client application information and the server's endpoint URLs.
// For the client credentials 2-legged OAuth2 flow, see the clientcredentials
// package (https://github.com/jfcote87/oauth2/clientcredentials).
type Config struct {
	// ClientID is the application's ID.
	ClientID string

	// ClientSecret is the application's secret.
	ClientSecret string

	// Endpoint contains the resource server's token endpoint
	// URLs. These are constants specific to each server
	// implementation so please refer to vendor documentation.
	// I have no intentions of maintaining a list of public endpoints.
	Endpoint Endpoint

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string

	// Scope specifies optional requested permissions.
	Scopes []string

	// ExpiryDelta determines how many seconds sooner a token should
	// expire than the retrieved expires_in setting.
	ExpiryDelta int64

	// HTTPClientFunc may be set to determine the *http.Client
	// used for Exchange and Refresh calls.  If not set, the default
	// for appengine applications is created via the urlfetch.Client
	// function.  Otherwise the http.DefaultClient is assumed.
	HTTPClientFunc ctxclient.Func
}

// Endpoint contains the OAuth 2.0 provider's authorization and token
// endpoint URLs.
type Endpoint struct {
	AuthURL  string
	TokenURL string
	// Set to true if server requires ClientID and Secret
	// in body rather than Basic Authentication.
	IDSecretInBody bool
}

var (
	// AccessTypeOnline and AccessTypeOffline are options passed
	// to the Options.AuthCodeURL method. They modify the
	// "access_type" field that gets sent in the URL returned by
	// AuthCodeURL.
	//
	// Online is the default if neither is specified.
	AccessTypeOnline = SetAuthURLParam("access_type", "online")
	// AccessTypeOffline is used an application needs to refresh
	// access tokens when the user is not present. This will
	// result in your application obtaining a refresh token the
	// first time your application exchanges an authorization
	// code for a user.
	AccessTypeOffline = SetAuthURLParam("access_type", "offline")

	// ApprovalForce forces the users to view the consent dialog
	// and confirm the permissions request at the URL returned
	// from AuthCodeURL, even if they've already done so.
	ApprovalForce = SetAuthURLParam("prompt", "consent")
)

// An AuthCodeOption is passed to Config.AuthCodeURL.
type AuthCodeOption interface {
	setValue(url.Values)
}

type setParam struct{ k, v string }

func (p setParam) setValue(m url.Values) { m.Set(p.k, p.v) }

// SetAuthURLParam builds an AuthCodeOption which passes key/value parameters
// to a provider's authorization endpoint.
func SetAuthURLParam(key, value string) AuthCodeOption {
	return setParam{key, value}
}

// AuthCodeURL returns a URL to OAuth 2.0 provider's consent page
// that asks for permissions for the required scopes explicitly.
//
// State is a token to protect the user from CSRF attacks. You must
// always provide a non-zero string and validate that it matches the
// the state query parameter on your redirect callback.
// See http://tools.ietf.org/html/rfc6749#section-10.12 for more info.
//
// Opts may include AccessTypeOnline or AccessTypeOffline, as well
// as ApprovalForce.
func (c *Config) AuthCodeURL(state string, opts ...AuthCodeOption) string {
	var buf = &bytes.Buffer{}
	buf.WriteString(c.Endpoint.AuthURL)
	v := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientID},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	if state != "" {
		// TODO(light): Docs say never to omit state; don't allow empty.
		v.Set("state", state)
	}
	for _, opt := range opts {
		opt.setValue(v)
	}
	if strings.Contains(c.Endpoint.AuthURL, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

// FromOptions returns a TokenSource that retrieves tokens using the
// parameters defined in opts.  Used by clientcredentials package.
func (c *Config) FromOptions(opts ...AuthCodeOption) TokenSource {
	v := make(url.Values)
	for _, o := range opts {
		o.setValue(v)
	}
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	return tokenRefreshFunc(func(ctx context.Context) (*Token, error) {
		return c.retrieveToken(ctx, v)
	})
}

// Exchange converts an authorization code into a token.
//
// It is used after a resource provider redirects the user back
// to the Redirect URI (the URL obtained from AuthCodeURL).
//
// The HTTP client to use is derived from the context.
// If a client is not provided via the context, http.DefaultClient is used.
//
// The code will be in the *http.Request.FormValue("code"). Before
// calling Exchange, be sure to validate FormValue("state").
func (c *Config) Exchange(ctx context.Context, code string, opts ...AuthCodeOption) (*Token, error) {
	v := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	if c.RedirectURL != "" {
		v.Set("redirect_uri", c.RedirectURL)
	}
	for _, opt := range opts {
		opt.setValue(v)
	}
	return c.retrieveToken(ctx, v)
}

// A TokenSource is anything that can return a token.  This package
// uses the passed context to determine/construct the appropriate
// http.Client for retrieving tokens and to allow for cancellation
// and timeouts of the Token request.
type TokenSource interface {
	// Token returns a token or an error.
	// Token must be safe for concurrent use by multiple goroutines.
	// The returned Token must not be modified.
	Token(context.Context) (*Token, error)
}

// RefreshToken retrieves a Token.  A developer may
// use this to construct custom caching TokenSources.
func (c *Config) RefreshToken(ctx context.Context, refreshToken string, opts ...AuthCodeOption) (*Token, error) {
	v := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	for _, o := range opts {
		o.setValue(v)
	}
	return c.retrieveToken(ctx, v)
}

type tokenRefreshFunc func(context.Context) (*Token, error)

func (trf tokenRefreshFunc) Token(ctx context.Context) (*Token, error) {
	if trf == nil {
		return nil, errors.New("nil tokenRefreshFunc")
	}
	return trf(ctx)
}

// cachedToken is a TokenSource that holds a single token in memory
// and validates its expiry before each call to retrieve it with
// Token. If it's expired, it will be auto-refreshed using the
// TokenSourc new.
type cachedToken struct {
	new TokenSource // called when t is expired.

	mu sync.Mutex // guards t
	t  *Token
}

// Token returns the current token if it's still valid, else will
// refresh the current token and return the new one.
func (s *cachedToken) Token(ctx context.Context) (*Token, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.t.Valid() {
		return s.t, nil
	}
	t, err := s.new.Token(ctx)
	if err != nil {
		return nil, err
	}
	s.t = t
	return t, nil
}

// ReuseTokenSource returns a TokenSource which repeatedly returns
// the same token as long as it's valid, starting with t.
// When its cached token is invalid, a new token is obtained from src.
//
// ReuseTokenSource is typically used to reuse tokens from a cache
// (such as a file on disk) between runs of a program, rather than
// obtaining new tokens unnecessarily.
//
// The initial token t may be nil, in which case the TokenSource is
// wrapped in a caching version if it isn't one already. This also
// means it's always safe to wrap ReuseTokenSource around any other
// TokenSource without adverse effects.
//
// ReuseTokenSource uses a mutex to allow the returned TokenSource
// to be used concurrently.
func ReuseTokenSource(t *Token, src TokenSource) TokenSource {
	// Don't wrap a reuseTokenSource in itself. That would work,
	// but cause an unnecessary number of mutex operations.
	// Just build the equivalent one.
	if rt, ok := src.(*cachedToken); ok {
		if t == nil || t == rt.t {
			// Just use it directly.
			return rt
		}
		src = rt.new
	}
	return &cachedToken{
		t:   t,
		new: src,
	}
}

// TokenSource returns a TokenSource that returns t until t expires,
// automatically refreshing it as necessary.  This tokensource is
// safe for use with different contexts.
// opts
func (c *Config) TokenSource(t *Token, opts ...AuthCodeOption) TokenSource {
	var refreshToken string
	if t != nil {
		refreshToken = t.RefreshToken
	}
	tfr := func(ctx context.Context) (*Token, error) {
		if c == nil {
			return nil, errors.New("nil config")
		}
		if refreshToken == "" {
			return nil, errors.New("empty refresh token")
		}
		tk, err := c.RefreshToken(ctx, refreshToken, opts...)
		if err != nil {
			return nil, err
		}
		if tk != nil && tk.RefreshToken > "" && tk.RefreshToken != refreshToken {
			refreshToken = tk.RefreshToken
		}
		return tk, nil
	}
	// ReuseTokenSource's mutex will protect refreshToken during concurrent operations
	return ReuseTokenSource(t, tokenRefreshFunc(tfr))
}

// StaticTokenSource returns a TokenSource that always returns the same token.
// Because the provided token t is never refreshed, StaticTokenSource is only
// useful for tokens that never expire.
func StaticTokenSource(t *Token) TokenSource {
	return staticTokenSource{t}
}

// staticTokenSource is a TokenSource that always returns the same Token.
type staticTokenSource struct {
	t *Token
}

func (s staticTokenSource) Token(ctx context.Context) (*Token, error) {
	return s.t, nil
}

// retrieveToken calls RetrieveToken taking notice of the Endopoint's IDSecretInBody flag
func (c *Config) retrieveToken(ctx context.Context, v url.Values) (*Token, error) {
	clientID, clientSecret := c.ClientID, c.ClientSecret
	if c.Endpoint.IDSecretInBody {
		v.Set("client_id", c.ClientID)
		v.Set("client_secret", c.ClientSecret)
		clientID, clientSecret = "", ""
	}
	req, err := http.NewRequest("POST", c.Endpoint.TokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if clientID > "" {
		req.SetBasicAuth(url.QueryEscape(clientID), url.QueryEscape(clientSecret))
	}
	var body []byte
	r, err := c.HTTPClientFunc.Do(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	body, err = ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: token body read error:  %v", err)
	}

	//var token *Token
	mappedValues := make(map[string]interface{})
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}
		for k := range vals {
			mappedValues[k] = vals.Get(k)
		}
	default:
		if err := json.NewDecoder(bytes.NewReader(body)).Decode(&mappedValues); err != nil {
			return nil, err
		}
	}
	// handle strings from x-www-form-urlencoded and for PayPayl
	if s, ok := mappedValues["expires_in"].(string); ok {
		if mappedValues["expires_in"], err = strconv.ParseInt(s, 10, 64); err != nil {
			return nil, fmt.Errorf("oauth2: unable to parse expires_in %v", err)
		}
	}
	return TokenFromMap(mappedValues, c.delta())
}

func (c *Config) delta() time.Duration {
	if c.ExpiryDelta > 0 {
		return time.Duration(c.ExpiryDelta) * time.Second
	}
	return time.Duration(DefaultExpiryDelta) * time.Second
}

// Client returns an HTTP client using the provided token.
// The token will auto-refresh as necessary. The underlying
// Client returns an HTTP client using the provided token.
// HTTP transport will be obtained using Config.HTTPClient.
// The returned client and its Transport should not be modified.
//
// The returned client uses the request's context to handle
// timeouts and cancellations.  It may be used concurrently
// as the token refresh is protected
func (c *Config) Client(t *Token) *http.Client {
	if c == nil {
		return nil
	}
	return &http.Client{
		Transport: &Transport{
			Source: c.TokenSource(t),
			Func:   c.HTTPClientFunc,
		},
	}
}
