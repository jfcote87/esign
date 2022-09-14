// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package oauth2

import (
	"context"
	"errors"
)

// PerRPCCredentials fulfills the grpc PerRPCCredentials interface allowing
// a developer to use a Tokensource for authorization
type PerRPCCredentials struct {
	TokenSource
}

// GetRequestMetadata returns authorization headers for a grpc credential
func (c *PerRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if c == nil {
		return nil, errors.New("nil credential")
	}
	tk, err := c.Token(ctx)
	if err != nil {
		return nil, err
	}
	return map[string]string{"authorization": tk.Type() + " " + tk.AccessToken}, nil
}

// RequireTransportSecurity needed for google.golang.org/grpc/credential
func (c *PerRPCCredentials) RequireTransportSecurity() bool {
	return true
}
