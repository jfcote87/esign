// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package oauth2

import (
	"errors"
	"net/http"

	"github.com/jfcote87/ctxclient"
)

// Client creates a client that sets an authorization
// header based on tokens created by the ctxtokensource ts.
func Client(ts TokenSource, f ctxclient.Func) *http.Client {
	if ts == nil {
		return &http.Client{
			Transport: &ctxclient.ErrorTransport{Err: errors.New("oauth2: nil tokensource specified")},
		}
	}
	return &http.Client{
		Transport: &Transport{
			Source: ts,
			Func:   f,
		},
	}
}

// Transport is an http.RoundTripper that makes OAuth 2.0 HTTP requests,
// wrapping a base RoundTripper and adding an Authorization header
// with a token from the supplied Sources.
//
// Transport is a low-level mechanism. Most code will use the
// higher-level Config.Client method instead.
type Transport struct {
	// Source supplies the token to add to outgoing requests'
	// Authorization headers.
	Source TokenSource
	ctxclient.Func
}

// RoundTrip authorizes and authenticates the request with an
// access token. If no token exists or token is expired,
// it fetches a new token passing along the request's context.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Source == nil {
		return ctxclient.RequestError(req, errors.New("oauth2: Transport's Source is nil"))
	}
	ctx := req.Context()
	tk, err := t.Source.Token(ctx)
	if err != nil {
		return ctxclient.RequestError(req, err)
	}
	req2 := cloneRequest(req) // per RoundTripper contract
	tk.SetAuthHeader(req2)
	return t.Func.Do(ctx, req2)
}

// deep copy header
func cloneRequest(r *http.Request) *http.Request {
	h := make(http.Header)
	for k, v := range r.Header {
		h[k] = v
	}
	newReq := *r

	newReq.Header = h
	return &newReq
}
