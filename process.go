package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
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
			if len(slice) > 1 && slice[1] == u.name {
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

func (u updater) writeUpdate(w *bufio.Writer) error {
	errWriter := NewErrorWriter(w)
	errWriter.WriteString(fmt.Sprintf("host %s\n", u.name))
	errWriter.WriteString(fmt.Sprintf("\tHostName %s\n", u.host))
	errWriter.WriteString(fmt.Sprintf("\tIdentityFile %s\n", u.identity))
	errWriter.WriteString("\tStrictHostKeyChecking no\n")
	errWriter.WriteString(fmt.Sprintf("\tUser %v\n", u.user))
	return errWriter.Err()
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
