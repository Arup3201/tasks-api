package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir(".")))

	err := http.ListenAndServe(":8000", nil)
	if err!=nil {
		fmt.Fprintf(os.Stderr, "http listen and server failed - %s", err)
	}
}
