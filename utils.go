package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

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

type Credentials struct {
	Token string
}

func GetCredentials() Credentials {
	file, err := ioutil.ReadFile("credentials.toml")
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

	err = ioutil.WriteFile(filename, []byte(content), 0666)
	CheckError(err)

	fmt.Printf("Your personal access token has been stored at %s/%s\n", currentDirectory, filename)
}

func GetToken() string {
	var token string

	if token = GetCredentials().Token; token == "" {
		log.Fatal(`You didn't auth yourself yet. Run lines-of-gode auth --token "your_token"`)
	}

	return token
}
