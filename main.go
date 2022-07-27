package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

var Token string

func main() {

	app := &cli.App{
		Name:  "lines-of-gode",
		Usage: "Check how many lines of code you have contributed.",
		Commands: []*cli.Command{
			&Auth,
			&Status,
			&Run,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
