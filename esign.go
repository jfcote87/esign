// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esign

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"net/http"
	"net/textproto"
	"net/url"
)

// ErrNilOp used to indicate a nil operation pointer
var ErrNilOp = errors.New("nil operation")

// APIv21  indicates that the url will resolve
// to /restapi/v2.1
var APIv21 APIVersion = &apiVersion{
	name:           "DocuSign REST API:v2.1",
	prefix:         "/restapi",
	accountReplace: true,
	versionPrefix:  "/v2.1",
	demoHost:       "demo.docusign.net",
}

// APIv2  indicates that the url will resolve
// to /restapi/v2
var APIv2 APIVersion = &apiVersion{
	name:           "DocuSign REST API:v2",
	prefix:         "/restapi",
	accountReplace: true,
	versionPrefix:  "/v2",
	demoHost:       "demo.docusign.net",
}

// AdminV2 handles calls for the admin api and urls will
// resolve to start with /management
var AdminV2 APIVersion = &apiVersion{
	name:     "DocuSign Admin API:v2.1",
	prefix:   "/Management",
	host:     "api.docusign.net",
	demoHost: "api-d.docusign.net",
}

// RoomsV2 resolves urls for monitor dataset calls
var RoomsV2 APIVersion = &apiVersion{
	name:           "DocuSign Rooms API - v2:v2",
	prefix:         "/restapi",
	accountReplace: true,
	versionPrefix:  "/v2",
	host:           "rooms.docusign.com",
	demoHost:       "demo.rooms.docusign.com",
}

// MonitorV2 resolves urls for monitor dataset calls
var MonitorV2 APIVersion = &apiVersion{
	name:     "Monitor API:v2.0",
	prefix:   "",
	host:     "lens.docusign.net",
	demoHost: "lens-d.docusign.net",
}

// ClickV1 defines url replacement for clickraps api
var ClickV1 APIVersion = &apiVersion{
	name:           "DocuSign Click API:v1",
	prefix:         "/clickapi",
	versionPrefix:  "/v1",
	accountReplace: true,
	demoHost:       "demo.docusign.net",
}

type apiVersion struct {
	name           string
	prefix         string
	host           string
	demoHost       string
	accountReplace bool
	versionPrefix  string
}

func (v *apiVersion) Name() string {
	return v.name
}

// APIVersion resolves the final op url by completing the partial path and host properties
// of u and returning a new URL.
type APIVersion interface {
	ResolveDSURL(u *url.URL, host string, accountID string, isDemo bool) *url.URL
	Name() string
}

// ResolveAPIHost determines the url's host based upon the version
func (v *apiVersion) resolveAPIHost(credentialHost string, isDemo bool) string {
	if isDemo {
		return v.demoHost
	}
	if v.host != "" {
		return v.host
	}
	return credentialHost
}

// ResolveDSURL updates the passed *url.URL's settings.
// https://developers.docusign.com/esign-rest-api/guides/authentication/user-info-endpoints#form-your-base-path
func (v *apiVersion) ResolveDSURL(u *url.URL, host string, accountID string, isDemo bool) *url.URL {
	if v == nil {
		return u
	}
	newURL := *u
	newURL.Scheme = "https"
	newURL.Host = v.resolveAPIHost(host, isDemo)

	if v.accountReplace && !strings.HasPrefix(u.Path, "/") {
		newURL.Path = v.prefix + v.versionPrefix + "/accounts/" + accountID + "/" + u.Path
		return &newURL
	}
	newURL.Path = v.prefix + u.Path
	return &newURL
}

// Credential adds an authorization header(s) for the http request,
// resolves the http client and finalizes the url.  Credentials may
// be created using the Oauth2Config and JWTConfig structs as well as
// legacy.Config.
type Credential interface {
	// AuthDo attaches an authorization header to a request, prepends
	// account and user ids to url, and sends request.
	AuthDo(context.Context, *Op) (*http.Response, error)
}

// Op contains all needed information to perform a DocuSign operation.
// Used in the sub packages, and may be used for testing or creating
// new/corrected operations.
type Op struct {
	// Used for authorization and for URL completion
	Credential Credential
	// POST,GET,PUT,DELETE
	Method string
	// If not prefixed with "/", Credential will prepend the accountId
	// /restapi/v2/accounts/{accountId}
	Path string
	// Payload will be marshalled into the request body
	Payload interface{}
	// Additional query parameters
	QueryOpts url.Values
	// Upload files for document creation
	Files []*UploadFile
	// Accept header value (usually json/application)
	Accept string
	// Leave nil for v2
	Version APIVersion
}

// ResponseError describes DocuSign's server error response.
// https://developers.docusign.com/esign-rest-api/guides/status-and-error-codes#general-error-response-handling
type ResponseError struct {
	ErrorCode   string `json:"errorCode,omitempty"`
	Description string `json:"message,omitempty"`
	Status      int    `json:"-"`
	Raw         []byte `json:"-"`
	OriginalErr error  `json:"-"`
}

// Error fulfills error interface
func (r ResponseError) Error() string {
	return fmt.Sprintf("Status: %d  %s: %s", r.Status, r.ErrorCode, r.Description)
}

// NewResponseError unmarshals buff, containing a DocuSign server error,
// into a ResponseError
func NewResponseError(buff []byte, status int) *ResponseError {
	re := ResponseError{
		Status: status,
		Raw:    buff,
	}
	_ = json.Unmarshal(buff, &re)
	return &re
}

