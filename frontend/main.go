package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type NameDay struct {
	Date string `json:"date"`
	Name string `json:"name"`
}

type DataRepo struct {
	Today     NameDay
	LastFetch time.Time
}

//go:embed template.html
var templateStr string

func main() {
	t, err := template.New("template").Parse(templateStr)
	if err != nil {
		log.Fatalln("Could not parse HTML template, ", err.Error())
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		repo, err := fetchData()
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "Could not connect to backend")
			return
		}
		t.Execute(w, struct {
			Name      string
			Date      string
			LastFetch string
		}{Name: repo.Today.Name,
			Date:      repo.Today.Date,
			LastFetch: repo.LastFetch.UTC().String()})
	})

	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchData() (DataRepo, error) {
	serverName := os.Getenv("BACKEND_HOST")
	endpoint := fmt.Sprintf("http://%s:3001/api/svatek", serverName)
	res, err := http.Get(endpoint)
	if err != nil {
		return DataRepo{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return DataRepo{}, err
	}
	var data DataRepo
	err = json.Unmarshal(body, &data)
	if err != nil {
		return DataRepo{}, err
	}
	return data, nil
}
