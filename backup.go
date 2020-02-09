package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
)

func (u updater) makeBackup() (source io.ReadCloser, err error) {
	errRename := os.Rename(u.configFileName, u.backupFileName)
	if errRename != nil {
		_, errStat := os.Stat(u.backupFileName)
		if errStat != nil {
			emptySource := ioutil.NopCloser(bytes.NewReader(nil))
			return emptySource, nil
			// return nil, fmt.Errorf("files not found: ['%s','%s']", config, backup)
		}
	}
	file, errOpen := os.Open(u.backupFileName)
	if errOpen != nil {
		return nil, errOpen
	}
	return file, nil
}

func (u updater) restoreBackup() error {
	return os.Rename(u.backupFileName, u.configFileName)
}
