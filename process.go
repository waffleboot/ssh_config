package main

import (
	"io"
	"text/template"
)

func (u updater) copyConfigWithUpdate(src io.Reader, dst io.Writer) error {
	buf := newWriter(dst)
	sourceScanner := newScanner(src)
	if errFind := sourceScanner.findServerName(u.ServerName, buf.writeWithNewLine); errFind != nil {
		return errFind
	}
	if errHandle := u.updateServerConfig(buf); errHandle != nil {
		return errHandle
	}
	if errWrite := sourceScanner.copyRest(buf.writeWithNewLine); errWrite != nil {
		return errWrite
	}
	return buf.Flush()
}

const serverTemplate = `host {{ .ServerName }}
{{if ne .Host "" }}{{printf "\tHostName %s" .Host }}{{else}}{{"\t# HostName"}}{{end}}
	IdentityFile {{ .Identity }}
	StrictHostKeyChecking no
	User {{ .User }}
`

func (u updater) updateServerConfig(w myWriter) error {
	tpl, tplError := template.New("update").Parse(serverTemplate)
	if tplError != nil {
		return tplError
	}
	return tpl.Execute(w, u)
}
