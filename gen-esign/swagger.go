// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// swagger provides structs and utilities for handling a swagger file.
// I created it for use with esign package. It is incomplete and not
// tested for other swagger implementations.

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var initialismMap map[string]bool

// create initialism map to speed ToGoName
func init() {
	initialismMap = make(map[string]bool)
	for _, v := range commonInitialisms {
		initialismMap[v] = true
	}
}

// Document organizes data from the Docusign swagger definition file
type Document struct {
	Version      string       `json:"swagger,omitempty"`
	Info         Info         `json:"info,omitempty"`
	Host         string       `json:"host,omitempty"`
	BasePath     string       `json:"base_path,omitempty"`
	Schemes      []string     `json:"schemes,omitempty"`
	Consumes     []string     `json:"consumes,omitempty"`
	Produces     []string     `json:"produces,omitempty"`
	ExternalDocs ExternalDocs `json:"external_docs,omitempty"`
	Operations   OpList       `json:"paths,omitempty"`
	Definitions  DefSlice     `json:"definitions,omitempty"`
	Tags         []Tag        `json:"tags,omitempty"`
	DSTags       []Tag        `json:"x-ds-categories,omitempty"`
}

// OpList provides custom json decoding for
// a slice of Options.
type OpList []Operation

// UnmarshalJSON reads a json map and decodes it into
// a slice of Operations.  It adds the path key and
// method keys to each operation.
func (ops *OpList) UnmarshalJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))
	tk, err := d.Token()
	for err == nil {
		switch tx := tk.(type) {
		case string:
			var newOps []Operation
			if newOps, err = unmarshalOps(d, tx); err == nil {
				*ops = append(*ops, newOps...)
			}
		}
		if tk, err = d.Token(); err == io.EOF {
			return nil
		}
	}
	return err
}

func unmarshalOps(d *json.Decoder, path string) ([]Operation, error) {
	var pathParams []Property
	var newOps []Operation
	tk, err := d.Token()
	for err == nil {
		switch tx := tk.(type) {
		case string:
			if tx == "parameters" {
				err = d.Decode(&pathParams)
			} else {
				op := Operation{Path: path, HTTPMethod: strings.ToUpper(tx)}
				if err = d.Decode(&op); err == nil {
					newOps = append(newOps, op)
				}
			}
		case json.Delim:
			if tx.String() == "}" {
				if len(pathParams) > 0 {
					for i := range newOps {
						newOps[i].Parameters = append(newOps[i].Parameters, pathParams...)
					}
				}
				return newOps, nil
			}
		}
		if err == nil {
			tk, err = d.Token()
		}
	}
	return nil, err
}

// Len defined to allow sort
func (ops OpList) Len() int { return len(ops) }

// Swap defined to allow sort
func (ops OpList) Swap(i, j int) { ops[i], ops[j] = ops[j], ops[i] }

// Less defined to allow sort
func (ops OpList) Less(i, j int) bool {
	if ops[i].Service == ops[j].Service {
		if len(ops[i].Tags) == 1 && len(ops[j].Tags) == 1 {
			if ops[i].Tags[0] == ops[j].Tags[0] {
				return ops[i].Method < ops[j].Method
			}
			return ops[i].Tags[0] < ops[j].Tags[0]
		}
		return ops[i].OperationID < ops[j].OperationID
	}
	return ops[i].Service < ops[j].Service
}

// Tag provides a full definition of an operation tag
type Tag struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Summary     string `json:"summary,omitempty"`
}

// Definition describes a structs needed for creating
// requests and reading responses.
type Definition struct {
	ID          string    `json:"-"`
	Description string    `json:"description,omitempty"`
	Type        string    `json:"type,omitempty"`
	Fields      FieldList `json:"properties,omitempty"`
	Name        string    `json:"x-ds-definition-name,omitempty"`
	Summary     string    `json:"x-ms-summary,omitempty"`
	Category    string    `json:"x-ds-category,omitempty"`
	Order       string    `json:"x-ds-order,omitempty"`
}

// CommentLines converts the lines of Description
// field into a string slice
func (d Definition) CommentLines() []string {
	if len(d.Description) == 0 {
		return nil
	}
	comments := strings.Split(d.Description, "\n")
	tspace := " "
	if strings.HasPrefix(comments[0], "A ") || strings.HasPrefix(comments[0], "An ") || strings.HasPrefix(comments[0], "The ") {
		tspace = " is "
	}
	comments[0] = d.StructName() + tspace + strings.ToLower(comments[0][0:1]) + comments[0][1:]
	return comments
}

// StructName returns the go formatted name
func (d Definition) StructName() string {
	return ToGoName(d.Name)
}

