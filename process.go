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
	if errFind := findServerName(sourceScanner, buf, u.ServerName); errFind != nil {
		return errFind
	}
	if errWrite := updateServerConfig(buf, u); errWrite != nil {
		return errWrite
	}
	if errWrite := copyRest(sourceScanner, buf); errWrite != nil {
		return errWrite
	}
	return buf.Flush()
}

func findServerName(sourceScanner *bufio.Scanner, w *bufio.Writer, serverName string) error {
	for sourceScanner.Scan() {
		textLine := sourceScanner.Text()
		if strings.HasPrefix(textLine, "host") {
			hostName := strings.Fields(textLine)
			if len(hostName) > 1 && hostName[1] == serverName {
				return nil
			}
		}
		if errWrite := writeWithNewLine(textLine, w); errWrite != nil {
			return errWrite
		}
	}
	return sourceScanner.Err()
}

const template = `host {{ .ServerName }}
{{if ne .Host "" }}{{printf "\tHostName %s" .Host }}{{else}}{{"\t# HostName"}}{{end}}
	IdentityFile {{ .Identity }}
	StrictHostKeyChecking no
	User {{ .User }}
`

func updateServerConfig(w *bufio.Writer, u updater) error {
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

func copyRest(sourceScanner *bufio.Scanner, w *bufio.Writer) error {
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
