// Copyright 2017 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package esign implements a service to use the version 2 Docusign
// rest api. Api documentation may be found at:
// https://docs.docusign.com/esign
package esign

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"strings"
	"sync"

	"github.com/jfcote87/esign/model"

	"net/http"
	"net/textproto"
	"net/url"

	"github.com/jfcote87/ctxclient"
	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

// Credential adds an authorization header(s) for a rest http request
// and completes the call URL host and accountId.
type Credential interface {
	// Authorize attaches an authorization header to a request and
	// and fixes the URL to the appropriate host.
	Authorize(context.Context, *http.Request) error
}

// Call describes a DocuSign function
type Call struct {
	// Used for authorization and for URL completion
	Credential Credential
	// POST,GET,PUT,DELETE
	Method string
	// If not prefixed with "/", Credential will
	// prepend the accountId path (i.e. /restapi/v2/accounts/{accountId})
	Path string
	// struct  will be JSON encoded
	Payload interface{}
	// Replacement parameters for the url
	// (i.e. "envelopeId": "")
	PathParameters map[string]string
	// Additional url query values
	QueryOpts url.Values
	// Upload files
	Files []*UploadFile
}

func (c *Call) getPath() string {
	path := c.Path
	// TODO: find a library to do the replacement
	if len(c.PathParameters) > 0 {
		px := strings.Split(path, "/")
		for i, s := range px {
			if match, ok := c.PathParameters[s]; ok {
				px[i] = match
			}
		}
		path = strings.Join(px, "/")
	}
	if len(c.QueryOpts) > 0 {
		return path + "?" + c.QueryOpts.Encode()
	}
	return path
}

// ResponseError is generated when docusign returns an http error.
//
// https://docs.docusign.com/esign/guide/appendix/status_and_error_codes.html#general-error-response-handling
type ResponseError struct {
	Err         string `json:"errorCode,omitempty"`
	Description string `json:"message,omitempty"`
	Status      int    `json:"-"`
	Raw         []byte `json:"-"`
}

// Error fulfills error interface
func (r ResponseError) Error() string {
	return fmt.Sprintf("Status: %d  %s: %s", r.Status, r.Err, r.Description)
}

// checkResponseStatus looks at the response for a 200 or 201.  If not it will
// decode the json into a Response Error.  Returns nil on  success.  Response
// Body is closed on error.
// https://docs.docusign.com/esign/guide/appendix/status_and_error_codes.html#general-error-response-handling
func checkResponseStatus(res *http.Response) *ResponseError {
	statusCode := res.StatusCode
	if statusCode != 200 && statusCode != 201 {
		re := &ResponseError{Status: statusCode}
		if res.ContentLength > 0 {
			var err error
			if re.Raw, err = ioutil.ReadAll(res.Body); err == nil {
				err = json.Unmarshal(re.Raw, re)
			}
			if err != nil {
				re.Description = err.Error()
			}
		}
		res.Body.Close()
		return re
	}
	return nil
}