// getGoFieldType translates from swagger type to go
func getGoFieldType(f Field) string {
	var fldType string
	switch f.Type {
	case "-":
		fldType = "-"
	case "":
		fldType = "*Self" //defMap[f.Ref].StructName()
	case "string":
		fldType = "string"
	case "integer":
		var validFormats = map[string]string{"int32": "int32", "int64": "int64"}
		var ok bool
		if fldType, ok = validFormats[f.Format]; !ok {
			fldType = "int"
		}
	case "boolean":
		fldType = "bool"
	case "number":
		fldType = "float64"
	case "object":
		if f.AdditionalProperties != nil {
			fldType = "map[string]" + getGoFieldType(Field{Type: f.AdditionalProperties.Type})
		}
	case "array":
		return handleArray(f.Items.Type)
	}
	return fldType
}

func handleArray(fldType string) string {
	switch fldType {
	case "":
		return "[]REF" //defMap[f.Items.Ref].StructName()
	case "string":
		return "[]string"
	case "integer":
		return "[]int"
	case "number":
		return "[]float64"
	}
	return ""
}

// StructFields returns all info need to create
// the struct field definition (go name, json name, comments, type).
// defMap is a map of all definitions and the overrides specify
// type overrides.
func (d Definition) StructFields(defMap map[string]Definition, overrides map[string]map[string]string) []StructField {
	// use x-definition-name for lookup
	overrideMap, ok := overrides[d.Name]
	if !ok {
		overrideMap = make(map[string]string)
	}
	var fldType string
	var fields []StructField
	if s, ok := overrideMap["TABS"]; ok {
		for _, nm := range strings.Split(s, ",") {
			fields = append(fields, StructField{
				Name: nm,
			})
		}
	}
	for _, f := range d.Fields {
		if fldType, ok = overrideMap[f.Name]; !ok {
			fldType = getGoFieldType(f)
			switch fldType {
			case "*Self":
				fldType = "*" + defMap[f.Ref].StructName()
			case "[]REF":
				fldType = "[]" + defMap[f.Items.Ref].StructName()
			}
		}
		if fldType != "-" {
			fields = append(fields, StructField{
				Name:     ToGoName(f.Name),
				JSON:     f.Name,
				Type:     fldType,
				Comments: strings.Split(f.Description, "\n"),
			})
		}
	}
	return fields
}

// DefSlice provides custom JSON processing to convert
// a map to a slice and adding the definition ID to
// the Definition.
type DefSlice []Definition

// UnmarshalJSON reads a map and coverts to a slice.
// Set the DefinititionID to the key.
func (ds *DefSlice) UnmarshalJSON(b []byte) error {
	d := json.NewDecoder(bytes.NewReader(b))
	tk, err := d.Token()
	for err == nil {
		switch tx := tk.(type) {
		case string:
			def := Definition{ID: tx}
			if err = d.Decode(&def); err != nil {
				continue
			}
			*ds = append(*ds, def)
		}
		tk, err = d.Token()

	}
	if err != io.EOF {
		return err
	}
	return nil
}

// Len defined to allow sort
func (ds DefSlice) Len() int { return len(ds) }

// Swap defined to allow sort
func (ds DefSlice) Swap(i, j int) { ds[i], ds[j] = ds[j], ds[i] }

// Less defined to allow sort
func (ds DefSlice) Less(i, j int) bool {
	return ds[i].Name < ds[j].Name
}

// StructField provides info to generate a struct definition
// ex:
// type <StructName> struct {
//     // <Comments>
//     <Name> <Type> `json:"<JSON>,omitempty"`
// }
type StructField struct {
	Name     string
	Comments []string
	Type     string
	JSON     string
}

