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

// VersionV21  indicates that the url will resolve
// to /restapi/v2.1
var VersionV21 = &APIVersion{Version: "v2.1"}

// ClickV1 defines url replacement for clickraps api
var ClickV1 = &APIVersion{
	Version: "v1",
	Prefix:  "clickapi",
}

// APIVersion defines the prefix used to resolve an operation's url.  If
// nil or blank, "v2" is assumed.
type APIVersion struct {
	Version string
	Prefix  string
}

// ResolveDSURL updates the passed *url.URL's settings.
// https://developers.docusign.com/esign-rest-api/guides/authentication/user-info-endpoints#form-your-base-path
func (v *APIVersion) ResolveDSURL(u *url.URL, host string, accountID string) *url.URL {
	newURL := *u
	newURL.Scheme = "https"
	newURL.Host = host
	var prefix = "/restapi"
	var version = "/v2" // default is v2 restapi
	if v != nil {
		if v.Prefix != "" {
			prefix = "/" + v.Prefix
		}
		if v.Version != "" {
			version = "/" + v.Version
		}
	}
	if strings.HasPrefix(u.Path, "/") {
		newURL.Path = prefix + u.Path
	} else {
		newURL.Path = prefix + version + "/accounts/" + accountID + "/" + u.Path
	}
	return &newURL
}

// Credential adds an authorization header(s) for the http request,
// resolves the http client and finalizes the url.  Credentials may
// be created using the Oauth2Config and JWTConfig structs as well as
// legacy.Config.
type Credential interface {
	// AuthDo attaches an authorization header to a request, prepends
	// account and user ids to url, and sends request.  This func must
	// always close the request Body.
	AuthDo(context.Context, *http.Request, *APIVersion) (*http.Response, error)
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
	// Set Accept to a mimeType if response will
	// not be application/json
	Accept string
	// Leave nil for v2
	Version *APIVersion
}

//type requestHandler interface {
//	Do(context.Context, *http.Request) (*http.Response, error)
//}

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
	json.Unmarshal(buff, &re)
	return &re
}

func getBodyFromPayload(payload interface{}, files []*UploadFile) (io.Reader, string, error) {
	var body io.Reader
	var ct string
	switch p := payload.(type) {
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
	if len(files) > 0 {
		if body != nil {
			files = append([]*UploadFile{{Reader: body, ContentType: ct}}, files...)
		}
		body, ct = multiPartBody(files)
	}
	return body, ct, nil
}

// createOpRequest prepares an http.Request and optionally logs the request body.
// UploadFiles will be closed on error.
func (op *Op) createOpRequest(ctx context.Context, accept string) (*http.Request, error) {

	body, ct, err := getBodyFromPayload(op.Payload, op.Files)
	if err != nil {
		op.closeFiles() // close any open files on error
		return nil, err
	}

	req, err := http.NewRequest(op.Method, op.Path, body)
	if err != nil {
		// close body
		if f, ok := body.(io.Closer); ok {
			f.Close()
		}
		return nil, err
	}
	if len(op.QueryOpts) > 0 {
		req.URL.RawQuery = op.QueryOpts.Encode()
	}
	if len(ct) > 0 {
		req.Header.Set("Content-Type", ct)
	}
	if len(accept) > 0 {
		req.Header.Set("Accept", accept)
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
				err = fmt.Errorf("Invalid upload file %v", f)
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

	acceptHdr := op.Accept
	if acceptHdr == "" {
		switch result.(type) {
		case **Download: // no accept header if **Download or nil
		case interface{}:
			acceptHdr = "application/json"
		}
	}

	// get request
	req, err := op.createOpRequest(ctx, acceptHdr)
	if err != nil {
		return err
	}

	res, err := op.Credential.AuthDo(ctx, req, op.Version)
	if err != nil {
		return err
	}

	switch f := result.(type) {
	case **Download: // return w/o closing response body
		*f = &Download{res.Body, res.ContentLength, res.Header.Get("Content-Type")}
		return nil
	case interface{}: // non-nil
		// parse response and check for context cancellation.
		done := make(chan error, 1) // buffered channel so go routine doesn't hang
		go func() {
			done <- json.NewDecoder(res.Body).Decode(result)
		}()
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-done:
		}
	}
	res.Body.Close()
	return err

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
		return
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
