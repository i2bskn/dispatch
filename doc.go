/*
Package tensile is net/http compatible HTTP request multiplexer.

For example:
	package main

	import (
		"fmt"
		"log"
		"net/http"

		"github.com/i2bskn/tensile"
	)

	func hello(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello %s!", tensile.Param(r, "name"))
	}

	func main() {
		mux := tensile.New()
		mux.HandleFunc("/hello/:name", hello)
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "Home!")
		})

		log.Fatal(http.ListenAndServe(":8080", mux))
	}
*/
package tensile
