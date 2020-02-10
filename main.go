package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli"
)

var version = ""

var gitCommit = ""

const usage = "usage: <.ssh-dir> <server-name> <hostname> <ssh-user> <identity file>"

func main() {

	app := cli.NewApp()
	app.Name = "ssh_config"
	app.Usage = usage

	var v []string
	if version != "" {
		v = append(v, version)
	}
	if gitCommit != "" {
		v = append(v, fmt.Sprintf("commit: %s", gitCommit))
	}
	app.Version = strings.Join(v, "\n")

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "print resulting config",
		},
	}

	app.Action = mainAction

	if err := app.Run(os.Args); err != nil {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
