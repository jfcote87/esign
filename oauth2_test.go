// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package esign_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jfcote87/oauth2"
	"github.com/jfcote87/testutils"
	"github.com/pwaterz/esign"
)

const tokenSuccessResponse = `{
	"access_token": "ISSUED_ACCESS_TOKEN",
	"token_type": "Bearer",
	"refresh_token": "ISSUED_REFRESH_TOKEN",
	"expires_in": 28800
  }`

const userInfoSuccessResponse = `{
	"sub": "50d89ab1-dad5-d00d-b410-92ee3110b970",
	"accounts": [
	  {
		"account_id": "fe0b61a3-3b9b-cafe-b7be-4592af32aa9b",
		"is_default": true,
		"account_name": "World Wide Co",
		"base_uri": "https://gotest.docusign.net"
	  },
	  {
		"account_id": "abcd61a3-3b9b-cafe-b7be-4592af32aa9b",
		"is_default": false,
		"account_name": "Account2",
		"base_uri": "https://gotest.docusign.net"
	  }
	],
	"name": "Susan Smart",
	"given_name": "Susan",
	"family_name": "Smart",
	"email": "susan.smart@example.com"
  }`

func getOAuth2ConfigTranspot() (*esign.OAuth2Config, *testutils.Transport) {
	testTransport := &testutils.Transport{}
	clx := &http.Client{Transport: testTransport}

	cfg := &esign.OAuth2Config{
		IntegratorKey: "KEY",
		Secret:        "SECRET",
		RedirURL:      "https://www.example.com/token",
		IsDemo:        true,
		HTTPClientFunc: func(ctx context.Context) (*http.Client, error) {
			return clx, nil
		},
	}
	return cfg, testTransport
}

func TestOuauth2Config_AuthURL(t *testing.T) {
	cfg, _ := getOAuth2ConfigTranspot()
	authURL := cfg.AuthURL("STATE")
	expectedURL := "https://account-d.docusign.com/oauth/auth?client_id=KEY&redirect_uri=https%3A%2F%2Fwww.example.com%2Ftoken&response_type=code&scope=signature&state=STATE"
	if authURL != expectedURL {
		t.Errorf("expected %s; got %s", expectedURL, authURL)
		return
	}

	// check for %20 replacement
	cfg.ExtendedLifetime = true
	cfg.Prompt = true
	cfg.UIlocales = []string{"en-us"}
	authURL = cfg.AuthURL("STATE")
	expectedURL = "https://account-d.docusign.com/oauth/auth?client_id=KEY&prompt=login&redirect_uri=https%3A%2F%2Fwww.example.com%2Ftoken&response_type=code&scope=signature%20extended&state=STATE&ui_locales=en-us"
	if authURL != expectedURL {
		t.Errorf("expected %s; got %s", expectedURL, authURL)
		return
	}

	cfg.UIlocales = nil
	cfg.Prompt = false
	authURL = cfg.AuthURL("STATE", "ASCOPE", "extended")
	expectedURL = "https://account-d.docusign.com/oauth/auth?client_id=KEY&redirect_uri=https%3A%2F%2Fwww.example.com%2Ftoken&response_type=code&scope=ASCOPE%20extended&state=STATE"
	if authURL != expectedURL {
		t.Errorf("expected %s; got %s", expectedURL, authURL)
		return
	}
	authURL = cfg.AuthURL("STATE", "ASCOPE")
	if authURL != expectedURL {
		t.Errorf("expected %s; got %s", expectedURL, authURL)
		return
	}
}

var exchangeResponseTest = &testutils.RequestTester{
	Host:    "account-d.docusign.com",
	Path:    "/oauth/token",
	Method:  "POST",
	Auth:    "Basic S0VZOlNFQ1JFVA==",
	Payload: []byte("code=CODE&grant_type=authorization_code&redirect_uri=https%3A%2F%2Fwww.example.com%2Ftoken"),
	ResponseFunc: func(r *http.Request) (*http.Response, error) {
		return testutils.MakeResponse(200, []byte(tokenSuccessResponse), nil), nil
	},
}
var userinfoResponseTest = &testutils.RequestTester{
	Host:   "account-d.docusign.com",
	Path:   "/oauth/userinfo",
	Method: "GET",
	Auth:   "Bearer ISSUED_ACCESS_TOKEN",
	ResponseFunc: func(r *http.Request) (*http.Response, error) {
		return testutils.MakeResponse(200, []byte(userInfoSuccessResponse), nil), nil
	},
}