// Do executes the call.  Response data is encoded into
// the call's Result.  If Result is a **File, the File
// ReadCloser is set to the *http.Response which will need
// to be closed by the calling function.
func (c *Call) Do(ctx context.Context, result interface{}) error {
	if c == nil {
		return errors.New("nil Call")
	}
	if ctx == nil {
		closeUploads(c.Files)
		return errors.New("context may not be nil")
	}
	// define now so may be used by deferred log function
	var responseBytes []byte
	var res *http.Response
	var isLogged bool

	// serialize payload into body	var body io.Reader
	var body io.Reader
	var ct string
	if len(c.Files) > 0 {
		// formatted body for file upload
		var cancelFunc func()
		body, ct, cancelFunc = multiBody(c.Payload, c.Files) // no error, errors will occur during read
		defer cancelFunc()                                   // cancelfunc will close all io.ReadCloses is c
	} else if c.Payload != nil {
		switch payload := c.Payload.(type) {
		case url.Values:
			body = strings.NewReader(payload.Encode())
			ct = "application/x-www-form-urlencoded"
		default:
			bx, err := json.Marshal(c.Payload)
			if err != nil {
				return err
			}
			body, ct = bytes.NewReader(bx), "application/json"
		}
	}
	//
	req, err := http.NewRequest(c.Method, c.getPath(), body)
	if err != nil {
		return err
	}
	if len(ct) > 0 {
		req.Header.Set("Content-Type", ct)
	}
	// Check for a raw file return and set for json result
	file, _ := result.(**File)
	if result != nil && file == nil {
		req.Header.Set("accept", "application/json")
	}

	// authorize request
	if err = c.Credential.Authorize(ctx, req); err != nil {
		return err
	}

	credInterface := c.Credential.(interface{})
	// set logging
	if logger, ok := credInterface.(loggerFunc); ok {
		isLogged = true
		defer logger(ctx, req, res, c.Payload, responseBytes)
	}

	// select client func and perform call.  nil ctxclient.Resolver calls ctxclient.Default
	clientResolver, _ := credInterface.(ctxclient.Func)
	if res, err = ctxhttp.Do(ctx, clientResolver.Client(ctx), req); err == nil {
		return err
	}
	// res.Body close on error
	if err := checkResponseStatus(res); err != nil {
		responseBytes = err.Raw
		return err
	}
	// raw file - do not close response body
	if file != nil {
		*file = &File{res.Body, res.ContentLength, res.Header.Get("Content-Type")}
		return nil
	}
	if result != nil {
		body = res.Body
		if isLogged { // if logging save body for log
			if responseBytes, err = ioutil.ReadAll(body); err == nil {
				err = json.Unmarshal(responseBytes, result)
			}
		} else { // No logging
			err = json.NewDecoder(body).Decode(result)
		}
	}
	res.Body.Close()
	return err
}

// multiBody is used to format calls containing files as a multipart/form-data body.
// Send payload and files thru a multipart writer to format multipart/form-data.  Use
// Use io.Pipe so we're not copying files into memory.
func multiBody(payload interface{}, files []*UploadFile) (io.Reader, string, func()) {
	pr, pw := io.Pipe()
	mpw := multipart.NewWriter(pw)
	var once sync.Once
	var err error
	cancelFunc := func() {
		// Wrap in a once so this is may be called in calling routine
		once.Do(func() {
			if err != nil { // On err close pipe reader to create error in reading routine.
				pr.CloseWithError(fmt.Errorf("batch: multiPart Error: %v", err))
			}
			mpw.Close() // close writers will create error in goroutine
			pw.Close()
			closeUploads(files)
		})
	}
	go func() {
		var ptw io.Writer
		defer cancelFunc()

		// write json payload first
		if payload != nil {
			mh := textproto.MIMEHeader{
				"Content-Type":        []string{"application/json"},
				"Content-Disposition": []string{"form-data"},
			}
			if ptw, err = mpw.CreatePart(mh); err == nil {
				err = json.NewEncoder(ptw).Encode(payload)
			}
			if err != nil {
				return
			}
		}

		for _, f := range files {
			mh := textproto.MIMEHeader{
				"Content-Type":        []string{f.ContentType},
				"Content-Disposition": []string{fmt.Sprintf("file; filename=\"%s\";documentid=%s", f.FileName, f.ID)},
			}
			if ptw, err = mpw.CreatePart(mh); err == nil {
				if _, err = io.Copy(ptw, f.Data); err != nil {
					break
				}
			}
		}
		return
	}()
	return pr, "multipart/form-data; boundary=" + mpw.Boundary(), cancelFunc
}

// File contains the body of an http response. Used to return
// image and pdf files from DocuSign. The developer needs to
// call Close() when finished reading.
type File struct {
	io.ReadCloser
	// ContentLength from response
	ContentLength int64
	// ContentType header value
	ContentType string
}

