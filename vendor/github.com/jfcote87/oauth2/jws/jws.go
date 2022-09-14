// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package jws provides a partial implementation
// of JSON Web Signature encoding and decoding.  It includes
// support for HS256, HS384, HS512, RS256, RS384, and RS512
// algorithms, although developers may extend this package
// by creating new Signer interfaces.
//
// See RFC 7515.
//
package jws // import "github.com/jfcote87/oauth2/jws"

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "crypto/sha256" // For HSXXX and RSXXX signer and verifier
	_ "crypto/sha512"
)

// ClaimSet contains information about the JWT signature including the
// permissions being requested (scopes), the target of the token, the issuer,
// the time the token was issued, and the lifetime of the token.
// see https://tools.ietf.org/html/rfc7519
type ClaimSet struct {
	Issuer    string // iss: client_id of the application making the access token request
	Audience  string // aud: descriptor of the intended target of the assertion (Optional).
	ExpiresAt int64  // exp: the expiration time of the assertion (seconds since Unix epoch)
	IssuedAt  int64  // iat: the time the assertion was issued (seconds since Unix epoch)
	NotBefore int64  // nbf: the time before which the JWT MUST NOT be accepted for processing (Optional)
	ID        string // jti: The "jti" (JWT ID) claim provides a unique identifier for the JWT (Optional)
	Subject   string // sub: Email/UserID for which the application is requesting delegated access (Optional).

	// See https://tools.ietf.org/html/rfc7519#section-4.3
	// This array is marshalled using custom code (see (c *ClaimSet) MarshalJSON()).
	PrivateClaims map[string]interface{}
}

// MarshalJSON flattens json output of PrivateClaims
func (c *ClaimSet) MarshalJSON() ([]byte, error) {
	pc := make(map[string]interface{})
	keys := []string{"iss", "aud", "jti", "sub"}
	for i, v := range []string{c.Issuer, c.Audience, c.ID, c.Subject} {
		if v > "" {
			pc[keys[i]] = v
		}
	}
	keys = []string{"exp", "iat", "nbf"}
	for i, v := range []int64{c.ExpiresAt, c.IssuedAt, c.NotBefore} {
		if v > 0 {
			pc[keys[i]] = v
		}
	}
	for k, v := range c.PrivateClaims {
		pc[k] = v
	}
	return json.Marshal(pc)
}

func (c *ClaimSet) setStringValues(k, v string) {
	switch k {
	case "iss":
		c.Issuer = v
	case "aud":
		c.Audience = v
	case "jti":
		c.ID = v
	case "sub":
		c.Subject = v
	default:
		c.PrivateClaims[k] = v
	}
}

func (c *ClaimSet) setNumericValues(k string, v float64) {
	switch k {
	case "exp":
		c.ExpiresAt = int64(v)
	case "iat":
		c.IssuedAt = int64(v)
	case "nbf":
		c.NotBefore = int64(v)
	default:
		c.PrivateClaims[k] = v
	}

}

// UnmarshalJSON places extra keys into PrivateClaims
func (c *ClaimSet) UnmarshalJSON(b []byte) error {
	pc := make(map[string]interface{})
	if err := json.Unmarshal(b, &pc); err != nil {
		return err
	}
	c.PrivateClaims = make(map[string]interface{})
	for k, v := range pc {
		switch val := v.(type) {
		case string:
			c.setStringValues(k, val)
		case float64:
			c.setNumericValues(k, val)
		default:
			c.PrivateClaims[k] = v
		}
	}
	return nil
}

// JWT creates a token using the signer
func (c *ClaimSet) JWT(signer Signer) (string, error) {
	payload, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	encodedPayload := make([]byte, base64.RawURLEncoding.EncodedLen(len(payload))+1)
	base64.RawURLEncoding.Encode(encodedPayload[1:], payload)
	encodedPayload[0] = '.'
	contentData := append(signer.Header(), encodedPayload...)
	sig, err := signer.Sign(contentData)
	if err != nil {
		return "", err
	}
	return string(contentData) + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}

// SetExpirationClaims sets the IssuedAt (iat) and ExpiresAt (exp) claims
func (c *ClaimSet) SetExpirationClaims(startOffset, tokenDuration time.Duration) error {
	if c == nil {
		return errors.New("nil Claim")
	}
	now := time.Now().Add(-startOffset)
	c.IssuedAt = now.Unix()
	c.ExpiresAt = now.Add(tokenDuration).Unix()
	if c.ExpiresAt <= c.IssuedAt {
		return fmt.Errorf("jws: invalid Exp = %v; must be later than Iat = %v", c.ExpiresAt, c.IssuedAt)
	}
	return nil
}

