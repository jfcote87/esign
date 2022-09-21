// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package legacy implements the deprecated legacy authentication
// methods for the version 2 Docusign rest api. Api documentation
// for these functions may be found at:
// https://docs.docusign.com/esign/guide/authentication/legacy_auth.html
package legacy

import (
	"context"
	"net/http"
	"net/url"

	"github.com/jfcote87/ctxclient"
	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/v2/model"
)

// Documentation: https://docs.docusign.com/esign/
//
// All REST API endpoints have the following base:
// https://{server}.docusign.net/restapi/v2
//
// DocuSign hosts multiple geo-dispersed ISO 27001-certified and SSAE 16-audited data centers. For example, account holders in North America might have the following baseUrl:
// https://na2.docusign.net/restapi/v2/accounts/{accountId}
// Whereas European users might access the following baseUrl:
// https://eu.docusign.net/restapi/v2/accounts/{accountId}
//
// EXAMPLES
// 	"https://www.docusign.net/restapi/v2"  (deprecated?)
// 	"https://na2.docusign.net/restapi/v2"   (north america)
// 	"https://na3.docusign.net/restapi/v2"   (north america)
// 	"https://eu.docusign.net/restapi/v2"   (europe)
// 	"https://demo.docusign.net/restapi/v2" (sandbox)
var (
	demoHost = "demo.docusign.net"
	baseHost = "www.docusign.net"
)

// OauthCredential provides authorization for rest request via
// docusign's oauth protocol.
//
// Documentation: https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#OAuth2/OAuth2 Authentication Support in DocuSign REST API.htm
// NOTE: Soon to be deprecated (https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#OAuth2/OAuth2 Token Request.htm)
// Use esign.Oauth2Credential in future.
type OauthCredential struct {
	// The docusign account used by the login user.  This may be
	// found using the LoginInformation call.
	AccountID     string `json:"account_id,omitempty"`
	AccessToken   string `json:"access_token,omitempty"`
	Scope         string `json:"scope,omitempty"`
	TokenType     string `json:"token_type,omitempty"`
	Host          string `json:"host,omitempty"`
	IsDemoAccount bool   `json:"isDemo,omitempty"`
	OnBehalfOf    string `json:"onBehalfOf,omitempty"`
	ctxclient.Func
}

// AuthDo updates request with authorization headers, adds account to URL and sends request
func (o OauthCredential) AuthDo(ctx context.Context, op *esign.Op) (*http.Response, error) {
	req, err := op.CreateRequest()
	if err != nil {
		return nil, err
	}
	req.URL = op.Version.ResolveDSURL(req.URL, getHost(o.IsDemoAccount, o.Host), o.AccountID, o.IsDemoAccount)

	var auth string
	if o.TokenType == "" {
		auth = "bearer " + o.AccessToken
	} else {
		auth = o.TokenType + " " + o.AccessToken
	}
	req.Header.Set("Authorization", auth)
	if o.OnBehalfOf != "" {
		req.Header.Set("X-DocuSign-Act-As-User", o.OnBehalfOf)
	}
	res, err := o.Func.Do(ctx, req)
	if nsErr, ok := err.(*ctxclient.NotSuccess); ok {
		return nil, esign.NewResponseError(nsErr.Body, nsErr.StatusCode)
	}
	return res, err
}

// Revoke invalidates the token ensuring that an error will occur on an subsequent uses.
func (o OauthCredential) Revoke(ctx context.Context) error {
	c := &esign.Op{
		Credential: o,
		Method:     "POST",
		Path:       "/oauth2/revoke",
		Payload: url.Values{
			"token": {o.AccessToken},
		},
		QueryOpts: make(url.Values),
	}
	return c.Do(ctx, nil)
}