type loggerFunc func(context.Context, *http.Request, *http.Response, interface{}, []byte)

func (l loggerFunc) log(ctx context.Context, req *http.Request, res *http.Response, payload interface{}, body []byte) {
	l(ctx, req, res, payload, body)
}

// WithLogger returns a credential that logs requests and responses to the
// logFunc function.
func WithLogger(credential Credential, logFunc func(ctx context.Context, req *http.Request, res *http.Response, payload interface{}, body []byte)) Credential {
	if logFunc == nil {
		return credential
	}
	if cx, ok := credential.(interface{}).(ctxclient.Func); ok {
		return struct {
			Credential
			ctxclient.Func
			loggerFunc
		}{credential, cx, logFunc}
	}
	return struct {
		Credential
		loggerFunc
	}{credential, logFunc}

}

// WithHTTPClientFunc specifies a function to determine the http.Client that
// calls will use.
func WithHTTPClientFunc(credential Credential, f func(context.Context) (*http.Client, error)) Credential {
	if f == nil {
		return credential
	}
	if l, ok := credential.(interface{}).(loggerFunc); ok {
		return struct {
			Credential
			loggerFunc
			ctxclient.Func
		}{credential, l, f}
	}
	return struct {
		Credential
		ctxclient.Func
	}{credential, f}
}

// WithHTTPClient specifies the http.Client the credential will use for calls.
func WithHTTPClient(credential Credential, client *http.Client) Credential {
	if client != nil {
		return WithHTTPClientFunc(credential,
			func(ctx context.Context) (*http.Client, error) {
				return client, nil
			})
	}
	return credential
}

// UploadFile describes an a document attachment for uploading
type UploadFile struct {
	// mime type of content
	ContentType string
	// file name to display in envelope
	FileName string
	// envelope documentId
	ID string
	// document order for envelope
	Order string
	// reader for creating file
	Data io.Reader
}

func closeUploads(files []*UploadFile) {
	for _, f := range files {
		if closer, ok := f.Data.(io.Closer); ok {
			closer.Close()
		}
	}
}

// GetTabValues returns a NameValue list of all entry tabs
func GetTabValues(tabs model.Tabs) []model.NameValue {
	results := make([]model.NameValue, 0)
	for _, v := range tabs.CheckboxTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, Value: fmt.Sprintf("%v", v.Selected)})
	}
	for _, v := range tabs.CompanyTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.DateSignedTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, Value: v.Value})
	}
	for _, v := range tabs.DateTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.EmailTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.FormulaTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.ListTabs {
		vals := make([]string, 0, len(v.ListItems))
		for _, x := range v.ListItems {
			if x.Selected {
				vals = append(vals, x.Value)
			}
		}
		results = append(results, model.NameValue{Name: v.TabLabel, Value: strings.Join(vals, ",")})
	}
	for _, v := range tabs.NoteTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, Value: v.Value})
	}
	for _, v := range tabs.NumberTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.RadioGroupTabs {
		for _, x := range v.Radios {
			if x.Selected {
				results = append(results, model.NameValue{Name: v.GroupName, Value: x.Value})
				break
			}
		}
	}
	for _, v := range tabs.SSNTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.TextTabs {
		results = append(results, model.NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}

	return results
}

// ResolveDSURL resolves a relative url.
// the host parameter determines which docusign server(s) to hit
//   EX: prod north america, prod europe, demo
// the accountID is used to finish the url's path.
func ResolveDSURL(ref *url.URL, host string, accountID string) {
	ref.Scheme = "https"
	ref.Host = host

	if strings.HasPrefix(ref.Path, "/") {
		ref.Path = "/restapi/v2" + ref.Path
	} else {
		ref.Path = "/restapi/v2/accounts/" + accountID + "/" + ref.Path
	}
}
