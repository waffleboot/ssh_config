package main

import (
	"bufio"
	"io"
	"strings"
	t "text/template"
)

func (u updater) process(r io.Reader, w io.Writer) error {
	buf := bufio.NewWriter(w)
	scanner := bufio.NewScanner(r)
	if errFind := u.findAndWrite(scanner, buf); errFind != nil {
		return errFind
	}
	if errWrite := u.writeUpdate(buf); errWrite != nil {
		return errWrite
	}
	if errWrite := restWrite(scanner, r, buf); errWrite != nil {
		return errWrite
	}
	return buf.Flush()
}

func (u updater) findAndWrite(scanner *bufio.Scanner, w *bufio.Writer) error {
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "host") {
			slice := strings.Fields(text)
			if len(slice) > 1 && slice[1] == u.Name {
				return scanner.Err()
			}
		}
		if _, errWrite := w.WriteString(text); errWrite != nil {
			return errWrite
		}
		if _, errWrite := w.WriteRune('\n'); errWrite != nil {
			return errWrite
		}
	}
	return scanner.Err()
}

const template = `host {{ .Name }}
	HostName {{ .Host }}
	IdentityFile {{ .Identity }}
	StrictHostKeyChecking no
	User {{ .User }}
`

func (u updater) writeUpdate(w *bufio.Writer) error {
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

func restWrite(scanner *bufio.Scanner, r io.Reader, w *bufio.Writer) error {
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "host") {
			if errWrite := writeWithNewLine(text, w); errWrite != nil {
				return errWrite
			}
			break
		}
	}
	for scanner.Scan() {
		if errWrite := writeWithNewLine(scanner.Text(), w); errWrite != nil {
			return errWrite
		}
	}
	return scanner.Err()
}
