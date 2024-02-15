package main

import (
	"encoding/json"
	"log"
	"os"
)

var targetDir string
var cb Config

func main() {
	args := os.Args
	if len(args) != 2 {
		log.Println("Usage: cb <dir>")
		return
	}

	targetDir = args[1]
	exDir, err := os.Executable()
	if err != nil {
		log.Fatal(err)
		return
	}

	cfgDir := exDir[:len(exDir)-6]

	_, err = os.Stat(cfgDir + "cb.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	b, err := os.ReadFile(cfgDir + "cb.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal(b, &cb)
	if err != nil {
		log.Fatal(err)
		return
	}

	traverseDirectories("")
	runBuild()

	// Traverse dir, read all files, for each line, if line == // @DEV, skip next line in tmp build
}
