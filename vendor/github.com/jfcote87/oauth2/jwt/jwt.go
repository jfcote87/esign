// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jwt implements the OAuth 2.0 JSON Web Token flow, commonly
// known as "two-legged OAuth 2.0".
//
// See: https://tools.ietf.org/html/draft-ietf-oauth-jwt-bearer-12
package jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/jfcote87/ctxclient"
	"github.com/jfcote87/oauth2"
	"github.com/jfcote87/oauth2/jws"
)

var (
	defaultGrantType = "urn:ietf:params:oauth:grant-type:jwt-bearer"
)

// Config is the configuration for using JWT to fetch tokens,
// commonly known as "two-legged OAuth 2.0".
type Config struct {
	// Signer is the func used to sign the JWT header and payload
	Signer jws.Signer

	// Issuer is the OAuth client identifier used when communicating with
	// the configured OAuth provider.
	Issuer string `json:"email,omitempty"`

	// Subject is the optional user to impersonate.
	Subject string `json:"subject,omitempty"`

	// TokenURL is the endpoint required to complete the 2-legged JWT flow.
	TokenURL string `json:"token_url,omitempty"`

	// Audience fills the claimset's aud parameter.  For Google Client API
	// this should be set to the TokenURL value.
	Audience string `json:"audience,omitempty"`

	// Scopes optionally specifies a list of requested permission scopes
	// which will be included as private claims named scopes in the
	// JWT claimset payload.
	Scopes []string `json:"scopes,omitempty"`

	// additional options for creating and sending claimset payload
	Options *ConfigOptions `json:"options,omitempty"`

	// HTTPClientFunc (Optional) specifies a function specifiying
	// the *http.Client used on Token calls to the oauth2
	// server.
	HTTPClientFunc ctxclient.Func `json:"-"`
}

// TokenSource returns a JWT TokenSource using the configuration
// in c and the HTTP client from the provided context.
func (c *Config) TokenSource(t *oauth2.Token) oauth2.TokenSource {
	return oauth2.ReuseTokenSource(t, c)
}

// Client returns an HTTP client wrapping the context's
// HTTP transport and adding Authorization headers with tokens
// obtained from c.
//
// The returned client and its Transport should not be modified.
func (c *Config) Client(t *oauth2.Token) (*http.Client, error) {
	return oauth2.Client(c.TokenSource(t), c.HTTPClientFunc), nil
}

// payload returns the body of a token request
func (c *Config) payload() (url.Values, error) {
	privateClaims := make(map[string]interface{})
	for k, v := range c.Options.getPrivateClaims() {
		privateClaims[k] = v
	}
	if len(c.Scopes) > 0 {
		privateClaims["scope"] = strings.Join(c.Scopes, " ")
	}

	claimSet := &jws.ClaimSet{
		Issuer:        c.Issuer,
		Audience:      c.Audience,
		Subject:       c.Subject,
		PrivateClaims: privateClaims,
	}
	if err := claimSet.SetExpirationClaims(c.Options.getIatOffset(), c.Options.getExpiresIn()); err != nil {
		return nil, err
	}

	tokenString, err := claimSet.JWT(c.Signer)
	if err != nil {
		return nil, err
	}
	formValues := url.Values{
		"grant_type": {defaultGrantType},
		"assertion":  {tokenString},
	}
	for k, v := range c.Options.getFormValues() {
		formValues[k] = v
	}
	return formValues, nil
}

// Token performs a signed JWT request to obtain a new token.
func (c *Config) Token(ctx context.Context) (*oauth2.Token, error) {
	payload, err := c.payload()
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClientFunc.PostForm(ctx, c.TokenURL, payload)
	if err != nil {
		return nil, fmt.Errorf("oauth2/jwt: cannot fetch token: %v", err)
	}
	defer resp.Body.Close()
	raw := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("oauth2/jwt: unable to decode token: %v", err)
	}
	tk, err := oauth2.TokenFromMap(raw, c.Options.getExpiryDelta())
	if err != nil {
		return nil, err
	}
	return tk, c.Options.postToken(ctx, tk, c)
}
