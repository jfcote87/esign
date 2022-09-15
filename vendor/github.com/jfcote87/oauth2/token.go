// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Token represents the crendentials used to authorize
// the requests to access protected resources on the OAuth 2.0
// provider's backend.
//
// Most users of this package should not access fields of Token
// directly. They're exported mostly for use by related packages
// implementing derivative OAuth2 flows.
type Token struct {
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`

	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`

	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`

	// Expiry is the optional expiration time of the access token.
	//
	// If zero, TokenSource implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time `json:"expiry,omitempty"`

	// raw optionally contains extra metadata from the server
	// when updating a token.
	raw interface{}
}

// DefaultExpiryDelta determines the number of seconds  a token should
// expire sooner than the delivered expiration time. This avoids late
// expirations due to client-server time mismatches and latency.
const DefaultExpiryDelta int64 = 10

// Type returns t.TokenType if non-empty, else "Bearer".
func (t *Token) Type() string {
	if t.TokenType == "" {
		return "Bearer"
	}
	switch strings.ToLower(t.TokenType) {
	case "bearer":
		return "Bearer"
	case "mac":
		return "MAC"
	case "basic":
		return "Basic"
	}
	return t.TokenType
}

// SetAuthHeader sets the Authorization header to r using the access
// token in t.
//
// This method is unnecessary when using Transport or an HTTP Client
// returned by this package.
func (t *Token) SetAuthHeader(r *http.Request) {
	r.Header.Set("Authorization", t.Type()+" "+t.AccessToken)
}

// WithExtra returns a new Token that's a clone of t, but using the
// provided raw extra map. This is only intended for use by packages
// implementing derivative OAuth2 flows.
func (t *Token) WithExtra(extra interface{}) *Token {
	t2 := new(Token)
	if t != nil { // nil check
		*t2 = *t
	}
	t2.raw = extra
	return t2
}

// Extra returns an extra field.
// Extra fields are key-value pairs returned by the server as a
// part of the token retrieval response.
func (t *Token) Extra(key string) interface{} {
	if t == nil {
		return nil
	}
	if raw, ok := t.raw.(map[string]interface{}); ok {
		return raw[key]
	}

	vals, ok := t.raw.(url.Values)
	if !ok {
		return nil
	}

	v := vals.Get(key)
	switch s := strings.TrimSpace(v); strings.Count(s, ".") {
	case 0: // Contains no "."; try to parse as int
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return i
		}
	case 1: // Contains a single "."; try to parse as float
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}

	return v
}

// expired reports whether the token is expired.
// t must be non-nil.
func (t *Token) expired() bool {
	if t.Expiry.IsZero() {
		return false
	}
	return t.Expiry.Before(time.Now())
}

// Valid reports whether t is non-nil, has an AccessToken, and is not expired.
func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

// TokenFromMap create a *Token from a map[string]interface{}. Expect the
// access_token, refresh_token and token_type values to be strings, expires_in
// may be string or a type convertible to int64.
func TokenFromMap(vals map[string]interface{}, expiryDelta time.Duration) (*Token, error) {
	t := &Token{raw: vals}
	var strValues = []struct {
		nm  string
		ptr *string
	}{
		{nm: "access_token", ptr: &t.AccessToken},
		{nm: "refresh_token", ptr: &t.RefreshToken},
		{nm: "token_type", ptr: &t.TokenType},
	}
	for _, fld := range strValues {
		switch v := vals[fld.nm].(type) {
		case nil:
		case string:
			*fld.ptr = v
		default:
			return nil, fmt.Errorf("%s must be a string", fld.nm)
		}
	}
	var numOfSeconds int64
	switch v := vals["expires_in"].(type) {
	case nil:
		return t, nil
	case int64:
		numOfSeconds = v
	case float64:
		numOfSeconds = int64(v)
	default:
		rv := reflect.Indirect(reflect.ValueOf(v))
		var intType = reflect.TypeOf(int64(0))
		if !rv.IsValid() || !rv.Type().ConvertibleTo(intType) {
			return nil, fmt.Errorf("unable to convert expires_in to int64")
		}
		numOfSeconds = rv.Convert(intType).Int()
	}
	t.Expiry = time.Now().Add((time.Duration(numOfSeconds) * time.Second) - expiryDelta)
	return t, nil
}
