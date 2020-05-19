// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package esign_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jfcote87/ctxclient"
	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/v2.1/folders"
	"github.com/jfcote87/esign/v2.1/templates"
	"github.com/jfcote87/oauth2"
	"github.com/jfcote87/testutils"
)

type TestCred struct {
	scheme string
	host   string
	acctID string
	ctxclient.Func
}

func (t *TestCred) AuthDo(ctx context.Context, req *http.Request, v *esign.APIVersion) (*http.Response, error) {
	req.URL = v.ResolveDSURL(req.URL, t.host, t.acctID)

	req.Header.Set("Authorization", "TESTAUTH")
	res, err := t.Func.Do(ctx, req)
	if nsErr, ok := err.(*ctxclient.NotSuccess); ok {
		return nil, esign.NewResponseError(nsErr.Body, nsErr.StatusCode)
	}
	return res, err
}

func (t *TestCred) SetClient(cl *http.Client) {
	t.Func = func(ctx context.Context) (*http.Client, error) {
		return cl, nil
	}
}

func TestOp_Do(t *testing.T) {
	// Setup test credential and transport to check for path substitution
	cx, testTransport := getTestCredentialClientTransport()

	var jsonResponse = []byte(`{"a": "val", "b": 9, "c": "X"}`)
	testTransport.Add(
		&testutils.RequestTester{ // test 0
			Path:     "/restapi/v2.1/accounts/1234/do/test0/testvar1/go",
			Auth:     "TESTAUTH",
			Method:   "GET",
			Response: testutils.MakeResponse(200, nil, nil),
		},
		&testutils.RequestTester{ // test 1
			Path:     "/restapi/v2/noaccount/test1/go",
			Auth:     "TESTAUTH",
			Method:   "GET",
			Query:    "a=B&c=D",
			Response: testutils.MakeResponse(200, nil, nil),
		},
		&testutils.RequestTester{ // test 2
			Path:     "/restapi/v2/noaccount/test2/go",
			Auth:     "TESTAUTH",
			Method:   "POST",
			Payload:  []byte("a=B&c=D"),
			Response: testutils.MakeResponse(200, nil, nil),
		},
		&testutils.RequestTester{ // test 3
			Path:     "/restapi/v2/noaccount/test3/go",
			Auth:     "TESTAUTH",
			Method:   "POST",
			Payload:  []byte("{\"a\":\"String\",\"b\":9}\n"),
			Response: testutils.MakeResponse(200, nil, nil),
		},
		&testutils.RequestTester{ // test 4
			Path:     "/restapi/v2/noaccount/test4/go",
			Auth:     "TESTAUTH",
			Method:   "POST",
			Payload:  []byte("0123456789"),
			Response: testutils.MakeResponse(200, nil, nil),
		},
		&testutils.RequestTester{ // test 5
			Auth:     "TESTAUTH",
			Method:   "GET",
			Response: testutils.MakeResponse(400, []byte(`No JSON`), nil),
		},
		&testutils.RequestTester{ // test 6
			Auth:     "TESTAUTH",
			Method:   "GET",
			Response: testutils.MakeResponse(400, []byte(`{"errorCode": "A", "message":"error desc"}`), nil),
		},
		&testutils.RequestTester{ // test 7
			Path:     "/restapi/v2/accounts/1234/do/test7/go",
			Auth:     "TESTAUTH",
			Method:   "POST",
			Response: testutils.MakeResponse(200, jsonResponse, nil),
		}) // Error Test

	ops := []esign.Op{
		{
			Credential: cx,
			Path:       strings.Join([]string{"do", "test0", "testvar1", "go"}, "/"),
			QueryOpts:  make(url.Values),
			Method:     "GET",
			Version:    esign.VersionV21,
		},
		{
			Credential: cx,
			Path:       "/v2/noaccount/test1/go",
			QueryOpts:  url.Values{"a": {"B"}, "c": {"D"}},
			Method:     "GET",
		},
		{
			Credential: cx,
			Path:       "/v2/noaccount/test2/go",
			QueryOpts:  make(url.Values),
			Method:     "POST",
			Payload:    url.Values{"a": {"B"}, "c": {"D"}},
		},
		{
			Credential: cx,
			Path:       "/v2/noaccount/test3/go",
			QueryOpts:  make(url.Values),
			Method:     "POST",
			Payload: map[string]interface{}{
				"a": "String", "b": 9,
			},
		},
		{
			Credential: cx,
			Path:       "/v2/noaccount/test4/go",
			QueryOpts:  make(url.Values),
			Method:     "POST",
			Payload: &esign.UploadFile{
				Reader:      bytes.NewReader([]byte("0123456789")),
				ContentType: "text/plain",
			},
		},
	}
	for i, op := range ops {
		if err := op.Do(context.Background(), nil); err != nil {
			t.Errorf("Error %d: %v", i, err)
		}
	}

	// check error handling for failed ResponseError
	op := &esign.Op{
		Credential: cx,
		Path:       "do/test5/go",
		QueryOpts:  make(url.Values),
		Method:     "GET",
	}
	err := op.Do(context.Background(), nil)
	if ex, ok := err.(*esign.ResponseError); !ok {
		t.Fatalf("test 5 expected *ResponseError; got %#v", err)
	} else if ex.Status != 400 || string(ex.Raw) != "No JSON" {
		t.Fatalf("test 5 expected status return of 400 and raw of \"No JSON\"; got %d, %q", ex.Status, string(ex.Raw))
	}

	// check error handling for properly formateed ResponseError
	op.Path = "do/test6/go"
	err = op.Do(context.Background(), nil)
	if ex, ok := err.(*esign.ResponseError); !ok {
		t.Fatalf("test 6 expected *ResponseError; got %#v", err)
	} else if ex.Status != 400 || ex.ErrorCode != "A" || ex.Description != "error desc" {
		t.Fatalf("test 6 expected status: 400, err: \"A\", description: \"error desc\"; got %#v", ex)
	}

	// check post and return value. expect success
	op.Path = "do/test7/go"
	op.Method = "POST"
	var result struct {
		A string
		B int64
		C string
	}
	if err := op.Do(context.Background(), &result); err != nil {
		t.Fatalf("JSON response failed: %#v", err)
	}
	if result.A != "val" || result.B != 9 || result.C != "X" {
		t.Fatalf("JSON response expected A=val, B=9, C=X; got %#v", result)
	}

}

