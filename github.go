package main

import (
	"context"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

func CreateClient(token string) (context.Context, *github.Client) {
	context := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	httpClient := oauth2.NewClient(context, tokenSource)
	client := github.NewClient(httpClient)
	return context, client
}


func GetUser(ctx context.Context, client *github.Client) *github.User {
	user, _, err := client.Users.Get(ctx, "")
	CheckError(err)

	return user
}

func GetRepositories(ctx context.Context, client *github.Client) []*github.Repository {
	repositories, _, err := client.Repositories.List(ctx, "", nil)
	CheckError(err)

	return repositories
}

type SetupGetCommits struct {
	ctx context.Context
	client *github.Client
	owner string
	user string
	repository string
}

func GetCommits(settings *SetupGetCommits) []*github.RepositoryCommit {
	options := github.CommitsListOptions{
		Author: settings.user,
	}

	commits, _, err := settings.client.Repositories.ListCommits(
		settings.ctx,
		settings.owner,
		settings.repository,
		&options,
	)
	CheckError(err)

	return commits
}

type SetupGetCommitInfo struct {
	ctx context.Context
	client *github.Client
	owner string
	repository string
	sha string
}

func GetCommitInfo(settings *SetupGetCommitInfo) *github.RepositoryCommit {
	commit, _, err := settings.client.Repositories.GetCommit(
		settings.ctx,
		settings.owner,
		settings.repository,
		settings.sha,
		nil,
	)

	CheckError(err)

	return commit
}
