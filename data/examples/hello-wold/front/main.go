package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	backHost := os.Getenv("BACKEND_HOST")
	url := fmt.Sprintf("http://%s:1234", backHost)
	log.Println("Starting frontend. Backend is ", url)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := ""
		log.Println("Calling ", url)
		resp, err := http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			data = err.Error()
		} else {
			buf, err := io.ReadAll(resp.Body)
			if err != nil {
				data = err.Error()
			} else {
				data = string(buf)
			}
		}
		log.Printf("From backend: %s", data)
		fmt.Fprintf(w, "Received backend data: %s\n", data)
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
