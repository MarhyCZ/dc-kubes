package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to my website!")
	})

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchData() {
	serverName := os.Getenv("BACKEND_HOST")
	endpoint := fmt.Sprintf("http://%s/api/", serverName)
	http.Get(endpoint)
}
