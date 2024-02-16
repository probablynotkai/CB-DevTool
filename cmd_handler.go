package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

var commands = []string{"init", "build"}

func handleArgs() {
	args := os.Args
	if len(args) != 2 || (len(args) > 1 && !includes(commands, args[1])) {
		log.Println("Usage: cb <init | build>")
		return
	}

	if len(args) == 2 && strings.EqualFold(args[1], "init") {
		initialiseProject()
		return
	}

	if len(args) == 2 && strings.EqualFold(args[1], "build") {
		loadConfig()
		buildProject()
		return
	}
}

func buildProject() {
	log.Println("Removing dev flags...")
	traverseDirectoriesAndScan("")

	log.Println("Building project...")
	runBuild()

	log.Println("Restoring original files...")
	traverseDirectoriesAndRestore("")

	err := os.RemoveAll(".\\tmp")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Done!")
}

func initialiseProject() {
	_, err := os.Stat(".\\cb.json")
	if err != nil {
		log.Println("Initialising cb-devtool...")

		tmpCfg := Config{}
		tmpCfg.Commands.Build = "ng build --configuration production"

		b, err := json.MarshalIndent(tmpCfg, "", "    ")
		if err != nil {
			log.Fatal(err)
			return
		}

		os.WriteFile(".\\cb.json", b, os.ModePerm)
		log.Println("This project has been initialised.")
		return
	}

	log.Println("This project has already been initialised.")
}
