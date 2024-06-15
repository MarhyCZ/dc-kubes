package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
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

func main() {
	repo, err := cacheData()
	if err != nil {
		log.Fatal("Could not start app " + err.Error())
	}
	http.HandleFunc("/api/svatek", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(repo)
	})

	err = http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func cacheData() (DataRepo, error) {
	time := time.Now()
	res, err := http.Get("https://svatky.adresa.info/json")
	if err != nil {
		return DataRepo{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return DataRepo{}, err
	}

	var data []NameDay
	err = json.Unmarshal(body, &data)
	if err != nil {
		return DataRepo{}, err
	}
	return DataRepo{
		Today:     data[0],
		LastFetch: time,
	}, nil
}
