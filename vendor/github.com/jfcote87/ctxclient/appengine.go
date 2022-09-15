// Copyright 2019 James Cote All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build appengine appenginevm

// App Engine hooks is based upon appengine code in golang.org/x/oauth2

package ctxclient // import "github.com/jfcote87/ctxclient"

import (
	"context"
	"errors"
	"net/http"

	"google.golang.org/appengine/urlfetch"
)

func init() {
	// set defaultContextClientFunc to return urlFetch client

	defaultFuncs = append(defaultFuncs, func(ctx context.Context) (*http.Client, error) {
		cl := urlfetch.Client(ctx)
		if cl == nil {
			return nil, errors.New("urlfetch returned nil client")
		}
		return cl, nil
	})
}