type jwtSection int

const (
	tokenHeader jwtSection = iota
	tokenPayload
	tokenSignature
)

func decodeSection(payload string, sectionType jwtSection, obj interface{}) error {
	sections := strings.Split(payload, ".")
	if len(sections) != 3 {
		// TODO(jbd): Provide more context about the error.
		return fmt.Errorf("jws: invalid token with %d sections", len(sections))
	}
	section := []byte(sections[sectionType])
	b := make([]byte, base64.RawURLEncoding.DecodedLen(len(section)))
	if _, err := base64.RawURLEncoding.Decode(b, section); err != nil {
		return err
	}
	return json.Unmarshal(b, obj)
}

// DecodePayload decodes a claim set from a JWT.
func DecodePayload(token string) (*ClaimSet, error) {
	c := &ClaimSet{}
	return c, decodeSection(token, tokenPayload, &c)
}

// DecodeHeader decodes the header from a JWT into hdr (usually
// a &map[string]interface{})
func DecodeHeader(token string, hdr interface{}) error {
	return decodeSection(token, tokenHeader, hdr)
}

// Signer provides a signature for a JWT as well as the Header
type Signer interface {
	Sign([]byte) ([]byte, error)
	Header() []byte
}

func encodeHeader(alg, keyID string) []byte {
	var hdr = struct {
		Alg string `json:"alg,omitempty"`
		Typ string `json:"typ,omitempty"`
		Kid string `json:"kid,omitempty"`
	}{alg, "JWT", keyID}
	hdrBytes, _ := json.Marshal(hdr)
	encodedHdr := make([]byte, base64.RawURLEncoding.EncodedLen(len(hdrBytes)))
	base64.RawURLEncoding.Encode(encodedHdr, hdrBytes)
	return encodedHdr
}

// RS256FromPEM creates a signer that implements the RS256 (RSA PKCS#1 with SHA-512)
// algorithm for the encoded key in pemBytes.  An error is returned if the pem encoding is
// invalid.  pemBytes should contain the contents of a PEM file using PKCS8 or PKCS1 encoding.
// PEM containers with a passphrase are not supported.
// Use the following command to convert a PKCS 12 file into a PEM.
//
//    $ openssl pkcs12 -in key.p12 -out key.pem -nodes
//
func RS256FromPEM(pemBytes []byte, keyID string) (Signer, error) {
	return rsaSignerFromPEM(pemBytes, keyID, crypto.SHA256, "RS256")
}

// RS256 creates a signer for the RS256 algorithm
func RS256(key *rsa.PrivateKey, keyID string) Signer {
	return &rsaSigner{
		key:    key,
		hash:   crypto.SHA256,
		header: encodeHeader("RS256", keyID),
	}
}

// RS384FromPEM creates a signer that implements the RS384 (RSA PKCS#1 with SHA-512)
// algorithm for the encoded key in pemBytes.  An error is returned if the pem encoding is
// invalid.  pemBytes should contain the contents of a PEM file using PKCS8 or PKCS1 encoding.
// PEM containers with a passphrase are not supported.
func RS384FromPEM(pemBytes []byte, keyID string) (Signer, error) {
	return rsaSignerFromPEM(pemBytes, keyID, crypto.SHA384, "RS384")
}

// RS384 creates a signer that implements the RS512 (RSA PKCS#1 with SHA-384)
// algorithm for the key.  keyID is the optional and will be used in the kid header claim.
func RS384(key *rsa.PrivateKey, keyID string) Signer {
	return &rsaSigner{
		key:    key,
		hash:   crypto.SHA384,
		header: encodeHeader("RS384", keyID),
	}
}

// RS512FromPEM creates a signer that implements the RS512 (RSA PKCS#1 with SHA-512)
// algorithm for the encoded key in pemBytes.  An error is returned if the pem encoding is
// invalid.  pemBytes should contain the contents of a PEM file using PKCS8 or PKCS1 encoding.
// PEM containers with a passphrase are not supported.
func RS512FromPEM(pemBytes []byte, keyID string) (Signer, error) {
	return rsaSignerFromPEM(pemBytes, keyID, crypto.SHA512, "RS512")
}

// RS512 creates a signer that implements the RS512 (RSA PKCS#1 with SHA-512)
// algorithm for the key.  keyID is the optional and will be used in the kid header claim.
func RS512(key *rsa.PrivateKey, keyID string) Signer {
	return &rsaSigner{
		key:    key,
		hash:   crypto.SHA512,
		header: encodeHeader("RS512", keyID),
	}
}

