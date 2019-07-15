// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//
//go:generate go run ./gen-esign/main.go ./gen-esign/swagger.go  ./gen-esign/overrides.go

/*
Package esign implements all interfaces of DocuSign's
eSignatature REST API as defined in the published  api.

Api documentation: https://developers.docusign.com/

*/
package esign // import "github.com/jfcote87/esign"
