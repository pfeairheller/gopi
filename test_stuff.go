package main

import (
	"fmt"
	"net/http"
	"log"
)


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
