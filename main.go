package main

import (
	"fmt"
	"io"
	"os"
	"path"
)

func newUpdater(args []string) updater {
	u := updater{}
	dir := args[1]
	u.ServerName = args[2]
	u.Host = args[3]
	u.User = args[4]
	u.configFileName = path.Join(dir, "config")
	u.backupFileName = path.Join(dir, "config.backup")
	u.Identity = args[5]
	return u
}

type updater struct {
	User           string
	ServerName     string
	Host           string
	Identity       string
	configFileName string
	backupFileName string
}

func (u updater) update() error {
	if errUpdate := u.tryUpdate(); errUpdate != nil {
		if errRestore := u.restoreBackup(); errRestore != nil {
			fmt.Fprintln(os.Stderr, errUpdate)
			return errRestore
		}
		return errUpdate
	}
	return nil
}

func (u updater) tryUpdate() error {
	src, errBackup := u.makeBackup()
	if errBackup != nil {
		return errBackup
	}
	defer close(src)
	dst, errConfig := os.Create(u.configFileName)
	if errConfig != nil {
		return errConfig
	}
	defer close(dst)
	return u.copyWithUpdate(src, dst)
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

const usage = "usage: <.ssh-dir> <server-name> <hostname> <ssh-user> <identity file>"

func main() {
	if len(os.Args) < 6 {
		fmt.Println(usage)
		os.Exit(1)
	}
	u := newUpdater(os.Args)
	if err := u.update(); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	u.printSSHConfig(os.Stdout)
}
