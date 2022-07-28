package main

import (
	"database/sql"
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

/*
"Run" is the main command of the CLI.
It will scrap all the commits who's author is the authed user, and check for the
additions and deletions of the non-excepted files (e.g.: package-lock.json)

If the commit was previously checked, it will be read from a SQLite database, for
performance improvements. If it's a new commit, it will be requested to GitHub and then
stored in the database.

This will be do to every repository where the user made a commit, getting a per-repository report
and finally a general report of the user's contributions.
*/
var Run cli.Command = cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Usage:   "Count your additions and deletions of your contributions.",
	Action: func(ctx *cli.Context) error {
		var additions, deletions, total int

		db, err := sql.Open("sqlite3", "linesofgode.db")
		CheckError(err)

		defer db.Close()

		CreateCommitsTable(db)

		token := GetToken()
		context, client := CreateClient(token)
		user := GetUser(context, client).GetLogin()

		fmt.Printf("Hello %v!, fetching your repositories now...\n", user)

		repositories := GetRepositories(context, client, user)

		for _, repository := range repositories {
			var repositoryStats = Stats{}

			commits := GetCommits(context, client, user, repository.GetName())

			for _, commit := range commits {
				var commitStats = Stats{}

				if CheckAuthor(user, commit.GetAuthor().GetLogin()) {
					sha := commit.GetSHA()

					storedCommit, err := GetCommitBySha(db, sha)
					if err != nil {
						files := GetCommitInfo(context, client, user, repository.GetName(), sha).Files

						for _, file := range files {
							if !CheckExceptions(file.GetFilename(), exceptions) {
								commitStats.Additions += file.GetAdditions()
								commitStats.Deletions += file.GetDeletions()
							}
						}

						defer InsertCommit(db, sha, commitStats.Additions, commitStats.Deletions)
					} else {
						commitStats.Additions += storedCommit.Additions
						commitStats.Deletions += storedCommit.Deletions
					}
				}
				repositoryStats.Additions += commitStats.Additions
				repositoryStats.Deletions += commitStats.Deletions
			}
			fmt.Printf(
				"Repository: %q had %d additions and %d deletions\n",
				repository.GetName(), repositoryStats.Additions, repositoryStats.Deletions,
			)
			additions += repositoryStats.Additions
			deletions += repositoryStats.Deletions
		}

		total = additions + deletions
		fmt.Printf(
			"You have added %d lines and deleted %d lines for a total of %d lines changed across %d repositories.",
			additions, deletions, total, len(repositories),
		)

		return nil
	},
}