type testFile struct {
	reader      *bytes.Reader
	closed      bool
	timesClosed int
	readFunc    func(*testFile, []byte) (int, error)
	m           sync.Mutex
}

func (t *testFile) Read(b []byte) (int, error) {
	if t.readFunc != nil {
		return t.readFunc(t, b)
	}
	return t.reader.Read(b)
}

func (t *testFile) IsClosed() bool {
	defer t.m.Unlock()
	t.m.Lock()
	return t.closed
}

func (t *testFile) Close() error {
	t.m.Lock()
	t.closed = true
	t.timesClosed++
	t.m.Unlock()
	return nil
}

func (t *testFile) ReOpen() {
	t.m.Lock()
	t.closed = false
	t.timesClosed = 0
	t.reader.Seek(0, io.SeekStart)
	t.m.Unlock()
}

type errReader struct{}

func (e errReader) Read(p []byte) (int, error) {
	return 0, errors.New("Read Error")
}

// testMultipart validates formatting of a multipart op
func testMultipart(req *http.Request) (*http.Response, error) {
	// Check multipart ops.  ensure fileUploads closed Do() routine.
	expectedContent := []string{"application/json", "text/plain", "application/octet-stream"}
	expectedFileName := []string{"", "file1.txt", "file2.txt"}
	defer req.Body.Close()

	_, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	mpr := multipart.NewReader(req.Body, params["boundary"])
	p, err := mpr.NextPart()
	for i := 0; err == nil && i < len(expectedContent); i++ {
		if expectedContent[i] != p.Header.Get("Content-Type") || expectedFileName[i] != p.FileName() {
			err = fmt.Errorf("multipart[%d] expected content-type: %s and file name: %s; got %s  %s", i, expectedContent[i], expectedFileName[i], p.Header.Get("Content-Type"), p.FileName())
		}
		p.Close()

		if err == nil {
			p, err = mpr.NextPart()
		}
	}
	if err != io.EOF {
		return testutils.MakeResponse(400, []byte(err.Error()), nil), nil
	}
	return testutils.MakeResponse(200, nil, nil), nil
}

func checkFilesClosed(t *testing.T, files ...*testFile) bool {
	var ok = true
	for i, f := range files {
		if !f.IsClosed() {
			t.Errorf("file%d not closed", i)
			ok = false
		}
		f.ReOpen()
	}
	return ok
}

