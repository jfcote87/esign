package esign_test

import (
	"bytes"
	"net/http"
	"net/url"
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
	esign.ResolveDSURL(req.URL, t.host, t.acctID)
	req.URL.Scheme = t.scheme

	req.Header.Set("Authorization", "TESTAUTH")
	return nil
}

func TestCallDo(t *testing.T) {

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
			Payload:        url.Values{"a": {"B"}, "c": {"D"}},
		},
		{
			Credential:     cx,
			Path:           "/noaccount/{TEST}/go",
			PathParameters: map[string]string{"{TEST}": "test3"},
			QueryOpts:      make(url.Values),
			Method:         "POST",
			Payload: map[string]interface{}{
				"a": "String", "b": 9,
			},
		},
	}

	testTransport.Add(&testutils.RequestTester{
		Path:     "/restapi/v2/accounts/1234/do/test0/testvar1/go",
		Auth:     "TESTAUTH",
		Method:   "GET",
		Response: testutils.MakeResponse(200, nil, make(http.Header)),
	})
	testTransport.Add(&testutils.RequestTester{
		Path:     "/restapi/v2/noaccount/test1/go",
		Auth:     "TESTAUTH",
		Method:   "GET",
		Query:    "a=B&c=D",
		Response: testutils.MakeResponse(200, nil, make(http.Header)),
	})
	testTransport.Add(&testutils.RequestTester{
		Path:     "/restapi/v2/noaccount/test2/go",
		Auth:     "TESTAUTH",
		Method:   "POST",
		Payload:  []byte("a=B&c=D"),
		Response: testutils.MakeResponse(200, nil, make(http.Header)),
	})
	testTransport.Add(&testutils.RequestTester{
		Path:     "/restapi/v2/noaccount/test3/go",
		Auth:     "TESTAUTH",
		Method:   "POST",
		Payload:  []byte(`{"a":"String","b":9}`),
		Response: testutils.MakeResponse(200, nil, make(http.Header)),
	})
	// Error Test
	testTransport.Add(&testutils.RequestTester{
		Auth:     "TESTAUTH",
		Method:   "GET",
		Response: testutils.MakeResponse(400, []byte(`No JSON`), make(http.Header)),
	})
	testTransport.Add(&testutils.RequestTester{
		Auth:     "TESTAUTH",
		Method:   "GET",
		Response: testutils.MakeResponse(400, []byte(`{"errorCode": "A", "message":"error desc"}`), make(http.Header)),
	})

	var jsonResponse = []byte(`{"a": "val", "b": 9, "c": "X"}`)
	testTransport.Add(&testutils.RequestTester{
		Path:     "/restapi/v2/accounts/1234/do/test6/go",
		Auth:     "TESTAUTH",
		Method:   "POST",
		Response: testutils.MakeResponse(200, jsonResponse, make(http.Header)),
	})

	for i, c := range calls {
		if err := c.Do(context.Background(), nil); err != nil {
			t.Errorf("Error %d: %v", i, err)
		}
	}

	// Check Error handling
	call := &esign.Call{
		Credential:     cx,
		Path:           "do/{TEST}/go",
		PathParameters: map[string]string{"{TEST}": "test4"},
		QueryOpts:      make(url.Values),
		Method:         "GET",
	}
	err := call.Do(context.Background(), nil)
	if ex, ok := err.(*esign.ResponseError); !ok {
		t.Fatalf("test 4 expected *ResponseError; got %#v", err)
	} else if ex.Status != 400 || string(ex.Raw) != "No JSON" {
		t.Fatalf("test 4 expected status return of 400 and raw of \"No JSON\"; got %d, %q", ex.Status, string(ex.Raw))
	}

	call.PathParameters = map[string]string{"{TEST}": "test5"}
	err = call.Do(context.Background(), nil)
	if ex, ok := err.(*esign.ResponseError); !ok {
		t.Fatalf("test 5 expected *ResponseError; got %#v", err)
	} else if ex.Status != 400 || ex.Err != "A" || ex.Description != "error desc" {
		t.Fatalf("test 5 expected status: 400, err: \"A\", description: \"error desc\"; got %#v", ex)
	}

	call.PathParameters = map[string]string{"{TEST}": "test6"}
	call.Method = "POST"
	var result struct {
		A string
		B int64
		C string
	}
	if err := call.Do(context.Background(), &result); err != nil {
		t.Fatalf("JSON response failed: %v", err)
	}
	if result.A != "val" || result.B != 9 || result.C != "X" {
		t.Fatalf("JSON response expected A=val, B=9, C=X; got %#v", result)
	}

}

type testFile struct {
	*bytes.Buffer
	isClosed    bool
	timesClosed int
}

func (t *testFile) Close() error {
	t.isClosed = true
	t.timesClosed++
	return nil
}

func TestFileUpload(t *testing.T) {
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
	_ = cx
	testTransport.Add(&testutils.RequestTester{
		Path:     "/restapi/v2/accounts/1234/multi",
		Auth:     "TESTAUTH",
		Method:   "POST",
		Response: testutils.MakeResponse(200, nil, make(http.Header)),
	})
	call := &esign.Call{
		Credential:     cx,
		Path:           "do/{TEST}/go",
		PathParameters: map[string]string{"{TEST}": "test4"},
		QueryOpts:      make(url.Values),
		Method:         "POST",
	}
	_ = call

}
