// Copyright 2017 James Cote and Liberty Fund, Inc.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Code generated by go-swagger; DO NOT EDIT.

// Package users implements the DocuSign SDK
// category Users.
// 
// Use the Users category to manage the users in your accounts.
// 
// You can:
// 
// * Set custom user settings.
// * Manage a users profile.
// * Add delete users.
// * Add and delete the intials and signature images for a user.
// Api documentation may be found at:
// https://docs.docusign.com/esign/restapi/Users
package users

import (
    "fmt"
    "net/url"
    "strings"
    
    "golang.org/x/net/context"
    
    "mystuff/esign"
    "mystuff/esign/model"
)

// Service generates DocuSign Users Category API calls
type Service struct {
    credential esign.Credential 
}

// New initializes a users service using cred to authorize calls.
func New(cred esign.Credential) *Service {
    return &Service{credential: cred}
}

// DeleteContactWithID replaces a particular contact associated with an account for the DocuSign service.
// SDK Method Users::deleteContactWithId
// https://docs.docusign.com/esign/restapi/Users/Contacts/delete
func (s *Service) DeleteContactWithID(contactID string) *DeleteContactWithIDCall {
    return &DeleteContactWithIDCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "DELETE",
            Path: "contacts/{contactId}",
            PathParameters: map[string]string{ 
                "{contactId}": contactID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// DeleteContactWithIDCall implements DocuSign API SDK Users::deleteContactWithId
type DeleteContactWithIDCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *DeleteContactWithIDCall) Do(ctx context.Context)  (*model.ContactUpdateResponse, error) {
    var res *model.ContactUpdateResponse
    return res, op.Call.Do(ctx, &res)
}

// DeleteContacts delete contacts associated with an account for the DocuSign service.
// SDK Method Users::deleteContacts
// https://docs.docusign.com/esign/restapi/Users/Contacts/deleteList
func (s *Service) DeleteContacts(contactModRequest *model.ContactModRequest) *DeleteContactsCall {
    return &DeleteContactsCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "DELETE",
            Path: "contacts",
            Payload: contactModRequest,
            QueryOpts: make(url.Values),
        },
    }
}

// DeleteContactsCall implements DocuSign API SDK Users::deleteContacts
type DeleteContactsCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *DeleteContactsCall) Do(ctx context.Context)  (*model.ContactUpdateResponse, error) {
    var res *model.ContactUpdateResponse
    return res, op.Call.Do(ctx, &res)
}

// GetContactByID gets a particular contact associated with the user's account.
// SDK Method Users::getContactById
// https://docs.docusign.com/esign/restapi/Users/Contacts/get
func (s *Service) GetContactByID(contactID string) *GetContactByIDCall {
    return &GetContactByIDCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "contacts/{contactId}",
            PathParameters: map[string]string{ 
                "{contactId}": contactID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// GetContactByIDCall implements DocuSign API SDK Users::getContactById
type GetContactByIDCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *GetContactByIDCall) Do(ctx context.Context)  (*model.ContactGetResponse, error) {
    var res *model.ContactGetResponse
    return res, op.Call.Do(ctx, &res)
}

// CloudProvider set the call query parameter cloud_provider
func (op *GetContactByIDCall) CloudProvider(val string) *GetContactByIDCall {
    op.QueryOpts.Set("cloud_provider", val)
    return op
}

// PostContacts imports multiple new contacts into the contacts collection from CSV, JSON, or XML (based on content type).
// SDK Method Users::postContacts
// https://docs.docusign.com/esign/restapi/Users/Contacts/create
func (s *Service) PostContacts(contactModRequest *model.ContactModRequest) *PostContactsCall {
    return &PostContactsCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "POST",
            Path: "contacts",
            Payload: contactModRequest,
            QueryOpts: make(url.Values),
        },
    }
}

// PostContactsCall implements DocuSign API SDK Users::postContacts
type PostContactsCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *PostContactsCall) Do(ctx context.Context)  (*model.ContactUpdateResponse, error) {
    var res *model.ContactUpdateResponse
    return res, op.Call.Do(ctx, &res)
}

