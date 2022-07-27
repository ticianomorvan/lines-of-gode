package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

/*
The auth command tries to authenticate an user against the GitHub API v3 using
their Personal Access Token. If it's approved, it will save their credentials for future
uses.
*/
var Auth cli.Command = cli.Command{
	Name:    "auth",
	Aliases: []string{"a"},
	Usage:   "Auth yourself with your personal access token and store it locally.",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Required:    true,
			Usage:       "Your personal access token",
			Value:       "",
			Destination: &Token,
		},
	},
	Action: func(ctx *cli.Context) error {
		context, client := CreateClient(Token)
		user := GetUser(context, client)

		fmt.Printf("Successfully authed as %q (%v)\n", user.GetLogin(), user.GetName())
		SaveCredentials(Token)
		return nil
	},
}
