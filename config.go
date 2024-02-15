package main

type Config struct {
	Commands Commands `json:"commands"`
}

type Commands struct {
	Build string `json:"build"`
}
