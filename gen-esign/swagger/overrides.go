// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package swagger

import (
	"strings"

	"github.com/jfcote87/esign"
)

// ServiceNameOverride provides map of new names x-ds-service
// value
var ServiceNameOverride = map[string]string{
	"Groups": "UserGroups",
}

// OperationSkipList contains operations ignore
// due to deprecation or incomplete definitions
var OperationSkipList = map[string]bool{
	"v2:OAuth2_PostRevoke":   true, // PostRevoke and PostToken are implemented in legacy package
	"v2:OAuth2_PostToken":    true,
	"v2.1:OAuth2_PostRevoke": true, // PostRevoke and PostToken are implemented in legacy package
	"v2.1:OAuth2_PostToken":  true,
}

// ResourceMaps provides links from the operation's tags to
// the service
var ResourceMaps = map[esign.APIVersion]map[string]string{
	esign.AdminV2: {
		"AccountSettingsExport":      "BulkOperations",
		"AccountSettingsImport":      "BulkOperations",
		"SingleAccountUserImport":    "BulkOperations",
		"UserExport":                 "BulkOperations",
		"UserImport":                 "BulkOperations",
		"IdentityProviders":          "IdentityProviders",
		"Organization":               "Organization",
		"ReservedDomains":            "ReservedDomains",
		"eSignUserManagement":        "UserManagement",
		"MultiProductUserManagement": "UserManagement",
		"Users":                      "UserManagement",
	},

	esign.RoomsV2: {
		"Accounts":                 "Accounts",
		"Documents":                "Documents",
		"ESignPermissionProfiles":  "ESignPermissionProfiles",
		"Fields":                   "Fields",
		"ExternalFormFillSessions": "Forms",
		"FormDetails":              "Forms",
		"FormGroups":               "Forms",
		"FormLibraries":            "Forms",
		"ClosingStatuses":          "GlobalResources",
		"ContactSides":             "GlobalResources",
		"Countries":                "GlobalResources",
		"Currencies":               "GlobalResources",
		"FinancingTypes":           "GlobalResources",
		"OriginsOfLeads":           "GlobalResources",
		"PropertyTypes":            "GlobalResources",
		"RoomContactTypes":         "GlobalResources",
		"SellerDecisionTypes":      "GlobalResources",
		"SpecialCircumstanceTypes": "GlobalResources",
		"States":                   "GlobalResources",
		"TaskDateTypes":            "GlobalResources",
		"TaskResponsibilityTypes":  "GlobalResources",
		"TaskStatuses":             "GlobalResources",
		"TimeZones":                "GlobalResources",
		"TransactionSides":         "GlobalResources",
		"Offices":                  "Offices",
		"Regions":                  "Regions",
		"Roles":                    "Roles",
		"RoomFolderss":             "Rooms",
		"Rooms":                    "Rooms",
		"RoomTemplates":            "RoomTemplates",
		"TaskLists":                "TaskLists",
		"TaskListTemplates":        "TaskLists",
		"Users":                    "Users",
	},
	esign.ClickV1: {
		"Uncategorized": "Click",
		"ClickWraps":    "Click",
	},
	esign.MonitorV2: {
		"DataSet": "Monitor",
	},
	esign.APIv2: {
		"AccountBrands":                         "Accounts",
		"AccountConsumerDisclosures":            "Accounts",
		"AccountCustomFields":                   "Accounts",
		"AccountPasswordRules":                  "Accounts",
		"AccountPermissionProfiles":             "Accounts",
		"Accounts":                              "Accounts",
		"AccountSealProviders":                  "Accounts",
		"AccountSignatureProviders":             "Accounts",
		"AccountTabSettings":                    "Accounts",
		"AccountWatermarks":                     "Accounts",
		"ConnectSecret":                         "Accounts",
		"ENoteConfigurations":                   "Accounts",
		"IdentityVerifications":                 "Accounts",
		"Authentication":                        "Authentication",
		"UserSocialAccountLogins":               "Authentication",
		"BillingPlans":                          "Billing",
		"Invoices":                              "Billing",
		"Payments":                              "Billing",
		"BulkEnvelopes":                         "BulkEnvelopes",
		"BulkSend":                              "BulkEnvelopes",
		"EnvelopeBulkRecipients":                "EnvelopeBulkRecipients",
		"CloudStorage":                          "CloudStorage",
		"CloudStorageProviders":                 "CloudStorage",
		"ConnectConfigurations":                 "Connect",
		"ConnectEvents":                         "Connect",
		"CustomTabs":                            "CustomTabs",
		"RequestLogs":                           "Diagnostics",
		"Resources":                             "Diagnostics",
		"Services":                              "Diagnostics",
		"BCCEmailArchive":                       "EmailArchive",
		"ChunkedUploads":                        "Envelopes",
		"Comments":                              "Envelopes",
		"DocumentResponsiveHtmlPreview":         "Envelopes",
		"EnvelopeAttachments":                   "Envelopes",
		"EnvelopeConsumerDisclosures":           "Envelopes",
		"EnvelopeCustomFields":                  "Envelopes",
		"EnvelopeDocumentFields":                "Envelopes",
		"EnvelopeDocumentHtmlDefinitions":       "Envelopes",
		"EnvelopeDocuments":                     "Envelopes",
		"EnvelopeDocumentTabs":                  "Envelopes",
		"EnvelopeDocumentVisibility":            "Envelopes",
		"EnvelopeEmailSettings":                 "Envelopes",
		"EnvelopeFormData":                      "Envelopes",
		"EnvelopeHtmlDefinitions":               "Envelopes",
		"EnvelopeLocks":                         "Envelopes",
		"EnvelopePublish":                       "Envelopes",
		"EnvelopeRecipients":                    "Envelopes",
		"EnvelopeRecipientTabs":                 "Envelopes",
		"Envelopes":                             "Envelopes",
		"EnvelopeTemplates":                     "Envelopes",
		"EnvelopeTransferRules":                 "Envelopes",
		"EnvelopeViews":                         "Envelopes",
		"NotaryJournals":                        "Envelopes",
		"ResponsiveHtmlPreview":                 "Envelopes",
		"TabsBlob":                              "Envelopes",
		"Folders":                               "Folders",
		"PaymentGatewayAccounts":                "Payments",
		"PowerFormData":                         "PowerForms",
		"PowerForms":                            "PowerForms",
		"SigningGroups":                         "SigningGroups",
		"SigningGroupUsers":                     "SigningGroups",
		"TemplateBulkRecipients":                "Templates",
		"TemplateCustomFields":                  "Templates",
		"TemplateDocumentFields":                "Templates",
		"TemplateDocumentHtmlDefinitions":       "Templates",
		"TemplateDocumentResponsiveHtmlPreview": "Templates",
		"TemplateDocuments":                     "Templates",
		"TemplateDocumentTabs":                  "Templates",
		"TemplateDocumentVisibility":            "Templates",
		"TemplateHtmlDefinitions":               "Templates",
		"TemplateLocks":                         "Templates",
		"TemplateRecipients":                    "Templates",
		"TemplateRecipientTabs":                 "Templates",
		"TemplateResponsiveHtmlPreview":         "Templates",
		"Templates":                             "Templates",
		"TemplateViews":                         "Templates",
		"GroupBrands":                           "UserGroups",
		"Groups":                                "UserGroups",
		"GroupUsers":                            "UserGroups",
		"Contacts":                              "Users",
		"UserCustomSettings":                    "Users",
		"UserProfiles":                          "Users",
		"Users":                                 "Users",
		"UserSignatures":                        "Users",
		"WorkspaceItems":                        "Workspaces",
		"Workspaces":                            "Workspaces",
	},
	esign.APIv21: {
		"AccountBrands":                         "Accounts",
		"AccountConsumerDisclosures":            "Accounts",
		"AccountCustomFields":                   "Accounts",
		"AccountPasswordRules":                  "Accounts",
		"AccountPermissionProfiles":             "Accounts",
		"Accounts":                              "Accounts",
		"AccountSealProviders":                  "Accounts",
		"AccountSignatureProviders":             "Accounts",
		"AccountSignatures":                     "Accounts",
		"AccountTabSettings":                    "Accounts",
		"AccountWatermarks":                     "Accounts",
		"ENoteConfigurations":                   "Accounts",
		"FavoriteTemplates":                     "Accounts",
		"IdentityVerifications":                 "Accounts",
		"BillingPlans":                          "Billing",
		"Invoices":                              "Billing",
		"Payments":                              "Billing",
		"BulkSend":                              "BulkEnvelopes",
		"CloudStorage":                          "CloudStorage",
		"CloudStorageProviders":                 "CloudStorage",
		"ConnectConfigurations":                 "Connect",
		"ConnectEvents":                         "Connect",
		"CustomTabs":                            "CustomTabs",
		"RequestLogs":                           "Diagnostics",
		"Resources":                             "Diagnostics",
		"Services":                              "Diagnostics",
		"BCCEmailArchive":                       "EmailArchive",
		"ChunkedUploads":                        "Envelopes",
		"Comments":                              "Envelopes",
		"DocumentResponsiveHtmlPreview":         "Envelopes",
		"EnvelopeAttachments":                   "Envelopes",
		"EnvelopeConsumerDisclosures":           "Envelopes",
		"EnvelopeCustomFields":                  "Envelopes",
		"EnvelopeDocumentFields":                "Envelopes",
		"EnvelopeDocumentHtmlDefinitions":       "Envelopes",
		"EnvelopeDocuments":                     "Envelopes",
		"EnvelopeDocumentTabs":                  "Envelopes",
		"EnvelopeDocumentVisibility":            "Envelopes",
		"EnvelopeEmailSettings":                 "Envelopes",
		"EnvelopeFormData":                      "Envelopes",
		"EnvelopeHtmlDefinitions":               "Envelopes",
		"EnvelopeLocks":                         "Envelopes",
		"EnvelopePublish":                       "Envelopes",
		"EnvelopeRecipients":                    "Envelopes",
		"EnvelopeRecipientTabs":                 "Envelopes",
		"Envelopes":                             "Envelopes",
		"EnvelopeTemplates":                     "Envelopes",
		"EnvelopeTransferRules":                 "Envelopes",
		"EnvelopeViews":                         "Envelopes",
		"EnvelopeWorkflowDefinition":            "Envelopes",
		"NotaryJournals":                        "Envelopes",
		"ResponsiveHtmlPreview":                 "Envelopes",
		"TabsBlob":                              "Envelopes",
		"Folders":                               "Folders",
		"NotaryJurisdiction":                    "Notary",
		"Notary":                                "Notary",
		"PaymentGatewayAccounts":                "Payments",
		"PowerFormData":                         "PowerForms",
		"PowerForms":                            "PowerForms",
		"SigningGroups":                         "SigningGroups",
		"SigningGroupUsers":                     "SigningGroups",
		"TemplateBulkRecipients":                "Templates",
		"TemplateCustomFields":                  "Templates",
		"TemplateDocumentFields":                "Templates",
		"TemplateDocumentHtmlDefinitions":       "Templates",
		"TemplateDocumentResponsiveHtmlPreview": "Templates",
		"TemplateDocuments":                     "Templates",
		"TemplateDocumentTabs":                  "Templates",
		"TemplateDocumentVisibility":            "Templates",
		"TemplateHtmlDefinitions":               "Templates",
		"TemplateLocks":                         "Templates",
		"TemplateRecipients":                    "Templates",
		"TemplateRecipientTabs":                 "Templates",
		"TemplateResponsiveHtmlPreview":         "Templates",
		"Templates":                             "Templates",
		"TemplateViews":                         "Templates",
		"GroupBrands":                           "UserGroups",
		"Groups":                                "UserGroups",
		"GroupUsers":                            "UserGroups",
		"Contacts":                              "Users",
		"UserCustomSettings":                    "Users",
		"UserProfiles":                          "Users",
		"Users":                                 "Users",
		"UserSignatures":                        "Users",
		"WorkspaceItems":                        "Workspaces",
		"Workspaces":                            "Workspaces",
	},
}

