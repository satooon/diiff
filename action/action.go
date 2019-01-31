package action

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

// Action interface
type Action interface {
	Action(ctx *cli.Context) error
}

type action struct{}

// NewAction return *action
func NewAction() Action {
	return &action{}
}

// action is cli action
func (a *action) Action(ctx *cli.Context) error {
	log.Println("CLI Args", ctx.Args())

	path, err := getPath(ctx.Args().Get(0))
	if err != nil {
		return err
	}
	log.Println("Path", path)

	filePath := NewFilePath()
	if err := filePath.Scan(path); err != nil {
		return err
	}

	return nil
}

func getPath(p string) (string, error) {
	if len(p) > 0 {
		return p, nil
	}
	return os.Getwd()
}
