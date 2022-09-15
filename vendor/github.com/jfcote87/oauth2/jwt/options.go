// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Copyright 2019 James F Cote All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jfcote87/oauth2"
	"github.com/jfcote87/oauth2/jws"
)

// ConfigOptions provide additional (rarely used options for the config)
type ConfigOptions struct {
	// ExpiresIn specifies how many seconds the token should be valid.  A server
	// may ignore the requested duration.  Default to 3600 (1 hour) if not set.
	// Use the SetExpiresIn func to set.
	ExpiresIn *int64
	// IatOffset is the number of seconds subtracted from the current time to
	// set the iat claim. Used for machines whose time is not perfectly in sync.
	// Google servers and others will not issue a token if the issued at time(iat)
	// is in the future. Nil value indicates a default of 10 secounds.  Use
	// the SetIatOffset func to set.
	IatOffset *int64
	// ExpiryDelta determines how many seconds sooner a token should
	// expire than the retrieved expires_in setting. Nil uses
	// oauth2.DefaultExpiryDelta.  Use the SetExpiryDelta func for setting.
	ExpiryDelta *int64
	// PrivateClaims(Optional) adds additional private claims to add
	// to the request
	PrivateClaims map[string]interface{} `json:"private_claims,omitempty"`
	// FormValues(Optional) adds addional form fields to request body
	FormValues url.Values `json:"form_values,omitempty"`
	// NewTokenFunc is called when after retrieving a new *oauth2.Token
	NewTokenFunc func(context.Context, *oauth2.Token, *Config) error
}

func (opts *ConfigOptions) getIatOffset() time.Duration {
	if opts == nil || opts.IatOffset == nil {
		return time.Duration(oauth2.DefaultExpiryDelta) * time.Second
	}
	return time.Duration(*opts.IatOffset) * time.Second
}

func (opts *ConfigOptions) getExpiryDelta() time.Duration {
	if opts == nil || opts.ExpiryDelta == nil {
		return time.Duration(oauth2.DefaultExpiryDelta) * time.Second
	}
	return time.Duration(*opts.ExpiryDelta) * time.Second
}

func (opts *ConfigOptions) getExpiresIn() time.Duration {
	if opts == nil || opts.ExpiresIn == nil {
		return time.Hour
	}
	return time.Duration(*opts.ExpiresIn) * time.Second
}

func (opts *ConfigOptions) getPrivateClaims() map[string]interface{} {
	if opts == nil {
		return nil
	}
	return opts.PrivateClaims
}

func (opts *ConfigOptions) getFormValues() url.Values {
	if opts == nil {
		return nil
	}
	return opts.FormValues
}

func (opts *ConfigOptions) postToken(ctx context.Context, tk *oauth2.Token, c *Config) error {
	if opts == nil || opts.NewTokenFunc == nil {
		return nil
	}
	return opts.NewTokenFunc(ctx, tk, c)
}

// SetExpiresIn sets the requested expiration time.  A new *ConfigOptions is
// returned if opts is nil.
func (opts *ConfigOptions) SetExpiresIn(numOfSeconds int64) *ConfigOptions {
	if opts == nil {
		opts = &ConfigOptions{}
	}
	opts.ExpiresIn = &numOfSeconds
	return opts
}

// SetIatOffset sets the IatOffset and returns a new *ConfigOptions if opts
// is nil
func (opts *ConfigOptions) SetIatOffset(numOfSeconds int64) *ConfigOptions {
	if opts == nil {
		opts = &ConfigOptions{}
	}
	opts.IatOffset = &numOfSeconds
	return opts
}

// SetExpiryDelta sets the ExpiryDelta and returns a new *ConfigOptions if
// opts is nil
func (opts *ConfigOptions) SetExpiryDelta(numOfSeconds int64) *ConfigOptions {
	if opts == nil {
		opts = &ConfigOptions{}
	}
	opts.ExpiryDelta = &numOfSeconds
	return opts
}

// SetPrivateClaims does exactly what its name says
func (opts *ConfigOptions) SetPrivateClaims(claims map[string]interface{}) *ConfigOptions {
	if opts == nil {
		opts = &ConfigOptions{}
	}
	opts.PrivateClaims = claims
	return opts
}

// SetFormValues does exactly what its name says
func (opts *ConfigOptions) SetFormValues(values url.Values) *ConfigOptions {
	if opts == nil {
		opts = &ConfigOptions{}
	}
	opts.FormValues = values
	return opts
}

// AddPrivateClaim add claim of key/value to the claimset
func (opts *ConfigOptions) AddPrivateClaim(key string, value interface{}) *ConfigOptions {
	if opts == nil {
		return &ConfigOptions{
			PrivateClaims: map[string]interface{}{key: value},
		}
	}
	if opts.PrivateClaims == nil {
		opts.PrivateClaims = map[string]interface{}{key: value}
	} else {
		opts.PrivateClaims[key] = value
	}
	return opts
}

// AddFormValue adds the key/value to the JWT urlform request
func (opts *ConfigOptions) AddFormValue(key string, value string) *ConfigOptions {
	if opts == nil {
		return &ConfigOptions{
			FormValues: url.Values{key: []string{value}},
		}
	}
	if opts.FormValues == nil {
		opts.FormValues = url.Values{key: []string{value}}
	} else {
		opts.FormValues.Add(key, value)
	}
	return opts
}

// DefaultCfgOptions returns ptr to an empty *ConfigOptions
func DefaultCfgOptions() *ConfigOptions {
	return &ConfigOptions{}
}

// IDTokenAsAccessToken returns a *ConfigOptions which sets the new Token
// access token to the id_token response
func IDTokenAsAccessToken() *ConfigOptions {
	return &ConfigOptions{
		NewTokenFunc: func(ctx context.Context, tk *oauth2.Token, c *Config) error {
			idToken, ok := tk.Extra("id_token").(string)
			if !ok {
				return errors.New("oauth2: response doesn't have JWT token")
			}
			tk.AccessToken = idToken
			return expiryFromIDToken(tk, idToken, c.Options.getExpiryDelta())
		},
	}
}

// IDTokenSetsExpiry returns a *ConfigOptions that sets new Token
// expiry based using id_token
func IDTokenSetsExpiry() *ConfigOptions {
	return &ConfigOptions{
		NewTokenFunc: func(ctx context.Context, tk *oauth2.Token, c *Config) error {
			idToken, ok := tk.Extra("id_token").(string)
			if !ok {
				return nil
			}
			return expiryFromIDToken(tk, idToken, c.Options.getExpiryDelta())
		},
	}
}

func expiryFromIDToken(tk *oauth2.Token, idToken string, delta time.Duration) error {
	claimSet, err := jws.DecodePayload(idToken)
	if err != nil {
		return fmt.Errorf("oauth2/jwt: error decoding JWT id_token: %v", err)
	}
	tk.Expiry = time.Unix(claimSet.ExpiresAt, 0).Add(-delta)
	return nil
}
