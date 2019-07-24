// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package click

import (
	"time"
)

// Listing contains the response to a clickwrap list method
type Listing struct {
	Clickwraps []Summary `json:"clickwraps,omitempty"`
}

// Summary provides summary data for a created/updated clickwrap
type Summary struct {
	ID          string     `json:"clickwrapId,omitempty"`
	Name        string     `json:"clickwrapName,omitempty"`
	CreatedTime *time.Time `json:"createdTime,omitempty"`
	// uuid of user that last modified clickwrap
	LastModifiedBy      string `json:"lastModifiedBy,omitempty"`
	RequireReacceptance bool   `json:"requireReacceptance,omitempty"`
	// Status may be active, inactive or deleted
	Status        string `json:"status,omitempty"`
	VersionNumber int32  `json:"versionNumber,omitempty"`
}

// Clickwrap fully describes a clickwrap version
type Clickwrap struct {
	ID                  string `json:"clickwrapId,omitempty"`
	Name                string `json:"name,omitempty"`
	RequireReacceptance bool   `json:"requireReacceptance,omitempty"`
	// Status may be active, inactive or deleted
	Status          string           `json:"status,omitempty"`
	UserID          string           `json:"userId,omitempty"`
	DisplaySettings *DisplaySettings `json:"displaySettings,omitempty"`
	Documents       []Document       `json:"documents,omitempty"`
}

// Settings are used by a  user agreement for determing how to display
// a clickwrap document
type Settings struct {
	DocumentDisplay string `json:"documentDisplay,omitempty"`
	Format          string `json:"format,omitempty"`
	HostOrigin      string `json:"hostOrigin,omitempty"`
	// "small", "medium" or "large"
	Size string `json:"size,omitempty"`
}

// DisplaySettings provide defaults for the clickwrap
type DisplaySettings struct {
	BrandID           string `json:"brandId,omitempty"`
	ConsentButtonText string `json:"consentButtonText,omitempty"`
	DeclineButtonText string `json:"declineButtonText,omitempty"`
	DisplayName       string `json:"displayName,omitempty"`
	DocumentDisplay   string `json:"documentDisplay,omitempty"`
	Downloadable      bool   `json:"bool,omitempty"`
	Format            string `json:"format,omitempty"`
	HasAccept         bool   `json:"hasAccept,omitempty"`
	HostOrigin        string `json:"hostOrigin,omitempty"`
	MustRead          bool   `json:"mustRead,omitempty"`
	MustView          bool   `json:"mustView,omitempty"`
	RequireAccept     bool   `json:"requireAccept,omitempty"`
	SendToEmail       bool   `json:"sendToEmail,omitempty"`
	Size              string `json:"size,omitempty"`
}

// Document used to update/create the original document for a clickwrap
type Document struct {
	Base64  []byte `json:"documentBase64,omitempty"`
	Name    string `json:"documentName,omitempty"`
	FileExt string `json:"fileExtension,omitempty"`
	Order   int32  `json:"order,omitempty"`
}

// UserAgreement describes a click from a clickwrap for a specific client
type UserAgreement struct {
	AccountID    string           `json:"accountId,omitempty"`
	AgreedOn     *time.Time       `json:"agreedOn,omitempty"`
	ID           string           `json:"agreementID,omitempty"`
	URL          string           `json:"agreementUrl,omitempty"`
	Clickwrap    *Clickwrap       `json:"clickwrap,omitempty"`
	ClickWrapID  string           `json:"clickwrapId,omitempty"`
	ClientUserID string           `json:"clientUserId,omitempty"`
	CreatedOn    *time.Time       `json:"createdOn,omitempty"`
	Settings     *DisplaySettings `json:"settings,omitempty"`
}

// AgreementList is DocuSign's response to a user agreement listing
type AgreementList struct {
	BeginCreatedOn        *time.Time      `json:"beginCreatedOn,omitempty"`
	MinimumPagesRemaining int32           `json:"minimumPagesRemaining,omitempty"`
	Page                  int32           `json:"page,omitempty"`
	PageSize              int32           `json:"pageSize,omitempty"`
	UserAgreements        []UserAgreement `json:"userAgreements,omitempty"`
}
