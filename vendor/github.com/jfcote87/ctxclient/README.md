# ctxclient package

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/jfcote87/ctxclient)
[![BSD 3-Clause License](https://img.shields.io/badge/license-bsd%203--clause-blue.svg)](https://github.com/jfcote87/ctxclient/blob/master/LICENSE)

Package ctxclient offers utilities for easing the making of http calls by handing boilerplate functions
and providing for the selection of an http.Client/Transport based upon the current context.  A ctxclient.Func
provides the basis of the functionality via the Client and Do methods.  A nil ctxclient simply returns the default
client while a Func that return the ErrUseDefault will also return the default on a Client call.

The Do func simplifies making an http request by handling context timeouts and checking for non 2xx response.  In the
case of a non 2xx response a NotSuccess error is returned which contains the body, header and status code of the response.

```go
var f ctxclient.Func

res, err = f.Do(ctx, req)
if err != nil {
    ...
}
```

Funcs may also be used to global defaults by registering a Func via the ctxclient.Register function.

Please review examples_test.go.

This package borrows from ideas found in
golang.org/x/oauth2.