// GetServicePrefixes provides a list of prefixes
// to remove from Tags for creating op FuncName
func GetServicePrefixes(service string) []string {
	var list = []string{service}
	if strings.HasSuffix(service, "s") {
		list = append(list, service[:len(service)-1])
	}
	if service == "BulkEnvelopes" {
		list = append(list, "EnvelopeBulk")
	}
	return list
}

// IsUploadFilesOperation checks whether the operation
// allow multipart file uploads
func IsUploadFilesOperation(opID string) bool {
	switch opID {
	case "v2:Envelopes_PostEnvelopes", "v2.1:Envelopes_PostEnvelopes":
		return true
	case "v2:Templates_PostTemplates", "v2.1:Templates_PostTemplates":
		return true
	case "v2:UserSignatures_PostUserSignatures", "v2.1:UserSignatures_PostUserSignatures":
		return true
	}
	return false
}

// GetParameterOverrides returns a map of all
// parameter Type overrides The returned map is
// map[<operationID>]map[<FieldName>]<GoType>
//
// I generated much of this list by searching the swagger file
// with the following rules:
// - if **true** found in description, set parameter type to bool.
// - if field name ends in "_date", set time.Time
// - if field names is count, start_position, etc, set to int
// - if description starts with "Comma separated list", set to ...string
//
// I eyeballed the doc as best I could so please let me know of any additions
// or corrections.
func GetParameterOverrides() map[string]map[string]string {
	return map[string]map[string]string{
		"AccountCustomFields_DeleteAccountCustomFields": {
			"apply_to_templates": "bool",
		},
		"AccountCustomFields_PostAccountCustomFields": {
			"apply_to_templates": "bool",
		},
		"AccountCustomFields_PutAccountCustomFields": {
			"apply_to_templates": "bool",
		},
		"Accounts_GetAccount": {
			"include_account_settings": "bool",
		},
		"Accounts_PostAccounts": {
			"preview_billing_plan": "bool",
		},
		"BillingInvoices_GetBillingInvoices": {
			"from_date": "time.Time",
			"to_date":   "time.Time",
		},
		"BillingPayments_GetPaymentList": {
			"from_date": "time.Time",
			"to_date":   "time.Time",
		},
		"BillingPlan_GetBillingPlan": {
			"include_credit_card_information": "bool",
			"include_metadata":                "bool",
			"include_successor_plans":         "bool",
		},
		"BillingPlan_PutBillingPlan": {
			"preview_billing_plan": "bool",
		},
		"Brand_GetBrand": {
			"include_external_references": "bool",
			"include_logos":               "bool",
		},
		"BrandResources_GetBrandResources": {
			"return_master": "bool",
		},
		"Brands_GetBrands": {
			"exclude_distributor_brand": "bool",
			"include_logos":             "bool",
		},
		"BulkEnvelopes_d": {
			"count":          "int",
			"include":        "...string",
			"start_position": "int",
		},
		"BulkEnvelopes_GetEnvelopesBulk": {
			"count":          "int",
			"include":        "...string",
			"start_position": "int",
		},
		"ChunkedUploads_GetChunkedUpload": {
			"include": "...string",
		},
		"CloudStorageFolder_GetCloudStorageFolder": {
			"count":          "int",
			"start_position": "int",
		},
		"CloudStorageFolder_GetCloudStorageFolderAll": {
			"cloud_storage_folder_path": "...string",
			"count":                     "int",
			"start_position":            "int",
		},
		"ConnectFailures_GetConnectLogs": {
			"from_date": "time.Time",
			"to_date":   "time.Time",
		},
		"ConnectLog_GetConnectLog": {
			"additional_info": "bool",
		},
		"ConnectLog_GetConnectLogs": {
			"from_date": "time.Time",
			"to_date":   "time.Time",
		},
		"Connect_GetConnectUsers": {
			"count":               "int",
			"list_included_users": "bool",
			"start_position":      "int",
			"status":              "...string",
		},
		"Documents_GetDocument": {
			"certificate":  "bool",
			"encrypt":      "bool",
			"show_changes": "bool",
			"watermark":    "bool",
		},
		"Documents_GetTemplateDocument": {
			"encrypt":      "bool",
			"show_changes": "bool",
		},
		"Documents_PutDocument": {
			"apply_document_fields": "bool",
		},
		"Documents_PutDocuments": {
			"apply_document_fields": "bool",
			"persist_tabs":          "bool",
		},
		"Documents_PutTemplateDocument": {
			"apply_document_fields":  "bool",
			"is_envelope_definition": "bool",
		},
		"Documents_PutTemplateDocuments": {
			"apply_document_fields": "bool",
			"persist_tabs":          "bool",
		},
		"Envelopes_GetEnvelope": {
			"advanced_update": "bool",
		},
		"Envelopes_GetEnvelopes": {
			"count":                   "int",
			"envelope_ids":            "...string",
			"folder_ids":              "...string",
			"from_date":               "time.Time",
			"intersecting_folder_ids": "...string",
			"powerformids":            "...string",
			"start_position":          "int",
			"status":                  "...string",
			"to_date":                 "time.Time",
			"transaction_ids":         "...string",
		},
		"Envelopes_PostEnvelopes": {
			"change_routing_order": "bool",
			"merge_roles_on_draft": "bool",
		},
		"Envelopes_PutEnvelope": {
			"advanced_update": "bool",
			"resend_envelope": "bool",
		},
		"Envelopes_PutStatus": {
			"from_date":      "time.Time",
			"start_position": "int",
			"to_date":        "time.Time",
		},
		"Folders_GetFolderItems": {
			"from_date":      "time.Time",
			"start_position": "int",
			"to_date":        "time.Time",
		},
		"Folders_GetFolders": {
			"start_position": "int",
		},
		"Groups_GetGroupUsers": {
			"count":          "int",
			"start_position": "int",
		},
		"Groups_GetGroups": {
			"count":          "int",
			"start_position": "int",
		},
		"LoginInformation_GetLoginInformation": {
			"include_account_id_guid": "bool",
		},
		"Pages_GetPageImage": {
			"dpi":          "int",
			"max_height":   "int",
			"max_width":    "int",
			"show_changes": "bool",
		},
		"Pages_GetPageImages": {
			"count":          "int",
			"dpi":            "int",
			"max_height":     "int",
			"max_width":      "int",
			"nocache":        "bool",
			"show_changes":   "bool",
			"start_position": "int",
		},
		"Pages_GetTemplatePageImage": {
			"dpi":          "int",
			"max_height":   "int",
			"max_width":    "int",
			"show_changes": "bool",
		},
		"Pages_GetTemplatePageImages": {
			"count":          "int",
			"dpi":            "int",
			"max_height":     "int",
			"max_width":      "int",
			"nocache":        "bool",
			"show_changes":   "bool",
			"start_position": "int",
		},
		"PermissionProfiles_GetPermissionProfile": {
			"include": "...string",
		},
		"PermissionProfiles_PostPermissionProfiles": {
			"include": "...string",
		},
		"PermissionProfiles_PutPermissionProfiles": {
			"include": "...string",
		},
		"PowerForms_GetPowerFormFormData": {
			"from_date": "time.Time",
			"to_date":   "time.Time",
		},
		"PowerForms_GetPowerFormsList": {
			"from_date": "time.Time",
			"to_date":   "time.Time",
		},
		"PowerForms_GetPowerFormsSenders": {
			"start_position": "int",
		},
		"Recipients_GetBulkRecipients": {
			"include_tabs":   "bool",
			"start_position": "int",
		},
		"Recipients_GetRecipientInitialsImage": {
			"include_chrome": "bool",
		},
		"Recipients_GetRecipientSignatureImage": {
			"include_chrome": "bool",
		},
		"Recipients_GetRecipientTabs": {
			"include_anchor_tab_locations": "bool",
		},
		"Recipients_GetRecipients": {
			"include_anchor_tab_locations": "bool",
			"include_extended":             "bool",
			"include_tabs":                 "bool",
		},
		"Recipients_GetTemplateBulkRecipients": {
			"include_tabs":   "bool",
			"start_position": "int",
		},
		"Recipients_GetTemplateRecipientTabs": {
			"include_anchor_tab_locations": "bool",
		},
		"Recipients_GetTemplateRecipients": {
			"include_anchor_tab_locations": "bool",
			"include_extended":             "bool",
			"include_tabs":                 "bool",
		},
		"Recipients_PostRecipients": {
			"resend_envelope": "bool",
		},
		"Recipients_PostTemplateRecipients": {
			"resend_envelope": "bool",
		},
		"Recipients_PutRecipients": {
			"resend_envelope": "bool",
		},
		"Recipients_PutTemplateRecipients": {
			"resend_envelope": "bool",
		},
		"SearchFolders_GetSearchFolderContents": {
			"all":                "bool",
			"count":              "int",
			"from_date":          "time.Time",
			"include_recipients": "bool",
			"start_position":     "int",
			"to_date":            "time.Time",
		},
		"SharedAccess_GetSharedAccess": {
			"count":          "int",
			"folder_ids":     "...string",
			"start_position": "int",
			"user_ids":       "...string",
		},
		"SharedAccess_PutSharedAccess": {
			"user_ids": "...string",
		},
		"SigningGroups_GetSigningGroups": {
			"include_users": "bool",
		},
		"Tabs_GetTabDefinitions": {
			"custom_tab_only": "bool",
		},
		"Templates_GetDocumentTemplates": {
			"include": "...string",
		},
		"Templates_GetTemplate": {
			"include": "...string",
		},
		"Templates_GetTemplates": {
			"count":              "int",
			"folder_ids":         "...string",
			"from_date":          "time.Time",
			"include":            "...string",
			"modified_from_date": "time.Time",
			"modified_to_date":   "time.Time",
			"start_position":     "int",
			"to_date":            "time.Time",
			"used_from_date":     "time.Time",
			"used_to_date":       "time.Time",
		},
		"UserSignatures_GetUserSignatureImage": {
			"include_chrome": "bool",
		},
		"UserSignatures_PutUserSignatureById": {
			"close_existing_signature": "bool",
		},
		"User_GetUser": {
			"additional_info": "bool",
		},
		"Users_GetUsers": {
			"additional_info":              "bool",
			"count":                        "int",
			"include_usersettings_for_csv": "bool",
			"start_position":               "int",
			"status":                       "...string",
		},
		"WorkspaceFilePages_GetWorkspaceFilePages": {
			"count":          "int",
			"dpi":            "int",
			"max_height":     "int",
			"max_width":      "int",
			"start_position": "int",
		},
		"WorkspaceFile_GetWorkspaceFile": {
			"is_download": "bool",
			"pdf_version": "bool",
		},
		"WorkspaceFolder_GetWorkspaceFolder": {
			"count":               "int",
			"include_files":       "bool",
			"include_sub_folders": "bool",
			"include_thumbnails":  "bool",
			"include_user_detail": "bool",
			"start_position":      "int",
		},
	}
}

