// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package click_test

import (
	"context"
	"fmt"
	"log"

	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/click"
)

func Example() {
	ctx := context.TODO()
	cfg := &esign.OAuth2Config{
		IntegratorKey:    "51d1a791-489c-4622-b743-19e0bd6f359e",
		Secret:           "f625e6f7-48e3-4226-adc5-66e434b21355",
		RedirURL:         "https://yourdomain.com/auth",
		AccountID:        "c23357a7-4f00-47f5-8802-94d2b1fb9a29",
		ExtendedLifetime: true,
		IsDemo:           true,
	}
	state := "SomeRandomStringOfSufficientSize"

	authURL := cfg.AuthURL(state, "click.manage")
	// Redirect user to consent page.
	fmt.Printf("Visit %s", authURL)

	// Enter code returned to redirectURL.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	credential, err := cfg.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}
	result, err := click.New(credential).List().Do(ctx)
	if err != nil {
		log.Fatalf("list clickwraps error: %v", err)
	}
	for _, cw := range result.Clickwraps {
		fmt.Printf("%s: %s", cw.ID, cw.Name)
	}

}
