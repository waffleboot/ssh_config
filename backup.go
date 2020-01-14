package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func makeBackup() (io.ReadCloser, error) {
	errRename := os.Rename(config, backup)
	if errRename != nil {
		_, errStat := os.Stat(backup)
		if errStat != nil {
			return ioutil.NopCloser(bytes.NewReader(nil)), nil
			// return nil, fmt.Errorf("files not found: ['%s','%s']", config, backup)
		}
	}
	file, errOpen := os.Open(backup)
	if errOpen != nil {
		return nil, errOpen
	}
	return file, nil
}

func restoreBackup() error {
	return os.Rename(backup, config)
}