// Field describes struct field for a definition
type Field struct {
	Name                 string              `json:"name,omitempty"`
	Description          string              `json:"description,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Items                *SchemaRef          `json:"items,omitempty"`
	Ref                  string              `json:"$ref,omitempty"`
	Format               string              `json:"format,omitempty"`
	AdditionalProperties *AdditionalProperty `json:"additionalProperties,omitempty"`
}

// AdditionalProperty defines the value type of a map
type AdditionalProperty struct {
	Type string `json:"type"`
}

// FieldList provides custom json decoding for
// a Definition propery map
type FieldList []Field

// UnmarshalJSON reads a json map of properties and
// converts it to a slice and adding the key as the
// Name.
func (f *FieldList) UnmarshalJSON(b []byte) error {
	var xmap = make(map[string]Field)
	if err := json.Unmarshal(b, &xmap); err != nil {
		return err
	}
	for k, v := range xmap {
		v.Name = k
		*f = append(*f, v)
	}
	sort.Slice(*f, func(i, j int) bool { return (*f)[i].Name < (*f)[j].Name })
	return nil
}

// Operation describes the endpoint, inputs and expected responses
// for a function
type Operation struct {
	HTTPMethod  string              `json:"-"`
	Path        string              `json:"-"`
	Tags        []string            `json:"tags,omitempty"`
	Summary     string              `json:"summary,omitempty"`
	Description string              `json:"description,omitempty"`
	OperationID string              `json:"operationId,omitempty"`
	Consumes    []string            `json:"consumes,omitempty"`
	Produces    []string            `json:"produces,omitempty"`
	Parameters  []Property          `json:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses,omitempty"`
	Deprecated  bool                `json:"deprecated,omitempty"`
	Examples    []Example           `json:"x-ds-example,omitempty"`
	MethodName  string              `json:"x-ds-methodname,omitempty"`
	Method      string              `json:"x-ds-method,omitempty"`
	Service     string              `json:"x-ds-service,omitempty"`
	InSDK       bool                `json:"x-ds-in-sdk,omitempty"`
	Status      string              `json:"x-ds-api-status,omitempty"`
}

// Accepts converts the Produces slice to a comma
// separated string to use for Accept Header.
func (o Operation) Accepts() string {
	return strings.Join(o.Produces, ", ")
}

// CommentLines returns a list of comments to annotate the operation.
func (o Operation) CommentLines(funcName string, docPrefix string, hasFileUploads bool, isMediaUpload bool) []string {

	comments := strings.Split(o.Summary, "\n")
	tspace := " "
	uncatFlag := (o.Service == "Uncategorized")

	if len(comments[0]) > 0 {
		if strings.HasPrefix(comments[0], "A ") || strings.HasPrefix(comments[0], "An ") || strings.HasPrefix(comments[0], "The ") {
			tspace = " is "
		}
		comments[0] = funcName + tspace + strings.ToLower(comments[0][0:1]) + comments[0][1:]
		if hasFileUploads {
			comments = append(comments, "If any uploads[x].Reader is an io.ReadCloser(s), Do() will always close Reader.")
		}
		if isMediaUpload {
			comments = append(comments, "If media is an io.ReadCloser, Do() will close media.")
		}
		if uncatFlag {
			comments = append(comments, "operation is uncategorized and subject to change.")
		} else {
			if len(o.Tags) > 0 {
				comments = append(comments, "", "https://developers.docusign.com/esign-rest-api/"+strings.ToLower(docPrefix+"reference/"+o.Service+"/"+o.Tags[0]+"/"+o.Method))
			}
			if o.InSDK {
				comments = append(comments, "", "SDK Method "+o.SDK())
			}
		}
		return comments
	}
	if uncatFlag {
		comments[0] = funcName + "is uncategorized and subject to change"
		return comments
	}
	if o.InSDK {
		comments[0] = funcName + " is SDK Method " + o.SDK()
		if len(o.Tags) > 0 {
			comments = append(comments, "", "https://developers.docusign.com/esign/restapi/"+o.Service+"/"+o.Tags[0]+"/"+o.Method)
		}
	}
	return comments
}

// OpPath removes the accountId prefix for the op path.  Allows
// for credential to fill in.
func (o Operation) OpPath(ver string) string {
	stdPrefix := "/" + ver + "/accounts/{accountId}/"
	stdPrefixLen := len(stdPrefix)
	if strings.HasPrefix(o.Path, stdPrefix) { //o.Path, "/" + ver + "/accounts/{accountId}") {
		if len(o.Path) < stdPrefixLen {
			return ""
		}
		return o.Path[stdPrefixLen:]
	}
	return o.Path
}

// OpPath2 creates a replacement string
func (o Operation) OpPath2(ver string, p []PathParam) string {
	path := o.OpPath(ver)
	if len(p) > 0 {
		parts := make([]string, 0)
		for _, part := range strings.Split(path, "/") {
			pathPart := `"` + part + `"`
			for _, k := range p {
				if part == "{"+k.Name+"}" {
					pathPart = k.GoName
					break
				}
			}
			parts = append(parts, pathPart)
		}
		return "strings.Join([]string{" + strings.Join(parts, ",") + `},"/")`
	}
	return `"` + path + `"`
}

// GoFuncName provides a go formatted name
func (o Operation) GoFuncName(prefixList []string) string {
	if len(o.Tags) == 1 {
		tag := o.Tags[0]
		method := strings.ToUpper(o.Method[0:1]) + o.Method[1:]
		for _, pre := range prefixList {
			if strings.HasPrefix(tag, pre) {
				return ToGoName(tag[len(pre):] + method)
			}
		}
		return ToGoName(tag + method)
	}
	return ToGoName(o.MethodName)
}

