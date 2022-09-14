# OAuth2 for Go

[![GoDoc](https://godoc.org/github.com/jfcote87/oauth2?status.svg)](https://godoc.org/github.com/jfcote87/oauth2)

oauth2 package is forked from golang.org/x/oauth2.  The major change is the handling of contexts.
Context for requesting tokens is taken from the http.request rather than stored in the
TokenSource.  This allows for a TokenSource to exist accross contexts and, when wrapped
by ReuseTokenSource func, to be safe for concurrent use.

The IDSecretInBody flag is added to the endpoint definition to replace "brokenAuthHeaderProviders"
logic.

Issues that are resolved in this implementation includ

https://github.com/golang/oauth2/issues/223 all: find a better way to set a custom HTTP client

https://github.com/golang/oauth2/issues/262 TokenSource. Token method should take in a Context

https://github.com/golang/oauth2/issues/198 google: no way to set jws.PrivateClaims (needed for Firebase)

https://github.com/golang/oauth2/issues/198 Token expiration tolerance should be configurable #249

https://github.com/golang/oauth2/issues/256 OAuth2: Ability to specify "audience" parameter to token refresh #256

https://github.com/golang/oauth2/issues/84 CacheToken/transport confusion
 

## Installation

~~~~
go get github.com/jfcote87/oauth2
~~~~


