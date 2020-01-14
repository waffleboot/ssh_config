package main

import "bufio"

type ErrorWriter struct {
	writer *bufio.Writer
	err    error
}

func NewErrorWriter(writer *bufio.Writer) *ErrorWriter {
	return &ErrorWriter{writer: writer}
}

func (w *ErrorWriter) WriteString(s string) {
	if w.err == nil {
		_, w.err = w.writer.WriteString(s)
	}
}

func (w *ErrorWriter) Err() error {
	return w.err
}
