package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var name = "master"
var config = "config"
var backup = "config.backup"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need hostname")
	}
	if err := runOrRestore(os.Args[1]); err != nil {
		log.Fatal(err)
	}
}

func runOrRestore(hostname string) error {
	if errRun := run(hostname); errRun != nil {
		if errRestore := restoreBackup(); errRestore != nil {
			fmt.Fprintln(os.Stderr, errRestore)
		}
		return errRun
	}
	return nil
}

func run(hostname string) error {
	r, errBackup := makeBackup()
	if errBackup != nil {
		return errBackup
	}
	defer close(r)
	w, errConfig := os.Create(config)
	if errConfig != nil {
		return errConfig
	}
	defer close(w)
	return process(hostname, r, w)
}

func close(file io.Closer) {
	errClose := file.Close()
	if errClose != nil {
		fmt.Fprintln(os.Stderr, errClose)
	}
}