// DownloadAddition describes an addiion to an operation that
// will return an esign.DownloadFile when the Accept header
// is set to the MimeType value.
type DownloadAddition struct {
	Name     string
	MimeType string
	Comments []string
}

// GetDownloadAdditions returns the special download funcs for
// an operation.  These are gleaned from the documentation
// not the swagger file.
func GetDownloadAdditions(opID string) []DownloadAddition {
	switch opID {
	case "BillingInvoices_GetBillingInvoice":
		return []DownloadAddition{
			{
				Name:     "PDF",
				MimeType: "application/pdf",
				Comments: []string{"PDF returns a pdf version of the invoice by setting", "the Accept header to application/pdf", "", "**not included in swagger definition"},
			},
		}
	case "APIRequestLog_GetRequestLogs":
		return []DownloadAddition{
			{
				Name:     "Zip",
				MimeType: "application/zip",
				Comments: []string{"Zip returns a zip file containing log files by setting", "the Accept header to application/zip", "", "**not included in swagger definition"},
			},
		}
	}
	return nil
}

// TabDefs return a list of embeded tab structs based upon the version
func TabDefs(apiname string, defMap map[string]Definition, overrides map[string]map[string]string) []Definition {
	switch apiname {
	case esign.APIv2.Name():
		return TabDefsV2(defMap, overrides)
	case esign.APIv21.Name():
		return TabDefsV21(defMap, overrides)
	}
	return make([]Definition, 0)
}

// TabDefsV2 creates a list definitions for embedded tab structs from the defMap parameter.
// overrides is updated with new override entries to allow tab definitions to generate.
func TabDefsV2(defMap map[string]Definition, overrides map[string]map[string]string) []Definition {
	// list of tab objects
	var tabObjects = []string{
		"approve",
		"checkbox",
		"company",
		"dateSigned",
		"date",
		"decline",
		"emailAddress",
		"email",
		"envelopeId",
		"firstName",
		"formulaTab",
		"fullName",
		"initialHere",
		"lastName",
		"list",
		"notarize",
		"note",
		"number",
		"radioGroup",
		"signerAttachment",
		"signHere",
		"ssn",
		"text",
		"tabGroup",
		"title",
		"view",
		"zip",
	}

	// list of types of tabs
	tabDefs := map[string]Definition{
		"Base": {
			ID:          "TabBase",
			Name:        "TabBase",
			Type:        "Object",
			Description: "contains common fields for all tabs",
			Summary:     "contains common fields for all tabs",
			Category:    "",
		},
		"Position": {
			ID:          "TabPosition",
			Name:        "TabPosition",
			Type:        "Object",
			Description: "contains common fields for all tabs that can position themselves on document",
			Summary:     "contains common fields for all tabs that can position themselves on document",
			Category:    "",
		},
		"Style": {
			ID:          "TabStyle",
			Name:        "TabStyle",
			Type:        "Object",
			Description: "contains common fields for all tabs that can set a display style",
			Summary:     "contains common fields for all tabs that can set a display style",
			Category:    "",
		},
		"Value": {
			ID:          "TabValue",
			Name:        "TabValue",
			Type:        "Object",
			Description: "add Value() func to tab",
			Summary:     "add Value() func to tab",
			Category:    "",
		},
	}
	// list of fields for each tab type
	tabFields := map[string][]string{
		"Base": {
			"conditionalParentLabel",
			"conditionalParentValue",
			"documentId",
			"recipientId",
		},
		"Position": {
			"anchorCaseSensitive",
			"anchorHorizontalAlignment",
			"anchorIgnoreIfNotPresent",
			"anchorMatchWholeWord",
			"anchorString",
			"anchorUnits",
			"anchorXOffset",
			"anchorYOffset",
			"customTabId",
			"errorDetails",
			"mergeField",
			"pageNumber",
			"status",
			"tabId",
			"tabLabel",
			"tabOrder",
			"templateLocked",
			"templateRequired",
			"xPosition",
			"yPosition",
		},
		"Style": {
			"bold",
			"font",
			"fontColor",
			"fontSize",
			"italic",
			"name",
			"underline",
		},
		"Value": {
			"value",
		},
	}
	// loop thru each tab definition
	for _, tabname := range tabObjects {
		dx := defMap["#/definitions/"+tabname]
		// create map of fields for easy lookup
		xmap := make(map[string]bool)
		// NOTE:  add tabLabel to notary and name to view.
		//This seems to be an error  in swagger definition.
		// TODO: remove when fixed in swagger file
		if tabname == "notarize" {
			xmap["tabLabel"] = true
		}
		if tabname == "view" {
			xmap["name"] = true
		}

		for _, f := range dx.Fields {
			xmap[f.Name] = true
		}
		// Get Overrides for this tab definition
		defOverrides, ok := overrides[dx.ID]
		if !ok {
			defOverrides = make(map[string]string)
			overrides[dx.ID] = defOverrides
		}

		if xmap["width"] {
			defOverrides["width"] = "string"
		}
		if xmap["height"] {
			defOverrides["height"] = "string"
		}

		memberOf := make([]string, 0) // tab types for this tab
		// Loop thru each tab type
		for _, nm := range []string{"Base", "Position", "Style", "Value"} {
			// check for match by checking for existence of each field
			isType := true
			for _, s := range tabFields[nm] {
				if isType = xmap[s]; !isType {
					break
				}
			}
			// if match, mark override for each field
			if isType {
				for _, f := range tabFields[nm] {
					// Definition.StructFields() will know to skip
					// outputting this field for the type def
					defOverrides[f] = "-"
				}
				memberOf = append(memberOf, "Tab"+nm)
			}
		}
		// create override.  Definition.StructFields will know to output
		// these embedded types.
		if len(memberOf) > 0 {
			defOverrides["TABS"] = strings.Join(memberOf, ",")
		}
	}

	// Get field definitions for embedded types.  Assume that
	// the Text tab meets all definitions so copy appropriate field
	// from its definition
	txtDef := defMap["#/definitions/text"]
	xmap := make(map[string]Field)
	for _, f := range txtDef.Fields {
		xmap[f.Name] = f

	}
	results := make([]Definition, 0)
	for _, tabDefName := range []string{"Base", "Position", "Style", "Value"} {
		ndef := tabDefs[tabDefName]
		for _, s := range tabFields[tabDefName] {
			fx, ok := xmap[s]
			if ok {
				ndef.Fields = append(ndef.Fields, fx)
			}
		}
		results = append(results, ndef)

	}
	// return list of embedded tab definitions
	return results
}

