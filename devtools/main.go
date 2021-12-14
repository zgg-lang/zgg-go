package main

import (
	"os"

	"github.com/urfave/cli"
)

var commands []*cli.Command

func main() {
	app := cli.NewApp()
	app.Commands = commands
	app.Run(os.Args)
}
