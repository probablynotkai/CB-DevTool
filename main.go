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

	log.Println("Loading configuration...")
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

	log.Println("Removing dev flags...")
	traverseDirectoriesAndScan("")

	log.Println("Building project...")
	runBuild()

	log.Println("Restoring original files...")
	traverseDirectoriesAndRestore("")

	err = os.RemoveAll(targetDir + "\\tmp")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Done!")
}
