// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

import (
	"net/http"
	"time"

	"github.com/jfcote87/oauth2"
	"github.com/jfcote87/oauth2/jws"
)

// ServiceAccount conforms to the json format of a Google Service Account Key and has
// the same structure as the golang.org/x/oauth2/jwt Config struct.
// Use ServiceAccount.Config() to convert to jwt.Config.
type ServiceAccount struct {
	// Email is the OAuth client identifier used when communicating with
	// the configured OAuth provider.
	Email string `json:"client_email,omitempty"`

	// PrivateKey contains the contents of an RSA private key or the
	// contents of a PEM file that contains a private key. The provided
	// private key is used to sign JWT payloads.
	// PEM containers with a passphrase are not supported.
	// Use the following command to convert a PKCS 12 file into a PEM.
	//
	//    $ openssl pkcs12 -in key.p12 -out key.pem -nodes
	//
	PrivateKey []byte `json:"private_key,omitempty"`

	// PrivateKeyID contains an optional hint indicating which key is being
	// used.
	PrivateKeyID string `json:"private_key_id,omitempty"`

	// Subject is the optional user to impersonate.
	Subject string `json:"subject,omitempty"`

	// Scopes optionally specifies a list of requested permission scopes.
	Scopes []string `json:"scopes,omitempty"`

	// TokenURL is the endpoint required to complete the 2-legged JWT flow.
	TokenURL string `json:"token_uri,omitempty"`

	// Expires optionally specifies how long the token is valid for.
	Expires time.Duration `json:"expires,omitempty"`
}

// Config translates the ServiceAccount settings to a jwt.Config
func (cfg ServiceAccount) Config() (*Config, error) {
	signer, err := jws.RS256FromPEM(cfg.PrivateKey, cfg.PrivateKeyID)
	if err != nil {
		return nil, err
	}
	opts := IDTokenSetsExpiry()
	if cfg.Expires != 0 {
		opts = opts.SetExpiresIn(int64(cfg.Expires / time.Second))
	}
	return &Config{
		Signer:   signer,
		Issuer:   cfg.Email,
		Subject:  cfg.Subject,
		TokenURL: cfg.TokenURL,
		Audience: cfg.TokenURL,
		Scopes:   cfg.Scopes,
		Options:  opts,
	}, nil
}

// Client returns a *jwt.Config client with the passed token.  Error
// returned if problems found with the private key
func (cfg ServiceAccount) Client(tk *oauth2.Token) (*http.Client, error) {
	jcfg, err := cfg.Config()
	if err != nil {
		return nil, err
	}
	return jcfg.Client(tk)
}