// Payload describes the body of an operation
type Payload struct {
	GoName string
	Type   string
}

// Payload formats a Payload struct from the operation's parameters
func (o Operation) Payload(defMap map[string]Definition) *Payload {
	for _, p := range o.Parameters {
		if p.In == "body" {
			var ifType = ""
			if p.Schema == nil {
				ifType = p.Type
			} else if p.Schema.Type > "" {
				if p.Schema.Type == "string" && p.Schema.Format == "binary" {
					ifType = "[]byte"
				} else {
					ifType = p.Schema.Type
				}
			} else {
				if def, ok := defMap[p.Schema.Ref]; ok {
					ifType = "*model." + ToGoName(def.Name)
				}
			}
			return &Payload{GoName: ToGoNameLC(p.Name), Type: ifType}
		}
	}
	// doesn't have body, check for empty post or put
	if o.HTTPMethod == "PUT" || o.HTTPMethod == "POST" {
		return &Payload{GoName: "upload", Type: "*esign.UploadFile"}
	}
	return nil
}

// PathParam provides a name/value pair for constructing
// a call url
type PathParam struct {
	Name   string
	GoName string
}

// PathParameters returns list of parameters used to
// construct a call url
func (o Operation) PathParameters() []PathParam {
	var params []PathParam
	for _, p := range o.Parameters {
		if p.In == "path" && p.Name != "accountId" {
			params = append(params, PathParam{
				Name:   p.Name,
				GoName: ToGoNameLC(p.Name),
			})
		}
	}
	return params
}

// QueryOpt describes a possibe url query parameters. Used to
// construct option funcs for a call.
type QueryOpt struct {
	Name     string
	GoName   string
	Type     string
	Value    string
	Comments []string
}

func queryComments(p Property) []string {
	l := strings.Split(strings.Trim(p.Description, " \n"), "\n")
	// remove excess date verbage
	if strings.HasPrefix(p.Description, "Specifies the date") {
		for j := 0; j < len(l); j++ {
			if l[j] == "" {
				l = l[0:j]
				break
			}
		}
	}
	// if blank default message
	if len(l) == 0 || (len(l) == 1 && len(l[0]) == 0) {
		return []string{ToGoName(p.Name) + " set the call query parameter " + p.Name}
	}
	tspace := " "
	if strings.HasPrefix(l[0], "A ") || strings.HasPrefix(l[0], "An ") || strings.HasPrefix(l[0], "The ") {
		tspace = " is "
	}
	l[0] = ToGoName(p.Name) + tspace + strings.ToLower(l[0][0:1]) + l[0][1:]
	return l
}

// QueryOpts returns list of all query parameters
func (o Operation) QueryOpts(overrides map[string]map[string]string) []QueryOpt {
	opOverrides, ok := overrides[o.OperationID]
	if !ok {
		opOverrides = make(map[string]string)
	}
	var params []QueryOpt
	for _, p := range o.Parameters {
		if p.In == "query" {
			ty, ok := opOverrides[p.Name]
			if !ok {
				ty = p.Type
			}
			/*comments := strings.Split(strings.TrimRight(p.Description, "\n"), "\n")
			if len(comments[0]) == 0 {
				comments = nil
			}*/
			params = append(params, QueryOpt{
				Name:     p.Name,
				GoName:   ToGoName(p.Name),
				Type:     ty,
				Value:    valueCode(ty),
				Comments: queryComments(p),
			})
		}
	}
	return params
}

// SDK returns the operation's DocuSign SDK ID
func (o Operation) SDK() string {
	return o.Service + "::" + o.MethodName
}

// Result defines the return value for an operation
func (o Operation) Result(structMap map[string]Definition) string {
	for k, v := range o.Responses {
		if k == "200" || k == "201" {
			if v.Schema != nil {
				if v.Schema.Ref != "" {
					if def, ok := structMap[v.Schema.Ref]; ok {
						return "*model." + ToGoName(def.Name)
					}
				}
				if v.Schema.Type == "file" {
					return "*esign.Download"
				}
				return v.Schema.Type
			}
		}
	}
	return ""
}

// Property provides custom
type Property struct {
	Name        string     `json:"name,omitempty"`
	In          string     `json:"in,omitempty"`
	Description string     `json:"description,omitempty"`
	Required    bool       `json:"required,omitempty"`
	Type        string     `json:"type,omitempty"`
	Schema      *SchemaRef `json:"schema,omitempty"`
}