// PutContacts replaces contacts associated with an account for the DocuSign service.
// SDK Method Users::putContacts
// https://docs.docusign.com/esign/restapi/Users/Contacts/update
func (s *Service) PutContacts(contactModRequest *model.ContactModRequest) *PutContactsCall {
    return &PutContactsCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "contacts",
            Payload: contactModRequest,
            QueryOpts: make(url.Values),
        },
    }
}

// PutContactsCall implements DocuSign API SDK Users::putContacts
type PutContactsCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *PutContactsCall) Do(ctx context.Context)  (*model.ContactUpdateResponse, error) {
    var res *model.ContactUpdateResponse
    return res, op.Call.Do(ctx, &res)
}

// DeleteCustomSettings deletes custom user settings for a specified user.
// SDK Method Users::deleteCustomSettings
// https://docs.docusign.com/esign/restapi/Users/UserCustomSettings/delete
func (s *Service) DeleteCustomSettings(userID string, userCustomSettings *model.CustomSettingsInformation) *DeleteCustomSettingsCall {
    return &DeleteCustomSettingsCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "DELETE",
            Path: "users/{userId}/custom_settings",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            Payload: userCustomSettings,
            QueryOpts: make(url.Values),
        },
    }
}

// DeleteCustomSettingsCall implements DocuSign API SDK Users::deleteCustomSettings
type DeleteCustomSettingsCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *DeleteCustomSettingsCall) Do(ctx context.Context)  (*model.CustomSettingsInformation, error) {
    var res *model.CustomSettingsInformation
    return res, op.Call.Do(ctx, &res)
}

// ListCustomSettings retrieves the custom user settings for a specified user.
// SDK Method Users::listCustomSettings
// https://docs.docusign.com/esign/restapi/Users/UserCustomSettings/list
func (s *Service) ListCustomSettings(userID string) *ListCustomSettingsCall {
    return &ListCustomSettingsCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users/{userId}/custom_settings",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// ListCustomSettingsCall implements DocuSign API SDK Users::listCustomSettings
type ListCustomSettingsCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *ListCustomSettingsCall) Do(ctx context.Context)  (*model.CustomSettingsInformation, error) {
    var res *model.CustomSettingsInformation
    return res, op.Call.Do(ctx, &res)
}

// UpdateCustomSettings adds or updates custom user settings for the specified user.
// SDK Method Users::updateCustomSettings
// https://docs.docusign.com/esign/restapi/Users/UserCustomSettings/update
func (s *Service) UpdateCustomSettings(userID string, userCustomSettings *model.CustomSettingsInformation) *UpdateCustomSettingsCall {
    return &UpdateCustomSettingsCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users/{userId}/custom_settings",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            Payload: userCustomSettings,
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateCustomSettingsCall implements DocuSign API SDK Users::updateCustomSettings
type UpdateCustomSettingsCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateCustomSettingsCall) Do(ctx context.Context)  (*model.CustomSettingsInformation, error) {
    var res *model.CustomSettingsInformation
    return res, op.Call.Do(ctx, &res)
}

// DeleteProfileImage deletes the user profile image for the specified user.
// SDK Method Users::deleteProfileImage
// https://docs.docusign.com/esign/restapi/Users/Users/deleteProfileImage
func (s *Service) DeleteProfileImage(userID string) *DeleteProfileImageCall {
    return &DeleteProfileImageCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "DELETE",
            Path: "users/{userId}/profile/image",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// DeleteProfileImageCall implements DocuSign API SDK Users::deleteProfileImage
type DeleteProfileImageCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *DeleteProfileImageCall) Do(ctx context.Context)  error {
    
    return op.Call.Do(ctx, nil)
}

// GetProfileImage retrieves the user profile image for the specified user.
// SDK Method Users::getProfileImage
// https://docs.docusign.com/esign/restapi/Users/Users/getProfileImage
func (s *Service) GetProfileImage(userID string) *GetProfileImageCall {
    return &GetProfileImageCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users/{userId}/profile/image",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// GetProfileImageCall implements DocuSign API SDK Users::getProfileImage
type GetProfileImageCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *GetProfileImageCall) Do(ctx context.Context)  (*esign.File, error) {
    var res *esign.File
    return res, op.Call.Do(ctx, &res)
}

