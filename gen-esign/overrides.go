// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file provides lists of type overrides for definition
// properties (struct fields) and for operation parameters found in
// Docusign's Rest API swagger definition.
package main

import (
	"strings"
)

// ServiceNameOverride provides map of new names x-ds-service
// value
var ServiceNameOverride = map[string]string{
	"Groups": "UserGroups",
}

// OperationSkipList contains operations ignore
// due to deprecation or incomplete definitions
var OperationSkipList = map[string]bool{
	"OAuth2_PostRevoke": true, // PostRevoke and PostToken are implemented in legacy package
	"OAuth2_PostToken":  true,
}

var serviceOverrides = map[string]string{
	"AccountSignatureProviders_GetSealProviders":                 "Accounts",
	"AccountIdentityVerification_GetAccountIdentityVerification": "Accounts",
	"NotaryJournals_GetNotaryJournals":                           "Envelopes",
	"Views_PostEnvelopeRecipientSharedView":                      "Envelopes",
}

type override struct {
	Object string
	Field  string
	Type   string
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
	case "Envelopes_PostEnvelopes":
		return true
	case "Templates_PostTemplates":
		return true
	case "UserSignatures_PostUserSignatures":
		return true
	}
	return false
}

// GetFieldOverrides returns a map of all
// field level type overrides for the esign
// api generation. The returned map is
// map[<structID>]map[<FieldName>]<GoType>
//
// In the specification, DocuSign lists every field as a string.
// I generated much of this list with the following rules.
// - definition properties with **true** in the description are
//   assumed to be bools set to DSBool
// - fields containing base64 in the name are assumed to be []byte
// - fields ending in DateTime are *time.Time//
// I eyeballed the doc as best I could so please let me know of any additions
// or corrections.
func GetFieldOverrides() map[string]map[string]string {
	return map[string]map[string]string{
		"witness": {
			"offlineAttributes": "interface{}",
		},
		"notaryRecipient": {
			"offlineAttributes": "interface{}",
		},
		"newAccountDefinition": {
			"socialAccountInformation": "*SocialAccountInformation",
		},
		"userSocialIdResult": {
			"socialAccountInformation": "[]SocialAccountInformation",
		},
		"AccountWatermarks": {
			"imageBase64": "[]byte",
		},
		"accountBillingPlan": {
			"isDowngrade":   "DSBool",
			"canUpgrade":    "DSBool",
			"enableSupport": "DSBool",
		},
		"accountBillingPlanResponse": {
			"billingAddressIsCreditCardAddress": "DSBool",
		},
		"accountInformation": {
			"isDowngrade":           "DSBool",
			"allowTransactionRooms": "DSBool",
			"canUpgrade":            "DSBool",
		},
		"accountRoleSettings": {
			"allowAccountManagement":                            "DSBool",
			"allowApiAccess":                                    "DSBool",
			"allowApiAccessToAccount":                           "DSBool",
			"allowApiSendingOnBehalfOfOthers":                   "DSBool",
			"allowApiSequentialSigning":                         "DSBool",
			"allowBulkSending":                                  "DSBool",
			"allowDocuSignDesktopClient":                        "DSBool",
			"allowedAddressBookAccess":                          "DSBool",
			"allowedTemplateAccess":                             "DSBool",
			"allowedToBeEnvelopeTransferRecipient":              "DSBool",
			"allowEnvelopeSending":                              "DSBool",
			"allowESealRecipients":                              "DSBool",
			"allowPowerFormsAdminToAccessAllPowerFormEnvelopes": "DSBool",
			"allowSendersToSetRecipientEmailLanguage":           "DSBool",
			"allowSignerAttachments":                            "DSBool",
			"allowSupplementalDocuments":                        "DSBool",
			"allowTaggingInSendAndCorrect":                      "DSBool",
			"allowVaulting":                                     "DSBool",
			"allowWetSigningOverride":                           "DSBool",
			"canCreateWorkspaces":                               "DSBool",
			"disableDocumentUpload":                             "DSBool",
			"disableOtherActions":                               "DSBool",
			"enableApiRequestLogging":                           "DSBool",
			"enableRecipientViewingNotifications":               "DSBool",
			"enableSequentialSigningInterface":                  "DSBool",
			"enableTransactionPointIntegration":                 "DSBool",
			"receiveCompletedSelfSignedDocumentsAsEmailLinks":   "DSBool",
			"supplementalDocumentsMustAccept":                   "DSBool",
			"supplementalDocumentsMustRead":                     "DSBool",
			"supplementalDocumentsMustView":                     "DSBool",
			"useNewDocuSignExperienceInterface":                 "DSBool",
			"useNewSendingInterface":                            "DSBool",
		},
		"accountSignatureRequired": {
			"isRequired": "DSBool",
		},
		"addressInformationInput": {
			"receiveInResponse": "DSBool",
		},
		"agent": {
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "DSBool",
			"inheritEmailNotificationConfiguration": "DSBool",
			"phoneAuthentication":                   "DSBool",
			"requireIdLookup":                       "DSBool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"templateLocked":                        "DSBool",
			"templateRequired":                      "DSBool",
		},
		"apiRequestLog": {
			"createdDateTime": "*time.Time",
		},
		"attachment": {
			"data": "[]byte",
		},
		"brand": {
			"isSendingDefault":        "DSBool",
			"isSigningDefault":        "DSBool",
			"isOverridingCompanyName": "DSBool",
		},
		"billingPlan": {
			"enableSupport": "DSBool",
		},
		"billingPlanInformation": {
			"enableSupport": "DSBool",
		},
		"bulkEnvelope": {
			"submittedDateTime": "*time.Time",
		},
		"carbonCopy": {
			"agentCanEditEmail":                     "DSBool",
			"agentCanEditName":                      "DSBool",
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "DSBool",
			"inheritEmailNotificationConfiguration": "DSBool",
			"phoneAuthentication":                   "DSBool",
			"requireIdLookup":                       "DSBool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"templateLocked":                        "DSBool",
			"templateRequired":                      "DSBool",
		},
		"certifiedDelivery": {
			"agentCanEditEmail":                     "DSBool",
			"agentCanEditName":                      "DSBool",
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "DSBool",
			"inheritEmailNotificationConfiguration": "DSBool",
			"phoneAuthentication":                   "DSBool",
			"requireIdLookup":                       "DSBool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"templateLocked":                        "DSBool",
			"templateRequired":                      "DSBool",
		},
		"checkbox": {
			"anchorCaseSensitive":          "DSBool",
			"anchorIgnoreIfNotPresent":     "DSBool",
			"anchorMatchWholeWord":         "DSBool",
			"locked":                       "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"selected":                     "DSBool",
			"shared":                       "DSBool",
			"templateLocked":               "DSBool",
			"templateRequired":             "DSBool",
		},
		"chunkedUploadRequest": {
			"data": "[]byte",
		},
		"chunkedUploadResponse": {
			"expirationDateTime": "*time.Time",
		},
		"company": {
			"concealValueOnDocument": "DSBool",
			"disableAutoSize":        "DSBool",
			"locked":                 "DSBool",
			"required":               "TabRequired",
		},
		"connectCustomConfiguration": {
			"allUsers":                          "DSBool",
			"allowEnvelopePublish":              "DSBool",
			"enableLog":                         "DSBool",
			"includeCertificateOfCompletion":    "DSBool",
			"includeDocumentFields":             "DSBool",
			"includeDocuments":                  "DSBool",
			"includeEnvelopeVoidReason":         "DSBool",
			"includeSenderAccountasCustomField": "DSBool",
			"includeTimeZoneInformation":        "DSBool",
			"requiresAcknowledgement":           "DSBool",
			"signMessageWithX509Certificate":    "DSBool",
			"useSoapInterface":                  "DSBool",
		},
		"connectDebugLog": {
			"eventDateTime": "*time.Time",
		},
		"connectLog": {
			"created": "*time.Time",
		},
		"consumerDisclosure": {
			"allowCDWithdraw":                    "DSBool",
			"useConsumerDisclosureWithinAccount": "DSBool",
			"withdrawByEmail":                    "DSBool",
			"withdrawByMail":                     "DSBool",
			"withdrawByPhone":                    "DSBool",
		},
		"contact": {
			"shared": "DSBool",
		},
		"currencyFeatureSetPrice": {
			"envelopeFee": "DSBool",
			"fixedFee":    "DSBool",
			"seatFee":     "DSBool",
		},
		"customField": {
			"required": "TabRequired",
			"show":     "DSBool",
		},
		"customField_v2": {
			"required": "TabRequired",
			"show":     "DSBool",
		},
		"date": {
			"concealValueOnDocument":       "DSBool",
			"disableAutoSize":              "DSBool",
			"locked":                       "DSBool",
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"senderRequired":               "DSBool",
			"shared":                       "DSBool",
		},
		"diagnosticsSettingsInformation": {
			"apiRequestLogging": "DSBool",
		},
		"dobInformationInput": {
			"receiveInResponse": "DSBool",
		},
		"document": {
			"authoritativeCopy":       "DSBool",
			"documentBase64":          "[]byte",
			"encryptedWithKeyManager": "DSBool",
			"includeInDownload":       "DSBool",
			"templateLocked":          "DSBool",
			"templateRequired":        "DSBool",
			"transformPdfFields":      "DSBool",
		},
		"documentHtmlDefinition": {
			"displayPageNumber":         "int32",
			"displayOrder":              "int32",
			"removeEmptyTags":           "DSBool",
			"showMobileOptimizedToggle": "DSBool",
		},
		"editor": {
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"inheritEmailNotificationConfiguration": "DSBool",
			"phoneAuthentication":                   "DSBool",
			"requireIdLookup":                       "DSBool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"templateLocked":                        "DSBool",
			"templateRequired":                      "DSBool",
		},
		"email": {
			"concealValueOnDocument":       "DSBool",
			"disableAutoSize":              "DSBool",
			"locked":                       "DSBool",
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"senderRequired":               "DSBool",
			"shared":                       "DSBool",
		},
		"envelope": {
			"allowMarkup":                 "DSBool",
			"allowReassign":               "DSBool",
			"allowViewHistory":            "DSBool",
			"asynchronous":                "DSBool",
			"authoritativeCopy":           "DSBool",
			"signerCanSignOnMobile":       "DSBool",
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"declinedDateTime":            "*time.Time",
			"deletedDateTime":             "*time.Time",
			"deliveredDateTime":           "*time.Time",
			"enableWetSign":               "DSBool",
			"enforceSignerVisibility":     "DSBool",
			"envelopeIdStamping":          "DSBool",
			"initialSentDateTime":         "*time.Time",
			"is21CFRPart11":               "DSBool",
			"isSignatureProviderEnvelope": "DSBool",
			"lastModifiedDateTime":        "*time.Time",
			"messageLock":                 "DSBool",
			"notification":                "DSBool",
			"recipientsLock":              "DSBool",
			"sentDateTime":                "*time.Time",
			"statusChangedDateTime":       "*time.Time",
			"useDisclosure":               "DSBool",
			"voidedDateTime":              "*time.Time",
		},
		"envelopeDefinition": {
			"authoritativeCopy":           "DSBool",
			"allowMarkup":                 "DSBool",
			"allowReassign":               "DSBool",
			"allowRecipientRecursion":     "DSBool",
			"allowViewHistory":            "DSBool",
			"asynchronous":                "DSBool",
			"signerCanSignOnMobile":       "DSBool",
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"declinedDateTime":            "*time.Time",
			"deletedDateTime":             "*time.Time",
			"deliveredDateTime":           "*time.Time",
			"enableWetSign":               "DSBool",
			"enforceSignerVisibility":     "DSBool",
			"envelopeIdStamping":          "DSBool",
			"initialSentDateTime":         "*time.Time",
			"is21CFRPart11":               "DSBool",
			"isSignatureProviderEnvelope": "DSBool",
			"lastModifiedDateTime":        "*time.Time",
			"messageLock":                 "DSBool",
			"recipientsLock":              "DSBool",
			"sentDateTime":                "*time.Time",
			"statusChangedDateTime":       "*time.Time",
			"useDisclosure":               "DSBool",
			"voidedDateTime":              "*time.Time",
		},
		"envelopeDocument": {
			"authoritativeCopy": "DSBool",
			"includeInDownload": "DSBool",
		},
		"envelopeFormData": {
			"sentDateTime": "*time.Time",
		},
		"envelopeNotificationRequest": {
			"useAccountDefaults": "DSBool",
		},
		"envelopeSummary": {
			"statusDateTime": "*time.Time",
		},
		"envelopeTemplate": {
			"allowMarkup":                 "DSBool",
			"allowReassign":               "DSBool",
			"allowViewHistory":            "DSBool",
			"asynchronous":                "DSBool",
			"authoritativeCopy":           "DSBool",
			"signerCanSignOnMobile":       "DSBool",
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"declinedDateTime":            "*time.Time",
			"deletedDateTime":             "*time.Time",
			"deliveredDateTime":           "*time.Time",
			"enableWetSign":               "DSBool",
			"enforceSignerVisibility":     "DSBool",
			"envelopeIdStamping":          "DSBool",
			"initialSentDateTime":         "*time.Time",
			"is21CFRPart11":               "DSBool",
			"isSignatureProviderEnvelope": "DSBool",
			"lastModifiedDateTime":        "*time.Time",
			"messageLock":                 "DSBool",
			"recipientsLock":              "DSBool",
			"sentDateTime":                "*time.Time",
			"statusChangedDateTime":       "*time.Time",
			"useDisclosure":               "DSBool",
			"voidedDateTime":              "*time.Time",
		},
		"envelopeTemplateDefinition": {
			"created":      "*time.Time",
			"lastModified": "*time.Time",
			"shared":       "DSBool",
		},
		"envelopeTemplateResult": {
			"authoritativeCopy":           "DSBool",
			"allowMarkup":                 "DSBool",
			"allowReassign":               "DSBool",
			"allowViewHistory":            "DSBool",
			"asynchronous":                "DSBool",
			"signerCanSignOnMobile":       "DSBool",
			"completedDateTime":           "*time.Time",
			"created":                     "*time.Time",
			"createdDateTime":             "*time.Time",
			"declinedDateTime":            "*time.Time",
			"deletedDateTime":             "*time.Time",
			"deliveredDateTime":           "*time.Time",
			"enableWetSign":               "DSBool",
			"enforceSignerVisibility":     "DSBool",
			"envelopeIdStamping":          "DSBool",
			"initialSentDateTime":         "*time.Time",
			"is21CFRPart11":               "DSBool",
			"isSignatureProviderEnvelope": "DSBool",
			"lastModified":                "*time.Time",
			"lastModifiedDateTime":        "*time.Time",
			"messageLock":                 "DSBool",
			"recipientsLock":              "DSBool",
			"sentDateTime":                "*time.Time",
			"shared":                      "DSBool",
			"statusChangedDateTime":       "*time.Time",
			"useDisclosure":               "DSBool",
			"voidedDateTime":              "*time.Time",
		},
		"eventNotification": {
			"includeCertificateOfCompletion":    "DSBool",
			"includeCertificateWithSoap":        "DSBool",
			"includeDocumentFields":             "DSBool",
			"includeDocuments":                  "DSBool",
			"includeEnvelopeVoidReason":         "DSBool",
			"includeSenderAccountAsCustomField": "DSBool",
			"includeTimeZone":                   "DSBool",
			"loggingEnabled":                    "DSBool",
			"requireAcknowledgment":             "DSBool",
			"signMessageWithX509Cert":           "DSBool",
			"useSoapInterface":                  "DSBool",
		},
		"expirations": {
			"expireEnabled": "DSBool",
		},
		"featureSet": {
			"isActive":      "DSBool",
			"is21CFRPart11": "DSBool",
			"isEnabled":     "DSBool",
		},
		"filter": {
			"fromDateTime": "*time.Time",
			"isTemplate":   "DSBool",
			"toDateTime":   "*time.Time",
		},
		"folderItem": {
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"is21CFRPart11":               "DSBool",
			"isSignatureProviderEnvelope": "DSBool",
			"sentDateTime":                "*time.Time",
			"shared":                      "DSBool",
		},
		"folderItem_v2": {
			"completedDateTime":           "*time.Time",
			"createdDateTime":             "*time.Time",
			"expireDateTime":              "*time.Time",
			"is21CFRPart11":               "DSBool",
			"isSignatureProviderEnvelope": "DSBool",
			"lastModifiedDateTime":        "*time.Time",
			"sentDateTime":                "*time.Time",
		},
		"formulaTab": {
			"concealValueOnDocument":       "DSBool",
			"disableAutoSize":              "DSBool",
			"isPaymentAmount":              "DSBool",
			"locked":                       "DSBool",
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"senderRequired":               "DSBool",
			"shared":                       "DSBool",
		},
		"inPersonSigner": {
			"canSignOffline":                        "DSBool",
			"declinedDateTime":                      "*time.Time",
			"defaultRecipient":                      "DSBool",
			"deliveredDateTime":                     "*time.Time",
			"inheritEmailNotificationConfiguration": "DSBool",
			"phoneAuthentication":                   "DSBool",
			"requireIdLookup":                       "DSBool",
			"requireSignOnPaper":                    "DSBool",
			"sentDateTime":                          "*time.Time",
			"signInEachLocation":                    "DSBool",
			"signedDateTime":                        "*time.Time",
			"templateLocked":                        "DSBool",
			"templateRequired":                      "DSBool",
		},
		"initialHere": {
			"optional": "DSBool",
		},
		"intermediary": {
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "DSBool",
			"inheritEmailNotificationConfiguration": "DSBool",
			"phoneAuthentication":                   "DSBool",
			"requireIdLookup":                       "DSBool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"templateLocked":                        "DSBool",
			"templateRequired":                      "DSBool",
		},
		"jurisdiction": {
			"allowSystemCreatedSeal": "DSBool",
			"allowUserUploadedSeal":  "DSBool",
			"commissionIdInSeal":     "DSBool",
			"countyInSeal":           "DSBool",
			"enabled":                "DSBool",
			"notaryPublicInSeal":     "DSBool",
			"stateNameInSeal":        "DSBool",
		},
		"list": {
			"locked":                       "DSBool",
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"senderRequired":               "DSBool",
			"shared":                       "DSBool",
		},
		"listCustomField": {
			"required": "TabRequired",
			"show":     "DSBool",
		},
		"listItem": {
			"selected": "DSBool",
		},
		"lockInformation": {
			"lockedUntilDateTime": "*time.Time",
		},
		"loginAccount": {
			"isDefault": "DSBool",
		},
		"memberGroupSharedItem": {
			"shared": "DSBool",
		},
		"mergeField": {
			"allowSenderToEdit": "DSBool",
		},
		"newUser": {
			"createdDateTime": "*time.Time",
		},
		"notarize": {
			"locked":   "DSBool",
			"required": "TabRequired",
		},
		"notaryHost": {
			"declinedDateTime":                      "*time.Time",
			"deliveredDateTime":                     "*time.Time",
			"inheritEmailNotificationConfiguration": "DSBool",
			"phoneAuthentication":                   "DSBool",
			"requireIdLookup":                       "DSBool",
			"sentDateTime":                          "*time.Time",
			"signedDateTime":                        "*time.Time",
			"templateLocked":                        "DSBool",
			"templateRequired":                      "DSBool",
		},
		"note": {
			"shared": "DSBool",
		},
		"notification": {
			"useAccountDefaults": "DSBool",
		},
		"number": {
			"concealValueOnDocument":       "DSBool",
			"disableAutoSize":              "DSBool",
			"isPaymentAmount":              "DSBool",
			"locked":                       "DSBool",
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"senderRequired":               "DSBool",
			"shared":                       "DSBool",
		},
		"permissionProfile": {
			"modifiedDateTime": "*time.Time",
		},
		"powerForm": {
			"createdDateTime": "*time.Time",
			"isActive":        "DSBool",
		},
		"purchasedEnvelopesInformation": {
			"receiptData": "[]byte",
		},
		"radio": {
			"anchorCaseSensitive":      "DSBool",
			"anchorIgnoreIfNotPresent": "DSBool",
			"anchorMatchWholeWord":     "DSBool",
			"locked":                   "DSBool",
			"required":                 "TabRequired",
			"selected":                 "DSBool",
		},
		"radioGroup": {
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"shared":                       "DSBool",
		},
		"recipientAttachment": {
			"data": "[]byte",
		},
		"recipientFormData": {
			"declinedTime":  "*time.Time",
			"deliveredTime": "*time.Time",
			"sentTime":      "*time.Time",
			"signedTime":    "*time.Time",
		},
		"recipientPhoneAuthentication": {
			"recipMayProvideNumber": "DSBool",
		},
		"recipientSignatureProvider": {
			"sealDocumentsWithTabsOnly": "DSBool",
		},
		"referralInformation": {
			"enableSupport": "DSBool",
		},
		"reminders": {
			"reminderEnabled": "DSBool",
		},
		"senderEmailNotifications": {
			"changedSigner":        "DSBool",
			"deliveryFailed":       "DSBool",
			"envelopeComplete":     "DSBool",
			"offlineSigningFailed": "DSBool",
			"recipientViewed":      "DSBool",
			"withdrawnConsent":     "DSBool",
		},
		"serviceInformation": {
			"buildBranchDeployedDateTime": "*time.Time",
		},
		"settingsMetadata": {
			"is21CFRPart11": "DSBool",
		},
		"sharedItem": {
			"shared": "DSBool",
		},
		"signHere": {
			"optional": "DSBool",
		},
		"signatureType": {
			"isDefault": "DSBool",
		},
		"signer": {
			"agentCanEditEmail":                     "DSBool",
			"agentCanEditName":                      "DSBool",
			"canSignOffline":                        "DSBool",
			"declinedDateTime":                      "*time.Time",
			"defaultRecipient":                      "DSBool",
			"deliveredDateTime":                     "*time.Time",
			"excludedDocuments":                     "DSBool",
			"inheritEmailNotificationConfiguration": "DSBool",
			"isBulkRecipient":                       "DSBool",
			"phoneAuthentication":                   "DSBool",
			"requireIdLookup":                       "DSBool",
			"requireSignOnPaper":                    "DSBool",
			"sentDateTime":                          "*time.Time",
			"signInEachLocation":                    "DSBool",
			"signedDateTime":                        "*time.Time",
			"templateLocked":                        "DSBool",
			"templateRequired":                      "DSBool",
		},
		"signerAttachment": {
			"optional": "DSBool",
		},
		"signerEmailNotifications": {
			"agentNotification":             "DSBool",
			"carbonCopyNotification":        "DSBool",
			"certifiedDeliveryNotification": "DSBool",
			"documentMarkupActivation":      "DSBool",
			"envelopeActivation":            "DSBool",
			"envelopeComplete":              "DSBool",
			"envelopeCorrected":             "DSBool",
			"envelopeDeclined":              "DSBool",
			"envelopeVoided":                "DSBool",
			"offlineSigningFailed":          "DSBool",
			"purgeDocuments":                "DSBool",
			"reassignedSigner":              "DSBool",
		},
		"signingGroup": {
			"created":  "*time.Time",
			"modified": "*time.Time",
		},
		"ssn": {
			"concealValueOnDocument":       "DSBool",
			"disableAutoSize":              "DSBool",
			"locked":                       "DSBool",
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"senderRequired":               "DSBool",
			"shared":                       "DSBool",
		},
		"ssn4InformationInput": {
			"receiveInResponse": "DSBool",
		},
		"tabAccountSettings": {
			"allowTabOrder":                       "DSBool",
			"approveDeclineTabsEnabled":           "DSBool",
			"calculatedFieldsEnabled":             "DSBool",
			"checkboxTabsEnabled":                 "DSBool",
			"dataFieldRegexEnabled":               "DSBool",
			"dataFieldSizeEnabled":                "DSBool",
			"firstLastEmailTabsEnabled":           "DSBool",
			"listTabsEnabled":                     "DSBool",
			"noteTabsEnabled":                     "DSBool",
			"radioTabsEnabled":                    "DSBool",
			"savingCustomTabsEnabled":             "DSBool",
			"senderToChangeTabAssignmentsEnabled": "DSBool",
			"sharedCustomTabsEnabled":             "DSBool",
			"tabDataLabelEnabled":                 "DSBool",
			"tabLocationEnabled":                  "DSBool",
			"tabLockingEnabled":                   "DSBool",
			"tabScaleEnabled":                     "DSBool",
			"tabTextFormattingEnabled":            "DSBool",
			"textTabsEnabled":                     "DSBool",
		},
		"tabMetadata": {
			"anchorCaseSensitive":      "DSBool",
			"anchorIgnoreIfNotPresent": "DSBool",
			"anchorMatchWholeWord":     "DSBool",
			"bold":                     "DSBool",
			"concealValueOnDocument":   "DSBool",
			"disableAutoSize":          "DSBool",
			"editable":                 "DSBool",
			"includedInEmail":          "DSBool",
			"italic":                   "DSBool",
			"lastModified":             "*time.Time",
			"locked":                   "DSBool",
			"required":                 "TabRequired",
			"shared":                   "DSBool",
			"underline":                "DSBool",
		},
		"TabPosition": {
			"anchorCaseSensitive":      "DSBool",
			"anchorIgnoreIfNotPresent": "DSBool",
			"anchorMatchWholeWord":     "DSBool",
			"templateLocked":           "DSBool",
			"templateRequired":         "DSBool",
		},
		"TabStyle": {
			"bold":      "DSBool",
			"italic":    "DSBool",
			"underline": "DSBool",
		},
		"templateNotificationRequest": {
			"useAccountDefaults": "DSBool",
		},
		"templateRole": {
			"defaultRecipient": "DSBool",
		},
		"templateSharedItem": {
			"shared": "DSBool",
		},
		"text": {
			"concealValueOnDocument":       "DSBool",
			"disableAutoSize":              "DSBool",
			"isPaymentAmount":              "DSBool",
			"locked":                       "DSBool",
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"senderRequired":               "DSBool",
			"shared":                       "DSBool",
		},
		"textCustomField": {
			"required": "TabRequired",
			"show":     "DSBool",
		},
		"title": {
			"concealValueOnDocument": "DSBool",
			"disableAutoSize":        "DSBool",
			"locked":                 "DSBool",
			"required":               "TabRequired",
		},
		"userAccountManagementGranularInformation": {
			"canManageAccountSecuritySettings": "DSBool",
			"canManageAccountSettings":         "DSBool",
			"canManageAdmins":                  "DSBool",
			"canManageGroups":                  "DSBool",
			"canManageReporting":               "DSBool",
			"canManageSharing":                 "DSBool",
			"canManageSigningGroups":           "DSBool",
			"canManageUsers":                   "DSBool",
		},
		"usageHistory": {
			"lastSentDateTime":   "*time.Time",
			"lastSignedDateTime": "*time.Time",
		},
		"userInformation": {
			"createdDateTime":              "*time.Time",
			"isActive":                     "DSBool",
			"sendActivationOnInvalidLogin": "DSBool",
		},
		"userProfile": {
			"displayOrganizationInfo": "DSBool",
			"displayPersonalInfo":     "DSBool",
			"displayProfile":          "DSBool",
			"displayUsageHistory":     "DSBool",
		},
		"userSharedItem": {
			"shared": "DSBool",
		},
		"userSignature": {
			"adoptedDateTime": "*time.Time",
			"createdDateTime": "*time.Time",
		},
		"view": {
			"required":     "TabRequired",
			"requiredRead": "DSBool",
		},
		"watermark": {
			"enabled": "DSBool",
		},
		"workspace": {
			"created":      "*time.Time",
			"lastModified": "*time.Time",
		},
		"workspaceItem": {
			"created":      "*time.Time",
			"isPublic":     "DSBool",
			"lastModified": "*time.Time",
		},
		"workspaceUser": {
			"activeSince":  "*time.Time",
			"created":      "*time.Time",
			"lastModified": "*time.Time",
		},
		"workspaceUserAuthorization": {
			"canDelete":   "DSBool",
			"canMove":     "DSBool",
			"canTransact": "DSBool",
			"canView":     "DSBool",
			"created":     "*time.Time",
			"modified":    "*time.Time",
		},
		"zip": {
			"concealValueOnDocument":       "DSBool",
			"disableAutoSize":              "DSBool",
			"locked":                       "DSBool",
			"requireAll":                   "DSBool",
			"requireInitialOnSharedChange": "DSBool",
			"required":                     "TabRequired",
			"senderRequired":               "DSBool",
			"shared":                       "DSBool",
		},
	}
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
		"BulkEnvelopes_GetBulkEnvelopesBatchId": {
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
func TabDefs(version string, defMap map[string]Definition, overrides map[string]map[string]string) []Definition {
	switch version {
	case "v2":
		return TabDefsV2(defMap, overrides)
	case "v2.1":
		return TabDefsV21(defMap, overrides)
	}
	return make([]Definition, 0, 0)
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
		"commentThread",
		"smartSection",
		"polyLineOverlay",
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

// CustomCode is lines of code to append to model.go
const CustomCode = `// GetTabValues returns a NameValue list of all entry tabs
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
			if item.Selected {
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
			if rb.Selected {
				results = append(results, NameValue{Name: rg.GroupName, Value: rb.Value})
			}
		}
	}
	return results
}`
