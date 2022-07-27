package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

type Credentials struct {
	Token string
}

// Check

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckExceptions(filename string, exceptions []string) bool {
	for _, exception := range exceptions {
		if strings.Contains(filename, exception) {
			return true
		}
	}

	return false
}

func CheckAuthor(user, author string) bool {
	return user == author
}

// Github

func CreateClient(token string) (context.Context, *github.Client) {
	context := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient := oauth2.NewClient(context, tokenSource)
	client := github.NewClient(httpClient)
	return context, client
}

// Credentials

func GetCredentials() Credentials {
	file, err := os.ReadFile("credentials.toml")
	CheckError(err)

	var credentials Credentials
	_, err = toml.Decode(string(file), &credentials)
	CheckError(err)

	return credentials
}

func SaveCredentials(token string) {
	filename := "credentials.toml"
	content := fmt.Sprintf("token = %q", token)

	currentDirectory, err := os.Getwd()
	CheckError(err)

	err = os.WriteFile(filename, []byte(content), 0666)
	CheckError(err)

	fmt.Printf("Your personal access token has been stored at %s/%s\n", currentDirectory, filename)
}

// Getters

func GetUser(ctx context.Context, client *github.Client) *github.User {
	user, _, err := client.Users.Get(ctx, "")
	CheckError(err)

	return user
}

func GetRepositories(ctx context.Context, client *github.Client, user string) []*github.Repository {
	repositories, _, err := client.Repositories.List(ctx, user, nil)
	CheckError(err)

	return repositories
}

func GetCommits(ctx context.Context, client *github.Client, user, repository string) []*github.RepositoryCommit {
	commits, _, err := client.Repositories.ListCommits(ctx, user, repository, nil)
	CheckError(err)

	return commits
}

func GetCommitInfo(ctx context.Context, client *github.Client, user, repository, sha string) *github.RepositoryCommit {
	commit, _, err := client.Repositories.GetCommit(ctx, user, repository, sha, nil)
	CheckError(err)

	return commit
}

func GetToken() string {
	var token string

	if token = GetCredentials().Token; token == "" {
		log.Fatal(`You didn't auth yourself yet. Run lines-of-gode auth --token "your_token"`)
	}

	return token
}