// Encoding set the call query parameter encoding
func (op *GetProfileImageCall) Encoding(val string) *GetProfileImageCall {
    op.QueryOpts.Set("encoding", val)
    return op
}

// UpdateProfileImage updates the user profile image for a specified user.
// SDK Method Users::updateProfileImage
// https://docs.docusign.com/esign/restapi/Users/Users/updateProfileImage
func (s *Service) UpdateProfileImage(userID string) *UpdateProfileImageCall {
    return &UpdateProfileImageCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users/{userId}/profile/image",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateProfileImageCall implements DocuSign API SDK Users::updateProfileImage
type UpdateProfileImageCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateProfileImageCall) Do(ctx context.Context)  error {
    
    return op.Call.Do(ctx, nil)
}

// GetProfile retrieves the user profile for a specified user.
// SDK Method Users::getProfile
// https://docs.docusign.com/esign/restapi/Users/UserProfiles/get
func (s *Service) GetProfile(userID string) *GetProfileCall {
    return &GetProfileCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users/{userId}/profile",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// GetProfileCall implements DocuSign API SDK Users::getProfile
type GetProfileCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *GetProfileCall) Do(ctx context.Context)  (*model.UserProfile, error) {
    var res *model.UserProfile
    return res, op.Call.Do(ctx, &res)
}

// UpdateProfile updates the user profile information for the specified user.
// SDK Method Users::updateProfile
// https://docs.docusign.com/esign/restapi/Users/UserProfiles/update
func (s *Service) UpdateProfile(userID string, userProfiles *model.UserProfile) *UpdateProfileCall {
    return &UpdateProfileCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users/{userId}/profile",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            Payload: userProfiles,
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateProfileCall implements DocuSign API SDK Users::updateProfile
type UpdateProfileCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateProfileCall) Do(ctx context.Context)  error {
    
    return op.Call.Do(ctx, nil)
}

// GetSettings gets the user account settings for a specified user.
// SDK Method Users::getSettings
// https://docs.docusign.com/esign/restapi/Users/Users/getSettings
func (s *Service) GetSettings(userID string) *GetSettingsCall {
    return &GetSettingsCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users/{userId}/settings",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// GetSettingsCall implements DocuSign API SDK Users::getSettings
type GetSettingsCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *GetSettingsCall) Do(ctx context.Context)  (*model.UserSettingsInformation, error) {
    var res *model.UserSettingsInformation
    return res, op.Call.Do(ctx, &res)
}

// UpdateSettings updates the user account settings for a specified user.
// SDK Method Users::updateSettings
// https://docs.docusign.com/esign/restapi/Users/Users/updateSettings
func (s *Service) UpdateSettings(userID string, userSettingsInformation *model.UserSettingsInformation) *UpdateSettingsCall {
    return &UpdateSettingsCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users/{userId}/settings",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            Payload: userSettingsInformation,
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateSettingsCall implements DocuSign API SDK Users::updateSettings
type UpdateSettingsCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateSettingsCall) Do(ctx context.Context)  error {
    
    return op.Call.Do(ctx, nil)
}