func TestOp_Do_FileUpload(t *testing.T) {
	cx, testTransport := getTestCredentialClientTransport()

	payload := struct {
		A string
		B int
	}{"STRING", 9}
	f1 := &testFile{
		reader: bytes.NewReader([]byte("12345678")),
	}
	f2 := &testFile{
		reader: bytes.NewReader([]byte("9876543210")),
	}

	// expect success
	testTransport.Add(&testutils.RequestTester{
		ResponseFunc: testMultipart,
	})
	op := &esign.Op{
		Credential: cx,
		Path:       "multipart/go",
		QueryOpts:  make(url.Values),
		Method:     "POST",
		Payload:    payload,
		Files: []*esign.UploadFile{
			{
				ContentType: "text/plain",
				FileName:    "file1.txt",
				ID:          "1",
				Reader:      f1,
			},
			{
				ContentType: "application/octet-stream",
				FileName:    "file2.txt",
				ID:          "2",
				Reader:      f2,
			},
		},
	}
	if err := op.Do(context.Background(), nil); err != nil {
		t.Fatalf("multipart test expected success; got %v", err)
	}
	if !checkFilesClosed(t, f1, f2) {
		t.Fatalf("multipart test success expected closed files")
	}

	f1.ReOpen()
	f2.ReOpen()
	// ensure files close on transport/network error
	cx.(*TestCred).SetClient(&http.Client{
		Transport: &ctxclient.ErrorTransport{Err: errors.New("ERROR")},
	})
	ctx := context.Background()
	if err := op.Do(ctx, nil); err == nil || err.Error() != "Post \"https://www.example.com/restapi/v2/accounts/1234/multipart/go\": ERROR" {
		t.Fatalf("multipart test expected post error; got %v", err)
	}
	time.Sleep(time.Second)
	if !checkFilesClosed(t, f1, f2) {
		t.Errorf("multipart network error expected closed files")
	}

	// ensure all files close if error reading previous file
	f1.ReOpen()
	f2.ReOpen()
	op.Files[0].Reader = errReader{}
	cx.(*TestCred).SetClient(&http.Client{
		Transport: testTransport,
	})
	testTransport.Add(&testutils.RequestTester{
		Payload:  []byte("Expect Error"),
		Response: testutils.MakeResponse(200, nil, nil),
	})

	switch err := op.Do(context.Background(), nil).(type) {
	case *url.Error:
	default:
		t.Errorf("multipart read expected io error; got %v", err)
	}

	if !f2.IsClosed() {
		t.Fatalf("multipart test post error expected closed files")
	}
}

func TestOp_Do_FileDownload(t *testing.T) {
	cx, testTransport := getTestCredentialClientTransport()
	op := &esign.Op{
		Credential: cx,
		Path:       "file",
		QueryOpts:  make(url.Values),
		Method:     "GET",
		Accept:     "text/plain",
	}
	f1 := &testFile{
		reader: bytes.NewReader([]byte("0123456789")),
	}

	var file *esign.Download
	testTransport.Add(&testutils.RequestTester{
		Header: http.Header{"Accept": []string{"text/plain"}},
		ResponseFunc: func(req *http.Request) (*http.Response, error) {
			res := testutils.MakeResponse(200, []byte("0123456789"), http.Header{"Content-Type": []string{"text/plain"}})
			res.Body = f1
			return res, nil
		},
	})
	if err := op.Do(context.Background(), &file); err != nil {
		t.Fatalf("expecte esign.Download; got error %v", err)
	}
	if file == nil {
		t.Errorf("expected *esign.Download; got nil")
	} else if file.ContentType != "text/plain" || file.ContentLength != 10 {
		t.Errorf("expected contentType of text/plain and ContentLength = 10; got %s, %d", file.ContentType, file.ContentLength)
	}
	if f1.IsClosed() {
		t.Errorf("expected open res.Body. got closed.")
	}
	file.Close()
}

func TestOp_FilesClosed(t *testing.T) {
	f1 := &testFile{
		reader: bytes.NewReader([]byte("12345678")),
	}
	f2 := &testFile{
		reader: bytes.NewReader([]byte("9876543210")),
	}
	op := &esign.Op{
		Credential: nil,
		Path:       "multipart/go",
		QueryOpts:  make(url.Values),
		Method:     "POST",
		Files: []*esign.UploadFile{
			{
				ContentType: "text/plain",
				FileName:    "file1.txt",
				ID:          "1",
				Reader:      f1,
			},
			{
				ContentType: "application/octet-stream",
				FileName:    "file2.txt",
				ID:          "2",
				Reader:      f2,
			},
		},
	}
	// ensure files close on nil context/invalid client/invalid credential error
	var ctx context.Context
	if err := op.Do(ctx, nil); err != nil && err.Error() != "nil context" {
		t.Errorf("expected nil context; got %v", err)
	}
	if !checkFilesClosed(t, f1, f2) {
		t.Fatalf("multipart test nil context expected closed files")
	}
	ctx = context.Background()
	if err := op.Do(ctx, nil); err != nil && err.Error() != "nil credential" {
		t.Errorf("expected nil credential; got %v", err)
	}
	if !checkFilesClosed(t, f1, f2) {
		t.Fatalf("multipart test nil credential expected closed files")
	}

	// create an error in Op.makeRequest
	cx, _ := getTestCredentialClientTransport()
	op.Credential = cx
	op.Method = "PO ST" //invalid method
	if err := op.Do(ctx, nil); err != nil && err.Error() != "net/http: invalid method \"PO ST\"" {
		t.Errorf("expected net/http: invalid method \"PO ST\"; got %v", err)
	}
	time.Sleep(time.Second)
	if !checkFilesClosed(t, f1, f2) {
		t.Fatalf("invalid method expected closed files")
	}

}