// Body creates an io.Reader marshalling the payload in the appropriate
// format and if files are available create a multipart form.
func (op *Op) Body() (io.Reader, string, error) {
	var body io.Reader
	var ct string
	switch p := op.Payload.(type) {
	case *UploadFile:
		return p.Reader, p.ContentType, nil
	case url.Values:
		body, ct = bytes.NewBufferString(p.Encode()), "application/x-www-form-urlencoded"
	case interface{}: // non-nil
		body, ct = bytes.NewBuffer(nil), "application/json"
		if err := json.NewEncoder(body.(io.Writer)).Encode(p); err != nil {
			return nil, "", err
		}
	}
	if len(op.Files) > 0 {
		var files = op.Files
		if body != nil {
			files = append([]*UploadFile{{Reader: body, ContentType: ct}}, op.Files...)
		}
		body, ct = multiPartBody(files)
	}
	return body, ct, nil
}

// CreateRequest prepares an http.Request and optionally logs the request body.
// UploadFiles will be closed on error.
func (op *Op) CreateRequest() (*http.Request, error) {
	body, ct, err := op.Body()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(op.Method, op.Path, body)
	if err != nil {
		if b, ok := body.(io.ReadCloser); ok {
			b.Close()
		}
		return nil, err
	}
	if len(op.QueryOpts) > 0 {
		req.URL.RawQuery = op.QueryOpts.Encode()
	}
	if body != nil {
		req.Header.Set("Content-Type", ct)
	}
	if op.Accept > "" {
		req.Header.Set("Accept", op.Accept)
	}
	return req, nil
}

func (op *Op) closeFiles() {
	if op != nil {
		for _, f := range op.Files {
			f.Close()
		}
		if f, ok := op.Payload.(*UploadFile); ok {
			f.Close()
		}
	}
}

// validate runs nil checks and returns an http.Client for the op
func (op *Op) validate(ctx context.Context) error {
	if op == nil {
		return errors.New("nil op")
	}
	var err error
	if ctx == nil {
		err = errors.New("nil context")
	}
	if err == nil && op.Credential == nil {
		err = errors.New("nil credential")
	}
	if err == nil {
		for _, f := range op.Files {
			if !f.Valid() {
				err = fmt.Errorf("invalid upload file %v", f)
				break
			}
		}
	}
	if err != nil {
		op.closeFiles()
	}
	return err
}

// Do sends a request to DocuSign.  Response data is decoded into
// result.  If result is a **Download, do sets the File.ReadCloser
// to the *http.Response.  The developer is responsible for closing
// the Download.ReadCloser.  Any non-2xx status code is returned as a
// *ResponseError.
func (op *Op) Do(ctx context.Context, result interface{}) error {
	// do nil checks and get client
	if err := op.validate(ctx); err != nil {
		return err
	}
	res, err := op.Credential.AuthDo(ctx, op)
	if err != nil {
		return err
	}
	// pass res.Body back if download
	if f, ok := result.(**Download); ok {
		*f = &Download{
			ReadCloser:         res.Body,
			ContentLength:      res.ContentLength,
			ContentType:        res.Header.Get("Content-Type"),
			ContentDisposition: res.Header.Get("Content-Disposition"),
		}
		return nil
	}
	defer res.Body.Close()
	if result != nil {
		return json.NewDecoder(res.Body).Decode(result)
	}
	return nil
}

// multiPartBody sends files thru a multipart writer. Using io.Pipe
// so we're not copying files into memory.
// https://developers.docusign.com/esign-rest-api/guides/requests-and-responses#multipart-form-requests
func multiPartBody(files []*UploadFile) (io.Reader, string) {
	pr, pw := io.Pipe()
	mpw := multipart.NewWriter(pw)
	go func() {
		var ptw io.Writer
		var err error
		// copy each file to multipart writer
		for _, f := range files {
			if err == nil {
				contentDisp := "form-data"
				if f.ID > "" {
					contentDisp = fmt.Sprintf("file; filename=\"%s\";documentid=%s", url.PathEscape(f.FileName), f.ID)
				}
				mh := textproto.MIMEHeader{
					"Content-Type":        []string{f.ContentType},
					"Content-Disposition": []string{contentDisp},
				}
				if ptw, err = mpw.CreatePart(mh); err == nil {
					_, err = io.Copy(ptw, f)
				}
			}
			f.Close()
		}
		if err == nil {
			mpw.Close()
		}
		pw.CloseWithError(err)
	}()
	return pr, "multipart/form-data; boundary=" + mpw.Boundary()
}

// Download is used to return image and pdf files from DocuSign. The developer
// needs to ensure to close when finished reading.
type Download struct {
	io.ReadCloser
	// ContentLength from response
	ContentLength int64
	// ContentType header value
	ContentType string
	// ContentDisposition header value
	ContentDisposition string
}

// UploadFile describes an a document attachment for uploading.
type UploadFile struct {
	// mime type of content
	ContentType string
	// file name to display in envelope or to identify signature
	FileName string
	// envelope documentId
	ID string
	// reader for creating file
	io.Reader
}

// Close closes the io.Reader if an io.Closer.
func (uf *UploadFile) Close() {
	if uf != nil {
		if closer, ok := uf.Reader.(io.Closer); ok {
			closer.Close()
		}
	}
}

// Valid ensures UploadFile.Reader is not nil.
func (uf *UploadFile) Valid() bool {
	return uf != nil && uf.Reader != nil
}
