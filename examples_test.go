// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package esign_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/jfcote87/esign"
	"github.com/jfcote87/esign/v2.1/envelopes"
	"github.com/jfcote87/esign/v2.1/folders"
	"github.com/jfcote87/esign/v2.1/model"
)

func ExampleOAuth2Config() {
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

	authURL := cfg.AuthURL(state)
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
	fl, err := folders.New(credential).List().Do(ctx)
	if err != nil {
		log.Fatalf("Folder list error: %v", err)
	}
	for _, fld := range fl.Folders {
		fmt.Printf("%s: %s", fld.Name, fld.FolderID)
	}

}

func ExampleJWTConfig() {
	ctx := context.TODO()
	cfg := &esign.JWTConfig{
		IntegratorKey: "51d1a791-489c-4622-b743-19e0bd6f359e",
		KeyPairID:     "1f57a66f-cc71-45cd-895e-9aaf39d5ccc4",
		PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
		MIICWwIBAAKBgGeDMVfH1+RBI/JMPfrIJDmBWxJ/wGQRjPUm6Cu2aHXMqtOjK3+u
		.........
		ZS1NWRHL8r7hdJL8lQYZPfNqyYcW7C1RW3vWbCRGMA==
		-----END RSFolderA PRIVATE KEY-----`,
		AccountID: "c23357a7-4f00-47f5-8802-94d2b1fb9a29",
		IsDemo:    true,
	}

	apiUserName := "78e5a047-f767-41f8-8dbd-10e3eed65c55"
	credential, err := cfg.Credential(apiUserName, nil, nil)
	if err != nil {
		log.Fatalf("create credential error: %v", err)
	}
	fl, err := folders.New(credential).List().Do(ctx)
	if err != nil {
		log.Fatalf("Folder list error: %v", err)
	}
	for _, fld := range fl.Folders {
		fmt.Printf("%s: %s", fld.Name, fld.FolderID)
	}
}

func ExampleJWTConfig_UserConsentURL() {
	// Prior to allowing jwt impersonation, each user must
	// provide consent.  Use the JWTConfig.UserConsentURL
	// to provide a url for individual to provide their
	// consent to DocuSign.
	apiUserID := "78e5a047-f767-41f8-8dbd-10e3eed65c55"

	cfg := &esign.JWTConfig{
		IntegratorKey: "51d1a791-489c-4622-b743-19e0bd6f359e",
		KeyPairID:     "1f57a66f-cc71-45cd-895e-9aaf39d5ccc4",
		PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
		MIICWwIBAAKBgGeDMVfH1+RBI/JMPfrIJDmBWxJ/wGQRjPUm6Cu2aHXMqtOjK3+u
		.........
		ZS1NWRHL8r7hdJL8lQYZPfNqyYcW7C1RW3vWbCRGMA==
		-----END RSA PRIVATE KEY-----`,
		AccountID: "c23357a7-4f00-47f5-8802-94d2b1fb9a29",
		IsDemo:    true,
	}

	// the redirURL must match one entered in DocuSign's Edit API
	// Integrator Key screen.
	consentURL := cfg.UserConsentURL("https://yourdomain.com/auth")
	_ = consentURL
	// redirect user to consentURL.  Once users given consent, a JWT
	// credential may be used to impersonate the user.  This only has
	// to be done once.
	ctx := context.TODO()
	cred, err := cfg.Credential(apiUserID, nil, nil)
	if err != nil {
		log.Fatalf("%v", err)
	}
	u, err := cred.UserInfo(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	log.Printf("%s", u.Name)

}

func getCredential(ctx context.Context, apiUser string) (*esign.OAuth2Credential, error) {
	cfg := &esign.JWTConfig{
		IntegratorKey: "51d1a791-489c-4622-b743-19e0bd6f359e",
		KeyPairID:     "1f57a66f-cc71-45cd-895e-9aaf39d5ccc4",
		PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
		MIICWwIBAAKBgGeDMVfH1+RBI/JMPfrIJDmBWxJ/wGQRjPUm6Cu2aHXMqtOjK3+u
		.........
		ZS1NWRHL8r7hdJL8lQYZPfNqyYcW7C1RW3vWbCRGMA==
		-----END RSA PRIVATE KEY-----`,
		AccountID: "c23357a7-4f00-47f5-8802-94d2b1fb9a29",
		IsDemo:    true,
	}

	return cfg.Credential(apiUser, nil, nil)

}

func ExampleTokenCredential() {
	// retrieve token frm DocuSign OAuthGenerator
	// https://developers.docusign.com/oauth-token-generator
	// or generate or
	var accessToken = `eyJ0eXAiOiJNVCIsImF...`
	credential := esign.TokenCredential(accessToken, true)
	fl, err := folders.New(credential).List().Do(context.Background())
	if err != nil {
		log.Fatalf("Folder list error: %v", err)
	}
	for _, fld := range fl.Folders {
		fmt.Printf("%s: %s", fld.Name, fld.FolderID)
	}
}

func Example_create_envelope() {
	ctx := context.TODO()
	apiUserID := "78e5a047-f767-41f8-8dbd-10e3eed65c55"
	cred, err := getCredential(ctx, apiUserID)
	if err != nil {
		log.Fatalf("credential error: %v", err)
	}
	sv := envelopes.New(cred)
	f1, err := ioutil.ReadFile("letter.pdf")
	if err != nil {
		log.Fatalf("file read error: %v", err)
	}

	env := &model.EnvelopeDefinition{
		EventNotification: &model.EventNotification{
			URL: "https://webhook.site/67ce4f73-7af7-4995-b37d-a1fab267e706",
			EnvelopeEvents: []model.EnvelopeEvent{
				{
					EnvelopeEventStatusCode: "Sent",
					IncludeDocuments:        "false",
				},
				{
					EnvelopeEventStatusCode: "Completed",
					IncludeDocuments:        "true",
				},
			},
			IncludeDocuments: true,
		},
		EmailSubject: "[Go eSignagure SDK] - Please sign this doc",
		EmailBlurb:   "Please sign this test document",
		Status:       "sent",
		Documents: []model.Document{
			{
				DocumentBase64: f1,
				Name:           "invite letter.pdf",
				DocumentID:     "1",
			},
		},
		Recipients: &model.Recipients{
			Signers: []model.Signer{
				{
					Email:             "user@example.com",
					EmailNotification: nil,
					Name:              "J F Cote",
					RecipientID:       "1",
					Tabs: &model.Tabs{
						SignHereTabs: []model.SignHere{
							{
								TabBase: model.TabBase{
									DocumentID:  "1",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "signature",
									XPosition:  "192",
									YPosition:  "160",
								},
							},
						},
						DateSignedTabs: []model.DateSigned{
							{
								TabBase: model.TabBase{
									DocumentID:  "1",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "dateSigned",
									XPosition:  "334",
									YPosition:  "179",
								},
							},
						},
					},
				},
			},
		},
	}
	envSummary, err := sv.Create(env).Do(ctx)
	if err != nil {
		log.Fatalf("create envelope error: %v", err)
	}
	log.Printf("New envelope id: %s", envSummary.EnvelopeID)
}

func Example_create_envelope_fileupload() {
	ctx := context.TODO()
	apiUserID := "78e5a047-f767-41f8-8dbd-10e3eed65c55"
	cred, err := getCredential(ctx, apiUserID)
	if err != nil {
		log.Fatalf("credential error: %v", err)
	}
	// open files for upload
	f1, err := os.Open("letter.pdf")
	if err != nil {
		log.Fatalf("file open letter.pdf: %v", err)
	}
	f2, err := os.Open("contract.pdf")
	if err != nil {
		log.Fatalf("file open contract.pdf: %v", err)
	}

	sv := envelopes.New(cred)

	env := &model.EnvelopeDefinition{
		EmailSubject: "[Go eSignagure SDK] - Please sign this doc",
		EmailBlurb:   "Please sign this test document",
		Status:       "sent",
		Documents: []model.Document{
			{
				Name:       "invite letter.pdf",
				DocumentID: "1",
			},
			{
				Name:       "contract.pdf",
				DocumentID: "2",
			},
		},
		Recipients: &model.Recipients{
			Signers: []model.Signer{
				{
					Email:             "user@example.com",
					EmailNotification: nil,
					Name:              "J F Cote",
					RecipientID:       "1",
					Tabs: &model.Tabs{
						SignHereTabs: []model.SignHere{
							{
								TabBase: model.TabBase{
									DocumentID:  "1",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "signature",
									XPosition:  "192",
									YPosition:  "160",
								},
							},
						},
						DateSignedTabs: []model.DateSigned{
							{
								TabBase: model.TabBase{
									DocumentID:  "1",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "dateSigned",
									XPosition:  "334",
									YPosition:  "179",
								},
							},
						},
						TextTabs: []model.Text{
							{
								TabBase: model.TabBase{
									DocumentID:  "2",
									RecipientID: "1",
								},
								TabPosition: model.TabPosition{
									PageNumber: "1",
									TabLabel:   "txtNote",
									XPosition:  "70",
									YPosition:  "564",
								},
								TabStyle: model.TabStyle{
									Name: "This is the tab tooltip",
								},
								Width:  "300",
								Height: "150",
							},
						},
					},
				},
			},
		},
	}
	// UploadFile.Reader is always closed on operation
	envSummary, err := sv.Create(env, &esign.UploadFile{
		ContentType: "application/pdf",
		FileName:    "invitation letter.pdf",
		ID:          "1",
		Reader:      f1,
	}, &esign.UploadFile{
		ContentType: "application/pdf",
		FileName:    "contract.pdf",
		ID:          "2",
		Reader:      f2,
	}).Do(ctx)

	if err != nil {
		log.Fatalf("create envelope error: %v", err)
	}
	log.Printf("New envelope id: %s", envSummary.EnvelopeID)

}

func Example_document_download() {

	ctx := context.TODO()
	envID := "78e5a047-f767-41f8-8dbd-10e3eed65c55"
	cred, err := getCredential(ctx, envID)
	if err != nil {
		log.Fatalf("credential error: %v", err)
	}
	sv := envelopes.New(cred)

	var dn *esign.Download
	if dn, err = sv.DocumentsGet("combined", envID).
		Certificate().
		Watermark().
		Do(ctx); err != nil {

		log.Fatalf("DL: %v", err)
	}
	defer dn.Close()
	log.Printf("File Type: %s  Length: %d", dn.ContentType, dn.ContentLength)
	f, err := os.Create("ds_env.pdf")
	if err != nil {
		log.Fatalf("file create: %v", err)
	}
	defer f.Close()
	if _, err := io.Copy(f, dn); err != nil {
		log.Fatalf("file create: %v", err)
	}

}
