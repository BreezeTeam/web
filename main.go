package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	myhandler := new(Myhandler)

	log.Fatal(http.ListenAndServe(":9999", myhandler))
}

/**
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
*/
type Myhandler struct{}

func (h *Myhandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", r.URL.Path)
	case "/hello":
		for k, v := range r.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}
