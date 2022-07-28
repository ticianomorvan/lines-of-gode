package main

import (
	"context"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/stretchr/testify/suite"
)

type UtilsTestSuite struct {
	suite.Suite
	token   string
	context context.Context
	client  *github.Client
	user    *github.User
}

func (suite *UtilsTestSuite) SetupSuite() {
	var currentToken = GetToken()
	context, client := CreateClient(currentToken)
	var currentUser = GetUser(context, client)

	suite.token = currentToken
	suite.context = context
	suite.client = client
	suite.user = currentUser
}

func (suite *UtilsTestSuite) TestCheckExceptions() {
	testFiles := map[string]bool{
		"package.json":       false,
		"App.tsx":            false,
		".gitignore":         true,
		"database.db":        false,
		"pnpm-lock.yaml":     true,
		"node_modules/babel": true,
		"vite.config.ts":     false,
		"main.go":            false,
	}

	for name, expected := range testFiles {
		suite.Equal(CheckExceptions(name, exceptions), expected)
	}
}

func (suite *UtilsTestSuite) TestCheckAuthor() {
	suite.True(CheckAuthor(suite.user.GetLogin(), suite.user.GetLogin()))
	suite.False(CheckAuthor(suite.user.GetLogin(), "Linus Torvalds"))
}

func (suite *UtilsTestSuite) TestSaveCredentials() {
	SaveCredentials(suite.token)
	suite.Contains(GetToken(), suite.token)
}

func (suite *UtilsTestSuite) TestGetCredentials() {
	suite.NotEmpty(GetCredentials())
}

func (suite *UtilsTestSuite) TestGetUser() {
	user := GetUser(suite.context, suite.client)
	suite.Equal(user, suite.user)
}

func (suite *UtilsTestSuite) TestGetRepositories() {
	repositories := GetRepositories(suite.context, suite.client, suite.user.GetLogin())
	suite.NotEmpty(repositories)
}

func (suite *UtilsTestSuite) TestGetCommits() {
	repository := GetRepositories(suite.context, suite.client, suite.user.GetLogin())
	commits := GetCommits(suite.context, suite.client, suite.user.GetLogin(), repository[0].GetName())
	suite.NotEmpty(commits)
}

func (suite *UtilsTestSuite) TestGetCommitInfo() {
	repository := GetRepositories(suite.context, suite.client, suite.user.GetLogin())
	commits := GetCommits(suite.context, suite.client, suite.user.GetLogin(), repository[0].GetName())
	info := GetCommitInfo(suite.context, suite.client, suite.user.GetLogin(), repository[0].GetName(), commits[0].GetSHA())
	suite.NotEmpty(info)
}

func (suite *UtilsTestSuite) TestGetToken() {
	token := GetToken()
	suite.NotEmpty(token)
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}
