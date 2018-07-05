// Package overrides provides lists of type overrides for definition
// properties (struct fields) and for operation parameters found in
// Docusign's Rest API swagger definition.  DocuSign defined all fields
// and parameters as strings, and well, I didn't really like that.
// So this package defines the overrides that generator
// uses.
package overrides

import (
	"strings"

	"github.com/jfcote87/esign/gen-esign/swagger"
)

type override struct {
	Object string
	Field  string
	Type   string
}

// GetOperationOverrides returns a map of func
// names that must be overridden.
func GetOperationOverrides() map[string]string {
	return map[string]string{
		"PasswordRules_GetPasswordRules": "GetCurrentUserPasswordRules",
	}
}

// GetParameterOverrides returns a map of
// all parameter Type overrides.
func GetParameterOverrides() map[string]map[string]string {
	defOverrides := make(map[string]map[string]string)
	for _, o := range propertyOverrides {
		m, ok := defOverrides[o.Object]
		if !ok {
			m = make(map[string]string)
			defOverrides[o.Object] = m
		}
		m[o.Field] = o.Type
	}
	return defOverrides
}

// GetFieldOverrides returns a map of all
// field level type overrides for the esign
// api generation. The returned map is
// map[<structID>]map[<FieldName>]<GoType>
func GetFieldOverrides() map[string]map[string]string {
	defOverrides := make(map[string]map[string]string)
	for _, o := range fieldOverrides {
		m, ok := defOverrides[o.Object]
		if !ok {
			m = make(map[string]string)
			defOverrides[o.Object] = m
		}
		m[o.Field] = o.Type
	}
	return defOverrides
}