// TabDefsV21 creates a list definitions for embedded tab structs from the defMap parameter.
// overrides is updated with new override entries to allow tab definitions to generate.
func TabDefsV21(defMap map[string]Definition, overrides map[string]map[string]string) []Definition {
	// list of tab objects
	var tabObjects = []string{
		"approve",
		"checkbox",
		"company",
		"commentThread",
		"commissionCounty",
		"commissionExpiration",
		"commissionNumber",
		"commissionState",
		"currency",
		"dateSigned",
		"date",
		"decline",
		"draw",
		"emailAddress",
		"email",
		"envelopeId",
		"firstName",
		"formulaTab",
		"fullName",
		"initialHere",
		"lastName",
		"list",
		"notarize",
		"notearySeal",
		"note",
		"number",
		"polyLineOverlay",
		"prefill",
		"radioGroup",
		"signerAttachment",
		"signHere",
		"smartSection",
		"ssn",
		"text",
		"tabGroup",
		"title",
		"view",
		"zip",
	}

	// list of types of tabs
	tabDefs := map[string]Definition{
		"Base": {
			ID:          "TabBase",
			Name:        "TabBase",
			Type:        "Object",
			Description: "contains common fields for all tabs",
			Summary:     "contains common fields for all tabs",
			Category:    "",
		},
		"GuidedForm": {
			ID:          "TabGuidedForm",
			Name:        "TabGuidedForm",
			Type:        "Object",
			Description: "contains common fields for all text box tabs",
			Summary:     "contains common fields for all text box tabs",
			Category:    "",
		},
		"Position": {
			ID:          "TabPosition",
			Name:        "TabPosition",
			Type:        "Object",
			Description: "contains common fields for all tabs that can position themselves on document",
			Summary:     "contains common fields for all tabs that can position themselves on document",
			Category:    "",
		},
		"Style": {
			ID:          "TabStyle",
			Name:        "TabStyle",
			Type:        "Object",
			Description: "contains common fields for all tabs that can set a display style",
			Summary:     "contains common fields for all tabs that can set a display style",
			Category:    "",
		},
		"Value": {
			ID:          "TabValue",
			Name:        "TabValue",
			Type:        "Object",
			Description: "add Value() func to tab",
			Summary:     "add Value() func to tab",
			Category:    "",
		},
	}
	// list of fields for each tab type
	tabFields := map[string][]string{
		"Base": {
			"conditionalParentLabel",
			"conditionalParentLabelMetadata",
			"conditionalParentValue",
			"conditionalParentValueMetadata",
			"documentId",
			"documentIdMetadata",
			"recipientId",
			"recipientIdMetadata",
			"recipientIdGuid",
			"recipientIdGuidMetadata",
			"tabGroupLabels",
			"tabGroupLabelsMetadata",
			"tabType",
			"tabTypeMetadata",
		},
		"GuidedForm": {
			"formOrder",
			"formOrderMetadata",
			"formPageLabel",
			"formPageLabelMetadata",
			"formPageNumber",
			"formPageNumberMetadata",
		},
		"Position": {
			"anchorCaseSensitive",
			"anchorCaseSensitiveMetadata",
			"anchorHorizontalAlignment",
			"anchorHorizontalAlignmentMetadata",
			"anchorIgnoreIfNotPresent",
			"anchorIgnoreIfNotPresentMetadata",
			"anchorMatchWholeWord",
			"anchorMatchWholeWordMetadata",
			"anchorString",
			"anchorStringMetadata",
			"anchorTabProcessorVersion",
			"anchorTabProcessorVersionMetadata",
			"anchorUnits",
			"anchorUnitsMetadata",
			"anchorXOffset",
			"anchorXOffsetMetadata",
			"anchorYOffset",
			"anchorYOffsetMetadata",
			"customTabId",
			"customTabIdMetadata",
			"errorDetails",
			"mergeField",
			"pageNumber",
			"pageNumberMetadata",
			"status",
			"statusMetadata",
			"tabId",
			"tabIdMetadata",
			"tabLabel",
			"tabLabelMetadata",
			"tabOrder",
			"tabOrderMetadata",
			"templateLocked",
			"templateLockedMetadata",
			"templateRequired",
			"templateRequiredMetadata",
			"xPosition",
			"xPositionMetadata",
			"yPosition",
			"yPositionMetadata",
		},
		"Style": {
			"bold",
			"boldMetadata",
			"font",
			"fontMetadata",
			"fontColor",
			"fontColorMetadata",
			"fontSize",
			"fontSizeMetadata",
			"italic",
			"italicMetadata",
			"name",
			"nameMetadata",
			"underline",
			"underlineMetadata",
		},
		"Value": {
			"value",
			"valueMetadata",
		},
	}
	// loop thru each tab definition
	for _, tabname := range tabObjects {
		dx := defMap["#/definitions/"+tabname]
		// create map of fields for easy lookup
		xmap := make(map[string]bool)
		// NOTE:  add tabLabel to notary and name to view.
		//This seems to be an error  in swagger definition.
		// TODO: remove when fixed in swagger file
		if tabname == "notarize" {
			xmap["tabLabel"] = true
		}
		if tabname == "view" {
			xmap["name"] = true
		}
		if tabname == "radioGroup" {
			xmap["tabGroupLabels"] = true
			xmap["tabGroupLabelsMetadata"] = true
		}

		for _, f := range dx.Fields {
			xmap[f.Name] = true
		}
		// Get Overrides for this tab definition
		defOverrides, ok := overrides[dx.ID]
		if !ok {
			defOverrides = make(map[string]string)
			overrides[dx.ID] = defOverrides
		}

		memberOf := make([]string, 0) // tab types for this tab
		// Loop thru each tab type
		for _, nm := range []string{"Base", "GuidedForm", "Position", "Style", "Value"} {
			// check for match by checking for existence of each field
			isType := true
			for _, s := range tabFields[nm] {
				if isType = xmap[s]; !isType {
					break
				}
			}
			// if match, mark override for each field
			if isType {
				for _, f := range tabFields[nm] {
					// Definition.StructFields() will know to skip
					// outputting this field for the type def
					defOverrides[f] = "-"
				}
				memberOf = append(memberOf, "Tab"+nm)
			}
		}
		// create override.  Definition.StructFields will know to output
		// these embedded types.
		if len(memberOf) > 0 {
			defOverrides["TABS"] = strings.Join(memberOf, ",")
		}
	}

	// Get field definitions for embedded types.  Assume that
	// the Text tab meets all definitions so copy appropriate field
	// from its definition
	txtDef := defMap["#/definitions/text"]
	xmap := make(map[string]Field)
	for _, f := range txtDef.Fields {
		xmap[f.Name] = f

	}
	results := make([]Definition, 0)
	for _, tabDefName := range []string{"Base", "GuidedForm", "Position", "Style", "Value"} {
		ndef := tabDefs[tabDefName]
		for _, s := range tabFields[tabDefName] {
			fx, ok := xmap[s]
			if ok {
				ndef.Fields = append(ndef.Fields, fx)
			}
		}
		results = append(results, ndef)

	}
	// return list of embedded tab definitions
	return results
}

// CustomCode provides additional specific code for a package
func CustomCode(apiname string) string {
	switch apiname {
	case esign.APIv2.Name():
		return esignTabCustomCode
	case esign.APIv21.Name():
		return esignTabCustomCode
	}
	return ""
}

// esignTabCustomCode is lines of code to append to esign v2 and v2.1 model.go
const esignTabCustomCode = `// Bool represents a DocuSign boolean value which is either a string "true" or a string "false".  This construct
// allows the setting of a false value that will not be omitted during a JSON Marshal.  Use the
// DSBool function to set the proper values
type Bool string

// True returns the bool value of the Bool string
func (b Bool) True() bool {
	return strings.ToLower(string(b)) == string(TRUE)
}

// DSBool converts the boolean value to a Bool
func DSBool(b bool) Bool {
	if b {
		return TRUE
	}
	return FALSE
}

const (
	// True is the standard true value for a Bool
	TRUE Bool = "true"
	// False is the standard false value for a Bool
	FALSE Bool = "false"

	// REQUIRED_DEFAULT sets the default value for (BOOL) Required field on a tab. This
	// constant is kept from a previous version where the Requried tab field was defined
	// as a integer.
	REQUIRED_DEFAULT Bool = ""
	// REQUIRED_FALSE sets the default value for (BOOL) Required field on a tab
	REQUIRED_FALSE Bool = "false"
	// REQUIRED_TRUE sets the default value for (BOOL) Required field on a tab
	REQUIRED_TRUE Bool = "true"
)


// GetTabValues returns a NameValue list of all entry tabs
func GetTabValues(tabs Tabs) []NameValue {
	results := make([]NameValue, 0)
	for _, v := range tabs.CheckboxTabs {
		results = append(results, NameValue{Name: v.TabLabel, Value: fmt.Sprintf("%v", v.Selected)})
	}
	for _, v := range tabs.CompanyTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.DateSignedTabs {
		results = append(results, NameValue{Name: v.TabLabel, Value: v.Value})
	}
	for _, v := range tabs.DateTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.EmailTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.FormulaTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	results = append(results, getListTabValues(tabs.ListTabs)...)
	for _, v := range tabs.NoteTabs {
		results = append(results, NameValue{Name: v.TabLabel, Value: v.Value})
	}
	for _, v := range tabs.NumberTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	results = append(results, getRadioTabValues(tabs.RadioGroupTabs)...)
	for _, v := range tabs.SSNTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.TextTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.ZipTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	return results
}

func getListTabValues(tabs []List) []NameValue {
	var results []NameValue
	for _, lt := range tabs {
		var vals []string
		for _, item := range lt.ListItems {
			if item.Selected.True() {
				vals = append(vals, item.Value)
			}
		}
		results = append(results, NameValue{Name: lt.TabLabel, Value: strings.Join(vals, ",")})
	}
	return results
}

func getRadioTabValues(tabs []RadioGroup) []NameValue {
	var results []NameValue
	for _, rg := range tabs {
		for _, rb := range rg.Radios {
			if rb.Selected.True() {
				results = append(results, NameValue{Name: rg.GroupName, Value: rb.Value})
			}
		}
	}
	return results
}`

// PackageScopes return definitions of scopes for the package
func PackageScopes(ver esign.APIVersion) string {
	switch ver {
	case esign.AdminV2:
		return scopesAdmin
	case esign.RoomsV2:
		return scopesRooms
	case esign.ClickV1:
		return scopesClick
	}
	return ""
}

const (
	scopesAdmin = `
// For more infomation on how to use scopes, see https://developers.docusign.com/docs/admin-api/admin101/auth/
const (
	// OAuthScopeOrganizationRead required to get lists of organizations and organization data.
	OAuthScopeOrganizationRead = "organization_read"
	// OAuthScopeGroupRead required to get lists of groups associated with an account.
	OAuthScopeGroupRead = "group_read"
	// OAuthScopePermissionRead equired to get lists of permission profiles associated with an account.
	OAuthScopePermissionRead = "permission_read"
	// OAuthScopeUserRead required to read user data.
	OAuthScopeUserRead = "user_read"
	// OAuthScopeUserWrite required to update, create, or delete users.
	OAuthScopeUserWrite = "user_write"
	// OAuthScopeAccountRead required to get account data.
	OAuthScopeAccountRead = "account_read"
	// OAuthScopeDomainRead required to get data on claimed domains for the organization.
	OAuthScopeDomainRead = "domain_read"
	// OAuthScopeIdentityProviderRead required to get data on already-defined identity providers for an organization.
	OAuthScopeIdentityProviderRead = "identity_provider_read"
)
`
	scopesRooms = `
// For more infomation on how to use scopes, see https://developers.docusign.com/docs/rooms-api/rooms101/auth/
const (
	// OAuthScopeRead authorizes reading DocuSign Rooms data
	OAuthScopeRead = "dtr.rooms.read"
	// OAuthScopeWrite authorizes updating DocuSign Room data
	OAuthScopeWrite = "dtr.rooms.write"
	// OAuthScopeDocumentsRead authorizes reading of documents from DocuSign Rooms
	OAuthScopeDocumentsRead = "dtr.documents.read"
	// OAuthScopeDocumentsWrite authorizes writing documents to DocuSign Rooms
	OAuthScopeDocumentsWrite = "dtr.documents.write"
	// OAuthScopeProfileRead authorizes reading of profile data for accounts or signers associated with your company
	OAuthScopeProfileRead = "dtr.profile.read"
	// OAuthScopeProfileWrite authorizes writing profile data to accounts or signers associated with your company
	OAuthScopeProfileWrite = "dtr.profile.write"
	// OAuthScopeCompanyRead authorizes reading information from all rooms and profiles associated with your company
	OAuthScopeCompanyRead = "dtr.company.read"
	// OAuthScopeCompanyWrite authorizes writing information to all rooms and profiles associated with your company
	OAuthScopeCompanyWrite = "dtr.company.write"
	// OAuthScopeForms authorizes use of endpoints related to the Forms feature
	OAuthScopeForms = "room_forms"
)
`
	scopesClick = `
// For more infomation on how to use scopes, see https://developers.docusign.com/docs/click-api/click101/auth/
const (
	// OAuthScopeManage enables all clickwrap operations, including creating, sending, and updating clickwraps;
	// getting a list of clickwraps, creating user agreements, getting a list of users, and retrieving responses.
	OAuthScopeManage = "click.manage"
	// OAuthScopeSend required to send a new clickwrap or check for a previously sent one.
	OAuthScopeSend = "click.send"
)
`
)

