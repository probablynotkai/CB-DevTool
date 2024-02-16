package main

import (
	"encoding/json"
	"log"
	"os"
)

var cb Config

type Config struct {
	Commands Commands `json:"commands"`
}

type Commands struct {
	Build string `json:"build"`
}

func loadConfig() {
	log.Println("Loading configuration...")
	_, err := os.Stat(".\\cb.json")
	if err != nil {
		log.Fatal("This project hasn't been initialised, use cb init.")
		return
	}

	b, err := os.ReadFile(".\\cb.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal(b, &cb)
	if err != nil {
		log.Fatal(err)
		return
	}
}
