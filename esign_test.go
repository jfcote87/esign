package esign_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/jfcote87/ctxclient"
	"github.com/jfcote87/esign"
	"github.com/jfcote87/testutils"

	"golang.org/x/net/context"
)

type TestCred struct {
	scheme string
	host   string
	acctID string
	ctxclient.Func
}

func (t *TestCred) Authorize(ctx context.Context, req *http.Request) error {
	log.Printf("OK I'm in Authorize")
	req.URL.Scheme = t.scheme
	req.URL.Host = t.host

	if !strings.HasPrefix(req.URL.Path, "/") {
		req.URL.Path = "/accounts/" + t.acctID + "/" + req.URL.Path
	}
	req.Header.Set("Authorization", "TESTAUTH")
	return nil
}

func TestCredential(t *testing.T) {

	testTransport := &testutils.Transport{}
	clx := &http.Client{Transport: testTransport}

	var cx esign.Credential
	cx = &TestCred{
		host:   "www.example.com",
		scheme: "https",
		acctID: "1234",
		Func: func(ctx context.Context) (*http.Client, error) {
			return clx, nil
		},
	}
	//cx = esign.WithLogger(cx, func(ctx context.Context, ix interface{}, u *url.URL, h http.Header) {}, nil)
	//t.Errorf("%s", esign.CredentialType(cx))
	//cx = esign.WithHTTPClientFunc(cx, func(ctx context.Context) (*http.Client, error) { return nil, nil })
	//t.Errorf("%s", esign.CredentialType(cx))
	//cx = esign.WithHTTPClientFunc(tx, func(ctx context.Context) (*http.Client, error) { return nil, nil })
	//t.Errorf("%s", esign.CredentialType(cx))

	/*cx = esign.WithHTTPClientFunc(cx, func(ctx context.Context) (*http.Client, error) {
		return clx, nil
	})*/
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
			Payload:        []byte("a=b&c=d"),
		},
		{
			Credential:     cx,
			Path:           "/noaccount/{TEST}/go",
			PathParameters: map[string]string{"{TEST}": "test3"},
			QueryOpts:      make(url.Values),
			Method:         "POST",
			Payload:        []byte("a=b&c=d"),
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
		t.Errorf("I am ehere %v !!!!!!!!!!!!1", r.URL)
		if r.URL.Path != "/accounts/1234/do/test0/testvar1/gco" {
			t.Errorf("OK I'm herer")
			return nil, fmt.Errorf("call0: expected path to be /accounts/1234/do/test0/testvar1/go; got %s", r.URL.Path)
		}
		t.Errorf("OK now what %v", r.Header)
		if r.Header.Get("Authorization") != "TESTAUTH" {
			return nil, fmt.Errorf("call0: expected authorization TESTAUTH; got %s.", r.Header.Get("Authorization"))
		}
		t.Errorf("OK now what")
		return testutils.MakeResponse(200, nil, make(http.Header)), nil
	})
	testTransport.Add(&testutils.RequestTester{
		Path:     "/accounts/1234/do/test0/testvar1/go",
		Auth:     "TESTAUTH",
		Method:   "GET",
		Response: testutils.MakeResponse(200, nil, make(http.Header)),
	})

	//t.Fatalf("Done")
	_ = calls
	//_, err := clx.Get("http://amagiweb.libertyfund.org/")
	//log.Printf("Err: %v", err)
	for i, c := range calls {
		if i == 0 {
			err := c.Do(context.Background(), nil)
			t.Errorf("%v", err)

		}
	}

}
