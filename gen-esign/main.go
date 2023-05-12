// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gen-esign creates the esign subpackages based upon DocuSign's
// esignature.rest.swagger.json definition file.

// Package main is the executable for gen-esign
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/gen-esign/swagger"
	"github.com/jfcote87/esign/gen-esign/templates"
)

const (
	openAPIdefinition = "https://github.com/docusign/eSign-OpenAPI-Specification"
	generatorVersion  = "20190720"
	pkgBaseName       = "github.com/jfcote87/esign"
)

var (
	definitionFileMap = map[string]esign.APIVersion{
		esign.APIv2.Name():     esign.APIv2,
		esign.APIv21.Name():    esign.APIv21,
		esign.AdminV2.Name():   esign.AdminV2,
		esign.ClickV1.Name():   esign.ClickV1,
		esign.MonitorV2.Name(): esign.MonitorV2,
		esign.RoomsV2.Name():   esign.RoomsV2,
	}
	apiParametersMap = map[esign.APIVersion]APIGenerateCfg{
		esign.APIv2: {
			DocPrefix:        "esign-rest-api/v2/",
			CallVersion:      "esign.APIv2",
			PackagePath:      "v2",
			ModelFile:        "v2/model/model.go",
			ModelPackage:     "model",
			ModelPackagePath: "v2/model",
			ModelIsPackage:   true,
			ModelImports:     []string{"fmt", "strings", "time"},
			fldOverrides:     swagger.GetFieldOverrides(),
			paramOverrides:   swagger.GetParameterOverrides(),
		},
		esign.APIv21: {
			DocPrefix:        "esign-rest-api/",
			CallVersion:      "esign.APIv21",
			PackagePath:      "v2.1",
			ModelFile:        "v2.1/model/model.go",
			ModelPackage:     "model",
			ModelPackagePath: "v2.1/model",
			ModelImports:     []string{"fmt", "strings", "time"},
			ModelIsPackage:   true,
			fldOverrides:     swagger.GetFieldOverrides(),
			paramOverrides:   swagger.GetParameterOverrides(),
		},
		esign.AdminV2: {
			DocPrefix:        "admin-api/",
			CallVersion:      "esign.AdminV2",
			PackagePath:      "admin",
			ModelFile:        "admin/admin.go",
			ModelPackage:     "admin",
			ModelPackagePath: "admin",
			ModelIsPackage:   true,
			UseMethodName:    true,
			fldOverrides:     make(map[string]map[string]string),
			paramOverrides:   make(map[string]map[string]string),
		},
		esign.RoomsV2: {
			DocPrefix:        "rooms-api/",
			CallVersion:      "esign.RoomsV2",
			PackagePath:      "rooms/",
			ModelFile:        "rooms/rooms.go",
			ModelPackage:     "rooms",
			ModelPackagePath: "rooms",
			ModelIsPackage:   true,
			UseMethodName:    true,
			fldOverrides:     make(map[string]map[string]string),
			paramOverrides:   make(map[string]map[string]string),
		},
		esign.ClickV1: {
			DocPrefix:        "click-api/",
			DocService:       "accounts",
			CallVersion:      "esign.ClickV1",
			PackagePath:      "",
			ModelFile:        "click/model.go",
			ModelPackage:     "click",
			ModelPackagePath: "click",
			ModelIsPackage:   false,
			UseMethodName:    true,
			fldOverrides:     make(map[string]map[string]string),
			paramOverrides:   make(map[string]map[string]string),
		},
		esign.MonitorV2: {
			DocPrefix:        "monitor-api/",
			DocService:       "monitor",
			CallVersion:      "esign.MonitorV2",
			PackagePath:      "",
			ModelFile:        "monitor/model.go",
			ModelPackage:     "monitor",
			ModelPackagePath: "monitor",
			ModelIsPackage:   false,
			UseMethodName:    true,
			fldOverrides:     make(map[string]map[string]string),
			paramOverrides:   make(map[string]map[string]string),
		},
	}

	baseDir     = flag.String("src", ".", "source directory")
	serviceTmpl = flag.String("service_templ", "", "override service package template")
	modelTmpl   = flag.String("model_templ", "", "api definitions template")
	specsFolder = flag.String("swaggerfiles", "gen-esign/specs", "directory containing swagger specification files")
	skipFormat  = flag.Bool("skip_format", false, "skip gofmt command on generated files")
)