// DeleteSignature removes removes signature information for the specified user.
// SDK Method Users::deleteSignature
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/delete
func (s *Service) DeleteSignature(signatureID string, userID string) *DeleteSignatureCall {
    return &DeleteSignatureCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "DELETE",
            Path: "users/{userId}/signatures/{signatureId}",
            PathParameters: map[string]string{ 
                "{signatureId}": signatureID,
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// DeleteSignatureCall implements DocuSign API SDK Users::deleteSignature
type DeleteSignatureCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *DeleteSignatureCall) Do(ctx context.Context)  error {
    
    return op.Call.Do(ctx, nil)
}

// DeleteSignatureImage deletes the user initials image or the  user signature image for the specified user.
// SDK Method Users::deleteSignatureImage
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/deleteImage
func (s *Service) DeleteSignatureImage(imageType string, signatureID string, userID string) *DeleteSignatureImageCall {
    return &DeleteSignatureImageCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "DELETE",
            Path: "users/{userId}/signatures/{signatureId}/{imageType}",
            PathParameters: map[string]string{ 
                "{imageType}": imageType,
                "{signatureId}": signatureID,
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// DeleteSignatureImageCall implements DocuSign API SDK Users::deleteSignatureImage
type DeleteSignatureImageCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *DeleteSignatureImageCall) Do(ctx context.Context)  (*model.UserSignature, error) {
    var res *model.UserSignature
    return res, op.Call.Do(ctx, &res)
}

// GetSignature gets the user signature information for the specified user.
// SDK Method Users::getSignature
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/get
func (s *Service) GetSignature(signatureID string, userID string) *GetSignatureCall {
    return &GetSignatureCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users/{userId}/signatures/{signatureId}",
            PathParameters: map[string]string{ 
                "{signatureId}": signatureID,
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// GetSignatureCall implements DocuSign API SDK Users::getSignature
type GetSignatureCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *GetSignatureCall) Do(ctx context.Context)  (*model.UserSignature, error) {
    var res *model.UserSignature
    return res, op.Call.Do(ctx, &res)
}

// GetSignatureImage retrieves the user initials image or the  user signature image for the specified user.
// SDK Method Users::getSignatureImage
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/getImage
func (s *Service) GetSignatureImage(imageType string, signatureID string, userID string) *GetSignatureImageCall {
    return &GetSignatureImageCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users/{userId}/signatures/{signatureId}/{imageType}",
            PathParameters: map[string]string{ 
                "{imageType}": imageType,
                "{signatureId}": signatureID,
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// GetSignatureImageCall implements DocuSign API SDK Users::getSignatureImage
type GetSignatureImageCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *GetSignatureImageCall) Do(ctx context.Context)  (*esign.File, error) {
    var res *esign.File
    return res, op.Call.Do(ctx, &res)
}

// IncludeChrome set the call query parameter include_chrome
func (op *GetSignatureImageCall) IncludeChrome() *GetSignatureImageCall {
    op.QueryOpts.Set("include_chrome", "true")
    return op
}

// ListSignatures retrieves a list of user signature definitions for a specified user.
// SDK Method Users::listSignatures
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/list
func (s *Service) ListSignatures(userID string) *ListSignaturesCall {
    return &ListSignaturesCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users/{userId}/signatures",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// ListSignaturesCall implements DocuSign API SDK Users::listSignatures
type ListSignaturesCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *ListSignaturesCall) Do(ctx context.Context)  (*model.UserSignaturesInformation, error) {
    var res *model.UserSignaturesInformation
    return res, op.Call.Do(ctx, &res)
}

// StampType set the call query parameter stamp_type
func (op *ListSignaturesCall) StampType(val string) *ListSignaturesCall {
    op.QueryOpts.Set("stamp_type", val)
    return op
}

// CreateSignatures adds user Signature and initials images to a Signature.
// SDK Method Users::createSignatures
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/create
func (s *Service) CreateSignatures(userID string, userSignaturesInformation *model.UserSignaturesInformation) *CreateSignaturesCall {
    return &CreateSignaturesCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "POST",
            Path: "users/{userId}/signatures",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            Payload: userSignaturesInformation,
            QueryOpts: make(url.Values),
        },
    }
}

// CreateSignaturesCall implements DocuSign API SDK Users::createSignatures
type CreateSignaturesCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *CreateSignaturesCall) Do(ctx context.Context)  (*model.UserSignaturesInformation, error) {
    var res *model.UserSignaturesInformation
    return res, op.Call.Do(ctx, &res)
}

// UpdateSignatures adds/updates a user signature.
// SDK Method Users::updateSignatures
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/updateList
func (s *Service) UpdateSignatures(userID string, userSignaturesInformation *model.UserSignaturesInformation) *UpdateSignaturesCall {
    return &UpdateSignaturesCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users/{userId}/signatures",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            Payload: userSignaturesInformation,
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateSignaturesCall implements DocuSign API SDK Users::updateSignatures
type UpdateSignaturesCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateSignaturesCall) Do(ctx context.Context)  (*model.UserSignaturesInformation, error) {
    var res *model.UserSignaturesInformation
    return res, op.Call.Do(ctx, &res)
}

// UpdateSignature updates the user signature for a specified user.
// SDK Method Users::updateSignature
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/update
func (s *Service) UpdateSignature(signatureID string, userID string, userSignatureDefinition *model.UserSignatureDefinition) *UpdateSignatureCall {
    return &UpdateSignatureCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users/{userId}/signatures/{signatureId}",
            PathParameters: map[string]string{ 
                "{signatureId}": signatureID,
                "{userId}": userID,
            },
            Payload: userSignatureDefinition,
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateSignatureCall implements DocuSign API SDK Users::updateSignature
type UpdateSignatureCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateSignatureCall) Do(ctx context.Context)  (*model.UserSignature, error) {
    var res *model.UserSignature
    return res, op.Call.Do(ctx, &res)
}

// CloseExistingSignature when set to **true**, closes the current signature.
func (op *UpdateSignatureCall) CloseExistingSignature() *UpdateSignatureCall {
    op.QueryOpts.Set("close_existing_signature", "true")
    return op
}

// UpdateSignatureImage updates the user signature image or user initials image for the specified user.
// SDK Method Users::updateSignatureImage
// https://docs.docusign.com/esign/restapi/Users/UserSignatures/updateImage
func (s *Service) UpdateSignatureImage(imageType string, signatureID string, userID string) *UpdateSignatureImageCall {
    return &UpdateSignatureImageCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users/{userId}/signatures/{signatureId}/{imageType}",
            PathParameters: map[string]string{ 
                "{imageType}": imageType,
                "{signatureId}": signatureID,
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateSignatureImageCall implements DocuSign API SDK Users::updateSignatureImage
type UpdateSignatureImageCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateSignatureImageCall) Do(ctx context.Context)  (*model.UserSignature, error) {
    var res *model.UserSignature
    return res, op.Call.Do(ctx, &res)
}

// GetInformation gets the user information for a specified user.
// SDK Method Users::getInformation
// https://docs.docusign.com/esign/restapi/Users/Users/get
func (s *Service) GetInformation(userID string) *GetInformationCall {
    return &GetInformationCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users/{userId}",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            QueryOpts: make(url.Values),
        },
    }
}

