package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

func traverseDirectories(path string) []string {
	_, err := os.Stat(targetDir + "\\" + path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	dir, err := os.ReadDir(targetDir + "\\" + path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	for _, v := range dir {
		if v.IsDir() {
			if v.Name() != "tmp" && v.Name() != "node_modules" && v.Name() != ".angular" && v.Name() != ".vscode" && v.Name() != ".idea" {
				if path == "" {
					traverseDirectories("\\" + v.Name() + "\\")
				} else {
					traverseDirectories("\\" + path + "\\" + v.Name() + "\\")
				}
			}
		} else {
			var filePath string
			if path == "" {
				filePath = v.Name()
			} else {
				filePath = path + "\\" + v.Name()
			}

			b, o := removeDevTags(targetDir + "\\" + filePath)
			if b != nil {
				log.Println(filePath)
				duplicateIntoTmp(v.Name(), path, o)
			}

			err = os.WriteFile(targetDir+"\\"+filePath, b, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

// mod, original
func removeDevTags(path string) ([]byte, []byte) {
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil, nil
	}

	var finalLines string
	var skipNext bool
	flagFound := false

	lines := strings.Split(string(b), "\n")
	for _, v := range lines {
		if strings.Contains(v, "// @DEV") {
			skipNext = true
			flagFound = true
			continue
		}

		if skipNext {
			skipNext = false
			continue
		}

		finalLines = finalLines + v + "\n"
	}

	if flagFound {
		return []byte(finalLines), b
	} else {
		return nil, b
	}
}

func duplicateIntoTmp(fileName string, path string, data []byte) {
	tmpDir := targetDir + "\\tmp\\"

	if path != "" {
		tmpDir = tmpDir + path + "\\"
	}

	_, err := os.Stat(tmpDir)
	if err != nil {
		os.MkdirAll(tmpDir, os.ModeDir)
	}

	os.WriteFile(tmpDir+fileName, data, os.ModePerm)
}

func runBuild() {
	var commands []string
	commands = strings.Split(cb.Commands.Build, "&&")

	commandArgs := map[string][]string{}
	for _, v := range commands {
		args := strings.Split(v, " ")
		commandArgs[args[0]] = args[1:]
	}

	for k, v := range commandArgs {
		cmd := exec.Command(k, v...)

		var serr bytes.Buffer
		cmd.Dir = targetDir
		cmd.Stderr = &serr

		err := cmd.Run()
		if err != nil {
			log.Fatal(serr.String())
			return
		}
	}

	// Traverse tmp again
	// Restore via WriteFile
	// Remove tmp

	os.Remove("tmp")
}

// run Config.Build.Command in tmp dir MAKE SURE THAT COMMANDS ARE SPLIT BY &&
// check for any new directories in tmp dir, if exists, likely the dist directory
// duplicate dist directory outside of tmp dir
// delete tmp dir

// IF file contains // @DEV then insert into TMP dir
// AFTER build, restore files from TMP back into root dir
// Best not to delete files, only update so dynamic in VSC
