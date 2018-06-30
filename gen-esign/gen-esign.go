// Copyright 2017 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// gen-esign creates the esign package based upon docusign's swagger file

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"text/template"

	"github.com/jfcote87/esign/gen-esign/overrides"
	"github.com/jfcote87/esign/gen-esign/swagger"
)

const (
	openAPIdefinition = "https://github.com/docusign/eSign-OpenAPI-Specification"
	generatorVersion  = "20171011"
)

var (
	basePkg     = flag.String("basepkg", "github.com/jfcote87/esign", "root package in gopath")
	baseDir     = flag.String("gopath", fmt.Sprintf("%s/src", build.Default.GOPATH), "GOPATH src directory")
	templDir    = flag.String("template", "gen-esign/templates", "directory containing output templates.")
	buildFlag   = flag.Bool("build", false, "Compile generated packages.")
	swaggerFile = flag.String("swagger_file", "esignature.rest.swagger.json", "If non-empty, the path to a local file on disk containing the API to generate. Exclusive with setting --api.")
)

// main program
func main() {
	if err := os.Chdir(getEsignDir()); err != nil {
		log.Fatalf("unable to chdir to %s: %v", getEsignDir(), err)
	}
	// Open swagger file and parse
	f, err := os.Open(*swaggerFile)
	if err != nil {
		log.Fatalf("Unable to open: %v", err)
	}
	var doc *swagger.Document
	if err = json.NewDecoder(f).Decode(&doc); err != nil {
		log.Fatalf("Unable to parse: %v", err)

	}
	f.Close()

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
	// create templates
	if !strings.HasPrefix(*templDir, "/") {
		*templDir = path.Join(*baseDir, *basePkg, *templDir)
	}
	genTemplates, err := template.ParseFiles(path.Join(*templDir, "service.tmpl"), path.Join(*templDir, "/model.tmpl"))
	if err != nil {
		log.Fatalf("Templates: %v", err)
	}
	resTempl := genTemplates.Lookup("service.tmpl")
	modelTempl := genTemplates.Lookup("model.tmpl")
	if resTempl == nil || modelTempl == nil {
		log.Fatalf("Nil templates")
	}

	// generate model.go
	if err := doModel(modelTempl, defList, defMap); err != nil {
		log.Fatalf("Generating model.go failed: %v", err)
	}

	sort.Sort(doc.Operations)
	ops := make(map[string][]swagger.Operation, 0) //**
	for _, op := range doc.Operations {
		opList := ops[op.Service]
		opList = append(opList, op)
		ops[op.Service] = opList
	}

	tagDescMap := make(map[string]string)
	for _, tag := range doc.DSTags {
		tagDescMap[tag.Name] = tag.Description
	}
	for k, v := range ops {
		log.Printf("Generating %s", k)
		descrip, _ := tagDescMap[k]
		if err := doPackage(resTempl, k, descrip, v, defMap); err != nil {
			log.Fatalf("generate %s.go failed: %v", k, err)
		}
	}
	log.Printf("Generating Tabs")
}

// doModel generates the model.go in the model package
func doModel(modelTempl *template.Template, defList []swagger.Definition, defMap map[string]swagger.Definition) error {
	// create model.go
	if err := os.Chdir(getEsignDir() + "/model"); err != nil {
		return err
	}
	f, err := os.Create("model.go")
	if err != nil {
		return err
	}
	defer f.Close()
	// get field overrides and tab overrides
	fldOverrides := overrides.GetFieldOverrides()
	tabDefs := overrides.TabDefs(defMap, fldOverrides)
	var data = struct {
		Definitions  []swagger.Definition
		DefMap       map[string]swagger.Definition
		FldOverrides map[string]map[string]string
		CustomCode   string
	}{
		Definitions:  append(tabDefs, defList...), // Prepend tab definitions
		DefMap:       defMap,
		FldOverrides: fldOverrides,
		CustomCode:   overrides.CustomCode,
	}
	return modelTempl.Execute(f, data)

}

// doPackage
func doPackage(resTempl *template.Template, serviceName, description string, ops []swagger.Operation, defMap map[string]swagger.Definition) error {
	var totOps = 0
	var totQry1 = 0
	var totQry = 0
	for _, xy := range ops {
		totOps++
		qx := xy.QueryOpts(overrides.GetParameterOverrides())
		if len(qx) > 0 {
			totQry++
			if len(qx) > 1 {
				totQry1++
			}
			//log.Printf("%s: %d", xy.GoFuncName(overrides.GetOperationOverrides()), len(qx))
		}
	}
	//log.Printf("%d  %d  %d", totOps, totQry, totQry1)

	packageName := strings.ToLower(serviceName)
	pkgDir := getEsignDir() + "/" + packageName
	if err := os.Chdir(pkgDir); err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(pkgDir, 0755); err == nil {
				os.Chdir(pkgDir)
			}
		}
		if err != nil {
			return err
		}
	}
	var data = struct {
		Service       string
		Package       string
		Operations    []swagger.Operation
		Comments      []string
		DefMap        map[string]swagger.Definition
		OpOverrides   map[string]string
		PropOverrides map[string]map[string]string
		Packages      []string
	}{
		Service:       serviceName,
		Package:       packageName,
		Operations:    ops,
		Comments:      strings.Split(strings.TrimRight(description, "\n"), "\n"),
		DefMap:        defMap,
		OpOverrides:   overrides.GetOperationOverrides(),
		PropOverrides: overrides.GetParameterOverrides(),
		Packages:      make([]string, 0, 10),
	}
	var useStrings bool
	var useFmt bool
	var useTime bool
	for _, o := range ops {
		for _, q := range o.QueryOpts(data.PropOverrides) {
			switch q.Type {
			case "...string":
				useStrings = true
			case "int":
				useFmt = true
			case "float64":
				useFmt = true
			case "time.Time":
				useTime = true
			}
		}
	}
	if useFmt {
		data.Packages = append(data.Packages, `"fmt"`)
	}
	data.Packages = append(data.Packages, `"net/url"`)
	if useStrings {
		data.Packages = append(data.Packages, `"strings"`)
	}
	if useTime {
		data.Packages = append(data.Packages, `"time"`)
	}
	data.Packages = append(data.Packages,
		"",
		`"golang.org/x/net/context"`,
		"",
		"\""+*basePkg+"\"",
		"\""+*basePkg+"/model\"")
	f, err := os.Create(packageName + ".go")
	if err != nil {
		return err
	}
	defer f.Close()

	return resTempl.Execute(f, data)
}

func getEsignDir() string {
	return path.Join(*baseDir, *basePkg)
}
