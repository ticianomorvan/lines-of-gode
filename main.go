package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/v45/github"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

type Credentials struct {
	Token string
}

func CreateClient(personalAccessToken string) (context.Context, *github.Client) {
	ctx := context.Background()
	staticToken := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: personalAccessToken},
	)

	httpClient := oauth2.NewClient(ctx, staticToken)
	client := github.NewClient(httpClient)
	return ctx, client
}

func GetCredentials() Credentials {
	file, err := os.ReadFile("credentials.toml")
	if err != nil {
		log.Fatal(err)
	}

	var credentials Credentials
	if _, err = toml.Decode(string(file), &credentials); err != nil {
		log.Fatal(err)
	}

	return credentials
}

func SaveCredentials(personalAccessToken string) {
	filename := "credentials.toml"
	content := fmt.Sprintf("token = %q", personalAccessToken)

	currentDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(filename, []byte(content), 0666); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Your personal access token has been stored successfully at %q\n", currentDirectory)
}

func CheckString(s string, exceptions []string) bool {
	for _, exception := range exceptions {
		if strings.Contains(s, exception) {
			return true
		}
	}

	return false
}

func main() {
	var personalToken string

	app := &cli.App{
		Name:  "lines-of-gode",
		Usage: "Check how many lines of code you have contributed.",
		Commands: []*cli.Command{
			{
				Name:    "auth",
				Aliases: []string{"a"},
				Usage:   "Auth yourself with your personal access token and save it for future tasks.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "token",
						Required:    true,
						Value:       "YOUR_ACCESS_TOKEN",
						Usage:       "Your personal access token. Refer to https://github.com/settings/tokens",
						Destination: &personalToken,
					},
				},
				Action: func(ctx *cli.Context) error {
					context, client := CreateClient(personalToken)
					user, _, err := client.Users.Get(context, "")
					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("Successfully authed as %q (%v)\n", user.GetLogin(), user.GetName())
					SaveCredentials(personalToken)
					return nil
				},
			},
			{
				Name:    "status",
				Aliases: []string{"s"},
				Usage:   "Check current authentication status.",
				Action: func(ctx *cli.Context) error {
					token := GetCredentials().Token
					if token == "" {
						log.Fatal(`You didn't auth yourself yet. Run lines-of-gode auth --token "your_token"`)
					}
					context, client := CreateClient(token)
					user, _, err := client.Users.Get(context, "")

					if err != nil {
						log.Fatal(err)
					}

					fmt.Printf("You are authed as %q (%v).\n", user.GetLogin(), user.GetName())

					return nil
				},
			},
			{
				Name:    "run",
				Aliases: []string{"r"},
				Usage:   "Count your additions and deletions on GitHub",
				Action: func(ctx *cli.Context) error {
					var additions int
					var deletions int
					var total int

					token := GetCredentials().Token
					if token == "" {
						log.Fatal(`You didn't auth yourself yet. Run lines-of-gode auth --token "your_token"`)
					}
					context, client := CreateClient(token)
					user, _, err := client.Users.Get(context, "")
					if err != nil {
						log.Fatal(err)
					}

					userName := user.GetLogin()

					fmt.Printf("Hello %v!, fetching your repositories now...\n", userName)

					repositories, _, err := client.Repositories.List(context, userName, nil)
					if err != nil {
						log.Fatal(err)
					}

					for _, repository := range repositories {
						var repositoryAdditions int
						var repositoryDeletions int
						commits, _, err := client.Repositories.ListCommits(context, userName, repository.GetName(), nil)
						if err != nil {
							log.Fatal(err)
						}

						for _, commit := range commits {
							if commit.GetAuthor().GetLogin() == userName {
								sha := commit.GetSHA()
								information, _, err := client.Repositories.GetCommit(context, userName, repository.GetName(), sha, nil)

								if err != nil {
									log.Fatal(err)
								}

								for _, file := range information.Files {
									exceptions := []string{
										"package-lock.json",
										"pnpm-lock.yaml",
										"gitignore",
										"eslint",
										"prettier",
										"yarn.lock",
										"node_modules",
									}

									if !CheckString(file.GetFilename(), exceptions) {
										repositoryAdditions += file.GetAdditions()
										repositoryDeletions += file.GetDeletions()
									}
								}
							}
						}
						additions += repositoryAdditions
						deletions += repositoryDeletions
						fmt.Printf("Repository: %q had %d additions and %d deletions\n", repository.GetName(), repositoryAdditions, repositoryDeletions)
					}

					total = additions + deletions

					fmt.Printf(
						"You have added %d lines and deleted %d lines for a total of %d lines changed across %d repositories.\n",
						additions, deletions, total, len(repositories),
					)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