// GetFieldOverrides returns a map of all
// field level type overrides for the esign
// api generation. The returned map is
// map[<structID>]map[<FieldName>]<GoType>
//
// In the specification, DocuSign lists every field as a string.
// I generated much of this list with the following rules.
// - definition properties with **true** or Boolean on the top lines' description are
//   assumed to be bools set to Bool
// - fields containing base64 in the name are assumed to be []byte
// - fields ending in DateTime are *time.Time//
// I eyeballed the doc as best I could so please let me know of any additions
// or corrections.
func GetFieldOverrides() map[string]map[string]string {
	return map[string]map[string]string{
		"AccountWatermarks": {
			"imageBase64": "[]byte",
		},
		"TabPosition": {
			"anchorCaseSensitive":      "Bool",
			"anchorIgnoreIfNotPresent": "Bool",
			"anchorMatchWholeWord":     "Bool",
			"templateLocked":           "Bool",
			"templateRequired":         "Bool",
		},
		"TabStyle": {
			"bold":      "Bool",
			"italic":    "Bool",
			"underline": "Bool",
		},
		"accessCodeFormat": {
			"formatRequired":           "Bool",
			"letterRequired":           "Bool",
			"numberRequired":           "Bool",
			"specialCharacterRequired": "Bool",
		},
		"accountBillingPlan": {
			"canUpgrade":    "Bool",
			"enableSupport": "Bool",
			"isDowngrade":   "Bool",
		},
		"accountBillingPlanResponse": {
			"billingAddressIsCreditCardAddress": "Bool",
		},
		"accountInformation": {
			"allowTransactionRooms":  "Bool",
			"canUpgrade":             "Bool",
			"envelopeSendingBlocked": "Bool",
			"isDowngrade":            "Bool",
		},
		"accountNotification": {
			"userOverrideEnabled": "Bool",
		},
		"accountPasswordRules": {
			"expirePassword":                         "Bool",
			"passwordIncludeDigit":                   "Bool",
			"passwordIncludeDigitOrSpecialCharacter": "Bool",
			"passwordIncludeLowerCase":               "Bool",
			"passwordIncludeSpecialCharacter":        "Bool",
			"passwordIncludeUpperCase":               "Bool",
		},
		"accountPasswordStrengthTypeOption": {
			"passwordIncludeDigit":                   "Bool",
			"passwordIncludeDigitOrSpecialCharacter": "Bool",
			"passwordIncludeLowerCase":               "Bool",
			"passwordIncludeSpecialCharacter":        "Bool",
			"passwordIncludeUpperCase":               "Bool",
		},
		"accountRoleSettings": {
			"allowAccountManagement":                            "Bool",
			"allowApiAccess":                                    "Bool",
			"allowApiAccessToAccount":                           "Bool",
			"allowApiSendingOnBehalfOfOthers":                   "Bool",
			"allowApiSequentialSigning":                         "Bool",
			"allowBulkSending":                                  "Bool",
			"allowDocuSignDesktopClient":                        "Bool",
			"allowESealRecipients":                              "Bool",
			"allowEnvelopeSending":                              "Bool",
			"allowPowerFormsAdminToAccessAllPowerFormEnvelopes": "Bool",
			"allowSendersToSetRecipientEmailLanguage":           "Bool",
			"allowSignerAttachments":                            "Bool",
			"allowSupplementalDocuments":                        "Bool",
			"allowTaggingInSendAndCorrect":                      "Bool",
			"allowVaulting":                                     "Bool",
			"allowWetSigningOverride":                           "Bool",
			"allowedAddressBookAccess":                          "Bool",
			"allowedTemplateAccess":                             "Bool",
			"allowedToBeEnvelopeTransferRecipient":              "Bool",
			"canCreateWorkspaces":                               "Bool",
			"disableDocumentUpload":                             "Bool",
			"disableOtherActions":                               "Bool",
			"enableApiRequestLogging":                           "Bool",
			"enableRecipientViewingNotifications":               "Bool",
			"enableSequentialSigningInterface":                  "Bool",
			"enableTransactionPointIntegration":                 "Bool",
			"receiveCompletedSelfSignedDocumentsAsEmailLinks":   "Bool",
			"supplementalDocumentsMustAccept":                   "Bool",
			"supplementalDocumentsMustRead":                     "Bool",
			"supplementalDocumentsMustView":                     "Bool",
			"useNewDocuSignExperienceInterface":                 "Bool",
			"useNewSendingInterface":                            "Bool",
		},
		"accountSettingsInformation": {
			"adoptSigConfig":                               "Bool",
			"advancedCorrect":                              "Bool",
			"allowAccessCodeFormat":                        "Bool",
			"allowAccountManagementGranular":               "Bool",
			"allowAccountMemberNameChange":                 "Bool",
			"allowAdvancedRecipientRoutingConditional":     "Bool",
			"allowBulkSend":                                "Bool",
			"allowCDWithdraw":                              "Bool",
			"allowConnectHttpListenerConfigs":              "Bool",
			"allowConsumerDisclosureOverride":              "Bool",
			"allowDataDownload":                            "Bool",
			"allowDocumentDisclosures":                     "Bool",
			"allowDocumentVisibility":                      "Bool",
			"allowDocumentsOnSignedEnvelopes":              "Bool",
			"allowEHankoStamps":                            "Bool",
			"allowEnvelopeCorrect":                         "Bool",
			"allowEnvelopePublishReporting":                "Bool",
			"allowExpressSignerCertificate":                "Bool",
			"allowExtendedSendingResourceFile":             "Bool",
			"allowExternalSignaturePad":                    "Bool",
			"allowIDVLevel1":                               "Bool",
			"allowInPerson":                                "Bool",
			"allowManagedStamps":                           "Bool",
			"allowMarkup":                                  "Bool",
			"allowMemberTimeZone":                          "Bool",
			"allowMergeFields":                             "Bool",
			"allowMultipleSignerAttachments":               "Bool",
			"allowOfflineSigning":                          "Bool",
			"allowOpenTrustSignerCertificate":              "Bool",
			"allowOrganizations":                           "Bool",
			"allowPaymentProcessing":                       "Bool",
			"allowPhoneAuthOverride":                       "Bool",
			"allowPhoneAuthentication":                     "Bool",
			"allowReminders":                               "Bool",
			"allowResourceFileBranding":                    "Bool",
			"allowSafeBioPharmaSignerCertificate":          "Bool",
			"allowSecurityAppliance":                       "Bool",
			"allowSendToCertifiedDelivery":                 "Bool",
			"allowSendToIntermediary":                      "Bool",
			"allowServerTemplates":                         "Bool",
			"allowSharedTabs":                              "Bool",
			"allowSignDocumentFromHomePage":                "Bool",
			"allowSignNow":                                 "Bool",
			"allowSignatureStamps":                         "Bool",
			"allowSignerReassign":                          "Bool",
			"allowSignerReassignOverride":                  "Bool",
			"allowSigningExtensions":                       "Bool",
			"allowSigningGroups":                           "Bool",
			"allowSigningRadioDeselect":                    "Bool",
			"allowSupplementalDocuments":                   "Bool",
			"attachCompletedEnvelope":                      "Bool",
			"autoProvisionSignerAccount":                   "Bool",
			"bccEmailArchive":                              "Bool",
			"bulkSend":                                     "Bool",
			"canSelfBrandSend":                             "Bool",
			"canSelfBrandSign":                             "Bool",
			"cfrUseWideImage":                              "Bool",
			"chromeSignatureEnabled":                       "Bool",
			"commentEmailShowMessageText":                  "Bool",
			"conditionalFieldsEnabled":                     "Bool",
			"convertPdfFields":                             "Bool",
			"disableMobileApp":                             "Bool",
			"disableMobilePushNotifications":               "Bool",
			"disableMobileSending":                         "Bool",
			"disableMultipleSessions":                      "Bool",
			"disableSignerCertView":                        "Bool",
			"disableSignerHistoryView":                     "Bool",
			"disableStyleSignature":                        "Bool",
			"disableUploadSignature":                       "Bool",
			"disableUserSharing":                           "Bool",
			"displayBetaSwitch":                            "Bool",
			"enableAccessCodeGenerator":                    "Bool",
			"enableAdvancedPayments":                       "Bool",
			"enableAdvancedPowerForms":                     "Bool",
			"enableAutoNav":                                "Bool",
			"enableCalculatedFields":                       "Bool",
			"enableClickwraps":                             "Bool",
			"enableCustomerSatisfactionMetricTracking":     "Bool",
			"enableEnvelopeStampingByAccountAdmin":         "Bool",
			"enableEnvelopeStampingByDSAdmin":              "Bool",
			"enablePaymentProcessing":                      "Bool",
			"enablePowerForm":                              "Bool",
			"enablePowerFormDirect":                        "Bool",
			"enableRequireSignOnPaper":                     "Bool",
			"enableReservedDomain":                         "Bool",
			"enableResponsiveSigning":                      "Bool",
			"enableSMSAuthentication":                      "Bool",
			"enableScheduledRelease":                       "Bool",
			"enableSendToAgent":                            "Bool",
			"enableSendToIntermediary":                     "Bool",
			"enableSendToManage":                           "Bool",
			"enableSendingTagsFontSettings":                "Bool",
			"enableSequentialSigningAPI":                   "Bool",
			"enableSequentialSigningUI":                    "Bool",
			"enableSignOnPaper":                            "Bool",
			"enableSignOnPaperOverride":                    "Bool",
			"enableSignWithNotary":                         "Bool",
			"enableSignerAttachments":                      "Bool",
			"enableSigningExtensionComments":               "Bool",
			"enableSigningExtensionConversations":          "Bool",
			"enableSigningOrderSettingsForAccount":         "Bool",
			"enableSmartContracts":                         "Bool",
			"enableStrikeThrough":                          "Bool",
			"enableVaulting":                               "Bool",
			"enforceTemplateNameUniqueness":                "Bool",
			"envelopeIntegrationEnabled":                   "Bool",
			"envelopeStampingDefaultValue":                 "Bool",
			"expressSend":                                  "Bool",
			"expressSendAllowTabs":                         "Bool",
			"faxOutEnabled":                                "Bool",
			"guidedFormsHtmlAllowed":                       "Bool",
			"hideAccountAddressInCoC":                      "Bool",
			"hidePricing":                                  "Bool",
			"inPersonSigningEnabled":                       "Bool",
			"inSessionEnabled":                             "Bool",
			"inSessionSuppressEmails":                      "Bool",
			"optInMobileSigningV02":                        "Bool",
			"optOutAutoNavTextAndTabColorUpdates":          "Bool",
			"optOutNewPlatformSeal":                        "Bool",
			"phoneAuthRecipientMayProvidePhoneNumber":      "Bool",
			"recipientSigningAutoNavigationControl":        "Bool",
			"recipientsCanSignOffline":                     "Bool",
			"require21CFRpt11Compliance":                   "Bool",
			"requireDeclineReason":                         "Bool",
			"requireExternalUserManagement":                "Bool",
			"selfSignedRecipientEmailDocumentUserOverride": "Bool",
			"senderCanSignInEachLocation":                  "Bool",
			"senderMustAuthenticateSigning":                "Bool",
			"setRecipEmailLang":                            "Bool",
			"setRecipSignLang":                             "Bool",
			"sharedTemplateFolders":                        "Bool",
			"showCompleteDialogInEmbeddedSession":          "Bool",
			"showConditionalRoutingOnSend":                 "Bool",
			"showInitialConditionalFields":                 "Bool",
			"showLocalizedWatermarks":                      "Bool",
			"showTutorials":                                "Bool",
			"signTimeShowAmPm":                             "Bool",
			"signerAttachCertificateToEnvelopePDF":         "Bool",
			"signerAttachConcat":                           "Bool",
			"signerCanCreateAccount":                       "Bool",
			"signerCanSignOnMobile":                        "Bool",
			"signerInSessionUseEnvelopeCompleteEmail":      "Bool",
			"signerMustHaveAccount":                        "Bool",
			"signerMustLoginToSign":                        "Bool",
			"signerShowSecureFieldInitialValues":           "Bool",
			"simplifiedSendingEnabled":                     "Bool",
			"singleSignOnEnabled":                          "Bool",
			"skipAuthCompletedEnvelopes":                   "Bool",
			"socialIdRecipAuth":                            "Bool",
			"specifyDocumentVisibility":                    "Bool",
			"startInAdvancedCorrect":                       "Bool",
			"supplementalDocumentsMustAccept":              "Bool",
			"supplementalDocumentsMustRead":                "Bool",
			"supplementalDocumentsMustView":                "Bool",
			"suppressCertificateEnforcement":               "Bool",
			"useConsumerDisclosure":                        "Bool",
			"useConsumerDisclosureWithinAccount":           "Bool",
			"useDocuSignExpressSignerCertificate":          "Bool",
			"useSAFESignerCertificates":                    "Bool",
			"useSignatureProviderPlatform":                 "Bool",
			"usesAPI":                                      "Bool",
			"validationsAllowed":                           "Bool",
			"validationsEnabled":                           "Bool",
			"waterMarkEnabled":                             "Bool",
			"writeReminderToEnvelopeHistory":               "Bool",
		},
		"accountSignature": {
			"disallowUserResizeStamp": "Bool",
			"isDefault":               "Bool",
		},
		"accountSignatureDefinition": {
			"disallowUserResizeStamp": "Bool",
			"isDefault":               "Bool",
		},
		"accountSignatureRequired": {
			"isRequired": "Bool",
		},
		"accountUISettings": {
			"hideUseATemplate": "Bool",
		},
		"addressInformationInput": {
			"receiveInResponse": "Bool",
		},
		"agent": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "Bool",
			"inheritEmailNotificationConfiguration": "Bool",
			"phoneAuthentication":                   "Bool",
			"requireIdLookup":                       "Bool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"apiRequestLog": {
			"createdDateTime": "*time.Time",
		},
		"approve": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"bold":                              "Bool",
			"italic":                            "Bool",
			"underline":                         "Bool",
		},
		"attachment": {
			"data": "[]byte",
		},
		"billingPaymentItem": {
			"paymentNumber": "Bool",
		},
		"billingPlan": {
			"enableSupport": "Bool",
		},
		"billingPlanInformation": {
			"enableSupport": "Bool",
		},
		"billingPlanPreview": {
			"isProrated": "Bool",
		},
		"brand": {
			"isOverridingCompanyName": "Bool",
			"isSendingDefault":        "Bool",
			"isSigningDefault":        "Bool",
		},
		"brandLink": {
			"showLink": "Bool",
		},
		"bulkEnvelope": {
			"submittedDateTime": "*time.Time",
		},
		"bulkSendEnvelopesInfo": {
			"authoritativeCopy": "Bool",
		},
		"bulkSendTestResponse": {
			"canBeSent": "Bool",
		},
		"carbonCopy": {
			"agentCanEditEmail":                     "Bool",
			"agentCanEditName":                      "Bool",
			"allowSystemOverrideForLockedRecipient": "Bool",
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "Bool",
			"inheritEmailNotificationConfiguration": "Bool",
			"phoneAuthentication":                   "Bool",
			"requireIdLookup":                       "Bool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"certifiedDelivery": {
			"agentCanEditEmail":                     "Bool",
			"agentCanEditName":                      "Bool",
			"allowSystemOverrideForLockedRecipient": "Bool",
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "Bool",
			"inheritEmailNotificationConfiguration": "Bool",
			"phoneAuthentication":                   "Bool",
			"requireIdLookup":                       "Bool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"checkbox": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorCaseSensitive":               "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"locked":                            "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"selected":                          "Bool",
			"shared":                            "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
		},
		"chunkedUploadRequest": {
			"data": "[]byte",
		},
		"chunkedUploadResponse": {
			"committed":          "Bool",
			"expirationDateTime": "*time.Time",
		},
		"commentThread": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"bold":                              "Bool",
			"italic":                            "Bool",
			"underline":                         "Bool",
		},
		"commissionCounty": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
		},
		"commissionExpiration": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
		},
		"commissionNumber": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
		},
		"commissionState": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
		},
		"company": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
		},
		"connectCustomConfiguration": {
			"allUsers":                          "Bool",
			"allowEnvelopePublish":              "Bool",
			"allowSalesforcePublish":            "Bool",
			"enableLog":                         "Bool",
			"includeCertificateOfCompletion":    "Bool",
			"includeDocumentFields":             "Bool",
			"includeDocuments":                  "Bool",
			"includeEnvelopeVoidReason":         "Bool",
			"includeSenderAccountasCustomField": "Bool",
			"includeTimeZoneInformation":        "Bool",
			"requireMutualTls":                  "Bool",
			"requiresAcknowledgement":           "Bool",
			"salesforceDocumentsAsContentFiles": "Bool",
			"signMessageWithX509Certificate":    "Bool",
			"useSoapInterface":                  "Bool",
		},
		"connectDebugLog": {
			"eventDateTime": "*time.Time",
		},
		"connectLog": {
			"created": "*time.Time",
		},
		"connectSalesforceObject": {
			"active":         "Bool",
			"onCompleteOnly": "Bool",
		},
		"connectUserObject": {
			"enabled": "Bool",
		},
		"consumerDisclosure": {
			"allowCDWithdraw":                    "Bool",
			"custom":                             "Bool",
			"mustAgreeToEsign":                   "Bool",
			"useBrand":                           "Bool",
			"useConsumerDisclosureWithinAccount": "Bool",
			"withdrawByEmail":                    "Bool",
			"withdrawByMail":                     "Bool",
			"withdrawByPhone":                    "Bool",
		},
		"contact": {
			"shared": "Bool",
		},
		"currency": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
		},
		"currencyFeatureSetPrice": {
			"envelopeFee": "Bool",
			"fixedFee":    "Bool",
			"seatFee":     "Bool",
		},
		"customField": {
			"required": "Bool",
			"show":     "Bool",
		},
		"customField_v2": {
			"required": "Bool",
			"show":     "Bool",
		},
		"date": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"requireAll":                        "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
		},
		"dateSigned": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
		},
		"decline": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"bold":                              "Bool",
			"italic":                            "Bool",
			"underline":                         "Bool",
		},
		"diagnosticsSettingsInformation": {
			"apiRequestLogging": "Bool",
		},
		"dobInformationInput": {
			"receiveInResponse": "Bool",
		},
		"document": {
			"authoritativeCopy":                      "Bool",
			"documentBase64":                         "[]byte",
			"encryptedWithKeyManager":                "Bool",
			"includeInDownload":                      "Bool",
			"signerMustAcknowledgeUseAccountDefault": "Bool",
			"templateLocked":                         "Bool",
			"templateRequired":                       "Bool",
			"transformPdfFields":                     "Bool",
		},
		"documentHtmlCollapsibleDisplaySettings": {
			"onlyArrowIsClickable": "Bool",
		},
		"documentHtmlDefinition": {
			"displayOrder":              "int32",
			"displayPageNumber":         "int32",
			"removeEmptyTags":           "Bool",
			"showMobileOptimizedToggle": "Bool",
		},
		"documentHtmlDisplayAnchor": {
			"caseSensitive":     "Bool",
			"removeEndAnchor":   "Bool",
			"removeStartAnchor": "Bool",
		},
		"documentHtmlDisplaySettings": {
			"hideLabelWhenOpened": "Bool",
		},
		"documentVisibility": {
			"visible": "Bool",
		},
		"draw": {
			"allowSignerUpload":                 "Bool",
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
			"shared":                            "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
		},
		"editor": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"inheritEmailNotificationConfiguration": "Bool",
			"phoneAuthentication":                   "Bool",
			"requireIdLookup":                       "Bool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"email": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"requireAll":                        "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
		},
		"emailAddress": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
		},
		"envelope": {
			"allowComments":               "Bool",
			"allowMarkup":                 "Bool",
			"allowReassign":               "Bool",
			"allowViewHistory":            "Bool",
			"asynchronous":                "Bool",
			"authoritativeCopy":           "Bool",
			"autoNavigation":              "Bool",
			"brandLock":                   "Bool",
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"declinedDateTime":            "*time.Time",
			"deletedDateTime":             "*time.Time",
			"deliveredDateTime":           "*time.Time",
			"disableResponsiveDocument":   "Bool",
			"enableWetSign":               "Bool",
			"enforceSignerVisibility":     "Bool",
			"envelopeIdStamping":          "Bool",
			"hasComments":                 "Bool",
			"hasWavFile":                  "Bool",
			"initialSentDateTime":         "*time.Time",
			"is21CFRPart11":               "Bool",
			"isDynamicEnvelope":           "Bool",
			"isSignatureProviderEnvelope": "Bool",
			"lastModifiedDateTime":        "*time.Time",
			"messageLock":                 "Bool",
			"notification":                "Bool",
			"recipientsLock":              "Bool",
			"sentDateTime":                "*time.Time",
			"signerCanSignOnMobile":       "Bool",
			"statusChangedDateTime":       "*time.Time",
			"useDisclosure":               "Bool",
			"voidedDateTime":              "*time.Time",
		},
		"envelopeDefinition": {
			"allowComments":               "Bool",
			"allowMarkup":                 "Bool",
			"allowReassign":               "Bool",
			"allowRecipientRecursion":     "Bool",
			"allowViewHistory":            "Bool",
			"asynchronous":                "Bool",
			"authoritativeCopy":           "Bool",
			"autoNavigation":              "Bool",
			"brandLock":                   "Bool",
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"declinedDateTime":            "*time.Time",
			"deletedDateTime":             "*time.Time",
			"deliveredDateTime":           "*time.Time",
			"disableResponsiveDocument":   "Bool",
			"enableWetSign":               "Bool",
			"enforceSignerVisibility":     "Bool",
			"envelopeIdStamping":          "Bool",
			"hasComments":                 "Bool",
			"hasFormDataChanged":          "Bool",
			"hasWavFile":                  "Bool",
			"initialSentDateTime":         "*time.Time",
			"is21CFRPart11":               "Bool",
			"isDynamicEnvelope":           "Bool",
			"isSignatureProviderEnvelope": "Bool",
			"lastModifiedDateTime":        "*time.Time",
			"messageLock":                 "Bool",
			"recipientsLock":              "Bool",
			"sentDateTime":                "*time.Time",
			"signerCanSignOnMobile":       "Bool",
			"statusChangedDateTime":       "*time.Time",
			"useDisclosure":               "Bool",
			"voidedDateTime":              "*time.Time",
		},
		"envelopeDocument": {
			"authoritativeCopy":     "Bool",
			"containsPdfFormFields": "Bool",
			"includeInDownload":     "Bool",
			"templateLocked":        "Bool",
			"templateRequired":      "Bool",
		},
		"envelopeEvent": {
			"includeDocuments": "Bool",
		},
		"envelopeFormData": {
			"sentDateTime": "*time.Time",
		},
		"envelopeId": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
		},
		"envelopeNotificationRequest": {
			"useAccountDefaults": "Bool",
		},
		"envelopePurgeConfiguration": {
			"purgeEnvelopes":                   "Bool",
			"redactPII":                        "Bool",
			"removeTabsAndEnvelopeAttachments": "Bool",
		},
		"envelopeSummary": {
			"statusDateTime": "*time.Time",
		},
		"envelopeTemplate": {
			"allowComments":               "Bool",
			"allowMarkup":                 "Bool",
			"allowReassign":               "Bool",
			"allowViewHistory":            "Bool",
			"asynchronous":                "Bool",
			"authoritativeCopy":           "Bool",
			"autoMatchSpecifiedByUser":    "Bool",
			"autoNavigation":              "Bool",
			"brandLock":                   "Bool",
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"declinedDateTime":            "*time.Time",
			"deletedDateTime":             "*time.Time",
			"deliveredDateTime":           "*time.Time",
			"disableResponsiveDocument":   "Bool",
			"enableWetSign":               "Bool",
			"enforceSignerVisibility":     "Bool",
			"envelopeIdStamping":          "Bool",
			"hasComments":                 "Bool",
			"hasWavFile":                  "Bool",
			"initialSentDateTime":         "*time.Time",
			"is21CFRPart11":               "Bool",
			"isDynamicEnvelope":           "Bool",
			"isSignatureProviderEnvelope": "Bool",
			"lastModifiedDateTime":        "*time.Time",
			"messageLock":                 "Bool",
			"passwordProtected":           "Bool",
			"recipientsLock":              "Bool",
			"sentDateTime":                "*time.Time",
			"shared":                      "Bool",
			"signerCanSignOnMobile":       "Bool",
			"statusChangedDateTime":       "*time.Time",
			"useDisclosure":               "Bool",
			"voidedDateTime":              "*time.Time",
		},
		"envelopeTemplateDefinition": {
			"created":      "*time.Time",
			"lastModified": "*time.Time",
			"shared":       "Bool",
		},
		"envelopeTemplateResult": {
			"allowMarkup":                 "Bool",
			"allowReassign":               "Bool",
			"allowViewHistory":            "Bool",
			"asynchronous":                "Bool",
			"authoritativeCopy":           "Bool",
			"completedDateTime":           "*time.Time",
			"created":                     "*time.Time",
			"createdDateTime":             "*time.Time",
			"declinedDateTime":            "*time.Time",
			"deletedDateTime":             "*time.Time",
			"deliveredDateTime":           "*time.Time",
			"enableWetSign":               "Bool",
			"enforceSignerVisibility":     "Bool",
			"envelopeIdStamping":          "Bool",
			"initialSentDateTime":         "*time.Time",
			"is21CFRPart11":               "Bool",
			"isSignatureProviderEnvelope": "Bool",
			"lastModified":                "*time.Time",
			"lastModifiedDateTime":        "*time.Time",
			"messageLock":                 "Bool",
			"recipientsLock":              "Bool",
			"sentDateTime":                "*time.Time",
			"shared":                      "Bool",
			"signerCanSignOnMobile":       "Bool",
			"statusChangedDateTime":       "*time.Time",
			"useDisclosure":               "Bool",
			"voidedDateTime":              "*time.Time",
		},
		"envelopeTransferRule": {
			"carbonCopyOriginalOwner": "Bool",
			"enabled":                 "Bool",
		},
		"envelopeTransferRuleRequest": {
			"carbonCopyOriginalOwner": "Bool",
			"enabled":                 "Bool",
		},
		"eventNotification": {
			"includeCertificateOfCompletion":    "Bool",
			"includeCertificateWithSoap":        "Bool",
			"includeDocumentFields":             "Bool",
			"includeDocuments":                  "Bool",
			"includeEnvelopeVoidReason":         "Bool",
			"includeHMAC":                       "Bool",
			"includeSenderAccountAsCustomField": "Bool",
			"includeTimeZone":                   "Bool",
			"loggingEnabled":                    "Bool",
			"requireAcknowledgment":             "Bool",
			"signMessageWithX509Cert":           "Bool",
			"useSoapInterface":                  "Bool",
		},
		"expirations": {
			"expireEnabled": "Bool",
		},
		"externalFile": {
			"supported": "Bool",
		},
		"featureSet": {
			"is21CFRPart11": "Bool",
			"isActive":      "Bool",
			"isEnabled":     "Bool",
		},
		"filter": {
			"actionRequired": "Bool",
			"isTemplate":     "Bool",
		},
		"firstName": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
		},
		"folder": {
			"hasAccess":     "Bool",
			"hasSubFolders": "Bool",
		},
		"folderItem": {
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"is21CFRPart11":               "Bool",
			"isSignatureProviderEnvelope": "Bool",
			"sentDateTime":                "*time.Time",
			"shared":                      "Bool",
		},
		"folderItem_v2": {
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"expireDateTime":              "*time.Time",
			"is21CFRPart11":               "Bool",
			"isSignatureProviderEnvelope": "Bool",
			"lastModifiedDateTime":        "*time.Time",
			"sentDateTime":                "*time.Time",
		},
		"formulaTab": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"isPaymentAmount":                   "Bool",
			"locked":                            "Bool",
			"requireAll":                        "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
		},
		"fullName": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
		},
		"idCheckConfiguration": {
			"isDefault": "Bool",
		},
		"inPersonSigner": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"autoNavigation":                        "Bool",
			"canSignOffline":                        "Bool",
			"declinedDateTime":                      "*time.Time",
			"defaultRecipient":                      "Bool",
			"deliveredDateTime":                     "*time.Time",
			"inheritEmailNotificationConfiguration": "Bool",
			"phoneAuthentication":                   "Bool",
			"recipientSuppliesTabs":                 "Bool",
			"requireIdLookup":                       "Bool",
			"requireSignOnPaper":                    "Bool",
			"requireUploadSignature":                "Bool",
			"sentDateTime":                          "*time.Time",
			"signInEachLocation":                    "Bool",
			"signedDateTime":                        "*time.Time",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"initialHere": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"optional":                          "Bool",
		},
		"intermediary": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "Bool",
			"inheritEmailNotificationConfiguration": "Bool",
			"phoneAuthentication":                   "Bool",
			"requireIdLookup":                       "Bool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"jurisdiction": {
			"allowSystemCreatedSeal": "Bool",
			"allowUserUploadedSeal":  "Bool",
			"commissionIdInSeal":     "Bool",
			"countyInSeal":           "Bool",
			"enabled":                "Bool",
			"notaryPublicInSeal":     "Bool",
			"stateNameInSeal":        "Bool",
		},
		"lastName": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
		},
		"list": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"bold":                              "Bool",
			"italic":                            "Bool",
			"locked":                            "Bool",
			"requireAll":                        "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
			"underline":                         "Bool",
		},
		"listCustomField": {
			"required": "Bool",
			"show":     "Bool",
		},
		"listItem": {
			"selected": "Bool",
		},
		"lockInformation": {
			"lockedUntilDateTime": "*time.Time",
			"useScratchPad":       "Bool",
		},
		"lockRequest": {
			"useScratchPad": "Bool",
		},
		"loginAccount": {
			"isDefault": "Bool",
		},
		"memberGroupSharedItem": {
			"shared": "Bool",
		},
		"mergeField": {
			"allowSenderToEdit": "Bool",
			"writeBack":         "Bool",
		},
		"newAccountDefinition": {
			"socialAccountInformation": "*SocialAccountInformation",
		},
		"newUser": {
			"createdDateTime": "*time.Time",
		},
		"notarize": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
		},
		"notaryHost": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"inheritEmailNotificationConfiguration": "Bool",
			"phoneAuthentication":                   "Bool",
			"requireIdLookup":                       "Bool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"notaryRecipient": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"autoNavigation":                        "Bool",
			"canSignOffline":                        "Bool",
			"defaultRecipient":                      "Bool",
			"offlineAttributes":                     "interface{}",
			"recipientSuppliesTabs":                 "Bool",
			"requireIdLookup":                       "Bool",
			"requireSignOnPaper":                    "Bool",
			"requireUploadSignature":                "Bool",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"notarySeal": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
		},
		"note": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"shared":                            "Bool",
		},
		"notification": {
			"useAccountDefaults": "Bool",
		},
		"number": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"isPaymentAmount":                   "Bool",
			"locked":                            "Bool",
			"requireAll":                        "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
		},
		"participant": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"requireIdLookup":                       "Bool",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"paymentGatewayAccount": {
			"allowCustomMetadata": "Bool",
			"isEnabled":           "Bool",
		},
		"permissionProfile": {
			"modifiedDateTime": "*time.Time",
		},
		"phoneNumber": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"bold":                              "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"italic":                            "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
			"underline":                         "Bool",
		},
		"polyLineOverlay": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"locked":                            "Bool",
			"shared":                            "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
		},
		"powerForm": {
			"createdDateTime":         "*time.Time",
			"isActive":                "Bool",
			"limitUseIntervalEnabled": "Bool",
			"maxUseEnabled":           "Bool",
		},
		"powerFormRecipient": {
			"accessCodeLocked":         "Bool",
			"accessCodeRequired":       "Bool",
			"emailLocked":              "Bool",
			"templateRequiresIdLookup": "Bool",
			"userNameLocked":           "Bool",
		},
		"purchasedEnvelopesInformation": {
			"receiptData": "[]byte",
		},
		"radio": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorCaseSensitive":               "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"bold":                              "Bool",
			"italic":                            "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
			"selected":                          "Bool",
			"underline":                         "Bool",
		},
		"radioGroup": {
			"requireAll":                   "Bool",
			"requireInitialOnSharedChange": "Bool",
			"shared":                       "Bool",
			"templateLocked":               "Bool",
			"templateRequired":             "Bool",
		},
		"recipientAttachment": {
			"data": "[]byte",
		},
		"recipientEvent": {
			"includeDocuments": "Bool",
		},
		"recipientFormData": {
			"declinedTime":  "*time.Time",
			"deliveredTime": "*time.Time",
			"sentTime":      "*time.Time",
			"signedTime":    "*time.Time",
		},
		"recipientNamesResponse": {
			"multipleUsers":          "Bool",
			"reservedRecipientEmail": "Bool",
		},
		"recipientPhoneAuthentication": {
			"recipMayProvideNumber": "Bool",
		},
		"recipientSignatureProvider": {
			"sealDocumentsWithTabsOnly": "Bool",
		},
		"referralInformation": {
			"enableSupport": "Bool",
		},
		"reminders": {
			"reminderEnabled": "Bool",
		},
		"sealSign": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"senderCompany": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"bold":                              "Bool",
			"italic":                            "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
			"underline":                         "Bool",
		},
		"senderEmailNotifications": {
			"changedSigner":                 "Bool",
			"commentsOnlyPrivateAndMention": "Bool",
			"commentsReceiveAll":            "Bool",
			"deliveryFailed":                "Bool",
			"envelopeComplete":              "Bool",
			"offlineSigningFailed":          "Bool",
			"purgeDocuments":                "Bool",
			"recipientViewed":               "Bool",
			"senderEnvelopeDeclined":        "Bool",
			"withdrawnConsent":              "Bool",
		},
		"senderName": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"bold":                              "Bool",
			"italic":                            "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
			"underline":                         "Bool",
		},
		"serviceInformation": {
			"buildBranchDeployedDateTime": "*time.Time",
		},
		"settingsMetadata": {
			"is21CFRPart11": "Bool",
		},
		"sharedItem": {
			"shared": "Bool",
		},
		"signHere": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"isSealSignTab":                     "Bool",
			"optional":                          "Bool",
		},
		"signatureType": {
			"isDefault": "Bool",
		},
		"signatureUser": {
			"isDefault": "Bool",
		},
		"signatureUserDef": {
			"isDefault": "Bool",
		},
		"signer": {
			"agentCanEditEmail":                     "Bool",
			"agentCanEditName":                      "Bool",
			"allowSystemOverrideForLockedRecipient": "Bool",
			"autoNavigation":                        "Bool",
			"canSignOffline":                        "Bool",
			"declinedDateTime":                      "*time.Time",
			"defaultRecipient":                      "Bool",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "Bool",
			"inheritEmailNotificationConfiguration": "Bool",
			"isBulkRecipient":                       "Bool",
			"phoneAuthentication":                   "Bool",
			"recipientSuppliesTabs":                 "Bool",
			"requireIdLookup":                       "Bool",
			"requireSignOnPaper":                    "Bool",
			"requireUploadSignature":                "Bool",
			"sentDateTime":                          "*time.Time",
			"signInEachLocation":                    "Bool",
			"signedDateTime":                        "*time.Time",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"signerAttachment": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"optional":                          "Bool",
		},
		"signerEmailNotifications": {
			"agentNotification":             "Bool",
			"carbonCopyNotification":        "Bool",
			"certifiedDeliveryNotification": "Bool",
			"commentsOnlyPrivateAndMention": "Bool",
			"commentsReceiveAll":            "Bool",
			"documentMarkupActivation":      "Bool",
			"envelopeActivation":            "Bool",
			"envelopeComplete":              "Bool",
			"envelopeCorrected":             "Bool",
			"envelopeDeclined":              "Bool",
			"envelopeVoided":                "Bool",
			"offlineSigningFailed":          "Bool",
			"purgeDocuments":                "Bool",
			"reassignedSigner":              "Bool",
			"whenSigningGroupMember":        "Bool",
		},
		"signingGroup": {
			"created":  "*time.Time",
			"modified": "*time.Time",
		},
		"smartSection": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"caseSensitive":                     "Bool",
			"locked":                            "Bool",
			"removeEndAnchor":                   "Bool",
			"removeStartAnchor":                 "Bool",
			"shared":                            "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
		},
		"smartSectionCollapsibleDisplaySettings": {
			"onlyArrowIsClickable": "Bool",
		},
		"smartSectionDisplaySettings": {
			"hideLabelWhenOpened": "Bool",
		},
		"ssn": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"requireAll":                        "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
		},
		"ssn4InformationInput": {
			"receiveInResponse": "Bool",
		},
		"stamp": {
			"disallowUserResizeStamp": "Bool",
		},
		"tabAccountSettings": {
			"allowTabOrder":                       "Bool",
			"approveDeclineTabsEnabled":           "Bool",
			"calculatedFieldsEnabled":             "Bool",
			"checkboxTabsEnabled":                 "Bool",
			"dataFieldRegexEnabled":               "Bool",
			"dataFieldSizeEnabled":                "Bool",
			"firstLastEmailTabsEnabled":           "Bool",
			"listTabsEnabled":                     "Bool",
			"noteTabsEnabled":                     "Bool",
			"radioTabsEnabled":                    "Bool",
			"savingCustomTabsEnabled":             "Bool",
			"senderToChangeTabAssignmentsEnabled": "Bool",
			"sharedCustomTabsEnabled":             "Bool",
			"tabDataLabelEnabled":                 "Bool",
			"tabLocationEnabled":                  "Bool",
			"tabLockingEnabled":                   "Bool",
			"tabScaleEnabled":                     "Bool",
			"tabTextFormattingEnabled":            "Bool",
			"textTabsEnabled":                     "Bool",
		},
		"tabGroup": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"anchorIgnoreIfNotPresent":          "Bool",
			"anchorMatchWholeWord":              "Bool",
			"templateLocked":                    "Bool",
			"templateRequired":                  "Bool",
		},
		"tabMetadata": {
			"anchorCaseSensitive":      "Bool",
			"anchorIgnoreIfNotPresent": "Bool",
			"anchorMatchWholeWord":     "Bool",
			"bold":                     "Bool",
			"concealValueOnDocument":   "Bool",
			"disableAutoSize":          "Bool",
			"editable":                 "Bool",
			"includedInEmail":          "Bool",
			"italic":                   "Bool",
			"lastModified":             "*time.Time",
			"locked":                   "Bool",
			"required":                 "Bool",
			"selected":                 "Bool",
			"shared":                   "Bool",
			"underline":                "Bool",
		},
		"templateNotificationRequest": {
			"useAccountDefaults": "Bool",
		},
		"templateRole": {
			"defaultRecipient": "Bool",
		},
		"templateSharedItem": {
			"shared": "Bool",
		},
		"text": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"isPaymentAmount":                   "Bool",
			"locked":                            "Bool",
			"requireAll":                        "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
		},
		"textCustomField": {
			"required": "Bool",
			"show":     "Bool",
		},
		"title": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"required":                          "Bool",
		},
		"usageHistory": {
			"lastSentDateTime":   "*time.Time",
			"lastSignedDateTime": "*time.Time",
		},
		"userAccountManagementGranularInformation": {
			"canManageAccountSecuritySettings": "Bool",
			"canManageAccountSettings":         "Bool",
			"canManageAdmins":                  "Bool",
			"canManageGroups":                  "Bool",
			"canManageReporting":               "Bool",
			"canManageSharing":                 "Bool",
			"canManageSigningGroups":           "Bool",
			"canManageUsers":                   "Bool",
		},
		"userInfo": {
			"loginStatus": "Bool",
		},
		"userInformation": {
			"createdDateTime":              "*time.Time",
			"enableConnectForUser":         "Bool",
			"isActive":                     "Bool",
			"isNAREnabled":                 "Bool",
			"sendActivationOnInvalidLogin": "Bool",
		},
		"userProfile": {
			"displayOrganizationInfo": "Bool",
			"displayPersonalInfo":     "Bool",
			"displayProfile":          "Bool",
			"displayUsageHistory":     "Bool",
		},
		"userSettingsInformation": {
			"allowAutoTagging":                      "Bool",
			"allowEnvelopeTransferTo":               "Bool",
			"allowEsealRecipients":                  "Bool",
			"allowRecipientLanguageSelection":       "Bool",
			"allowSendOnBehalfOf":                   "Bool",
			"allowSupplementalDocuments":            "Bool",
			"apiAccountWideAccess":                  "Bool",
			"apiCanExportAC":                        "Bool",
			"bulkSend":                              "Bool",
			"canManageAccount":                      "Bool",
			"canManageTemplates":                    "Bool",
			"canSendAPIRequests":                    "Bool",
			"canSendEnvelope":                       "Bool",
			"canSignEnvelope":                       "Bool",
			"canUseScratchpad":                      "Bool",
			"disableDocumentUpload":                 "Bool",
			"disableOtherActions":                   "Bool",
			"enableSequentialSigningAPI":            "Bool",
			"enableSequentialSigningUI":             "Bool",
			"enableSignOnPaperOverride":             "Bool",
			"enableSignerAttachments":               "Bool",
			"enableVaulting":                        "Bool",
			"manageClickwrapsMode":                  "Bool",
			"recipientViewedNotification":           "Bool",
			"supplementalDocumentIncludeInDownload": "Bool",
			"supplementalDocumentsMustAccept":       "Bool",
			"supplementalDocumentsMustRead":         "Bool",
			"supplementalDocumentsMustView":         "Bool",
			"templateActiveCreation":                "Bool",
			"templateApplyNotify":                   "Bool",
			"templateAutoMatching":                  "Bool",
			"templatePageLevelMatching":             "Bool",
		},
		"userSharedItem": {
			"shared": "Bool",
		},
		"userSignature": {
			"adoptedDateTime":         "*time.Time",
			"createdDateTime":         "*time.Time",
			"disallowUserResizeStamp": "Bool",
			"isDefault":               "Bool",
		},
		"userSignatureDefinition": {
			"disallowUserResizeStamp": "Bool",
			"isDefault":               "Bool",
		},
		"userSocialIdResult": {
			"socialAccountInformation": "[]SocialAccountInformation",
		},
		"view": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"bold":                              "Bool",
			"italic":                            "Bool",
			"required":                          "Bool",
			"requiredRead":                      "Bool",
			"underline":                         "Bool",
		},
		"watermark": {
			"enabled": "Bool",
		},
		"witness": {
			"allowSystemOverrideForLockedRecipient": "Bool",
			"autoNavigation":                        "Bool",
			"canSignOffline":                        "Bool",
			"defaultRecipient":                      "Bool",
			"offlineAttributes":                     "interface{}",
			"recipientSuppliesTabs":                 "Bool",
			"requireIdLookup":                       "Bool",
			"requireSignOnPaper":                    "Bool",
			"requireUploadSignature":                "Bool",
			"suppressEmails":                        "Bool",
			"templateLocked":                        "Bool",
			"templateRequired":                      "Bool",
		},
		"workspace": {
			"created":      "*time.Time",
			"lastModified": "*time.Time",
		},
		"workspaceItem": {
			"created":      "*time.Time",
			"isPublic":     "Bool",
			"lastModified": "*time.Time",
		},
		"workspaceSettings": {
			"commentsAllowed": "Bool",
		},
		"workspaceUser": {
			"activeSince":  "*time.Time",
			"created":      "*time.Time",
			"lastModified": "*time.Time",
		},
		"workspaceUserAuthorization": {
			"canDelete":   "Bool",
			"canMove":     "Bool",
			"canTransact": "Bool",
			"canView":     "Bool",
			"created":     "*time.Time",
			"modified":    "*time.Time",
		},
		"zip": {
			"anchorAllowWhiteSpaceInCharacters": "Bool",
			"concealValueOnDocument":            "Bool",
			"disableAutoSize":                   "Bool",
			"locked":                            "Bool",
			"requireAll":                        "Bool",
			"requireInitialOnSharedChange":      "Bool",
			"required":                          "Bool",
			"senderRequired":                    "Bool",
			"shared":                            "Bool",
			"useDash4":                          "Bool",
		},
	}
}