// APIGenerateCfg contains parameters for generating an eSignature version
type APIGenerateCfg struct {
	esign.APIVersion
	Templates        *template.Template // templates
	BaseDir          string             // source directory
	BasePkg          string
	SkipFormat       bool
	Name             string
	Version          string
	DocPrefix        string
	DocService       string
	CallVersion      string
	PackagePath      string
	ModelFile        string
	ModelPackage     string
	ModelPackagePath string
	ModelImports     []string
	ModelIsPackage   bool
	UseMethodName    bool
	//ResourceMap      map[string]string
	fldOverrides   map[string]map[string]string
	paramOverrides map[string]map[string]string
}

// ResourceMap returns a map of tags to service name
func (api APIGenerateCfg) ResourceMap() map[string]string {
	return swagger.ResourceMaps[api.APIVersion]
}

func main() {
	flag.Parse()

	pkgBaseDir, pkgSwaggerDir, skipFormatting := *baseDir, *specsFolder, *skipFormat
	codeTmpl, err := parseTemplates(*serviceTmpl, *modelTmpl)
	if err != nil {
		log.Fatalf("Templates: %v", err)
	}
	docmap, err := decodeSwaggerDocs(pkgSwaggerDir)
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err := os.Chdir(pkgBaseDir); err != nil {
		log.Fatalf("unable to set directory to %s: %v", pkgBaseDir, err)
	}
	if pkgBaseDir, err = os.Getwd(); err != nil {
		log.Fatalf("unable to retrieve working diretory: %v", err)
	}

	for v, doc := range docmap {
		cfg, ok := apiParametersMap[v]
		if !ok {
			log.Printf("skipping %s has no parameters entry", v.Name())
			continue
		}
		cfg.APIVersion = v
		cfg.Name = v.Name()
		cfg.Version = doc.Info.Version
		cfg.BaseDir = pkgBaseDir
		cfg.BasePkg = pkgBaseName
		cfg.Templates = codeTmpl
		cfg.SkipFormat = skipFormatting

		if err := cfg.genVersion(&doc); err != nil {
			log.Printf("%s %v", cfg.Name, err)
			return
		}
	}
}

func parseTemplates(serviceTmplFile, modelTmplFile string) (*template.Template, error) {
	var err error
	svc, model := templates.Service, templates.Model
	if serviceTmplFile > "" {
		b, err := ioutil.ReadFile(serviceTmplFile)
		if err != nil {
			return nil, fmt.Errorf("%s read %w", serviceTmplFile, err)
		}
		svc = string(b)
	}
	tmpl, err := template.New("service.tmpl").Parse(svc)
	if err != nil {
		return nil, fmt.Errorf("service.tmpl: %w", err)
	}
	if modelTmplFile > "" {
		b, err := ioutil.ReadFile(modelTmplFile)
		if err != nil {
			return nil, fmt.Errorf("%s read %w", modelTmplFile, err)
		}
		model = string(b)
	}
	if _, err = tmpl.New("model.tmpl").Parse(model); err != nil {
		return nil, fmt.Errorf("model.tmpl: %w", err)
	}

	return tmpl, err
}

func decodeSwaggerDocs(folderName string) (map[esign.APIVersion]swagger.Document, error) {
	fis, err := ioutil.ReadDir(folderName)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	if !strings.HasSuffix(folderName, "/") {
		folderName += "/"
	}
	var results = make(map[esign.APIVersion]swagger.Document)
	for _, f := range fis {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		b, err := ioutil.ReadFile(folderName + f.Name())
		if err != nil {
			return nil, err
		}
		var doc *swagger.Document
		if err = json.Unmarshal(b, &doc); err != nil {
			return nil, fmt.Errorf("%s decode %w", f.Name(), err)
		}
		apikey := doc.Info.Title + ":" + doc.Info.Version
		apiVersion, ok := definitionFileMap[apikey]
		if !ok {
			return nil, fmt.Errorf("no matching api version for %s", apikey)
		}
		results[apiVersion] = *doc
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no definition files specified in %s", folderName)
	}
	return results, nil
}

