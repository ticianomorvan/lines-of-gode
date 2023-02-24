package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Stats struct {
	Additions int
	Deletions int
}

type Commit struct {
	gorm.Model
	ID        int64
	Sha       string
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

		db, err := gorm.Open(sqlite.Open("linesofgode.db"), &gorm.Config{})
		CheckError(err)

		db.AutoMigrate(&Commit{})

		token := GetToken()
		context, client := CreateClient(token)
		user := GetUser(context, client)

		fmt.Printf("Hello %v!, fetching your repositories now...\n", user.GetLogin())

		repositories := GetRepositories(context, client)

		for _, repository := range repositories {
			var repositoryStats = Stats{}

			commits := GetCommits(
				context,
				client,
				repository.GetOwner().GetLogin(),
				user.GetLogin(),
				repository.GetName(),
			)

			for _, commit := range commits {
				var commitStats = Stats{}

				sha := commit.GetSHA()

				files := GetCommitInfo(
					context,
					client,
					repository.GetOwner().GetLogin(),
					repository.GetName(),
					sha).Files
				for _, file := range files {
					if !CheckExceptions(file.GetFilename(), exceptions) {
						commitStats.Additions += file.GetAdditions()
						commitStats.Deletions += file.GetDeletions()
						}
					}
					repositoryStats.Additions += commitStats.Additions
					repositoryStats.Deletions += commitStats.Deletions
				}

				additions += repositoryStats.Additions
				deletions += repositoryStats.Deletions
				fmt.Printf("%v: +%v -%v\n", repository.GetFullName(), repositoryStats.Additions, repositoryStats.Deletions)
			}
		
		total = additions + deletions

		fmt.Printf(
			"You added around %v lines and deleted around %v lines for a total of %v lines changed.\n",
			additions, deletions, total,
		)

		return nil
	},
}

		/*

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
			"You have added %d lines and deleted %d lines for a total of %d lines changed across %d repositories.\n",
			additions, deletions, total, len(repositories),
		)

		return nil
	},
}
*/
