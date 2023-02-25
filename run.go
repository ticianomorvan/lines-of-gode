package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Commit struct {
	gorm.Model
	Sha       string
	Additions int
	Deletions int
}

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

			commits := GetCommits(&SetupGetCommits{
				ctx: context,
				client: client,
				owner: repository.GetOwner().GetLogin(),
				user: user.GetLogin(),
				repository: repository.GetName(),
			})

			for _, commit := range commits {
				var commitStats = Stats{}

				sha := commit.GetSHA()

				files := GetCommitInfo(&SetupGetCommitInfo{
					ctx: context,
					client: client,
					owner: repository.GetOwner().GetLogin(),
					repository: repository.GetName(),
					sha: sha,
				}).Files

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