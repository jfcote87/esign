// Copyright 2022 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ratelimit provides tools for reporting on DocuSign's rate limits.  Documentation
// may be found at https://developers.docusign.com/docs/esign-soap-api/esign101/security/call-limits/
package ratelimit

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/jfcote87/esign"
)

type reportCtxKey struct{}

// ReportPtrContextKey is used to store a **Report pointer into a
// context to retrieve a Report for a single call.  When used in a
// context passing through a RateLimit.Credential, the most recent report
// is set to the **Report pointer.
// See ExampleCredential_context() in example.go
var ReportPtrContextKey *reportCtxKey

// Report contains the rate limit values after an
// API request
type Report struct {
	RateLimit      int64
	RateRemaining  int64
	rateReset      int64
	BurstLimit     int64
	BurstRemaining int64
}

// ResetAt returns the time when the rate remaining
// will reset to the the initial value
func (r Report) ResetAt() time.Time {
	return time.Unix(r.rateReset, 0)
}

// IsEmpty returns true if all fields are empty suggesting that docusign
// returned no rate limit data
func (r Report) IsEmpty() bool {
	return r.RateLimit == 0 && r.RateRemaining == 0 &&
		r.rateReset == 0 && r.BurstLimit == 0 && r.BurstRemaining == 0
}

// ReportHandler is used by the ratelimit.Credential store, analyze or process the
// new report value.  The Handle method should be thread safe as credentials are designed
// to be concurrent.
type ReportHandler interface {
	Handle(context.Context, *http.Response) error
}

// Credential allows for the analyzing of rate limit reports.  Simply
// set the Credential to an existing credential and the Handler.
type Credential struct {
	esign.Credential
	ReportHandler
}

// AuthDo authorizes the request and creates a Report from the response.  The Handler
// then processes the Report before returning the response.
func (rlc *Credential) AuthDo(ctx context.Context, op *esign.Op) (*http.Response, error) {
	if rlc == nil || rlc.Credential == nil {
		return nil, fmt.Errorf("ratelimitcredential no child credential specified")
	}
	res, err := rlc.Credential.AuthDo(ctx, op)
	if err != nil {
		return nil, err
	}
	if ptr, ok := ctx.Value(ReportPtrContextKey).(**Report); ok {
		*ptr = New(res.Header)
	}

	if rlc.ReportHandler != nil {
		if err := rlc.Handle(ctx, res); err != nil {
			res.Body.Close()
			return nil, err
		}
	}
	return res, nil
}

// New creates a Report from an http.Header.  See docusign documentation
// for explanation of header values.
func New(hdr http.Header) *Report {
	var keys = []string{
		"X-RateLimit-Limit",
		"X-RateLimit-Remaining",
		"X-RateLimit-Reset",
		"X-BurstLimit-Limit",
		"X-BurstLimit-Remaining",
	}
	var rpt Report
	var rptFields = []*int64{
		&rpt.RateLimit,
		&rpt.RateRemaining,
		&rpt.rateReset,
		&rpt.BurstLimit,
		&rpt.BurstRemaining,
	}
	for i, k := range keys {
		*rptFields[i], _ = strconv.ParseInt(hdr.Get(k), 10, 64)
	}
	return &rpt
}
