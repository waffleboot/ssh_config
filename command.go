package main

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/urfave/cli"
)

type updater struct {
	User           string
	ServerName     string
	Host           string
	Identity       string
	configFileName string
	backupFileName string
}

func newUpdater(context *cli.Context) updater {
	u := updater{}
	dir := context.Args().Get(0)
	u.ServerName = context.Args().Get(1)
	u.Host = context.Args().Get(2)
	u.User = context.Args().Get(3)
	u.configFileName = path.Join(dir, "config")
	u.backupFileName = path.Join(dir, "config.backup")
	u.Identity = context.Args().Get(4)
	return u
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
	return u.copyConfigWithUpdate(src, dst)
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

func mainAction(context *cli.Context) error {
	if err := checkArgs(context, 5, exactArgs); err != nil {
		return err
	}
	u := newUpdater(context)
	if err := u.update(); err != nil {
		return err
	}
	if context.Bool("verbose") {
		u.printSSHConfig(os.Stdout)
	}
	return nil

}