func rsaSignerFromPEM(pemBytes []byte, keyID string, hash crypto.Hash, alg string) (Signer, error) {
	key, err := ParseRSAKey(pemBytes)
	if err != nil {
		return nil, err
	}
	return &rsaSigner{
		key:    key,
		hash:   hash,
		header: encodeHeader(alg, keyID),
	}, nil
}

type rsaSigner struct {
	key    *rsa.PrivateKey
	hash   crypto.Hash
	header []byte
}

func (rs *rsaSigner) Sign(data []byte) ([]byte, error) {
	h := rs.hash.New()
	h.Write(data)
	return rsa.SignPKCS1v15(rand.Reader, rs.key, rs.hash, h.Sum(nil))
}

func (rs *rsaSigner) Header() []byte {
	return rs.header
}

type hmacSigner struct {
	secret []byte
	hash   crypto.Hash
	header []byte
}

func (h *hmacSigner) Sign(data []byte) ([]byte, error) {
	hm := hmac.New(h.hash.New, h.secret)
	hm.Write(data)
	return hm.Sum(nil), nil
}

func (h *hmacSigner) Header() []byte {
	return h.header
}

// HS256 returns a signer implementing the HMAC with SHA-256
// algorithm with the passed secret.
func HS256(secret []byte) Signer {
	return &hmacSigner{
		secret: secret,
		hash:   crypto.SHA256,
		header: encodeHeader("HS256", ""),
	}
}

// HS384 returns a signer implementing the HMAC with SHA-384
// algorithm with the passed secret.
func HS384(secret []byte) Signer {
	return &hmacSigner{
		secret: secret,
		hash:   crypto.SHA384,
		header: encodeHeader("HS384", ""),
	}
}

// HS512 returns a signer implementing the HMAC with SHA-512
// algorithm with the passed secret.
func HS512(secret []byte) Signer {
	return &hmacSigner{
		secret: secret,
		hash:   crypto.SHA512,
		header: encodeHeader("HS512", ""),
	}
}

// Verifier is a funct that verifies the signature of a specific content
type Verifier func(signature, content []byte) error

// RS256Verifier verifies the signature using PKCS1v15 using key
func RS256Verifier(key *rsa.PublicKey) Verifier {
	return rsaVerify(key, crypto.SHA256)
}

// RS384Verifier verifies the signature using PKCS1v15 using key
func RS384Verifier(key *rsa.PublicKey) Verifier {
	return rsaVerify(key, crypto.SHA384)
}

// RS512Verifier verifies the signature using PKCS1v15 using key
func RS512Verifier(key *rsa.PublicKey) Verifier {
	return rsaVerify(key, crypto.SHA512)
}

// HS256Verifier verifies the signature using SHA256 hmac using secret
func HS256Verifier(secret []byte) Verifier {
	return hmacVerify(secret, crypto.SHA256, "HS256")
}

// HS384Verifier verifies the signature using SHA384 hmac using secret
func HS384Verifier(secret []byte) Verifier {
	return hmacVerify(secret, crypto.SHA384, "HS384")
}

// HS512Verifier verifies the signature using SHA384 hmac using secret
func HS512Verifier(secret []byte) Verifier {
	return hmacVerify(secret, crypto.SHA512, "HS512")
}

// Verify tests whether the provided JWT token's signature is valid
func Verify(token string, v Verifier) error {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("jws: invalid token received, token must have 3 parts")
	}
	signatureString, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return err
	}
	signedContent := parts[0] + "." + parts[1]
	return v([]byte(signatureString), []byte(signedContent))
}

func rsaVerify(key *rsa.PublicKey, hash crypto.Hash) Verifier {
	return func(signature, signedContent []byte) error {
		h := hash.New()
		h.Write(signedContent)
		return rsa.VerifyPKCS1v15(key, hash, h.Sum(nil), signature)
	}
}

func hmacVerify(secret []byte, hash crypto.Hash, alg string) Verifier {
	return func(signature, signedContent []byte) error {
		h := hmac.New(hash.New, secret)
		h.Write(signedContent)
		if bytes.Equal(h.Sum(nil), signature) {
			return nil
		}
		return fmt.Errorf("invalid %s signature", alg)
	}
}

// ParseRSAKey converts the binary contents of a private key file
// to an *rsa.PrivateKey. It detects whether the private key is in a
// PEM container or not. If so, it extracts the the private key
// from PEM container before conversion. It only supports PEM
// containers with no passphrase.
func ParseRSAKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block != nil {
		key = block.Bytes
	}
	parsedKey, err := x509.ParsePKCS8PrivateKey(key)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("private key should be a PEM or plain PKSC1 or PKCS8; parse error: %v", err)
		}
	}
	parsed, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key is invalid")
	}
	return parsed, nil
}