func TestOp_Do_ContextCancel(t *testing.T) {
	cx, testTransport := getTestCredentialClientTransport()
	ctx, cancelFunc := context.WithCancel(context.Background())
	// First run without context cancel for baseline
	op := &esign.Op{
		Credential: cx,
		Path:       "multipart/go",
		QueryOpts:  make(url.Values),
		Method:     "POST",
	}
	testTransport.Add(&testutils.RequestTester{
		Response: testutils.MakeResponse(200, []byte("0123456789"), nil),
	})
	var result *TokenCache
	switch err := op.Do(ctx, &result).(type) {
	case *json.UnmarshalTypeError:
	default:
		t.Fatalf("expected *json.UnmarshalTypeError; got %#v", err)
	}

	testTransport.Add(&testutils.RequestTester{
		ResponseFunc: func(*http.Request) (*http.Response, error) {
			cancelFunc()
			return nil, errors.New("Should not see this error Error")
		},
	})
	testTransport.Add(&testutils.RequestTester{
		Response: testutils.MakeResponse(200, []byte("0123456789"), nil),
	})
	op = &esign.Op{
		Credential: cx,
		Path:       "multipart/go",
		QueryOpts:  make(url.Values),
		Method:     "POST",
	}
	if err := op.Do(ctx, &result); err == nil || err.Error() != "context canceled" {
		t.Errorf("expected context canceled; got %v", err)
	}
}

func getTestCredentialClientTransport() (esign.Credential, *testutils.Transport) {
	testTransport := &testutils.Transport{}
	clx := &http.Client{Transport: testTransport}

	var cx esign.Credential = &TestCred{
		host:   "www.example.com",
		scheme: "https",
		acctID: "1234",
		Func: func(ctx context.Context) (*http.Client, error) {
			return clx, nil
		},
	}
	return cx, testTransport
}

type TokenCache struct {
	Token *oauth2.Token   `json:"token"`
	User  *esign.UserInfo `json:"user_info"`
}

// TestGenerateOps creates an OAuth2Credential using environment
// variables DOCUSIGN_Token DOCUSIGN_AcccountID and or
// DOCUSIGN_JWTConfig and DOCUSIGN_JWTAPIUser.  If neither of these
// variables are set, skip the test.
func TestGeneratedOps(t *testing.T) {
	cred, err := getLocalCredential()
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if cred == nil {
		t.Skip()
	}
	ctx := context.Background()
	// Read throught  all folders
	sv := folders.New(cred)
	l, err := sv.List().Do(ctx)
	if err != nil {
		t.Errorf("List: %v", err)
		return
	}
	if len(l.Folders) < 1 {
		t.Errorf("expecting multiple folders")
	}

	// Read through all templates
	svT := templates.New(cred)
	tList, err := svT.List().Do(context.Background())
	if err != nil {
		t.Errorf("Template List: %v", err)
	}

	for _, tmpl := range tList.EnvelopeTemplates {
		t.Logf("Getting: %s", tmpl.TemplateID)
		tx, err := svT.Get(tmpl.TemplateID).Include("recipients").Do(context.Background())
		if err != nil {
			t.Errorf("unable to open template %s: %v", tmpl.Name, err)
			continue
		}
		t.Logf("Got: %s", tx.TemplateID)
	}
}

func getLocalCredential() (*esign.OAuth2Credential, error) {
	if tk, ok := os.LookupEnv("DOCUSIGN_Token"); ok {
		acctID, _ := os.LookupEnv("DOCUSIGN_AccountID")
		return esign.TokenCredential(tk, true).WithAccountID(acctID), nil
	}

	if jwtConfigJSON, ok := os.LookupEnv("DOCUSIGN_JWTConfig"); ok {
		jwtAPIUserName, ok := os.LookupEnv("DOCUSIGN_JWTAPIUser")
		if !ok {
			return nil, fmt.Errorf("expected DOCUSIGN_JWTAPIUser environment variable with DOCUSIGN_JWTConfig=%s", jwtConfigJSON)
		}

		buffer, err := ioutil.ReadFile(jwtConfigJSON)
		if err != nil {
			return nil, fmt.Errorf("%s open: %v", jwtConfigJSON, err)
		}
		var cfg *esign.JWTConfig
		if err = json.Unmarshal(buffer, &cfg); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		return cfg.Credential(jwtAPIUserName, nil, nil)
	}
	return nil, nil
}