// Config provides methods to authenticate via a user/password combination.  It may also
// be used to generate an OauthCredential.  If Host is empty, the IsDemoAccount
// is used to determine the host.
//
// Documentation:  https://www.docusign.com/p/RESTAPIGuide/RESTAPIGuide.htm#SOBO/Send On Behalf Of Functionality in the DocuSign REST API.htm
type Config struct {
	// The docusign account used by the login user.  This may be
	// found using the LoginInformation call.
	AccountID     string `json:"acctId,omitempty"`
	IntegratorKey string `json:"key"`
	UserName      string `json:"user"`
	Password      string `json:"pwd"`
	Host          string `json:"host,omitempty"`
	// Deprecated - use Host
	IsDemoAccount bool   `json:"isDemo,omitempty"`
	OnBehalfOf    string `json:"onBehalfOf,omitempty"`
	ctxclient.Func
}

// OauthCredential retrieves an OauthCredential  from docusign
// using the username and password from Config. The returned
// token does not have a expiration although it may be revoked
// via OauthCredential.Revoke()
func (c *Config) OauthCredential(ctx context.Context) (*OauthCredential, error) {
	call := &esign.Op{
		Credential: c,
		Method:     "POST",
		Path:       "/oauth2/token",
		Payload: url.Values{
			"grant_type": []string{"password"},
			"client_id":  []string{c.IntegratorKey},
			"username":   []string{c.UserName},
			"password":   []string{c.Password},
			"scope":      []string{"api"},
		},
		QueryOpts: make(url.Values),
	}
	var ret *model.OauthAccess
	if err := call.Do(ctx, &ret); err != nil {
		return nil, err
	}
	return &OauthCredential{
		AccountID:     c.AccountID,
		AccessToken:   ret.AccessToken,
		Scope:         ret.Scope,
		TokenType:     ret.TokenType,
		Host:          c.Host,
		IsDemoAccount: c.IsDemoAccount,
	}, nil
}

// OnBehalfOfCredential returns an *OauthCredential for the user name specied by nm.  oauthCred
// must be a credential for a user with administrative rights on the account.
func (o *OauthCredential) OnBehalfOfCredential(ctx context.Context, integratorKey, nm string) (*OauthCredential, error) {
	call := &esign.Op{
		Credential: o,
		Method:     "POST",
		Path:       "/oauth2/token",
		Payload: url.Values{
			"grant_type": []string{"password"},
			"client_id":  []string{integratorKey},
			"username":   []string{nm},
			"scope":      []string{"api"},
		},
		QueryOpts: make(url.Values),
	}
	var ret *model.OauthAccess
	if err := call.Do(ctx, &ret); err != nil {
		return nil, err
	}
	return &OauthCredential{
		AccountID:     o.AccountID,
		AccessToken:   ret.AccessToken,
		Scope:         ret.Scope,
		TokenType:     ret.TokenType,
		Host:          o.Host,
		IsDemoAccount: o.IsDemoAccount,
	}, nil
}

// AuthDo adds authorization headers, adds accountID to url and sends request
func (c Config) AuthDo(ctx context.Context, op *esign.Op) (*http.Response, error) {
	req, err := op.CreateRequest()
	req.URL = op.Version.ResolveDSURL(req.URL, getHost(c.IsDemoAccount, c.Host), c.AccountID, c.IsDemoAccount)

	var onBehalfOf string
	if c.OnBehalfOf != "" {
		onBehalfOf = "<SendOnBehalfOf>" + c.OnBehalfOf + "</SendOnBehalfOf>"
	}
	authString := "<DocuSignCredentials>" + onBehalfOf +
		"<Username>" + c.UserName + "</Username><Password>" +
		c.Password + "</Password><IntegratorKey>" +
		c.IntegratorKey + "</IntegratorKey></DocuSignCredentials>"
	req.Header.Set("X-DocuSign-Authentication", authString)
	res, err := c.Func.Do(ctx, req)
	if nsErr, ok := err.(*ctxclient.NotSuccess); ok {
		return nil, esign.NewResponseError(nsErr.Body, nsErr.StatusCode)
	}
	return res, err
}

func getHost(isDemo bool, host string) string {
	if host > "" {
		return host
	}
	if isDemo {
		return demoHost
	}
	return baseHost
}