func (api *APIGenerateCfg) genVersion(doc *swagger.Document) error {
	// Put the Definitions (structs) in order
	sort.Sort(doc.Definitions)

	defMap := make(map[string]swagger.Definition)
	structName := ""
	defList := make([]swagger.Definition, 0, len(doc.Definitions))
	// create defMap for lookup later field and parameter
	// lookups.  Make certain defList has only unique names.
	for _, def := range doc.Definitions {
		defMap["#/definitions/"+def.ID] = def
		if structName != def.Name {
			defList = append(defList, def)
		}
		structName = def.Name
	}

	// generate model.go first
	if err := api.doModel(defList, defMap); err != nil {
		return fmt.Errorf("%v Generating model.go failed: %v", api.Version, err)
	}

	sort.Sort(doc.Operations)
	ops := make(map[string][]swagger.Operation, 0)
	for _, op := range doc.Operations {
		if op.Status == "restricted" {
			log.Printf("Skipping: %s %s", op.Status, op.OperationID)
			continue
		}
		if op.Service == "" {
			log.Printf("No service specified: %s", op.OperationID)
			continue
		}
		fullService := api.Version + ":" + op.Service
		if newServiceName, ok := swagger.ServiceNameOverride[fullService]; ok {
			op.Service = newServiceName
		}
		fullOpName := api.Version + ":" + op.OperationID
		if !swagger.OperationSkipList[fullOpName] {
			serviceName, ok := api.ResourceMap()[op.Tags[0]]
			if ok {
				op.Service = serviceName
			} else {
				log.Printf("No service found for %s", op.Tags[0])
			}
			ops[op.Service] = append(ops[op.Service], op)
		}
	}
	tagDescMap := make(map[string]string)
	for _, tag := range doc.DSTags {
		tagDescMap[tag.Name] = tag.Description
	}

	serviceTmpl := api.Templates.Lookup("service.tmpl")
	for k, v := range ops {
		log.Printf("Generating %s", k)
		descrip := tagDescMap[k]

		if err := api.doPackage(serviceTmpl, k, descrip, v, defMap); err != nil {
			return fmt.Errorf("%s generate %s.go failed: %v", api.Version, k, err)
		}
	}
	return nil
}

// doModel generates the model.go in the model package
func (api *APIGenerateCfg) doModel(defList []swagger.Definition, defMap map[string]swagger.Definition) error {
	modelTempl := api.Templates.Lookup("model.tmpl")
	// create model.go
	// get field overrides and tab overrides
	tabDefs := swagger.TabDefs(api.Name, defMap, api.fldOverrides)
	var data = struct {
		Definitions      []swagger.Definition
		DefMap           map[string]swagger.Definition
		FldOverrides     map[string]map[string]string
		CustomCode       string
		DocPrefix        string
		VersionID        string
		IsPackage        bool
		ModelPackage     string
		ModelPackagePath string
		ModelImports     []string
		Scopes           string
	}{
		Definitions:  append(tabDefs, defList...), // Prepend tab definitions
		DefMap:       defMap,
		FldOverrides: api.fldOverrides,
		CustomCode:   swagger.CustomCode(api.Name),
		DocPrefix:    api.DocPrefix,

		VersionID:        api.Version,
		IsPackage:        api.ModelIsPackage,
		ModelPackage:     api.ModelPackage,
		ModelPackagePath: api.ModelPackagePath,
		ModelImports:     api.ModelImports,
		Scopes:           swagger.PackageScopes(api.APIVersion),
	}
	modelBuffer := &bytes.Buffer{}
	if err := modelTempl.Execute(modelBuffer, data); err != nil {
		return err
	}

	if *skipFormat {
		return api.makePackageFile(api.ModelFile, modelBuffer.Bytes())
	}
	fmtBytes, err := format.Source(modelBuffer.Bytes())
	if err != nil {
		log.Printf("Source Error: %v", err)
		return err
	}
	return api.makePackageFile(api.ModelFile, fmtBytes)
}

// ExtOperation contains all needed info
// for the template merge
type ExtOperation struct {
	swagger.Operation
	OpPayload         *swagger.Payload
	HasUploads        bool
	IsMediaUpload     bool
	PathParams        []swagger.PathParam
	FuncName          string
	QueryOptions      []swagger.QueryOpt
	Result            string
	DownloadAdditions []swagger.DownloadAddition
	JSONResponse      bool
}

