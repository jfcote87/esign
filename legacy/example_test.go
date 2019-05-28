// Copyright 2019 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package legacy_test

import (
	"context"
	"fmt"
	"log"

	"github.com/jfcote87/esign/folders"
	"github.com/jfcote87/esign/legacy"
)

func Example_config() {
	cfg := legacy.Config{
		AccountID:     "c23357a7-4f00-47f5-8802-94d2b1fb9a29",
		IntegratorKey: "51d1a791-489c-4622-b743-19e0bd6f359e",
		UserName:      "sample User",
		Password:      "****************",
		Host:          "na1.docusign.com",
	}
	sv := folders.New(cfg)
	fl, err := sv.List().Do(context.Background())
	if err != nil {
		log.Fatalf("Folder list error: %v", err)
	}
	for _, fld := range fl.Folders {
		fmt.Printf("%s: %s", fld.Name, fld.FolderID)
	}
}

func Example_oauth() {
	cfg := legacy.Config{
		AccountID:     "c23357a7-4f00-47f5-8802-94d2b1fb9a29",
		IntegratorKey: "51d1a791-489c-4622-b743-19e0bd6f359e",
		UserName:      "sample User",
		Password:      "****************",
		Host:          "na1.docusign.com",
	}
	cred, err := cfg.OauthCredential(context.Background())
	if err != nil {
		log.Fatalf("oauth err: %v", err)
	}
	// save cred for future use
	sv := folders.New(cred)
	fl, err := sv.List().Do(context.Background())
	if err != nil {
		log.Fatalf("Folder list error: %v", err)
	}
	for _, fld := range fl.Folders {
		fmt.Printf("%s: %s", fld.Name, fld.FolderID)
	}
}
