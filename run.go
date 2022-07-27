package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

type Stats struct {
	Additions int
	Deletions int
}

var exceptions []string = []string{
	"package-lock.json",
	"pnpm-lock.yaml",
	"gitignore",
	"eslint",
	"prettier",
	"yarn.lock",
	"node_modules",
}

var Run cli.Command = cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "Count your additions and deletions of your contributions.",
	Action: func(ctx *cli.Context) error {
		var additions, deletions, total int

		token := GetToken()
		context, client := CreateClient(token)
		user := GetUser(context, client).GetLogin()

		fmt.Printf("Hello %v!, fetching your repositories now...\n", user)

		repositories := GetRepositories(context, client, user)

		for _, repository := range repositories {
			var stats = Stats{}

			commits := GetCommits(context, client, user, repository.GetName())

			for _, commit := range commits {
				if CheckAuthor(user, commit.GetAuthor().GetLogin()) {
					sha := commit.GetSHA()
					files := GetCommitInfo(context, client, user, repository.GetName(), sha).Files

					for _, file := range files {
						if !CheckExceptions(file.GetFilename(), exceptions) {
							stats.Additions += file.GetAdditions()
							stats.Deletions += file.GetDeletions()
						}
					}
				}
			}
			additions += stats.Additions
			deletions += stats.Deletions
			fmt.Printf("Repository: %q had %d additions and %d deletions\n", repository.GetName(), stats.Additions, stats.Deletions)
		}

		total = additions + deletions

		fmt.Printf(
			"You have added %d lines and deleted %d lines for a total of %d lines changed across %d repositories.",
			additions, deletions, total, len(repositories),
		)

		return nil
	},
}
