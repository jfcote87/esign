# DocuSign eSignature RestApi v2 and v2.1(Beta) for Go

[![GoDoc](https://godoc.org/github.com/jfcote87/esign?status.svg)](https://godoc.org/github.com/jfcote87/esign)

esign provides Go packages for interacting with DocuSign's eSignature RestApi and
has been created using the lastest published swagger definition.
Definitions for call payloads and return values are found in  {VERSIONID}/model/model.go file.  
Docusign has split the API into major categories, and each category has a corresponding package 
in the v2 and v2.1 directory.  

The package requires Go 1.7 or above and has been tested with Go 1.9-1.12.

The previous package github.com/jfcote87/docusign is now deprecated.

## Announcements

eSignature v2.1 and [click api](https://developers.docusign.com/click-api) are now available.

## Resources

Official documentation: [https://developers.docusign.com/](https://developers.docusign.com/)

## Package Updates

All packages, excepte for gen-esign and legacy, are generated from DocuSign's [OpenAPI(swagger) specification](https://github.com/docusign/eSign-OpenAPI-Specification) using the gen-esign package.

Corrections to field names and definitions are documented in gen-esign/overrides/overrides.go.

The package names correspond to the API categories listed on the
[DocuSign REST API Reference](https://developers.docusign.com/esign-rest-api/reference) page.
Each operation contains an SDK Method which corresponds to the package operation which you will
find in operation definition.

## Authentication

Authentication is handled by the esign.Credential interface.  OAuth2Credentials may be created
via the OAuth2Config struct for 3-legged oauth and the JWTConfig struct for 2-legged oauth. Examples
are shown for each in the esign examples.

UserID/Password login is available via the legacy.Config which may also be used to create non-expiring
oauth tokens (legacy.OauthCredential).  Examples are shown in the legacy package.

## Models

The model package describes the structure of all data passed to and received from API calls.

## Operations

Each package has a service object which is initialized via the <packagename>.New(<credential>) call.
The service methods define all operations for the package with corresponding options.  An operation is
executed via a Do(context.Context) function.  A context must be passwed for all operation

## Example

Create envelope

```go
    import "github.com/jfcote87/esign"
    import "github.com/jfcote87/esign/v2/envelope"
    import "github.com/jfcote87/esign/v2/model"

    sv := envelopes.New(credential)

    f1, err := ioutil.ReadFile("letter.pdf")
    if err != nil {
        return nil, err
    }
    f2, err := ioutil.ReadFile("contract.pdf")
    if err != nil {
        return nil, err
    }

    env := &model.EnvelopeDefinition{
        EmailSubject: "[Go eSignagure SDK] - Please sign this doc",
        EmailBlurb:   "Please sign this test document",
        Status:       "sent",
        Documents: []model.Document{
            {
                DocumentBase64: f1,
                Name:           "invite letter.pdf",
                DocumentID:     "1",
            },
            {
                DocumentBase64: f2,
                Name:           "contract.pdf",
                DocumentID:     "2",
            },
        },
        Recipients: &model.Recipients{
            Signers: []model.Signer{
                {
                    Email:             email,
                    EmailNotification: nil,
                    Name:              name,
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
                                Width:  300,
                                Height: 150,
                            },
                        },
                    },
                },
            },
        },
    }
    envSummary, err := sv.Create(env).Do(context.Background())
```

## Testing

To test the package with your DocuSign Sandbox, obtain a token from
the [DocuSign OAuth2 Generator](https://developers.docusign.com/oauth-token-generator)
prior to running go test.  Set the environment variable DOCUSIGN_Token
to the generated access token and run go test.  This will read through
a list of user folders and all templates.

```sh
$ export DOCUSIGN_Token=eyJ0eXAiO...
$ go test
```