# Lines of Gode

lines-of-gode is a CLI tool to count how many lines of code you have changed across your repositories (or those where you contributed). It's inspired by [Jothin-kumar/lines-of-code](https://github.com/Jothin-kumar/lines-of-code).

## Using it

To use it, you can follow two ways: **running a binary** or **building from source**.

NOTE: You _have_ to create a Personal Access Token for using `lines-of-gode`. It's mostly GitHub's security reasons and you can do it in less than five minutes, head up to [Settings](https://github.com/settings/tokens) and create a new token with some permissions (you may use it for other things, but `lines-of-gode` only requires access to your repositories and your user).

### Running a binary

This is the easiest one, just head to the **Releases** section and download the executable file. Once you have it, run it from your terminal with `./lines-of-gode`. Running it without arguments will display the help section.

### Building from source

The hardest way, but maybe the one you will have to follow if the binary is not available. You're going to need a Go compiler (preferably 1.13+, I use 1.18 for developing but anything over 1.13 _should_ be fine) and a little bit of experience with the terminal.

So, once you have the requirements, clone this repository with

```
git clone https://github.com/Ti7oyan/lines-of-gode
```

or if you use `github-cli`.

```
gh repo clone Ti7oyan/lines-of-code
```

Then, open the folder with a terminal (or just run `cd` to it), build the CLI with `go build` and... that's it!

But, if you want to use it across your OS, you will have to install it with `go install` and, if you didn't have the `GOPATH` set, run `export PATH=$PATH:$HOME/go/bin`.

With all of this, you can now run `lines of gode` from anywhere.

## Contributing

I'm just starting learning Go so, any tips for code style, refactor and general observations are really appreciated, as well as direct contribution!

Copyright (c) 2022 Ticiano Morvan
