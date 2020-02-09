package main

import (
	"bufio"
	"io"
	"strings"
)

type myScanner struct {
	*bufio.Scanner
}

func newScanner(src io.Reader) myScanner {
	return myScanner{bufio.NewScanner(src)}
}

func (sourceScanner myScanner) findServerName(serverName string, other func(string) error) error {
	for sourceScanner.Scan() {
		textLine := sourceScanner.Text()
		if strings.HasPrefix(textLine, "host") {
			hostName := strings.Fields(textLine)
			if len(hostName) > 1 && hostName[1] == serverName {
				return nil
			}
		}
		if errWrite := other(textLine); errWrite != nil {
			return errWrite
		}
	}
	return sourceScanner.Err()
}

func (sourceScanner myScanner) copyRest(other func(string) error) error {
	for sourceScanner.Scan() {
		textLine := sourceScanner.Text()
		if strings.HasPrefix(textLine, "host") {
			if errWrite := other(textLine); errWrite != nil {
				return errWrite
			}
			break
		}
	}
	for sourceScanner.Scan() {
		if errWrite := other(sourceScanner.Text()); errWrite != nil {
			return errWrite
		}
	}
	return sourceScanner.Err()
}
