// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package click implements DocuSign's click api. Api documentation
// for these functions may be found at:
// https://docs.docusign.com/esign/guide/authentication/legacy_auth.htmlpackage
package click // import "github.com/jfcote87/esign/click"

import (
	"context"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jfcote87/esign"
)

var clickV1 = &esign.APIVersion{
	Prefix:  "clickapi",
	Version: "v1",
}

// Service implements DocuSign Clickwrap API operations
type Service struct {
	credential esign.Credential
}

// New initializes a accounts service using cred to authorize ops.
func New(cred esign.Credential) *Service {
	return &Service{credential: cred}
}

// CreateOp creates a clickwrap for an account
type CreateOp esign.Op

// Create builds a CreateOp using the passed clickwrap
func (s *Service) Create(clickwrap *Clickwrap) *CreateOp {
	return &CreateOp{
		Credential: s.credential,
		Method:     "POST",
		Path:       "clickwraps",
		Payload:    clickwrap,
		QueryOpts:  make(url.Values),
		Version:    clickV1,
	}
}

// Do executes the operation
func (op *CreateOp) Do(ctx context.Context) (*Summary, error) {
	var res *Summary
	return res, ((*esign.Op)(op)).Do(ctx, &res)
}

// UpdateOp updates an existing clickwrap and creates a new
// version.
type UpdateOp esign.Op

// Update builds an UpdateOp to create a new version
func (s *Service) Update(id string, clickwrap *Clickwrap) *UpdateOp {
	return &UpdateOp{
		Credential: s.credential,

		Method:    "POST",
		Path:      strings.Join([]string{"clickwraps", id, "versions"}, "/"),
		Payload:   clickwrap,
		QueryOpts: make(url.Values),
		Version:   clickV1,
	}
}

// Do executes the operation
func (op *UpdateOp) Do(ctx context.Context) (*Summary, error) {
	var res *Summary
	return res, ((*esign.Op)(op)).Do(ctx, &res)
}

// ListOp gets all the clickwraps for an account
type ListOp esign.Op

// List builds the ListOp
func (s *Service) List() *ListOp {
	return &ListOp{
		Credential: s.credential,
		Method:     "GET",
		Path:       "clickwraps",
		QueryOpts:  make(url.Values),
		Version:    clickV1,
	}
}

// FromDate (optional) sets the date from which created
// clickwaps will be returned (optional)
func (op *ListOp) FromDate(val time.Time) *ListOp {
	if op != nil {
		op.QueryOpts.Set("from_date", val.Format(time.RFC3339))
	}
	return op
}

// ToDate sets the date  up to which created clickwraps
// will be returned  (optional)
func (op *ListOp) ToDate(val time.Time) *ListOp {
	if op != nil {
		op.QueryOpts.Set("to_date", val.Format(time.RFC3339))
	}
	return op
}

// Name sets the name of clickwraps to return (optional)
func (op *ListOp) Name(val string) *ListOp {
	if op != nil {
		op.QueryOpts.Set("name", val)
	}
	return op
}

// Page sets the page number of the clickwraps to return
func (op *ListOp) Page(val int32) *ListOp {
	if op != nil {
		op.QueryOpts.Set("page_number", strconv.Itoa(int(val)))
	}
	return op
}

// Status filters the statuses of the clickwrapts to return.
// Valid values are active, inactive and deleted.
func (op *ListOp) Status(vals ...string) *ListOp {
	if op != nil {
		for _, val := range vals {
			op.QueryOpts.Add("status", val)
		}
	}
	return op
}

// VersionNumber sets the version number of the clickwraps to return
func (op *ListOp) VersionNumber(val int32) *ListOp {
	if op != nil {
		op.QueryOpts.Set("version_number", strconv.Itoa(int(val)))
	}
	return op
}

// Do executes the operation
func (op *ListOp) Do(ctx context.Context) (*Listing, error) {
	var res *Listing
	return res, ((*esign.Op)(op)).Do(ctx, &res)
}

// SetAgreementOp creates/update as user agreement for a client
type SetAgreementOp esign.Op

// SetAgreement builds a SetAgreementOp for creating/updating
func (s *Service) SetAgreement(clickwrapID string, agreement *UserAgreement) *SetAgreementOp {
	return &SetAgreementOp{
		Credential: s.credential,
		Method:     "POST",
		Path:       strings.Join([]string{"clickwraps", clickwrapID, "versions"}, "/"),
		Payload:    agreement,
		QueryOpts:  make(url.Values),
		Version:    clickV1,
	}
}

// Do executes the operation
func (op *SetAgreementOp) Do(ctx context.Context) (*UserAgreement, error) {
	var res *UserAgreement
	return res, ((*esign.Op)(op)).Do(ctx, &res)
}

// GetAgreementsOp lists user agreements for a specific clickwrap
type GetAgreementsOp esign.Op

// GetAgreements builds a GetAgreementsOp for the passed clickwrapID
func (s *Service) GetAgreements(clickwrapID string) *GetAgreementsOp {
	return &GetAgreementsOp{
		Credential: s.credential,
		Method:     "GET",
		Path:       strings.Join([]string{"clickwrap", clickwrapID, "users"}, "/"),
		QueryOpts:  make(url.Values),
		Version:    clickV1,
	}
}

// ClientUserID adds a filter to limit agreements for the client
func (op *GetAgreementsOp) ClientUserID(val string) *GetAgreementsOp {
	if op != nil {
		op.QueryOpts.Set("client_user_id", val)
	}
	return op
}

// Page filters agreements to those with the passed page number
func (op *GetAgreementsOp) Page(val int32) *GetAgreementsOp {
	if op != nil {
		op.QueryOpts.Set("page_number", strconv.Itoa(int(val)))
	}
	return op
}

// Do executes the operation
func (op *GetAgreementsOp) Do(ctx context.Context) (*AgreementList, error) {
	var res *AgreementList
	return res, ((*esign.Op)(op)).Do(ctx, &res)
}
