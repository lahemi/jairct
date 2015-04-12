package main

import "os"

func mustMakeDir(dir string) {
	if fi, err := os.Stat(dir); err != nil {
		if err := os.MkdirAll(dir, 0777); err != nil {
			die("Unable to create " + dir)
		}
	} else if !fi.Mode().IsDir() {
		die("There is a file with the same name as the " + dir + " already present")
	}
}

func mustHaveFile(file string) {
	if fi, err := os.Stat(file); err == nil && fi.Mode().IsRegular() {
		return
	}
	die("No file: " + file + " found")
}

func setupEnvironment() {
	mustMakeDir(configPath)
	mustMakeDir(dataPath)

	mustHaveFile(dbFile)

	setupDB()
}
