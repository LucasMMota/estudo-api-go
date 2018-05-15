package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

/*
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/

func main() {

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome!")
}
