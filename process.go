package main

import (
	"bufio"
	"io"
	"strings"
	t "text/template"
)

func (u updater) copyWithUpdate(src io.Reader, dst io.Writer) error {
	buf := bufio.NewWriter(dst)
	sourceScanner := bufio.NewScanner(src)
	if errFind := u.findServerName(sourceScanner, buf); errFind != nil {
		return errFind
	}
	if errWrite := u.updateServerConfig(buf); errWrite != nil {
		return errWrite
	}
	if errWrite := restWrite(sourceScanner, buf); errWrite != nil {
		return errWrite
	}
	return buf.Flush()
}

func (u updater) findServerName(sourceScanner *bufio.Scanner, w *bufio.Writer) error {
	for sourceScanner.Scan() {
		textLine := sourceScanner.Text()
		if strings.HasPrefix(textLine, "host") {
			hostName := strings.Fields(textLine)
			if len(hostName) > 1 && hostName[1] == u.ServerName {
				return nil
			}
		}
		if errWrite := writeWithNewLine(textLine, w); errWrite != nil {
			return errWrite
		}
	}
	return sourceScanner.Err()
}

const template = `host {{ .Name }}
	HostName {{ .Host }}
	IdentityFile {{ .Identity }}
	StrictHostKeyChecking no
	User {{ .User }}
`

func (u updater) updateServerConfig(w *bufio.Writer) error {
	tpl, tplError := t.New("update").Parse(template)
	if tplError != nil {
		return tplError
	}
	return tpl.Execute(w, u)
}

func writeWithNewLine(s string, w *bufio.Writer) error {
	if _, errWrite := w.WriteString(s); errWrite != nil {
		return errWrite
	}
	if _, errWrite := w.WriteRune('\n'); errWrite != nil {
		return errWrite
	}
	return nil
}

func restWrite(sourceScanner *bufio.Scanner, w *bufio.Writer) error {
	for sourceScanner.Scan() {
		textLine := sourceScanner.Text()
		if strings.HasPrefix(textLine, "host") {
			if errWrite := writeWithNewLine(textLine, w); errWrite != nil {
				return errWrite
			}
			break
		}
	}
	for sourceScanner.Scan() {
		if errWrite := writeWithNewLine(sourceScanner.Text(), w); errWrite != nil {
			return errWrite
		}
	}
	return sourceScanner.Err()
}