// propertyOverrides lists operation parameter overrides.  I generated
// much of this list by searching the swagger file with the following rules:
// - if **true** found in description, set parameter type to bool.
// - if field name ends in "_date", set time.Time
// - if field names is count, start_position, etc, set to int
// - if description starts with "Comma separated list", set to ...string
//
// I eyeballed the doc as best I could so please let me know of any additions
// or corrections.
var propertyOverrides = []override{
	{Object: "Accounts_PostAccounts", Field: "preview_billing_plan", Type: "bool"},
	{Object: "Accounts_GetAccount", Field: "include_account_settings", Type: "bool"},
	{Object: "BillingPlan_GetBillingPlan", Field: "include_credit_card_information", Type: "bool"},
	{Object: "BillingPlan_GetBillingPlan", Field: "include_metadata", Type: "bool"},
	{Object: "BillingPlan_GetBillingPlan", Field: "include_successor_plans", Type: "bool"},
	{Object: "BillingPlan_PutBillingPlan", Field: "preview_billing_plan", Type: "bool"},
	{Object: "Brands_GetBrands", Field: "exclude_distributor_brand", Type: "bool"},
	{Object: "Brands_GetBrands", Field: "include_logos", Type: "bool"},
	{Object: "Envelopes_PostEnvelopes", Field: "merge_roles_on_draft", Type: "bool"},
	{Object: "Envelopes_PutEnvelope", Field: "advanced_update", Type: "bool"},
	{Object: "Envelopes_PutEnvelope", Field: "resend_envelope", Type: "bool"},
	{Object: "Documents_PutDocuments", Field: "apply_document_fields", Type: "bool"},
	{Object: "Documents_GetDocument", Field: "encrypt", Type: "bool"},
	{Object: "Documents_GetDocument", Field: "show_changes", Type: "bool"},
	{Object: "Documents_GetDocument", Field: "watermark", Type: "bool"},
	{Object: "Documents_PutDocument", Field: "apply_document_fields", Type: "bool"},
	{Object: "Recipients_GetRecipients", Field: "include_anchor_tab_locations", Type: "bool"},
	{Object: "Recipients_GetRecipients", Field: "include_extended", Type: "bool"},
	{Object: "Recipients_GetRecipients", Field: "include_tabs", Type: "bool"},
	{Object: "Recipients_PutRecipients", Field: "resend_envelope", Type: "bool"},
	{Object: "Recipients_PostRecipients", Field: "resend_envelope", Type: "bool"},
	{Object: "Recipients_GetBulkRecipients", Field: "include_tabs", Type: "bool"},
	{Object: "Recipients_GetRecipientSignatureImage", Field: "include_chrome", Type: "bool"},
	{Object: "Recipients_GetRecipientTabs", Field: "include_anchor_tab_locations", Type: "bool"},
	{Object: "SearchFolders_GetSearchFolderContents", Field: "include_recipients", Type: "bool"},
	{Object: "SigningGroups_GetSigningGroups", Field: "include_users", Type: "bool"},
	{Object: "Tabs_GetTabDefinitions", Field: "custom_tab_only", Type: "bool"},
	{Object: "Documents_PutTemplateDocuments", Field: "apply_document_fields", Type: "bool"},
	{Object: "Documents_PutTemplateDocument", Field: "apply_document_fields", Type: "bool"},
	{Object: "Recipients_GetTemplateRecipients", Field: "include_anchor_tab_locations", Type: "bool"},
	{Object: "Recipients_GetTemplateRecipients", Field: "include_extended", Type: "bool"},
	{Object: "Recipients_GetTemplateRecipients", Field: "include_tabs", Type: "bool"},
	{Object: "Recipients_PutTemplateRecipients", Field: "resend_envelope", Type: "bool"},
	{Object: "Recipients_PostTemplateRecipients", Field: "resend_envelope", Type: "bool"},
	{Object: "Recipients_GetTemplateBulkRecipients", Field: "include_tabs", Type: "bool"},
	{Object: "Recipients_GetTemplateRecipientTabs", Field: "include_anchor_tab_locations", Type: "bool"},
	{Object: "Users_GetUsers", Field: "additional_info", Type: "bool"},
	{Object: "User_GetUser", Field: "additional_info", Type: "bool"},
	{Object: "UserSignatures_PutUserSignatureById", Field: "close_existing_signature", Type: "bool"},
	{Object: "WorkspaceFolder_GetWorkspaceFolder", Field: "include_files", Type: "bool"},
	{Object: "WorkspaceFolder_GetWorkspaceFolder", Field: "include_sub_folders", Type: "bool"},
	{Object: "WorkspaceFolder_GetWorkspaceFolder", Field: "include_thumbnails", Type: "bool"},
	{Object: "WorkspaceFolder_GetWorkspaceFolder", Field: "include_user_detail", Type: "bool"},
	{Object: "WorkspaceFile_GetWorkspaceFile", Field: "is_download", Type: "bool"},
	{Object: "WorkspaceFile_GetWorkspaceFile", Field: "pdf_version", Type: "bool"},
	{Object: "LoginInformation_GetLoginInformation", Field: "include_account_id_guid", Type: "bool"},
	{Object: "BillingInvoices_GetBillingInvoices", Field: "from_date", Type: "time.Time"},
	{Object: "BillingInvoices_GetBillingInvoices", Field: "to_date", Type: "time.Time"},
	{Object: "BillingPayments_GetPaymentList", Field: "from_date", Type: "time.Time"},
	{Object: "BillingPayments_GetPaymentList", Field: "to_date", Type: "time.Time"},
	{Object: "ConnectFailures_GetConnectLogs", Field: "from_date", Type: "time.Time"},
	{Object: "ConnectFailures_GetConnectLogs", Field: "to_date", Type: "time.Time"},
	{Object: "ConnectLog_GetConnectLogs", Field: "from_date", Type: "time.Time"},
	{Object: "ConnectLog_GetConnectLogs", Field: "to_date", Type: "time.Time"},
	{Object: "Envelopes_GetEnvelopes", Field: "from_date", Type: "time.Time"},
	{Object: "Envelopes_GetEnvelopes", Field: "to_date", Type: "time.Time"},
	{Object: "Envelopes_PutStatus", Field: "from_date", Type: "time.Time"},
	{Object: "Envelopes_PutStatus", Field: "to_date", Type: "time.Time"},
	{Object: "Folders_GetFolderItems", Field: "from_date", Type: "time.Time"},
	{Object: "Folders_GetFolderItems", Field: "to_date", Type: "time.Time"},
	{Object: "PowerForms_GetPowerFormsList", Field: "from_date", Type: "time.Time"},
	{Object: "PowerForms_GetPowerFormsList", Field: "to_date", Type: "time.Time"},
	{Object: "PowerForms_GetPowerFormFormData", Field: "from_date", Type: "time.Time"},
	{Object: "PowerForms_GetPowerFormFormData", Field: "to_date", Type: "time.Time"},
	{Object: "SearchFolders_GetSearchFolderContents", Field: "from_date", Type: "time.Time"},
	{Object: "SearchFolders_GetSearchFolderContents", Field: "to_date", Type: "time.Time"},
	{Object: "Templates_GetTemplates", Field: "from_date", Type: "time.Time"},
	{Object: "Templates_GetTemplates", Field: "to_date", Type: "time.Time"},
	{Object: "Templates_GetTemplates", Field: "used_from_date", Type: "time.Time"},
	{Object: "Templates_GetTemplates", Field: "used_to_date", Type: "time.Time"},
	{Object: "BulkEnvelopes_GetEnvelopesBulk", Field: "count", Type: "int"},
	{Object: "BulkEnvelopes_GetBulkEnvelopesBatchId", Field: "count", Type: "int"},
	{Object: "Connect_GetConnectUsers", Field: "count", Type: "int"},
	{Object: "Envelopes_GetEnvelopes", Field: "count", Type: "int"},
	{Object: "Pages_GetPageImages", Field: "count", Type: "int"},
	{Object: "Groups_GetGroups", Field: "count", Type: "int"},
	{Object: "Groups_GetGroupUsers", Field: "count", Type: "int"},
	{Object: "SearchFolders_GetSearchFolderContents", Field: "count", Type: "int"},
	{Object: "SharedAccess_GetSharedAccess", Field: "count", Type: "int"},
	{Object: "Templates_GetTemplates", Field: "count", Type: "int"},
	{Object: "Pages_GetTemplatePageImages", Field: "count", Type: "int"},
	{Object: "Users_GetUsers", Field: "count", Type: "int"},
	{Object: "CloudStorageFolder_GetCloudStorageFolderAll", Field: "count", Type: "int"},
	{Object: "CloudStorageFolder_GetCloudStorageFolder", Field: "count", Type: "int"},
	{Object: "WorkspaceFolder_GetWorkspaceFolder", Field: "count", Type: "int"},
	{Object: "WorkspaceFilePages_GetWorkspaceFilePages", Field: "count", Type: "int"},
	{Object: "BulkEnvelopes_GetEnvelopesBulk", Field: "start_position", Type: "int"},
	{Object: "BulkEnvelopes_GetBulkEnvelopesBatchId", Field: "start_position", Type: "int"},
	{Object: "Connect_GetConnectUsers", Field: "start_position", Type: "int"},
	{Object: "Envelopes_GetEnvelopes", Field: "start_position", Type: "int"},
	{Object: "Pages_GetPageImages", Field: "start_position", Type: "int"},
	{Object: "Recipients_GetBulkRecipients", Field: "start_position", Type: "int"},
	{Object: "Envelopes_PutStatus", Field: "start_position", Type: "int"},
	{Object: "Folders_GetFolders", Field: "start_position", Type: "int"},
	{Object: "Folders_GetFolderItems", Field: "start_position", Type: "int"},
	{Object: "Groups_GetGroups", Field: "start_position", Type: "int"},
	{Object: "Groups_GetGroupUsers", Field: "start_position", Type: "int"},
	{Object: "PowerForms_GetPowerFormsSenders", Field: "start_position", Type: "int"},
	{Object: "SearchFolders_GetSearchFolderContents", Field: "start_position", Type: "int"},
	{Object: "SharedAccess_GetSharedAccess", Field: "start_position", Type: "int"},
	{Object: "Templates_GetTemplates", Field: "start_position", Type: "int"},
	{Object: "Pages_GetTemplatePageImages", Field: "start_position", Type: "int"},
	{Object: "Recipients_GetTemplateBulkRecipients", Field: "start_position", Type: "int"},
	{Object: "Users_GetUsers", Field: "start_position", Type: "int"},
	{Object: "CloudStorageFolder_GetCloudStorageFolderAll", Field: "start_position", Type: "int"},
	{Object: "CloudStorageFolder_GetCloudStorageFolder", Field: "start_position", Type: "int"},
	{Object: "WorkspaceFolder_GetWorkspaceFolder", Field: "start_position", Type: "int"},
	{Object: "WorkspaceFilePages_GetWorkspaceFilePages", Field: "start_position", Type: "int"},
	{Object: "WorkspaceFilePages_GetWorkspaceFilePages", Field: "dpi", Type: "int"},
	{Object: "WorkspaceFilePages_GetWorkspaceFilePages", Field: "max_height", Type: "int"},
	{Object: "WorkspaceFilePages_GetWorkspaceFilePages", Field: "max_width", Type: "int"},
	{Object: "Users_GetUsers", Field: "include_usersettings_for_csv", Type: "bool"},
	{Object: "Pages_GetTemplatePageImages", Field: "dpi", Type: "int"},
	{Object: "Pages_GetTemplatePageImages", Field: "max_height", Type: "int"},
	{Object: "Pages_GetTemplatePageImages", Field: "max_width", Type: "int"},
	{Object: "Pages_GetTemplatePageImages", Field: "nocache", Type: "bool"},
	{Object: "Pages_GetTemplatePageImages", Field: "show_changes", Type: "bool"},
	{Object: "Pages_GetTemplatePageImage", Field: "dpi", Type: "int"},
	{Object: "Pages_GetTemplatePageImage", Field: "max_height", Type: "int"},
	{Object: "Pages_GetTemplatePageImage", Field: "max_width", Type: "int"},
	{Object: "Pages_GetTemplatePageImage", Field: "show_changes", Type: "bool"},
	{Object: "Documents_GetTemplateDocument", Field: "encrypt", Type: "bool"},
	{Object: "Documents_GetTemplateDocument", Field: "show_changes", Type: "bool"},
	{Object: "SearchFolders_GetSearchFolderContents", Field: "all", Type: "bool"},
	{Object: "Pages_GetPageImages", Field: "dpi", Type: "int"},
	{Object: "Pages_GetPageImages", Field: "max_height", Type: "int"},
	{Object: "Pages_GetPageImages", Field: "max_width", Type: "int"},
	{Object: "Pages_GetPageImages", Field: "nocache", Type: "bool"},
	{Object: "Pages_GetPageImages", Field: "show_changes", Type: "bool"},
	{Object: "Pages_GetPageImage", Field: "dpi", Type: "int"},
	{Object: "Pages_GetPageImage", Field: "max_height", Type: "int"},
	{Object: "Pages_GetPageImage", Field: "max_width", Type: "int"},
	{Object: "Pages_GetPageImage", Field: "show_changes", Type: "bool"},
	{Object: "Brand_GetBrand", Field: "include_external_references", Type: "bool"},
	{Object: "Brand_GetBrand", Field: "include_logos", Type: "bool"},
	{Object: "Connect_GetConnectUsers", Field: "list_included_users", Type: "bool"},
	{Object: "ConnectLog_GetConnectLog", Field: "additional_info", Type: "bool"},
	{Object: "AccountCustomFields_PostAccountCustomFields", Field: "apply_to_templates", Type: "bool"},
	{Object: "AccountCustomFields_PutAccountCustomFields", Field: "apply_to_templates", Type: "bool"},
	{Object: "AccountCustomFields_DeleteAccountCustomFields", Field: "apply_to_templates", Type: "bool"},
	{Object: "Envelopes_GetEnvelope", Field: "advanced_update", Type: "bool"},
	{Object: "Documents_PutTemplateDocument", Field: "is_envelope_definition", Type: "bool"},
	{Object: "Recipients_GetRecipientInitialsImage", Field: "include_chrome", Type: "bool"},
	{Object: "UserSignatures_GetUserSignatureImage", Field: "include_chrome", Type: "bool"},
	{Object: "BulkEnvelopes_GetEnvelopesBulk", Field: "include", Type: "...string"},
	{Object: "BulkEnvelopes_GetBulkEnvelopesBatchId", Field: "include", Type: "...string"},
	{Object: "ChunkedUploads_GetChunkedUpload", Field: "include", Type: "...string"},
	{Object: "Connect_GetConnectUsers", Field: "status", Type: "...string"},
	{Object: "Envelopes_GetEnvelopes", Field: "envelope_ids", Type: "...string"},
	{Object: "Envelopes_GetEnvelopes", Field: "folder_ids", Type: "...string"},
	{Object: "Envelopes_GetEnvelopes", Field: "intersecting_folder_ids", Type: "...string"},
	{Object: "Envelopes_GetEnvelopes", Field: "powerformids", Type: "...string"},
	{Object: "Envelopes_GetEnvelopes", Field: "status", Type: "...string"},
	{Object: "Envelopes_GetEnvelopes", Field: "transaction_ids", Type: "...string"},
	{Object: "Templates_GetDocumentTemplates", Field: "include", Type: "...string"},
	{Object: "PermissionProfiles_PostPermissionProfiles", Field: "include", Type: "...string"},
	{Object: "PermissionProfiles_GetPermissionProfile", Field: "include", Type: "...string"},
	{Object: "PermissionProfiles_PutPermissionProfiles", Field: "include", Type: "...string"},
	{Object: "SharedAccess_GetSharedAccess", Field: "folder_ids", Type: "...string"},
	{Object: "SharedAccess_GetSharedAccess", Field: "user_ids", Type: "...string"},
	{Object: "SharedAccess_PutSharedAccess", Field: "user_ids", Type: "...string"},
	{Object: "Templates_GetTemplates", Field: "folder_ids", Type: "...string"},
	{Object: "Templates_GetTemplates", Field: "include", Type: "...string"},
	{Object: "Templates_GetTemplate", Field: "include", Type: "...string"},
	{Object: "Users_GetUsers", Field: "status", Type: "...string"},
	{Object: "CloudStorageFolder_GetCloudStorageFolderAll", Field: "cloud_storage_folder_path", Type: "...string"},
}

