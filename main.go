package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"monot/commands"
)

func main() {
	app := &cli.App{
		Name:  "monot",
		Usage: "Simple Monorepo Manager",
		Action: func(*cli.Context) error {
			fmt.Println("Welcome to Monot!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "init",
				Action:  commands.InitCmd,
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "run run-local",
				Action:  commands.RunCmd,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