// GetInformationCall implements DocuSign API SDK Users::getInformation
type GetInformationCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *GetInformationCall) Do(ctx context.Context)  (*model.UserInformation, error) {
    var res *model.UserInformation
    return res, op.Call.Do(ctx, &res)
}

// AdditionalInfo when set to **true**, the full list of user information is returned for each user in the account.
func (op *GetInformationCall) AdditionalInfo() *GetInformationCall {
    op.QueryOpts.Set("additional_info", "true")
    return op
}

// Email set the call query parameter email
func (op *GetInformationCall) Email(val string) *GetInformationCall {
    op.QueryOpts.Set("email", val)
    return op
}

// UpdateUser updates the specified user information.
// SDK Method Users::updateUser
// https://docs.docusign.com/esign/restapi/Users/Users/update
func (s *Service) UpdateUser(userID string, users *model.UserInformation) *UpdateUserCall {
    return &UpdateUserCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users/{userId}",
            PathParameters: map[string]string{ 
                "{userId}": userID,
            },
            Payload: users,
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateUserCall implements DocuSign API SDK Users::updateUser
type UpdateUserCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateUserCall) Do(ctx context.Context)  (*model.UserInformation, error) {
    var res *model.UserInformation
    return res, op.Call.Do(ctx, &res)
}

// Delete removes users account privileges.
// SDK Method Users::delete
// https://docs.docusign.com/esign/restapi/Users/Users/delete
func (s *Service) Delete(userInfoList *model.UserInfoList) *DeleteCall {
    return &DeleteCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "DELETE",
            Path: "users",
            Payload: userInfoList,
            QueryOpts: make(url.Values),
        },
    }
}

// DeleteCall implements DocuSign API SDK Users::delete
type DeleteCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *DeleteCall) Do(ctx context.Context)  (*model.UsersResponse, error) {
    var res *model.UsersResponse
    return res, op.Call.Do(ctx, &res)
}

// List retrieves the list of users for the specified account.
// SDK Method Users::list
// https://docs.docusign.com/esign/restapi/Users/Users/list
func (s *Service) List() *ListCall {
    return &ListCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "GET",
            Path: "users",
            QueryOpts: make(url.Values),
        },
    }
}

