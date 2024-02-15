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
			if v.Name() != "tmp" {
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

			b, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatal(err)
				return nil
			}
			duplicateIntoTmp(v.Name(), path, b)

			b = removeDevTags(targetDir + "\\" + filePath)
		}
	}

	return nil
}

func removeDevTags(path string) []byte {
	_, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var finalLines string
	var skipNext bool

	lines := strings.Split(string(b), "\n")
	for _, v := range lines {
		if strings.Contains(v, "// @DEV") {
			skipNext = true
			continue
		}

		if skipNext {
			skipNext = false
			continue
		}

		finalLines = finalLines + v + "\n"
	}

	return []byte(finalLines)
}

func duplicateIntoTmp(fileName string, path string, data []byte) {
	tmpDir := targetDir + "\\tmp\\"

	if path != "" {
		tmpDir = tmpDir + path + "\\"
	}

	_, err := os.Stat(tmpDir)
	if err != nil {
		os.Mkdir(tmpDir, os.ModeDir)
	}

	os.WriteFile(tmpDir+fileName, data, os.ModePerm)
}

func runBuild() {
	_, err := os.Stat("tmp")
	if err != nil {
		log.Fatal(err)
	}

	var existingFiles []string
	d, err := os.ReadDir("tmp")
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, f := range d {
		existingFiles = append(existingFiles, f.Name())
	}

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
		cmd.Dir = targetDir + "\\tmp"
		cmd.Stderr = &serr

		err := cmd.Run()
		if err != nil {
			log.Fatal(serr.String())
			return
		}
	}

	d, err = os.ReadDir("tmp")
	if err != nil {
		log.Fatal(err)
		return
	}

	var toExtract []os.DirEntry
	for _, f := range d {
		found := false

		for _, e := range existingFiles {
			if e == f.Name() {
				found = true
			}
		}

		if !found {
			toExtract = append(toExtract, f)
		}
	}

	for _, e := range toExtract {
		b, err := os.ReadFile("tmp\\" + e.Name())
		if err != nil {
			log.Fatal(err)
			return
		}

		os.WriteFile(targetDir+"\\"+e.Name(), b, os.ModePerm)
	}

	os.Remove("tmp")
}

// run Config.Build.Command in tmp dir MAKE SURE THAT COMMANDS ARE SPLIT BY &&
// check for any new directories in tmp dir, if exists, likely the dist directory
// duplicate dist directory outside of tmp dir
// delete tmp dir