// valueCode generates code for updating a call's
// query options.
func valueCode(ty string) string {
	switch ty {
	case "bool":
		return `"true"`
	case "int":
		return "fmt.Sprintf(\"%d\", val )"
	case "float64":
		return "fmt.Sprintf(\"%f\", val )"
	case "time.Time":
		return "val.Format(time.RFC3339)"
	case "...string":
		return `strings.Join(val,",")`
	}
	return "val"
}

// Example contains an Operation's example code.  Not
// currently used in generation.
type Example struct {
	Description string                 `json:"description,omitempty"`
	Direction   string                 `json:"direction,omitempty"`
	Format      string                 `json:"format,omitempty"`
	Response    map[string]interface{} `json:"response,omitempty"`
	Request     map[string]interface{} `json:"request,omitempty"`
	Style       string                 `json:"style,omitempty"`
	Title       string                 `json:"title,omitempty"`
}

// SchemaRef provides a key to identify a field/property/parameter's
// data type.
type SchemaRef struct {
	Ref    string `json:"$ref"`
	Type   string `json:"type"`
	Format string `json:"format"`
}

// Response defines the datatype of a call's response.
type Response struct {
	Description string
	Schema      *SchemaRef
}

// Info describes data from the swagger file.  Not used
// in generation
type Info struct {
	Version        string  `json:"version,omitempty"`
	Title          string  `json:"title,omitempty"`
	Description    string  `json:"description,omitempty"`
	TermsOfService string  `json:"terms_of_service,omitempty"`
	Contact        Contact `json:"contact,omitempty"`
}

// Contact is not used in generation
type Contact struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// ExternalDocs provides links to external web documentation
type ExternalDocs struct {
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

// ToGoNameLC formats a
func ToGoNameLC(name string) string {
	nm := ToGoName(name)
	return strings.ToLower(nm[0:1]) + nm[1:]
}

var registeredGoNames = map[string]bool{
	"TabGuidedForm": true,
}

// ToGoName translates a swagger name which can be underscored or camel cased to a name that golint likes
func ToGoName(name string) string {
	if _, ok := registeredGoNames[name]; ok {
		return name
	}
	var out []string
	for _, w := range split(name) {
		uw := strings.ToUpper(w)
		mod := int(math.Min(float64(len(uw)), 2))
		if !initialismMap[uw] && !initialismMap[uw[:len(uw)-mod]] {
			uw = strings.ToUpper(w[:1]) + strings.ToLower(w[1:])
		}
		out = append(out, uw)
	}

	result := strings.Join(out, "")
	if len(result) > 0 {
		ud := strings.ToUpper(result[:1])
		ru := []rune(ud)
		if unicode.IsUpper(ru[0]) {
			result = ud + result[1:]
		} else {
			result = "X" + ud + result[1:]
		}
	}
	return result
}

// Prepares strings by splitting by caps, spaces, dashes, and underscore
func split(str string) (words []string) {
	repl := strings.NewReplacer(
		"@", "At ",
		"&", "And ",
		"|", "Pipe ",
		"$", "Dollar ",
		"!", "Bang ",
		"-", " ",
		"_", " ",
	)

	rex1 := regexp.MustCompile(`(\p{Lu})`)
	rex2 := regexp.MustCompile(`(\pL|\pM|\pN|\p{Pc})+`)

	str = strings.Trim(str, " ")

	// Convert dash and underscore to spaces
	str = repl.Replace(str)

	// Split when uppercase is found (needed for Snake)
	str = rex1.ReplaceAllString(str, " $1")
	// check if consecutive single char things make up an initialism

	for _, k := range commonInitialisms {
		str = strings.Replace(str, rex1.ReplaceAllString(k, " $1"), " "+k, -1)
	}
	// Get the final list of words
	words = rex2.FindAllString(str, -1)

	return
}

// Taken from https://github.com/golang/lint/blob/3390df4df2787994aea98de825b964ac7944b817/lint.go#L732-L769
var commonInitialisms = []string{
	"ACL",
	"API",
	"ASCII",
	"CPU",
	"CSS",
	"DNS",
	"EOF",
	"GUID",
	"HTML",
	"HTTPS",
	"HTTP",
	"ID",
	"IP",
	"JSON",
	"LHS",
	"QPS",
	"RAM",
	"RHS",
	"RPC",
	"SLA",
	"SMTP",
	"SQL",
	"SSH",
	"SSN", // added by JFC
	"TCP",
	"TLS",
	"TTL",
	"UDP",
	"UI",
	"UID",
	"UUID",
	"URI",
	"URL",
	"UTF8",
	"VM",
	"XML",
	"XMPP",
	"XSRF",
	"XSS",
}
