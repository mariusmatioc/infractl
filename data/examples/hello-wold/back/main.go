package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting backend")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request from:", r.RemoteAddr)
		fmt.Fprintf(w, "Hello World!")
	})

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		log.Fatal(err)
	}
}
