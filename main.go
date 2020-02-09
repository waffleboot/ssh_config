package main

import (
	"fmt"
	"io"
	"log"
	"os"
	pth "path"
)

func newUpdater() updater {
	u := updater{}
	path := os.Args[1]
	u.Name = os.Args[2]
	u.Host = os.Args[3]
	u.User = os.Args[4]
	u.config = pth.Join(path, "config")
	u.backup = pth.Join(path, "config.backup")
	u.Identity = os.Args[5]
	return u
}

type updater struct {
	User     string
	Name     string
	Host     string
	config   string
	backup   string
	Identity string
}

func (u updater) runOrRestore() error {
	if errRun := u.run(); errRun != nil {
		if errRestore := u.restoreBackup(); errRestore != nil {
			fmt.Fprintln(os.Stderr, errRun)
			return errRestore
		}
		return errRun
	}
	return nil
}

func (u updater) run() error {
	r, errBackup := u.makeBackup()
	if errBackup != nil {
		return errBackup
	}
	defer close(r)
	w, errConfig := os.Create(u.config)
	if errConfig != nil {
		return errConfig
	}
	defer close(w)
	return u.process(r, w)
}

func close(file io.Closer) {
	errClose := file.Close()
	if errClose != nil {
		fmt.Fprintln(os.Stderr, errClose)
	}
}

func (u updater) dump() {
	file, errOpen := os.Open(u.config)
	if errOpen != nil {
		return
	}
	defer file.Close()
	io.Copy(os.Stdout, file)
}

func main() {
	if len(os.Args) < 6 {
		fmt.Println("usage: <.ssh-dir> <name> <host> <user> <identity file>")
		os.Exit(1)
	}
	u := newUpdater()
	if err := u.runOrRestore(); err != nil {
		log.Fatal(err)
	}
	u.dump()
}