var refreshResponseTest = &testutils.RequestTester{
	Path:    "/oauth/token",
	Payload: []byte("grant_type=refresh_token&refresh_token=refresh"),
	ResponseFunc: func(r *http.Request) (*http.Response, error) {
		return testutils.MakeResponse(200, []byte(tokenSuccessResponse), nil), nil
	},
}

func TestOAuth2Config_Exchange(t *testing.T) {
	// Test OAuth2Credential flow
	cfg, testTransport := getOAuth2ConfigTranspot()

	testTransport.Add(exchangeResponseTest, userinfoResponseTest)
	ctx := context.Background()

	var savedToken *oauth2.Token
	var savedUserInfo *esign.UserInfo

	cfg.CacheFunc = func(cx context.Context, tk oauth2.Token, ui esign.UserInfo) {
		savedToken = &tk
		savedUserInfo = &ui
	}

	ocr, err := cfg.Exchange(ctx, "CODE")
	if err != nil {
		t.Fatalf("expected successful code exchage; got %v", err)
	}
	u, err := ocr.UserInfo(ctx)
	if err != nil {
		t.Fatalf("expected userInfo for Susan Smart; got error %v", err)
	}
	if u.Name != "Susan Smart" {
		t.Fatalf("expected user name Susan Smart; got %s", u.Name)
	}
	if savedToken == nil || savedUserInfo == nil {
		t.Fatalf("token and userinfo should be cached; got savedToken is nil %v and savedUserInfo is nil %v", (savedToken == nil), (savedUserInfo == nil))
	}

	tk, err := ocr.Token(ctx)
	if err != nil {
		t.Fatalf("expected token; got %v", err)
	}
	cfg.AccountID = "INVALID ACCOUNT"
	if _, err = cfg.Credential(tk, u); err == nil || err.Error() != "no account INVALID ACCOUNT for susan.smart@example.com" {
		t.Fatalf("expected no account INVALID ACCOUNT for susan.smart@example.com; got %v", err)
	}
	cfg.AccountID = "fe0b61a3-3b9b-cafe-b7be-4592af32aa9b"
	if _, err = cfg.Credential(tk, u); err != nil {
		t.Fatalf("expected successful credential; got %v", err)
	}
	if _, err = cfg.Credential(nil, nil); err == nil || err.Error() != "token may not be nil" {
		t.Fatalf("expected \"token may not be nil\"; got %v", err)
	}
}

func TestOAuth2Config_Refresh(t *testing.T) {
	cfg, testTransport := getOAuth2ConfigTranspot()

	var savedToken *oauth2.Token
	var savedUserInfo *esign.UserInfo

	cfg.CacheFunc = func(cx context.Context, tk oauth2.Token, ui esign.UserInfo) {
		savedToken = &tk
		savedUserInfo = &ui
	}

	testTransport.Add(refreshResponseTest, userinfoResponseTest)

	var tk *oauth2.Token
	ctx := context.Background()

	ocr, err := cfg.Credential(&oauth2.Token{RefreshToken: "refresh"}, nil)
	if err != nil {
		t.Fatalf("expected successful credential create; got %v", err)
	}
	if tk, err = ocr.Token(ctx); err != nil {
		t.Fatalf("expected token; got %v", err)
	}
	if tk.AccessToken != "ISSUED_ACCESS_TOKEN" {
		t.Fatalf("expected token ISSUED_ACCESS_TOKEN; got %s", tk.AccessToken)
	}

	testTransport.Add(refreshResponseTest, userinfoResponseTest)
	ocr, err = cfg.Credential(&oauth2.Token{RefreshToken: "refresh"}, nil)
	if err != nil {
		t.Fatalf("expected successful credential create; got %v", err)
	}
	u, err := ocr.UserInfo(ctx)
	if err != nil {
		t.Fatalf("expecte userinfo success; got %v", err)
	}
	if u.Email != "susan.smart@example.com" {
		t.Fatalf("expected email susan.smart@example.com; got %s", u.Email)
	}

	if savedToken == nil || savedUserInfo == nil {
		t.Fatalf("token and userinfo should be cached; got savedToken is nil %v and savedUserInfo is nil %v", (savedToken == nil), (savedUserInfo == nil))
	}

	testTransport.Add(&testutils.RequestTester{
		Path:   "/restapi/v2/accounts/" + u.Accounts[0].AccountID + "/abc/def",
		Header: http.Header{"Authorization": {"Bearer ISSUED_ACCESS_TOKEN"}},
		Host:   "gotest.docusign.net",
	}, &testutils.RequestTester{
		Path:   "/restapi/v2.1/accounts/" + u.Accounts[0].AccountID + "/abc/def",
		Header: http.Header{"Authorization": {"Bearer ISSUED_ACCESS_TOKEN"}},
		Host:   "gotest.docusign.net",
	})
	req, _ := http.NewRequest("GET", "abc/def", nil)
	if res, err := ocr.AuthDo(ctx, req, nil); err != nil {
		_ = res
		t.Errorf("%v", err)
	} else {
		res.Body.Close()
	}
	req, _ = http.NewRequest("GET", "abc/def", nil)
	if res, err := ocr.AuthDo(ctx, req, esign.VersionV21); err != nil {
		t.Errorf("%v", err)
	} else {
		res.Body.Close()
	}
}

