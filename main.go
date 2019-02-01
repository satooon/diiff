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

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "debug",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		a, err := action.NewAction(ctx.Bool("debug"))
		if err != nil {
			return err
		}
		return a.Action(ctx)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
