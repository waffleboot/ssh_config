package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func (u updater) makeBackup() (io.ReadCloser, error) {
	errRename := os.Rename(u.config, u.backup)
	if errRename != nil {
		_, errStat := os.Stat(u.backup)
		if errStat != nil {
			return ioutil.NopCloser(bytes.NewReader(nil)), nil
			// return nil, fmt.Errorf("files not found: ['%s','%s']", config, backup)
		}
	}
	file, errOpen := os.Open(u.backup)
	if errOpen != nil {
		return nil, errOpen
	}
	return file, nil
}

func (u updater) restoreBackup() error {
	return os.Rename(u.backup, u.config)
}