func TestJWTConfig(t *testing.T) {
	var testPK = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAyki3KNQlqFYHQOg+uywV1GNbi/Zvgs2MLYVMiJ/NYeBIZgMm
STDW8mtiR1kLSMq/glzvQdFWPZTzbxkIqiYESoUsErIbZVsMzDNgneDy3XZqXYAS
qT5X2QH1vsCP6Cni4T7Ooj6aFqAsq/7ERGoudP4CO8he82QlcWNMupoWrNZw12AB
J4HSqGT6ebi2YaPXCPCVMr3NqBc8AJGkaFG+RokhRCqSUZUboVQ52vLt7f4Xn4FI
0HAWYegA3kEsCTVQmsNSX/3pUGoCtg4kAOKDUfyPHPCWjA94M8OAU5qnXg/HnZTP
1uP5XnaNhd+po/LklqxMY2tCUf6VUhilUNyw0QIDAQABAoIBAQCh0oIT+4MUo52x
4xksCxx7h/CYi1Cxx1W4pMaRFaXsAsxoL2TVcGjEDfvVL/rDBM8nrskIUjs3kI0d
91zjIP6VzutvGWSpNKmMQh2sr2QanryAiBBlrCYCyHqbWtjE1Z1WrDQJvyLtrr2N
6oWAZaE8nmeTA7xR4W/CwbmEHfi90nB9xxtb6iJNMJAguMsvQ+oBxN4tQYCeNUGo
r88wd8vQyQjFCuU7Jzt8oSzcrP7D/pCgR4XhpU4ODsif8KMaAXS6H7Pt0QfLTkST
AaIq9NBjBvQ5VqkpwWvGHzE2oZ2cfVBu3+sfhi3bmNCkHnmoPlOhfortVDDObwpw
FA4+f71BAoGBAP80L/WseRIOqDkQ+wKbdMOwmyk8p6AlqnDiiGNXe2OsOarImTNn
U2L4xr8MpmOjkDr1aF7e6lIXvtDWyqrIaqmlMf/8xNGMNu24kFTRNxqlII9Yq3fP
sB0LGygnm1aEznK3uKzEIPFdHG0liOdsI3O6TF0PZXPFDFkJV+ERaRFFAoGBAMrq
Q9MjCYrVX2hlyYnv8l2EhQA3AtUXcQhM2JoH1pY/0QwLjloPrUnHSsWuRxf3vuA0
jkSzaoqOu2g/RyVEIPfhaLSptSs82vnLytsE+oPOKfQB28EyfJZcddbONmnCuJY1
4QKYVOzZBqDArD1U5JMZu3UotL2QmXDZDzamtIwdAoGBAMtU0UF0gaIZe368QMH7
CjVAaN+aLBQ07m+yjehYsz7e4bNo0GdcU9vvSqq9cXTBxRC0psuv4BI4SRgrip43
wIQZ0pSa2FX82WbePmDVsInSNvb/Nt7m4vLA/oonxGRSvAo6xzEfsv+bqCJuXX3F
cxmpvV4H/lUXEpd+Ej6ImKXhAoGBALBQ0tJ5lWcPdLGQEIlM97oO1kqTgmCK1+qw
a12cBffUR99Bg1X6XUbIZs5SWvAWk8LZp+1GQQNYdrtkkHtvMX5yXLru479IR7Xa
QNADCXLSB15A5yR+rAczHCmkUV+glSfgdT3+A30yLzIreP5p75tqNprc3gABz3Jh
CXkhbax5AoGAMrZdtA8h9gTdQfqo7QTpUHVP7sFm1Cv/JVDR+iIguF9inLPA/jqN
LHOH+9K3mKx8s6FIuSKsB9it1xCBx5PcP5lBE/9E0z72HC4S7eVVZJEQU2YxfLyS
ZhC2gm1mAAZF9SBYwxTJ7vIcXRWi8uOB6yM7QQhuUpduK236a1lJZao=
-----END RSA PRIVATE KEY-----`
	testTransport := &testutils.Transport{}
	clx := &http.Client{Transport: testTransport}

	cfg := esign.JWTConfig{
		IntegratorKey: "KEY",
		PrivateKey:    testPK,
		KeyPairID:     "1234567890123",
		IsDemo:        true,
		HTTPClientFunc: func(ctx context.Context) (*http.Client, error) {
			return clx, nil
		},
	}
	var expectedConsentURL = "https://account-d.docusign.com/oauth/auth?client_id=KEY&redirect_uri=https%3A%2F%2Fwww.docusign.com&response_type=code&scope=signature%20impersonation"
	if userConsentURL := cfg.UserConsentURL("https://www.docusign.com"); userConsentURL != expectedConsentURL {
		t.Fatalf("expected %s; got %s", expectedConsentURL, userConsentURL)
	}
	ocr, _ := cfg.Credential("50d89ab1-dad5-d00d-b410-92ee3110b970", nil, nil)

	_ = ocr
	var exchangeResponseTest = &testutils.RequestTester{
		Host:   "account-d.docusign.com",
		Path:   "/oauth/token",
		Method: "POST",
		ResponseFunc: func(r *http.Request) (*http.Response, error) {
			return testutils.MakeResponse(200, []byte(tokenSuccessResponse), nil), nil
		},
	}
	var userinfoResponseTest = &testutils.RequestTester{
		Host:   "account-d.docusign.com",
		Path:   "/oauth/userinfo",
		Method: "GET",
		Auth:   "Bearer ISSUED_ACCESS_TOKEN",
		ResponseFunc: func(r *http.Request) (*http.Response, error) {
			return testutils.MakeResponse(200, []byte(userInfoSuccessResponse), nil), nil
		},
	}
	testTransport.Add(exchangeResponseTest, userinfoResponseTest)
	ctx := context.Background()
	tk, err := ocr.Token(ctx)
	if err != nil {
		t.Errorf("expected token; got error %v", err)
	}

	ocr, _ = cfg.Credential("50d89ab1-dad5-d00d-b410-92ee3110b970", tk, nil)

	testTransport.Add(userinfoResponseTest)

	u, err := ocr.UserInfo(ctx)
	if err != nil {
		t.Errorf("userinf error: %v", err)
	}

	testTransport.Add(&testutils.RequestTester{
		Path:   "/restapi/v2/accounts/" + u.Accounts[0].AccountID + "/abc/def",
		Header: http.Header{"Authorization": {"Bearer ISSUED_ACCESS_TOKEN"}},
		Host:   "gotest.docusign.net",
	}, &testutils.RequestTester{
		Path:   "/restapi/v2.1/accounts/" + u.Accounts[0].AccountID + "/abc/def",
		Header: http.Header{"Authorization": {"Bearer ISSUED_ACCESS_TOKEN"}},
		Host:   "gotest.docusign.net",
	})
	req, _ := http.NewRequest("GET", "abc/def", nil)
	if res, err := ocr.AuthDo(ctx, req, nil); err != nil {
		_ = res
		t.Errorf("%v", err)
	} else {
		res.Body.Close()
	}
	req, _ = http.NewRequest("GET", "abc/def", nil)
	if res, err := ocr.AuthDo(ctx, req, esign.VersionV21); err != nil {
		t.Errorf("%v", err)
	} else {
		res.Body.Close()
	}
}

func TestTokenCredential(t *testing.T) {
	ctx := context.Background()
	testTransport := &testutils.Transport{}
	cred := esign.TokenCredential("ABCDEF", true).
		SetClientFunc(func(ctx context.Context) (*http.Client, error) {
			return &http.Client{Transport: testTransport}, nil
		})
	testOp := &esign.Op{
		Credential: cred,
		Method:     "GET",
		Path:       "testcmd",
	}
	_ = testOp
	expectedAuthHeader := http.Header{
		"Authorization": []string{"Bearer ABCDEF"},
	}
	testTransport.Add(
		&testutils.RequestTester{
			Path:     "/oauth/userinfo",
			Header:   expectedAuthHeader,
			Response: testutils.MakeResponse(400, []byte("invalid token"), nil),
		},
		&testutils.RequestTester{
			Path:     "/oauth/userinfo",
			Header:   expectedAuthHeader,
			Response: testutils.MakeResponse(200, []byte(userInfoSuccessResponse), nil),
		},
		&testutils.RequestTester{
			Path:     "/restapi/v2/accounts/fe0b61a3-3b9b-cafe-b7be-4592af32aa9b/testcmd",
			Header:   expectedAuthHeader,
			Response: testutils.MakeResponse(200, []byte("{}"), nil),
		},
		&testutils.RequestTester{
			Path:     "/restapi/v2/accounts/abcd61a3-3b9b-cafe-b7be-4592af32aa9b/testcmd",
			Header:   expectedAuthHeader,
			Response: testutils.MakeResponse(200, []byte("{}"), nil),
		},
	)
	// check for userinfo fail
	switch err := testOp.Do(ctx, nil).(type) {
	case nil:
		t.Errorf("invalid token expected 400 status; got success")
		return
	case *esign.ResponseError:
	default:
		t.Errorf("%v", err)
		return
	}

	if err := testOp.Do(ctx, nil); err != nil {
		t.Errorf("%v", err)
	}
	testOp.Credential = cred.WithAccountID("BAD_ACCT")
	if err := testOp.Do(ctx, nil); err == nil || err.Error() != "no account BAD_ACCT for susan.smart@example.com" {
		t.Errorf("expected no account BAD_ACCT; got %v", err)
		return
	}
	testOp.Credential = cred.WithAccountID("abcd61a3-3b9b-cafe-b7be-4592af32aa9b")
	if err := testOp.Do(ctx, nil); err != nil {
		t.Errorf("expected success; got %v", err)
	}
}

func TestJWTExternalAdminConsentURL(t *testing.T) {
	jwtCfg := esign.JWTConfig{
		IntegratorKey: "INT_KEY",
		IsDemo:        false,
	}
	// invalid authType
	_, err := jwtCfg.ExternalAdminConsentURL("https://www.example.com", "a", "", false)
	if err == nil {
		t.Errorf("expected error; got success")
	}
	// scopes empty
	_, err = jwtCfg.ExternalAdminConsentURL("https://www.example.com", "code", "STATE", false)
	if err == nil {
		t.Errorf("expected error; got success")
	}

	authURL, _ := jwtCfg.ExternalAdminConsentURL("https://www.example.com", "code", "STATE", false, "signature", "impersonation")
	expectedURL := "https://account.docusign.com/oauth/auth?admin_consent_scope=signature%20impersonation&client_id=INT_KEY&redirect_uri=https%3A%2F%2Fwww.example.com&response_type=code&scope=openid&state=STATE"
	if authURL != expectedURL {
		t.Errorf("expected %s; got %s", expectedURL, authURL)
		return
	}
	authURL, _ = jwtCfg.ExternalAdminConsentURL("https://www.example.com", "token", "STATE", true, "signature", "impersonation")
	expectedURL = "https://account.docusign.com/oauth/auth?admin_consent_scope=signature%20impersonation&client_id=INT_KEY&prompt=login&redirect_uri=https%3A%2F%2Fwww.example.com&response_type=token&scope=openid&state=STATE"
	if authURL != expectedURL {
		t.Errorf("expected %s; got %s", expectedURL, authURL)
		return
	}
}