// doPackage creates a subpackage go file
func (api *APIGenerateCfg) doPackage(resTempl *template.Template, serviceName string, description string,
	ops []swagger.Operation, defMap map[string]swagger.Definition) error {
	packageName := strings.ToLower(serviceName)
	packageFile := packageName + "/" + packageName + ".go"
	if api.PackagePath > "" {
		packageFile = api.PackagePath + "/" + packageFile
	}
	comments := strings.Split(strings.TrimRight(description, "\n"), "\n")
	if packageName == "uncategorized" {
		comments = append(comments, "Uncategorized calls may change or move to other packages.")
	}

	extOps := make([]ExtOperation, 0, len(ops))
	modelPkg := ""
	if api.ModelIsPackage {
		modelPkg = api.ModelPackage
	}
	for _, op := range ops {
		modelPackage := api.ModelPackage
		if !api.ModelIsPackage {
			modelPackage = ""
		}
		payload := op.Payload(defMap, modelPackage)
		extOps = append(extOps, ExtOperation{
			Operation:         op,
			OpPayload:         payload,
			HasUploads:        swagger.IsUploadFilesOperation(api.Version + ":" + op.OperationID),
			IsMediaUpload:     payload != nil && payload.Type == "*esign.UploadFile",
			PathParams:        op.PathParameters(),
			FuncName:          op.GoFuncName(api.UseMethodName, swagger.GetServicePrefixes(op.Service)),
			QueryOptions:      op.QueryOpts(api.paramOverrides),
			Result:            op.Result(defMap, modelPkg),
			DownloadAdditions: swagger.GetDownloadAdditions(api.Version + ":" + op.OperationID),
			JSONResponse:      op.ReturnsJSON(),
		})
	}
	docService := serviceName
	if api.DocService > "" {
		docService = api.DocService
	}
	var data = struct {
		Service          string
		Package          string
		Directory        string
		Operations       []ExtOperation
		Comments         []string
		Packages         []string
		PackagePath      string
		ModelPackage     string
		ModelPackagePath string
		ModelIsPackage   bool
		DocPrefix        string
		DocService       string
		VersionID        string
		CallVersion      string
		AddDocLinks      bool
		Accept           string
	}{
		Service:          serviceName,
		Package:          packageName,
		Directory:        api.BasePkg,
		Operations:       extOps,
		Comments:         comments,
		Packages:         []string{`"context"`, `"net/url"`},
		PackagePath:      api.PackagePath,
		ModelPackage:     api.ModelPackage,
		ModelPackagePath: api.ModelPackagePath,
		ModelIsPackage:   api.ModelIsPackage,
		DocPrefix:        api.DocPrefix,
		DocService:       docService,
		VersionID:        api.Version,
		CallVersion:      api.CallVersion,
		AddDocLinks:      (serviceName != "Uncategorized"),
	}
	importMap := make(map[string]bool)
	for _, op := range extOps {
		if len(op.PathParameters()) > 0 {
			importMap[`"strings"`] = true
			break
		}
	}
	for _, o := range extOps {
		for _, q := range o.QueryOptions {
			switch q.Type {
			case "...string":
				importMap[`"strings"`] = true
			case "int", "int32", "int64":
				importMap[`"fmt"`] = true
			case "float64":
				importMap[`"fmt"`] = true
			case "time.Time":
				importMap[`"time"`] = true
			}
		}
		if o.IsMediaUpload {
			importMap[`"io"`] = true
		}
	}

	for k, v := range importMap {
		if v {
			data.Packages = append(data.Packages, k)
		}
	}
	sort.Strings(data.Packages)
	data.Packages = append(data.Packages,
		"",
		"\""+api.BasePkg+"\"")
	if api.ModelIsPackage {
		data.Packages = append(data.Packages,
			"\""+api.BasePkg+"/"+api.ModelPackagePath+"\"")
	}

	pkgBuffer := &bytes.Buffer{}
	if err := resTempl.Execute(pkgBuffer, data); err != nil {
		return err
	}
	if !*skipFormat {
		pkgBytes, err := format.Source(pkgBuffer.Bytes())
		if err == nil {
			pkgBuffer = bytes.NewBuffer(pkgBytes)
		}
	}
	return api.makePackageFile(packageFile, pkgBuffer.Bytes())

}

func (api *APIGenerateCfg) makePackageFile(fileName string, content []byte) error {

	if err := os.MkdirAll(path.Dir(fileName), 0755); err != nil {
		return err
	}
	return os.WriteFile(fileName, content, 0755)
}
