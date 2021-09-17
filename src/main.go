package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("gobase test server: starting on port 2000")
	http.ListenAndServe(":2000", nil)
}
