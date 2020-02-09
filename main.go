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
	dir := os.Args[1]
	u.Name = os.Args[2]
	u.Host = os.Args[3]
	u.User = os.Args[4]
	u.configFileName = pth.Join(dir, "config")
	u.backupFileName = pth.Join(dir, "config.backup")
	u.Identity = os.Args[5]
	return u
}

type updater struct {
	User           string
	Name           string
	Host           string
	configFileName string
	backupFileName string
	Identity       string
}

func (u updater) tryUpdate() error {
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
	w, errConfig := os.Create(u.configFileName)
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

func (u updater) printSSHConfig(out io.Writer) error {
	file, errOpen := os.Open(u.configFileName)
	if errOpen != nil {
		return errOpen
	}
	defer file.Close()
	_, errCopy := io.Copy(out, file)
	return errCopy
}

const usage = "usage: <.ssh-dir> <name> <host> <user> <identity file>"

func main() {
	if len(os.Args) < 6 {
		fmt.Println(usage)
		os.Exit(1)
	}
	u := newUpdater()
	if err := u.tryUpdate(); err != nil {
		log.Fatal(err)
	}
	u.printSSHConfig(os.Stdout)
}