// ListCall implements DocuSign API SDK Users::list
type ListCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *ListCall) Do(ctx context.Context)  (*model.UserInformationList, error) {
    var res *model.UserInformationList
    return res, op.Call.Do(ctx, &res)
}

// AdditionalInfo when set to **true**, the full list of user information is returned for each user in the account.
func (op *ListCall) AdditionalInfo() *ListCall {
    op.QueryOpts.Set("additional_info", "true")
    return op
}

// Count number of records to return. The number must be greater than 0 and less than or equal to 100.
func (op *ListCall) Count(val int) *ListCall {
    op.QueryOpts.Set("count", fmt.Sprintf("%d", val ))
    return op
}

// Email set the call query parameter email
func (op *ListCall) Email(val string) *ListCall {
    op.QueryOpts.Set("email", val)
    return op
}

// EmailSubstring filters the returned user records by the email address or a sub-string of email address.
func (op *ListCall) EmailSubstring(val string) *ListCall {
    op.QueryOpts.Set("email_substring", val)
    return op
}

// GroupID filters user records returned by one or more group Id's.
func (op *ListCall) GroupID(val string) *ListCall {
    op.QueryOpts.Set("group_id", val)
    return op
}

// IncludeUsersettingsForCsv set the call query parameter include_usersettings_for_csv
func (op *ListCall) IncludeUsersettingsForCsv() *ListCall {
    op.QueryOpts.Set("include_usersettings_for_csv", "true")
    return op
}

// LoginStatus set the call query parameter login_status
func (op *ListCall) LoginStatus(val string) *ListCall {
    op.QueryOpts.Set("login_status", val)
    return op
}

// NotGroupID set the call query parameter not_group_id
func (op *ListCall) NotGroupID(val string) *ListCall {
    op.QueryOpts.Set("not_group_id", val)
    return op
}

// StartPosition starting value for the list.
func (op *ListCall) StartPosition(val int) *ListCall {
    op.QueryOpts.Set("start_position", fmt.Sprintf("%d", val ))
    return op
}

// Status filters the results by user status.
// You can specify a comma-separated
// list of the following statuses:
// 
// * ActivationRequired 
// * ActivationSent 
// * Active
// * Closed 
// * Disabled
func (op *ListCall) Status(val ...string) *ListCall {
    op.QueryOpts.Set("status", strings.Join(val,","))
    return op
}

// UserNameSubstring filters the user records returned by the user name or a sub-string of user name.
func (op *ListCall) UserNameSubstring(val string) *ListCall {
    op.QueryOpts.Set("user_name_substring", val)
    return op
}

// Create adds news user to the specified account.
// SDK Method Users::create
// https://docs.docusign.com/esign/restapi/Users/Users/create
func (s *Service) Create(newUsersDefinition *model.NewUsersDefinition) *CreateCall {
    return &CreateCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "POST",
            Path: "users",
            Payload: newUsersDefinition,
            QueryOpts: make(url.Values),
        },
    }
}

// CreateCall implements DocuSign API SDK Users::create
type CreateCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *CreateCall) Do(ctx context.Context)  (*model.NewUsersSummary, error) {
    var res *model.NewUsersSummary
    return res, op.Call.Do(ctx, &res)
}

// UpdateUsers change one or more user in the specified account.
// SDK Method Users::updateUsers
// https://docs.docusign.com/esign/restapi/Users/Users/updateList
func (s *Service) UpdateUsers(userInformationList *model.UserInformationList) *UpdateUsersCall {
    return &UpdateUsersCall{
        &esign.Call{
            Credential: s.credential,
    		Method:  "PUT",
            Path: "users",
            Payload: userInformationList,
            QueryOpts: make(url.Values),
        },
    }
}

// UpdateUsersCall implements DocuSign API SDK Users::updateUsers
type UpdateUsersCall struct {
    *esign.Call
}

// Do executes the call.  A nil context will return error.
func (op *UpdateUsersCall) Do(ctx context.Context)  (*model.UserInformationList, error) {
    var res *model.UserInformationList
    return res, op.Call.Do(ctx, &res)
}
