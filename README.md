# dispatch

[![GoDoc](https://godoc.org/github.com/i2bskn/dispatch?status.svg)](https://godoc.org/github.com/i2bskn/dispatch)
[![Build Status](https://travis-ci.org/i2bskn/dispatch.svg?branch=master)](https://travis-ci.org/i2bskn/dispatch)
[![codecov](https://codecov.io/gh/i2bskn/dispatch/branch/master/graph/badge.svg)](https://codecov.io/gh/i2bskn/dispatch)

dispatch is HTTP request multiplexer compatible with `ServeMux` of `net/http`.

## Dependencies

- [Go](https://golang.org/) 1.7 or lator

No dependency on the third party library.

## Installation

```
go get -u github.com/i2bskn/dispatch
```

## Usage

```Go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/i2bskn/dispatch"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", dispatch.Param(r, "name"))
}

func main() {
	mux := dispatch.New()
	mux.HandleFunc("/hello/:name", hello)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Home!")
	})
	log.Fatal(http.ListenAndServe(":8080", mux))
}
```

See also [GoDoc](https://godoc.org/github.com/i2bskn/dispatch).

## License

dispatch is available under the MIT.