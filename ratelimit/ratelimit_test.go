// Copyright 2022 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ratelimit_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jfcote87/ctxclient"
	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/ratelimit"
)

type RateLimitServer struct {
	u         string
	rdr       *bytes.Reader
	m         sync.Mutex
	RLR       ratelimit.Report
	resetTime int64
}

func (rls *RateLimitServer) Handle(ctx context.Context, res *http.Response) error {
	if rls == nil {
		return errors.New("nil handler")
	}
	rpt := ratelimit.New(res.Header)
	if rpt.IsEmpty() {
		return errors.New("empty report")
	}
	rls.m.Lock()
	rls.RLR = *rpt
	rls.m.Unlock()
	return nil
}

func (rls *RateLimitServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("x-msg") == "noheader" {
		_, _ = w.Write([]byte("{}"))
		return
	}
	if r.Header.Get("x-msg") == "404" {
		http.Error(w, "path not found", 404)
		return
	}
	tm := time.Now().Add(time.Hour)
	rls.resetTime = tm.Unix()
	w.Header().Set("X-RateLimit-Limit", "30000")
	w.Header().Set("X-RateLimit-Remaining", "1")
	w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(tm.Unix(), 10))
	w.Header().Set("X-BurstLimit-Limit", "100")
	w.Header().Set("X-BurstLimit-Remaining", "32")
}

var xmsgHeaderKey = &RateLimitServer{}

func (rls *RateLimitServer) AuthDo(ctx context.Context,
	op *esign.Op) (*http.Response, error) {

	r, _ := http.NewRequest(op.Method, rls.u, rls)
	if xmsg, ok := ctx.Value(xmsgHeaderKey).(string); ok {
		r.Header.Set("x-msg", xmsg)
	}
	return (ctxclient.Func)(nil).Do(ctx, r)
}

func (rls *RateLimitServer) Close() error {
	return nil
}

func (rls *RateLimitServer) Read(b []byte) (int, error) {
	return rls.rdr.Read(b)
}

func TestRateLimitCredential_AuthDo(t *testing.T) {
	rls := &RateLimitServer{rdr: bytes.NewReader([]byte(""))}
	ts := httptest.NewServer(rls)
	rls.u = ts.URL

	stndCred := &ratelimit.Credential{
		Credential:    rls,
		ReportHandler: rls,
	}
	tests := []struct {
		name       string
		credential *ratelimit.Credential
		msg        string
		wantErr    bool
		errmsg     string
		want       *ratelimit.Report
	}{
		{name: "test00", msg: "", wantErr: true, errmsg: "ratelimitcredential no child credential specified"},
		{name: "test01", credential: stndCred, msg: "404", wantErr: true, errmsg: "response returned 404 404 Not Found: path not found"},
		{name: "test02", credential: stndCred, msg: "noheader", wantErr: true, errmsg: "empty report"},
		{name: "test03", credential: stndCred, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.msg != "" {
				ctx = context.WithValue(ctx, xmsgHeaderKey, tt.msg)
			}
			op := &esign.Op{
				Method: "POST",
				Path:   "/a",
			}
			res, err := tt.credential.AuthDo(ctx, op)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("%s expected success; got %v", tt.name, err)
					return
				}
				if !strings.HasPrefix(err.Error(), tt.errmsg) {
					t.Errorf("%sX", tt.errmsg)
					t.Errorf("%sX", err.Error())
					t.Errorf("%s expected %s; got %v", tt.name, tt.errmsg, err.Error())
				}
				return
			}
			if tt.wantErr {
				t.Errorf("%s expected %s; got success", tt.name, tt.errmsg)
				return
			}
			res.Body.Close()

			if rls.RLR.RateLimit != 30000 ||
				rls.RLR.RateRemaining != 1 ||
				rls.RLR.BurstLimit != 100 ||
				rls.RLR.BurstRemaining != 32 ||
				rls.RLR.ResetAt().Unix() != rls.resetTime {
				t.Errorf("expected 30000, 1, 100, 32, %s; got %d %d %d %d %s",
					time.Unix(rls.resetTime, 0).Format("15:04:05"), rls.RLR.RateLimit, rls.RLR.RateRemaining,
					rls.RLR.BurstLimit, rls.RLR.BurstRemaining,
					rls.RLR.ResetAt().Format("15:04:05"))
			}
		})
	}
}
