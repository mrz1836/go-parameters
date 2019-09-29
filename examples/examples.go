/*
Package main is an example package showing the use of the parameters package
*/
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mrz1836/go-parameters"
)

// Index is a basic request
func Index(w http.ResponseWriter, req *http.Request) {
	_, _ = fmt.Fprint(w, `{"Hello":"World"}`)
}

// Hello is an example of using parameters from the request
func Hello(w http.ResponseWriter, req *http.Request) {

	params := parameters.GetParams(req)

	name, ok := params.GetStringOk("name")
	if !ok {
		name = "unknown"
	}

	_, _ = fmt.Fprintf(w, `{"hello":"%s"}`, name)
}

// main starts the router and http server
func main() {
	router := httprouter.New()
	router.GET("/", parameters.GeneralJSONResponse(Index))
	router.GET("/hello/:name", parameters.GeneralJSONResponse(Hello))

	log.Println("Running examples on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
