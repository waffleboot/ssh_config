package main

import "io"
import "bufio"

type myWriter struct {
	*bufio.Writer
}

func newWriter(dst io.Writer) myWriter {
	return myWriter{bufio.NewWriter(dst)}
}

func (w myWriter) writeWithNewLine(s string) error {
	if _, errWrite := w.WriteString(s); errWrite != nil {
		return errWrite
	}
	if _, errWrite := w.WriteRune('\n'); errWrite != nil {
		return errWrite
	}
	return nil
}
