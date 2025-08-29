package main

import (
	"fmt"
	"log"
	"net/http"
)

const PORT = ":8000"

func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!\n")
	})


	log.Printf("Starting the server at %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