// fieldOverrides provides a list of fields and their new type. In
// the specification, DocuSign lists every field as a string.   I generated much
// of this list with the following rules.
// - definition properties with **true** in the description are
//   assumed to be bools set to DSBool
// - fields containing base64 in the name are assumed to be []byte
// - fields ending in DateTime are *time.Time
var fieldOverrides = []override{
	{Object: "recipientAttachment", Field: "data", Type: "[]byte"},
	{Object: "attachment", Field: "data", Type: "[]byte"},
	{Object: "chunkedUploadRequest", Field: "data", Type: "[]byte"},
	{Object: "document", Field: "documentBase64", Type: "[]byte"},
	{Object: "AccountWatermarks", Field: "imageBase64", Type: "[]byte"},
	{Object: "purchasedEnvelopesInformation", Field: "receiptData", Type: "[]byte"},
	{Object: "accountBillingPlan", Field: "canUpgrade", Type: "DSBool"},
	{Object: "accountBillingPlan", Field: "enableSupport", Type: "DSBool"},
	{Object: "accountBillingPlanResponse", Field: "billingAddressIsCreditCardAddress", Type: "DSBool"},
	{Object: "accountInformation", Field: "allowTransactionRooms", Type: "DSBool"},
	{Object: "accountInformation", Field: "canUpgrade", Type: "DSBool"},
	{Object: "addressInformationInput", Field: "receiveInResponse", Type: "DSBool"},
	{Object: "agent", Field: "clientUserId", Type: "DSBool"},
	{Object: "agent", Field: "excludedDocuments", Type: "DSBool"},
	{Object: "agent", Field: "inheritEmailNotificationConfiguration", Type: "DSBool"},
	{Object: "agent", Field: "phoneAuthentication", Type: "DSBool"},
	{Object: "agent", Field: "requireIdLookup", Type: "DSBool"},
	{Object: "agent", Field: "signingGroupId", Type: "DSBool"},
	{Object: "agent", Field: "templateLocked", Type: "DSBool"},
	{Object: "agent", Field: "templateRequired", Type: "DSBool"},
	{Object: "approve", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "approve", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "approve", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "approve", Field: "bold", Type: "DSBool"},
	{Object: "approve", Field: "italic", Type: "DSBool"},
	{Object: "approve", Field: "templateLocked", Type: "DSBool"},
	{Object: "approve", Field: "templateRequired", Type: "DSBool"},
	{Object: "approve", Field: "underline", Type: "DSBool"},
	{Object: "billingPaymentItem", Field: "paymentNumber", Type: "DSBool"},
	{Object: "billingPlan", Field: "enableSupport", Type: "DSBool"},
	{Object: "billingPlanInformation", Field: "enableSupport", Type: "DSBool"},
	{Object: "captiveRecipient", Field: "clientUserId", Type: "DSBool"},
	{Object: "carbonCopy", Field: "clientUserId", Type: "DSBool"},
	{Object: "carbonCopy", Field: "excludedDocuments", Type: "DSBool"},
	{Object: "carbonCopy", Field: "inheritEmailNotificationConfiguration", Type: "DSBool"},
	{Object: "carbonCopy", Field: "phoneAuthentication", Type: "DSBool"},
	{Object: "carbonCopy", Field: "requireIdLookup", Type: "DSBool"},
	{Object: "carbonCopy", Field: "signingGroupId", Type: "DSBool"},
	{Object: "carbonCopy", Field: "templateLocked", Type: "DSBool"},
	{Object: "carbonCopy", Field: "templateRequired", Type: "DSBool"},
	{Object: "certifiedDelivery", Field: "clientUserId", Type: "DSBool"},
	{Object: "certifiedDelivery", Field: "excludedDocuments", Type: "DSBool"},
	{Object: "certifiedDelivery", Field: "inheritEmailNotificationConfiguration", Type: "DSBool"},
	{Object: "certifiedDelivery", Field: "phoneAuthentication", Type: "DSBool"},
	{Object: "certifiedDelivery", Field: "requireIdLookup", Type: "DSBool"},
	{Object: "certifiedDelivery", Field: "signingGroupId", Type: "DSBool"},
	{Object: "certifiedDelivery", Field: "templateLocked", Type: "DSBool"},
	{Object: "certifiedDelivery", Field: "templateRequired", Type: "DSBool"},
	{Object: "checkbox", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "checkbox", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "checkbox", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "checkbox", Field: "locked", Type: "DSBool"},
	{Object: "checkbox", Field: "required", Type: "DSBool"},
	{Object: "checkbox", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "checkbox", Field: "selected", Type: "DSBool"},
	{Object: "checkbox", Field: "shared", Type: "DSBool"},
	{Object: "checkbox", Field: "templateLocked", Type: "DSBool"},
	{Object: "checkbox", Field: "templateRequired", Type: "DSBool"},
	{Object: "company", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "company", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "company", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "company", Field: "bold", Type: "DSBool"},
	{Object: "company", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "company", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "company", Field: "italic", Type: "DSBool"},
	{Object: "company", Field: "locked", Type: "DSBool"},
	{Object: "company", Field: "required", Type: "DSBool"},
	{Object: "company", Field: "templateLocked", Type: "DSBool"},
	{Object: "company", Field: "templateRequired", Type: "DSBool"},
	{Object: "company", Field: "underline", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "allowEnvelopePublish", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "allUsers", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "enableLog", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "includeCertificateOfCompletion", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "includeDocumentFields", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "includeDocuments", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "includeEnvelopeVoidReason", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "includeSenderAccountasCustomField", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "includeTimeZoneInformation", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "requiresAcknowledgement", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "signMessageWithX509Certificate", Type: "DSBool"},
	{Object: "connectCustomConfiguration", Field: "useSoapInterface", Type: "DSBool"},
	{Object: "contact", Field: "shared", Type: "DSBool"},
	{Object: "currencyFeatureSetPrice", Field: "envelopeFee", Type: "DSBool"},
	{Object: "currencyFeatureSetPrice", Field: "fixedFee", Type: "DSBool"},
	{Object: "currencyFeatureSetPrice", Field: "seatFee", Type: "DSBool"},
	{Object: "customField", Field: "required", Type: "DSBool"},
	{Object: "customField_v2", Field: "required", Type: "DSBool"},
	{Object: "date", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "date", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "date", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "date", Field: "bold", Type: "DSBool"},
	{Object: "date", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "date", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "date", Field: "italic", Type: "DSBool"},
	{Object: "date", Field: "locked", Type: "DSBool"},
	{Object: "date", Field: "requireAll", Type: "DSBool"},
	{Object: "date", Field: "required", Type: "DSBool"},
	{Object: "date", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "date", Field: "senderRequired", Type: "DSBool"},
	{Object: "date", Field: "shared", Type: "DSBool"},
	{Object: "date", Field: "templateLocked", Type: "DSBool"},
	{Object: "date", Field: "templateRequired", Type: "DSBool"},
	{Object: "date", Field: "underline", Type: "DSBool"},
	{Object: "dateSigned", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "dateSigned", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "dateSigned", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "dateSigned", Field: "bold", Type: "DSBool"},
	{Object: "dateSigned", Field: "italic", Type: "DSBool"},
	{Object: "dateSigned", Field: "templateLocked", Type: "DSBool"},
	{Object: "dateSigned", Field: "templateRequired", Type: "DSBool"},
	{Object: "dateSigned", Field: "underline", Type: "DSBool"},
	{Object: "decline", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "decline", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "decline", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "decline", Field: "bold", Type: "DSBool"},
	{Object: "decline", Field: "italic", Type: "DSBool"},
	{Object: "decline", Field: "templateLocked", Type: "DSBool"},
	{Object: "decline", Field: "templateRequired", Type: "DSBool"},
	{Object: "decline", Field: "underline", Type: "DSBool"},
	{Object: "diagnosticsSettingsInformation", Field: "apiRequestLogging", Type: "DSBool"},
	{Object: "dobInformationInput", Field: "receiveInResponse", Type: "DSBool"},
	{Object: "document", Field: "encryptedWithKeyManager", Type: "DSBool"},
	{Object: "document", Field: "includeInDownload", Type: "DSBool"},
	{Object: "document", Field: "templateLocked", Type: "DSBool"},
	{Object: "document", Field: "templateRequired", Type: "DSBool"},
	{Object: "document", Field: "transformPdfFields", Type: "DSBool"},
	{Object: "editor", Field: "clientUserId", Type: "DSBool"},
	{Object: "editor", Field: "inheritEmailNotificationConfiguration", Type: "DSBool"},
	{Object: "editor", Field: "phoneAuthentication", Type: "DSBool"},
	{Object: "editor", Field: "requireIdLookup", Type: "DSBool"},
	{Object: "editor", Field: "signingGroupId", Type: "DSBool"},
	{Object: "editor", Field: "templateLocked", Type: "DSBool"},
	{Object: "editor", Field: "templateRequired", Type: "DSBool"},
	{Object: "email", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "email", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "email", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "email", Field: "bold", Type: "DSBool"},
	{Object: "email", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "email", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "email", Field: "italic", Type: "DSBool"},
	{Object: "email", Field: "locked", Type: "DSBool"},
	{Object: "email", Field: "requireAll", Type: "DSBool"},
	{Object: "email", Field: "required", Type: "DSBool"},
	{Object: "email", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "email", Field: "senderRequired", Type: "DSBool"},
	{Object: "email", Field: "shared", Type: "DSBool"},
	{Object: "email", Field: "templateLocked", Type: "DSBool"},
	{Object: "email", Field: "templateRequired", Type: "DSBool"},
	{Object: "email", Field: "underline", Type: "DSBool"},
	{Object: "emailAddress", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "emailAddress", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "emailAddress", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "emailAddress", Field: "bold", Type: "DSBool"},
	{Object: "emailAddress", Field: "italic", Type: "DSBool"},
	{Object: "emailAddress", Field: "templateLocked", Type: "DSBool"},
	{Object: "emailAddress", Field: "templateRequired", Type: "DSBool"},
	{Object: "emailAddress", Field: "underline", Type: "DSBool"},
	{Object: "envelope", Field: "allowMarkup", Type: "DSBool"},
	{Object: "envelope", Field: "allowReassign", Type: "DSBool"},
	{Object: "envelope", Field: "asynchronous", Type: "DSBool"},
	{Object: "envelope", Field: "enableWetSign", Type: "DSBool"},
	{Object: "envelope", Field: "enforceSignerVisibility", Type: "DSBool"},
	{Object: "envelope", Field: "envelopeIdStamping", Type: "DSBool"},
	{Object: "envelope", Field: "is21CFRPart11", Type: "DSBool"},
	{Object: "envelope", Field: "messageLock", Type: "DSBool"},
	{Object: "envelope", Field: "notification", Type: "DSBool"},
	{Object: "envelope", Field: "recipientsLock", Type: "DSBool"},
	{Object: "envelope", Field: "useDisclosure", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "allowMarkup", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "allowReassign", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "allowRecipientRecursion", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "asynchronous", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "enableWetSign", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "enforceSignerVisibility", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "envelopeIdStamping", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "is21CFRPart11", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "messageLock", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "recipientsLock", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "templateRoles", Type: "DSBool"},
	{Object: "envelopeDefinition", Field: "useDisclosure", Type: "DSBool"},
	{Object: "envelopeDocument", Field: "includeInDownload", Type: "DSBool"},
	{Object: "envelopeId", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "envelopeId", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "envelopeId", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "envelopeId", Field: "bold", Type: "DSBool"},
	{Object: "envelopeId", Field: "italic", Type: "DSBool"},
	{Object: "envelopeId", Field: "templateLocked", Type: "DSBool"},
	{Object: "envelopeId", Field: "templateRequired", Type: "DSBool"},
	{Object: "envelopeId", Field: "underline", Type: "DSBool"},
	{Object: "envelopeNotificationRequest", Field: "useAccountDefaults", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "allowMarkup", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "allowReassign", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "asynchronous", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "enableWetSign", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "enforceSignerVisibility", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "envelopeIdStamping", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "envelopeTemplateDefinition", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "is21CFRPart11", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "messageLock", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "recipientsLock", Type: "DSBool"},
	{Object: "envelopeTemplate", Field: "useDisclosure", Type: "DSBool"},
	{Object: "envelopeTemplateDefinition", Field: "shared", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "allowMarkup", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "allowReassign", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "asynchronous", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "enableWetSign", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "enforceSignerVisibility", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "envelopeIdStamping", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "is21CFRPart11", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "messageLock", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "recipientsLock", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "shared", Type: "DSBool"},
	{Object: "envelopeTemplateResult", Field: "useDisclosure", Type: "DSBool"},
	{Object: "eventNotification", Field: "includeCertificateOfCompletion", Type: "DSBool"},
	{Object: "eventNotification", Field: "includeCertificateWithSoap", Type: "DSBool"},
	{Object: "eventNotification", Field: "includeDocumentFields", Type: "DSBool"},
	{Object: "eventNotification", Field: "includeDocuments", Type: "DSBool"},
	{Object: "eventNotification", Field: "includeEnvelopeVoidReason", Type: "DSBool"},
	{Object: "eventNotification", Field: "includeSenderAccountAsCustomField", Type: "DSBool"},
	{Object: "eventNotification", Field: "includeTimeZone", Type: "DSBool"},
	{Object: "eventNotification", Field: "loggingEnabled", Type: "DSBool"},
	{Object: "eventNotification", Field: "requireAcknowledgment", Type: "DSBool"},
	{Object: "eventNotification", Field: "signMessageWithX509Cert", Type: "DSBool"},
	{Object: "eventNotification", Field: "useSoapInterface", Type: "DSBool"},
	{Object: "expirations", Field: "expireEnabled", Type: "DSBool"},
	{Object: "featureSet", Field: "is21CFRPart11", Type: "DSBool"},
	{Object: "firstName", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "firstName", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "firstName", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "firstName", Field: "bold", Type: "DSBool"},
	{Object: "firstName", Field: "italic", Type: "DSBool"},
	{Object: "firstName", Field: "templateLocked", Type: "DSBool"},
	{Object: "firstName", Field: "templateRequired", Type: "DSBool"},
	{Object: "firstName", Field: "underline", Type: "DSBool"},
	{Object: "folderItem", Field: "is21CFRPart11", Type: "DSBool"},
	{Object: "folderItem", Field: "shared", Type: "DSBool"},
	{Object: "folderItem_v2", Field: "is21CFRPart11", Type: "DSBool"},
	{Object: "formulaTab", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "formulaTab", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "formulaTab", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "formulaTab", Field: "bold", Type: "DSBool"},
	{Object: "formulaTab", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "formulaTab", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "formulaTab", Field: "isPaymentAmount", Type: "DSBool"},
	{Object: "formulaTab", Field: "italic", Type: "DSBool"},
	{Object: "formulaTab", Field: "locked", Type: "DSBool"},
	{Object: "formulaTab", Field: "requireAll", Type: "DSBool"},
	{Object: "formulaTab", Field: "required", Type: "DSBool"},
	{Object: "formulaTab", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "formulaTab", Field: "senderRequired", Type: "DSBool"},
	{Object: "formulaTab", Field: "shared", Type: "DSBool"},
	{Object: "formulaTab", Field: "templateLocked", Type: "DSBool"},
	{Object: "formulaTab", Field: "templateRequired", Type: "DSBool"},
	{Object: "formulaTab", Field: "underline", Type: "DSBool"},
	{Object: "fullName", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "fullName", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "fullName", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "fullName", Field: "bold", Type: "DSBool"},
	{Object: "fullName", Field: "italic", Type: "DSBool"},
	{Object: "fullName", Field: "templateLocked", Type: "DSBool"},
	{Object: "fullName", Field: "templateRequired", Type: "DSBool"},
	{Object: "fullName", Field: "underline", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "canSignOffline", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "clientUserId", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "defaultRecipient", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "inheritEmailNotificationConfiguration", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "phoneAuthentication", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "requireIdLookup", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "requireSignOnPaper", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "signInEachLocation", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "signingGroupId", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "templateLocked", Type: "DSBool"},
	{Object: "inPersonSigner", Field: "templateRequired", Type: "DSBool"},
	{Object: "initialHere", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "initialHere", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "initialHere", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "initialHere", Field: "templateLocked", Type: "DSBool"},
	{Object: "initialHere", Field: "templateRequired", Type: "DSBool"},
	{Object: "intermediary", Field: "clientUserId", Type: "DSBool"},
	{Object: "intermediary", Field: "excludedDocuments", Type: "DSBool"},
	{Object: "intermediary", Field: "inheritEmailNotificationConfiguration", Type: "DSBool"},
	{Object: "intermediary", Field: "phoneAuthentication", Type: "DSBool"},
	{Object: "intermediary", Field: "requireIdLookup", Type: "DSBool"},
	{Object: "intermediary", Field: "signingGroupId", Type: "DSBool"},
	{Object: "intermediary", Field: "templateLocked", Type: "DSBool"},
	{Object: "intermediary", Field: "templateRequired", Type: "DSBool"},
	{Object: "lastName", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "lastName", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "lastName", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "lastName", Field: "bold", Type: "DSBool"},
	{Object: "lastName", Field: "italic", Type: "DSBool"},
	{Object: "lastName", Field: "templateLocked", Type: "DSBool"},
	{Object: "lastName", Field: "templateRequired", Type: "DSBool"},
	{Object: "lastName", Field: "underline", Type: "DSBool"},
	{Object: "list", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "list", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "list", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "list", Field: "bold", Type: "DSBool"},
	{Object: "list", Field: "italic", Type: "DSBool"},
	{Object: "list", Field: "locked", Type: "DSBool"},
	{Object: "list", Field: "requireAll", Type: "DSBool"},
	{Object: "list", Field: "required", Type: "DSBool"},
	{Object: "list", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "list", Field: "senderRequired", Type: "DSBool"},
	{Object: "list", Field: "shared", Type: "DSBool"},
	{Object: "list", Field: "templateLocked", Type: "DSBool"},
	{Object: "list", Field: "templateRequired", Type: "DSBool"},
	{Object: "list", Field: "underline", Type: "DSBool"},
	{Object: "listCustomField", Field: "required", Type: "DSBool"},
	{Object: "listItem", Field: "selected", Type: "DSBool"},
	{Object: "memberGroupSharedItem", Field: "shared", Type: "DSBool"},
	{Object: "mergeField", Field: "allowSenderToEdit", Type: "DSBool"},
	{Object: "notaryHost", Field: "clientUserId", Type: "DSBool"},
	{Object: "notaryHost", Field: "inheritEmailNotificationConfiguration", Type: "DSBool"},
	{Object: "notaryHost", Field: "phoneAuthentication", Type: "DSBool"},
	{Object: "notaryHost", Field: "requireIdLookup", Type: "DSBool"},
	{Object: "notaryHost", Field: "templateLocked", Type: "DSBool"},
	{Object: "notaryHost", Field: "templateRequired", Type: "DSBool"},
	{Object: "note", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "note", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "note", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "note", Field: "bold", Type: "DSBool"},
	{Object: "note", Field: "italic", Type: "DSBool"},
	{Object: "note", Field: "shared", Type: "DSBool"},
	{Object: "note", Field: "templateLocked", Type: "DSBool"},
	{Object: "note", Field: "templateRequired", Type: "DSBool"},
	{Object: "note", Field: "underline", Type: "DSBool"},
	{Object: "notification", Field: "useAccountDefaults", Type: "DSBool"},
	{Object: "number", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "number", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "number", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "number", Field: "bold", Type: "DSBool"},
	{Object: "number", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "number", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "number", Field: "isPaymentAmount", Type: "DSBool"},
	{Object: "number", Field: "italic", Type: "DSBool"},
	{Object: "number", Field: "locked", Type: "DSBool"},
	{Object: "number", Field: "requireAll", Type: "DSBool"},
	{Object: "number", Field: "required", Type: "DSBool"},
	{Object: "number", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "number", Field: "senderRequired", Type: "DSBool"},
	{Object: "number", Field: "shared", Type: "DSBool"},
	{Object: "number", Field: "templateLocked", Type: "DSBool"},
	{Object: "number", Field: "templateRequired", Type: "DSBool"},
	{Object: "number", Field: "underline", Type: "DSBool"},
	{Object: "radio", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "radio", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "radio", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "radio", Field: "locked", Type: "DSBool"},
	{Object: "radio", Field: "required", Type: "DSBool"},
	{Object: "radio", Field: "selected", Type: "DSBool"},
	{Object: "radioGroup", Field: "requireAll", Type: "DSBool"},
	{Object: "radioGroup", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "radioGroup", Field: "shared", Type: "DSBool"},
	{Object: "recipientPhoneAuthentication", Field: "recipMayProvideNumber", Type: "DSBool"},
	{Object: "referralInformation", Field: "enableSupport", Type: "DSBool"},
	{Object: "reminders", Field: "reminderEnabled", Type: "DSBool"},
	{Object: "senderEmailNotifications", Field: "changedSigner", Type: "DSBool"},
	{Object: "senderEmailNotifications", Field: "deliveryFailed", Type: "DSBool"},
	{Object: "senderEmailNotifications", Field: "envelopeComplete", Type: "DSBool"},
	{Object: "senderEmailNotifications", Field: "offlineSigningFailed", Type: "DSBool"},
	{Object: "senderEmailNotifications", Field: "recipientViewed", Type: "DSBool"},
	{Object: "senderEmailNotifications", Field: "withdrawnConsent", Type: "DSBool"},
	{Object: "settingsMetadata", Field: "is21CFRPart11", Type: "DSBool"},
	{Object: "sharedItem", Field: "shared", Type: "DSBool"},
	{Object: "signHere", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "signHere", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "signHere", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "signHere", Field: "templateLocked", Type: "DSBool"},
	{Object: "signHere", Field: "templateRequired", Type: "DSBool"},
	{Object: "signer", Field: "canSignOffline", Type: "DSBool"},
	{Object: "signer", Field: "clientUserId", Type: "DSBool"},
	{Object: "signer", Field: "defaultRecipient", Type: "DSBool"},
	{Object: "signer", Field: "excludedDocuments", Type: "DSBool"},
	{Object: "signer", Field: "inheritEmailNotificationConfiguration", Type: "DSBool"},
	{Object: "signer", Field: "isBulkRecipient", Type: "DSBool"},
	{Object: "signer", Field: "phoneAuthentication", Type: "DSBool"},
	{Object: "signer", Field: "requireIdLookup", Type: "DSBool"},
	{Object: "signer", Field: "requireSignOnPaper", Type: "DSBool"},
	{Object: "signer", Field: "signInEachLocation", Type: "DSBool"},
	{Object: "signer", Field: "signingGroupId", Type: "DSBool"},
	{Object: "signer", Field: "templateLocked", Type: "DSBool"},
	{Object: "signer", Field: "templateRequired", Type: "DSBool"},
	{Object: "signerAttachment", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "signerAttachment", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "signerAttachment", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "signerAttachment", Field: "templateLocked", Type: "DSBool"},
	{Object: "signerAttachment", Field: "templateRequired", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "agentNotification", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "carbonCopyNotification", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "certifiedDeliveryNotification", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "documentMarkupActivation", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "envelopeActivation", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "envelopeComplete", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "envelopeCorrected", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "envelopeDeclined", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "envelopeVoided", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "offlineSigningFailed", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "purgeDocuments", Type: "DSBool"},
	{Object: "signerEmailNotifications", Field: "reassignedSigner", Type: "DSBool"},
	{Object: "signingGroup", Field: "signingGroupId", Type: "DSBool"},
	{Object: "ssn", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "ssn", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "ssn", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "ssn", Field: "bold", Type: "DSBool"},
	{Object: "ssn", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "ssn", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "ssn", Field: "italic", Type: "DSBool"},
	{Object: "ssn", Field: "locked", Type: "DSBool"},
	{Object: "ssn", Field: "requireAll", Type: "DSBool"},
	{Object: "ssn", Field: "required", Type: "DSBool"},
	{Object: "ssn", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "ssn", Field: "senderRequired", Type: "DSBool"},
	{Object: "ssn", Field: "shared", Type: "DSBool"},
	{Object: "ssn", Field: "templateLocked", Type: "DSBool"},
	{Object: "ssn", Field: "templateRequired", Type: "DSBool"},
	{Object: "ssn", Field: "underline", Type: "DSBool"},
	{Object: "ssn4InformationInput", Field: "receiveInResponse", Type: "DSBool"},
	{Object: "tabMetadata", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "tabMetadata", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "tabMetadata", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "tabMetadata", Field: "bold", Type: "DSBool"},
	{Object: "tabMetadata", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "tabMetadata", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "tabMetadata", Field: "editable", Type: "DSBool"},
	{Object: "tabMetadata", Field: "includedInEmail", Type: "DSBool"},
	{Object: "tabMetadata", Field: "italic", Type: "DSBool"},
	{Object: "tabMetadata", Field: "locked", Type: "DSBool"},
	{Object: "tabMetadata", Field: "required", Type: "DSBool"},
	{Object: "tabMetadata", Field: "shared", Type: "DSBool"},
	{Object: "tabMetadata", Field: "underline", Type: "DSBool"},
	{Object: "templateNotificationRequest", Field: "useAccountDefaults", Type: "DSBool"},
	{Object: "templateRole", Field: "clientUserId", Type: "DSBool"},
	{Object: "templateRole", Field: "defaultRecipient", Type: "DSBool"},
	{Object: "templateRole", Field: "signingGroupId", Type: "DSBool"},
	{Object: "templateSharedItem", Field: "shared", Type: "DSBool"},
	{Object: "text", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "text", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "text", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "text", Field: "bold", Type: "DSBool"},
	{Object: "text", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "text", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "text", Field: "isPaymentAmount", Type: "DSBool"},
	{Object: "text", Field: "italic", Type: "DSBool"},
	{Object: "text", Field: "locked", Type: "DSBool"},
	{Object: "text", Field: "requireAll", Type: "DSBool"},
	{Object: "text", Field: "required", Type: "DSBool"},
	{Object: "text", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "text", Field: "senderRequired", Type: "DSBool"},
	{Object: "text", Field: "shared", Type: "DSBool"},
	{Object: "text", Field: "templateLocked", Type: "DSBool"},
	{Object: "text", Field: "templateRequired", Type: "DSBool"},
	{Object: "text", Field: "underline", Type: "DSBool"},
	{Object: "textCustomField", Field: "required", Type: "DSBool"},
	{Object: "title", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "title", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "title", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "title", Field: "bold", Type: "DSBool"},
	{Object: "title", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "title", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "title", Field: "italic", Type: "DSBool"},
	{Object: "title", Field: "locked", Type: "DSBool"},
	{Object: "title", Field: "required", Type: "DSBool"},
	{Object: "title", Field: "templateLocked", Type: "DSBool"},
	{Object: "title", Field: "templateRequired", Type: "DSBool"},
	{Object: "title", Field: "underline", Type: "DSBool"},
	{Object: "userInformation", Field: "sendActivationOnInvalidLogin", Type: "DSBool"},
	{Object: "userProfile", Field: "displayOrganizationInfo", Type: "DSBool"},
	{Object: "userProfile", Field: "displayPersonalInfo", Type: "DSBool"},
	{Object: "userProfile", Field: "displayProfile", Type: "DSBool"},
	{Object: "userProfile", Field: "displayUsageHistory", Type: "DSBool"},
	{Object: "userSharedItem", Field: "shared", Type: "DSBool"},
	{Object: "view", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "view", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "view", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "view", Field: "bold", Type: "DSBool"},
	{Object: "view", Field: "italic", Type: "DSBool"},
	{Object: "view", Field: "required", Type: "DSBool"},
	{Object: "view", Field: "templateLocked", Type: "DSBool"},
	{Object: "view", Field: "templateRequired", Type: "DSBool"},
	{Object: "view", Field: "underline", Type: "DSBool"},
	{Object: "zip", Field: "anchorCaseSensitive", Type: "DSBool"},
	{Object: "zip", Field: "anchorIgnoreIfNotPresent", Type: "DSBool"},
	{Object: "zip", Field: "anchorMatchWholeWord", Type: "DSBool"},
	{Object: "zip", Field: "bold", Type: "DSBool"},
	{Object: "zip", Field: "concealValueOnDocument", Type: "DSBool"},
	{Object: "zip", Field: "disableAutoSize", Type: "DSBool"},
	{Object: "zip", Field: "italic", Type: "DSBool"},
	{Object: "zip", Field: "locked", Type: "DSBool"},
	{Object: "zip", Field: "requireAll", Type: "DSBool"},
	{Object: "zip", Field: "required", Type: "DSBool"},
	{Object: "zip", Field: "requireInitialOnSharedChange", Type: "DSBool"},
	{Object: "zip", Field: "senderRequired", Type: "DSBool"},
	{Object: "zip", Field: "shared", Type: "DSBool"},
	{Object: "zip", Field: "templateLocked", Type: "DSBool"},
	{Object: "zip", Field: "templateRequired", Type: "DSBool"},
	{Object: "zip", Field: "underline", Type: "DSBool"},
	{Object: "agent", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "agent", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "agent", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "agent", Field: "signedDateTime", Type: "*time.Time"},
	{Object: "apiRequestLog", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "bulkEnvelope", Field: "submittedDateTime", Type: "*time.Time"},
	{Object: "carbonCopy", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "carbonCopy", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "carbonCopy", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "carbonCopy", Field: "signedDateTime", Type: "*time.Time"},
	{Object: "certifiedDelivery", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "certifiedDelivery", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "certifiedDelivery", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "certifiedDelivery", Field: "signedDateTime", Type: "*time.Time"},
	{Object: "chunkedUploadResponse", Field: "expirationDateTime", Type: "*time.Time"},
	{Object: "connectDebugLog", Field: "eventDateTime", Type: "*time.Time"},
	{Object: "editor", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "editor", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "editor", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "editor", Field: "signedDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "completedDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "deletedDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "initialSentDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "lastModifiedDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "statusChangedDateTime", Type: "*time.Time"},
	{Object: "envelope", Field: "voidedDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "completedDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "deletedDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "initialSentDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "lastModifiedDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "statusChangedDateTime", Type: "*time.Time"},
	{Object: "envelopeDefinition", Field: "voidedDateTime", Type: "*time.Time"},
	{Object: "envelopeFormData", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "envelopeSummary", Field: "statusDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "completedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "deletedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "initialSentDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "lastModifiedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "statusChangedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplate", Field: "voidedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "completedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "deletedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "initialSentDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "lastModifiedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "statusChangedDateTime", Type: "*time.Time"},
	{Object: "envelopeTemplateResult", Field: "voidedDateTime", Type: "*time.Time"},
	{Object: "filter", Field: "fromDateTime", Type: "*time.Time"},
	{Object: "filter", Field: "toDateTime", Type: "*time.Time"},
	{Object: "folderItem", Field: "completedDateTime", Type: "*time.Time"},
	{Object: "folderItem", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "folderItem", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "folderItem_v2", Field: "completedDateTime", Type: "*time.Time"},
	{Object: "folderItem_v2", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "folderItem_v2", Field: "expireDateTime", Type: "*time.Time"},
	{Object: "folderItem_v2", Field: "lastModifiedDateTime", Type: "*time.Time"},
	{Object: "folderItem_v2", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "inPersonSigner", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "inPersonSigner", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "inPersonSigner", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "inPersonSigner", Field: "signedDateTime", Type: "*time.Time"},
	{Object: "intermediary", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "intermediary", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "intermediary", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "intermediary", Field: "signedDateTime", Type: "*time.Time"},
	{Object: "lockInformation", Field: "lockedUntilDateTime", Type: "*time.Time"},
	{Object: "lockInformation", Field: "lockedUntilDateTime", Type: "*time.Time"},
	{Object: "newUser", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "notaryHost", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "notaryHost", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "notaryHost", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "notaryHost", Field: "signedDateTime", Type: "*time.Time"},
	{Object: "permissionProfile", Field: "modifiedDateTime", Type: "*time.Time"},
	{Object: "powerForm", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "serviceInformation", Field: "buildBranchDeployedDateTime", Type: "*time.Time"},
	{Object: "signer", Field: "declinedDateTime", Type: "*time.Time"},
	{Object: "signer", Field: "deliveredDateTime", Type: "*time.Time"},
	{Object: "signer", Field: "sentDateTime", Type: "*time.Time"},
	{Object: "signer", Field: "signedDateTime", Type: "*time.Time"},
	{Object: "usageHistory", Field: "lastSentDateTime", Type: "*time.Time"},
	{Object: "usageHistory", Field: "lastSignedDateTime", Type: "*time.Time"},
	{Object: "userInformation", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "userSignature", Field: "adoptedDateTime", Type: "*time.Time"},
	{Object: "userSignature", Field: "createdDateTime", Type: "*time.Time"},
	{Object: "recipientFormData", Field: "declinedTime", Type: "*time.Time"},
	{Object: "recipientFormData", Field: "deliveredTime", Type: "*time.Time"},
	{Object: "recipientFormData", Field: "sentTime", Type: "*time.Time"},
	{Object: "recipientFormData", Field: "signedTime", Type: "*time.Time"},
}

// TabDefs creates a list definitions for embedded tab structs from the defMap parameter.
// overrides is updated with new override entries to allow tab definitions to generate.
func TabDefs(defMap map[string]swagger.Definition, overrides map[string]map[string]string) []swagger.Definition {
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
		"note",
		"number",
		"radioGroup",
		"signerAttachment",
		"signHere",
		"ssn",
		"text",
		"title",
		"view",
		"zip",
	}

	// list of types of tabs
	tabDefs := map[string]swagger.Definition{
		"Base": swagger.Definition{
			ID:          "TabBase",
			Name:        "TabBase",
			Type:        "Object",
			Description: "contains common fields for all tabs",
			Summary:     "contains common fields for all tabs",
			Category:    "",
		},
		"Position": swagger.Definition{
			ID:          "TabPosition",
			Name:        "TabPosition",
			Type:        "Object",
			Description: "contains common fields for all tabs that can position themselves on document",
			Summary:     "contains common fields for all tabs that can position themselves on document",
			Category:    "",
		},
		"Style": swagger.Definition{
			ID:          "TabStyle",
			Name:        "TabStyle",
			Type:        "Object",
			Description: "contains common fields for all tabs that can set a display style",
			Summary:     "contains common fields for all tabs that can set a display style",
			Category:    "",
		},
		"Value": swagger.Definition{
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
		"Base": []string{
			"conditionalParentLabel",
			"conditionalParentValue",
			"documentId",
			"recipientId",
		},
		"Position": []string{
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
		"Style": []string{
			"bold",
			"font",
			"fontColor",
			"fontSize",
			"italic",
			"name",
			"underline",
		},
		"Value": []string{
			"value",
		},
	}
	// loop thru each tab definition
	for _, tabname := range tabObjects {
		dx := defMap["#/definitions/"+tabname]
		// create map of fields for easy lookup
		xmap := make(map[string]bool)
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
	xmap := make(map[string]swagger.Field)
	for _, f := range txtDef.Fields {
		xmap[f.Name] = f
	}
	results := make([]swagger.Definition, 0)
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

// CustomCode is lines of code to append to model.go
const CustomCode = `// GetTabValues returns a NameValue list of all entry tabs
func GetTabValues(tabs Tabs) []NameValue {
	results := make([]NameValue, 0)
	for _, v := range tabs.CheckboxTabs {
        fmt.Println("oK")
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
	for _, v := range tabs.ListTabs {
		vals := make([]string, 0, len(v.ListItems))
		for _, x := range v.ListItems {
			if x.Selected {
				vals = append(vals, x.Value)
			}
		}
		results = append(results, NameValue{Name: v.TabLabel, Value: strings.Join(vals, ",")})
	}
	for _, v := range tabs.NoteTabs {
		results = append(results, NameValue{Name: v.TabLabel, Value: v.Value})
	}
	for _, v := range tabs.NumberTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.RadioGroupTabs {
		for _, x := range v.Radios {
			if x.Selected {
				results = append(results, NameValue{Name: v.GroupName, Value: x.Value})
				break
			}
		}
	}
	for _, v := range tabs.SSNTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	for _, v := range tabs.TextTabs {
		results = append(results, NameValue{Name: v.TabLabel, OriginalValue: v.OriginalValue, Value: v.Value})
	}
	return results
}`
