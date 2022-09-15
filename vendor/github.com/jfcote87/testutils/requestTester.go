// Copyright 2019 James F Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testutils contains routines to help with
// creating tests
package testutils // import "github.com/jfcote87/testutils"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

// RequestTester contains expected values of the request
// and returns either the passed Response or the
// results of the ResponseFunc. If both are nil,
// an empty 200 OK response is returned.
type RequestTester struct {
	Path        string      // checks req.URL.Path
	Auth        string      // req.Header.Get("Authorization")
	Method      string      // req.Method
	Query       string      // req.URL.RawQuery
	Host        string      // req.URL.Host
	ContentType string      // req.Header.Get("Content-type")
	Header      http.Header // req.Header
	Payload     []byte      // req.Body
	// returned response if all tests pass and
	// ResponseFunc is nil. You may use the MakeResponse
	// function to create an *http.Response.
	Response *http.Response
	// ResponseFunc used to add custom tests of the Request.
	// No need to close request body as the RoundTripper function
	// handles this.
	ResponseFunc func(*http.Request) (*http.Response, error)
}

// Check compares expected values with the req parameter
func (r RequestTester) Check(req *http.Request) error {
	if r.Path > "" && r.Path != req.URL.Path {
		return fmt.Errorf("expected request path %s; got %s", r.Path, req.URL.Path)
	}
	if r.Auth > "" && r.Auth != req.Header.Get("Authorization") {
		return fmt.Errorf("expecte auth header %s; got %s", r.Auth, req.Header.Get("Authorization"))
	}
	if r.Method > "" && r.Method != req.Method {
		return fmt.Errorf("expected method %s; got %s", r.Method, req.Method)
	}
	if r.Query > "" && r.Query != req.URL.RawQuery {
		return fmt.Errorf("expected query args %s; got %s", r.Query, req.URL.RawQuery)
	}
	if r.Host > "" && r.Host != req.URL.Host {
		return fmt.Errorf("expected host %s; got %s", r.Host, req.URL.Host)
	}
	if r.ContentType > "" && r.ContentType != req.Header.Get("ContentType") {
		return fmt.Errorf("expected content-type %s; got %s", r.ContentType, req.Header.Get("ContentType"))
	}
	for k := range r.Header {
		if r.Header.Get(k) != req.Header.Get(k) {
			return fmt.Errorf("expected header %s = %s; got %s", k, r.Header.Get(k), req.Header.Get(k))
		}
	}
	if len(r.Payload) > 0 {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("unable to read request body: %v", err)
		}
		if bytes.Compare(b, r.Payload) != 0 {
			return fmt.Errorf("expected body %s; got %s", string(r.Payload), string(b))
		}

	}
	return nil
}

// Transport contains an array of http.Response, handler funcs and errors
// that help create an http request test
type Transport struct {
	Queue []*RequestTester
}

// RoundTrip fulfills the http.Transport interface{} by creating
// http responses and errors from the Transport queue
func (tx *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	defer func() {
		if req.Body != nil {
			req.Body.Close()
		}
	}()
	var t *RequestTester
	if len(tx.Queue) > 0 {
		defer func() {
			tx.Queue[0] = nil
			tx.Queue = tx.Queue[1:]
		}()
		t = tx.Queue[0]
	}
	if err := t.Check(req); err != nil {
		return nil, err
	}
	if t.ResponseFunc != nil {
		return t.ResponseFunc(req)
	}
	if t.Response == nil {
		return MakeResponse(200, nil, nil), nil
	}
	return t.Response, nil

}

// Add adds a new response to the queue.
func (tx *Transport) Add(val ...*RequestTester) {
	tx.Queue = append(tx.Queue, val...)
	return
}

// MakeResponse creates an *http.Response for later processing
func MakeResponse(status int, body []byte, header http.Header) *http.Response {
	if header == nil {
		header = make(http.Header)
	}
	return &http.Response{
		Body:          ioutil.NopCloser(bytes.NewReader(body)),
		StatusCode:    status,
		Status:        http.StatusText(status),
		ContentLength: int64(len(body)),
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
	}
}
