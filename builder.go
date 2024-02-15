package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
)

var ignoredDir = []string{"tmp", "node_modules", ".angular", ".vscode", ".idea", ".git"}

func traverseDirectoriesAndScan(path string) {
	_, err := os.Stat(targetDir + "\\" + path)
	if err != nil {
		log.Fatal(err)
		return
	}

	dir, err := os.ReadDir(targetDir + "\\" + path)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, v := range dir {
		if v.IsDir() {
			if !includes(ignoredDir, v.Name()) {
				if path == "" {
					traverseDirectoriesAndScan("\\" + v.Name() + "\\")
				} else {
					traverseDirectoriesAndScan("\\" + path + "\\" + v.Name() + "\\")
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
				duplicateIntoTmp(v.Name(), path, o)

				err = os.WriteFile(targetDir+"\\"+filePath, b, os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}
}

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

	var firstIter string
	var skipNext bool
	flagFound := false

	lines := strings.Split(string(b), "\n")
	for _, v := range lines {
		if strings.Contains(strings.ToUpper(v), "// @DEV") {
			skipNext = true
			flagFound = true
			continue
		}

		if skipNext {
			skipNext = false
			continue
		}

		firstIter = firstIter + v + "\n"
	}

	var finalLines string
	devOpen := false
	lines = strings.Split(firstIter, "\n")
	for _, v := range lines {
		if strings.Contains(strings.ToUpper(v), "// @START-DEV") {
			devOpen = true
			flagFound = true
			continue
		}

		if strings.Contains(strings.ToUpper(v), "// @END-DEV") {
			if !devOpen {
				log.Println("[WARN] @END-DEV tag found but no opening tag.")
			} else {
				devOpen = false
			}
			continue
		}

		if devOpen {
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
	commands := strings.Split(cb.Commands.Build, "&&")
	commandArgs := map[string][]string{}

	for _, v := range commands {
		args := strings.Split(v, " ")
		commandArgs[args[0]] = args[1:]
	}

	for k, v := range commandArgs {
		log.Printf("Executing command starting %s...\n", k)
		cmd := exec.Command(k, v...)

		var serr bytes.Buffer
		cmd.Dir = targetDir
		cmd.Stderr = &serr

		b, err := cmd.Output()
		if err != nil {
			log.Fatal(serr.String())
			return
		}

		log.Println(string(b))
	}
}

func traverseDirectoriesAndRestore(path string) {
	_, err := os.Stat(targetDir + "\\tmp\\" + path)
	if err != nil {
		log.Fatal(err)
		return
	}

	dir, err := os.ReadDir(targetDir + "\\tmp\\" + path)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, v := range dir {
		if v.IsDir() {
			if !includes(ignoredDir, v.Name()) {
				if path == "" {
					traverseDirectoriesAndRestore(v.Name() + "\\")
				} else {
					traverseDirectoriesAndRestore(path + "\\" + v.Name() + "\\")
				}
			}
		} else {
			var filePath string
			if path == "" {
				filePath = v.Name()
			} else {
				filePath = path + "\\" + v.Name()
			}

			b, err := os.ReadFile(targetDir + "\\tmp\\" + filePath)
			if err != nil {
				log.Println("here?")
				log.Fatal(err)
				return
			}

			log.Printf("Restored %s!\n", v.Name())
			os.WriteFile(filePath, b, os.ModePerm)
		}
	}
}

func includes(arr []string, find string) bool {
	for _, v := range arr {
		if strings.EqualFold(v, find) {
			return true
		}
	}
	return false
}

// run Config.Build.Command in tmp dir MAKE SURE THAT COMMANDS ARE SPLIT BY &&
// check for any new directories in tmp dir, if exists, likely the dist directory
// duplicate dist directory outside of tmp dir
// delete tmp dir

// IF file contains // @DEV then insert into TMP dir
// AFTER build, restore files from TMP back into root dir
// Best not to delete files, only update so dynamic in VSC
