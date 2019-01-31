package main

import (
	"log"
	"os"

	"github.com/satooon/diiff/action"

	"github.com/urfave/cli"
)

const (
	name    = "diiff"
	author  = "satooon"
	version = "0.0.1"
	usage   = "file diff"
)

func main() {
	app := cli.NewApp()
	app.Name = name
	app.Author = author
	app.Version = version
	app.Usage = usage

	app.Action = action.NewAction().Action

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
