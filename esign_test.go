package esign_test

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jfcote87/esign"
	"github.com/jfcote87/testutils"

	"golang.org/x/net/context"
)

type TestCred struct {
	scheme string
	host   string
	acctID string
}

func (t *TestCred) Authorize(ctx context.Context, req *http.Request) error {
	req.URL.Scheme = t.scheme
	req.URL.Host = t.host

	if !strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = "/accounts/" + t.acctID + "/" + req.URL.Path
	}
	req.Header.Set("Authorization", "TESTAUTH")
	return nil
}

func TestCredential(t *testing.T) {
	var cx esign.Credential
	cx = &TestCred{
		host:   "www.example.com",
		scheme: "https",
		acctID: "1234",
	}
	//cx = esign.WithLogger(cx, func(ctx context.Context, ix interface{}, u *url.URL, h http.Header) {}, nil)
	//t.Errorf("%s", esign.CredentialType(cx))
	//cx = esign.WithHTTPClientFunc(cx, func(ctx context.Context) (*http.Client, error) { return nil, nil })
	//t.Errorf("%s", esign.CredentialType(cx))
	//cx = esign.WithHTTPClientFunc(tx, func(ctx context.Context) (*http.Client, error) { return nil, nil })
	//t.Errorf("%s", esign.CredentialType(cx))

	testTransport := &testutils.Transport{}
	cx = esign.WithHTTPClientFunc(cx, func(ctx context.Context) (*http.Client, error) {
		return &http.Client{Transport: testTransport}, nil
	})
	calls := []esign.Call{
		{
			Credential:     cx,
			Path:           "do/{TEST}/{VAR1}/go",
			PathParameters: map[string]string{"{TEST}": "test0", "{VAR1}": "testvar1"},
			QueryOpts:      make(url.Values),
			Method:         "GET",
		},
		{
			Credential:     cx,
			Path:           "/noaccount/{TEST}/go",
			PathParameters: map[string]string{"{TEST}": "test1"},
			QueryOpts:      url.Values{"a": {"B"}, "c": {"D"}},
			Method:         "GET",
		},
		{
			Credential:     cx,
			Path:           "/noaccount/{TEST}/go",
			PathParameters: map[string]string{"{TEST}": "test2"},
			QueryOpts:      make(url.Values),
			Method:         "POST",
			Payload:        []byte("a=b?c=d"),
		},
		{
			Credential:     cx,
			Path:           "/noaccount/{TEST}/go",
			PathParameters: map[string]string{"{TEST}": "test3"},
			QueryOpts:      make(url.Values),
			Method:         "POST",
			Payload:        []byte("a=b?c=d"),
		},
		{
			Credential:     cx,
			Path:           "/noaccount/{TEST}/go",
			PathParameters: map[string]string{"{TEST}": "test4"},
			QueryOpts:      make(url.Values),
			Method:         "POST",
			Payload:        nil,
		},
	}
	testTransport.Add(func(r *http.Request) (*http.Response, error) {

		if r.URL.Path != "/accounts/1234/do/test0/testvar1/go" {
			return nil, fmt.Errorf("call0: expected path to be /accounts/1234/do/test0/testvar1/go; got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "TESTAUTH" {
			return nil, fmt.Errorf("call0: expected authorization TESTAUTH; got %s.", r.Header.Get("Authorization"))
		}
		return testutils.MakeResponse(200, nil, make(http.Header)), nil
	})
	testTransport.Add(&testutils.RequestTester{
		Path:     "/accounts/1234/do/test0/testvar1/go",
		Auth:     "TESTAUTH",
		Method:   "GET",
		Response: testutils.MakeResponse(200, nil, make(http.Header)),
	})
	_ = calls

}
