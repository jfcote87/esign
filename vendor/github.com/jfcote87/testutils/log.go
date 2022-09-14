// Copyright 2019 James Cote All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testutils

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// LogTransport provides a transport to log http requests.  If Base
// is nil then http.DefaultTransport is assumed to complete the round trip.LogTransport
// SaveFunc must be set for logging to occur.
type LogTransport struct {
	Base     http.RoundTripper
	SaveFunc func(ctx context.Context, rl *RequestLog)
}

// RoundTrip logs values from request and response
func (lt *LogTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rt := lt.Base
	if rt == nil {
		rt = http.DefaultTransport
	}
	if lt.SaveFunc == nil {
		return lt.RoundTrip(req)
	}

	rl, err := logRequest(req)
	defer lt.SaveFunc(req.Context(), rl)
	if err != nil {
		return nil, err
	}
	res, err := rt.RoundTrip(req)
	if err != nil {
		rl.RespErr = err
		return nil, err
	}
	if err = rl.logResponse(res); err != nil {
		return nil, err
	}
	return res, nil
}

// RequestLog contains values from the request and response
type RequestLog struct {
	URL        *url.URL
	Headers    http.Header
	Method     string
	Proto      string
	Body       []byte
	Error      error
	StatusCode int
	RespHeader http.Header
	RespBody   []byte
	RespErr    error
}

func logRequest(req *http.Request) (*RequestLog, error) {
	rlog := &RequestLog{
		URL:     req.URL,
		Headers: req.Header,
		Method:  req.Method,
		Proto:   req.Proto,
	}
	if req.Body != nil {
		body := &bytes.Buffer{}
		defer req.Body.Close()
		if _, err := io.Copy(body, req.Body); err != nil {
			rlog.Error = err
			return nil, err
		}
		rlog.Body = body.Bytes()
		req.Body = ioutil.NopCloser(body)
	}
	return rlog, nil
}

func (rl *RequestLog) logResponse(res *http.Response) error {
	rl.StatusCode = res.StatusCode
	rl.RespHeader = res.Header
	if res.Body != nil {
		body := &bytes.Buffer{}
		defer res.Body.Close()
		if _, err := io.Copy(body, res.Body); err != nil {
			rl.RespErr = err
			return err
		}
		rl.RespBody = body.Bytes()
		res.Body = ioutil.NopCloser(body)
	}
	return nil
}
