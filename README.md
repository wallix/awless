`awless` is a nice, easy-to-use command line interface (CLI) to manage Amazon Web Services.
It is designed to keep simple and secure the display, creation, update and deletion of virtual resources.

<!-- TOC depthFrom:1 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Overview](#overview)
- [Install](#install)
- [Getting started](#getting-started)
- [Usage](#usage)
- [Bash and Zsh autocomplete](#bash-and-zsh-autocomplete)
	- [Bash](#bash)
	- [Zsh](#zsh)
- [Contributing](#contributing)
	- [Source install](#source-install)
	- [Test](#test)
	- [Build and Run](#build-and-run)

<!-- /TOC -->

# Overview

Using Amazon Web Services through CLI is often discouraging, namely because of its complexity and API variety.
`awless` has been created with the idea to run the most frequent actions much easier, by using simple and intuitive commands or pre-defined scripts, without editing any line of JSON or dealing with policies.
It brings a new approach to manage virtualized infrastructures through CLI.

Awless provides:

- A simple command to list virtualized resources (subnets, instances, groups, users, etc.): `awless list`
- Several output formats either human (tables, trees) or machine readable (csv, json): `--format`
- The ability to create a local snapshot of the infrastructure deployed in the remote cloud: `awless sync`
- The analysis of what has changed on the cloud since the last local snapshot: `awless diff`
- A git-based versioning of what has been deployed in the cloud: `awless show revisions`
- The simple and secure creation of virtual resources using pre-defined scenarios: `awless create`
- Commands autocomplete for Unix/Linux's bash and zsh `awless completion`

# Install

Download the latest awless executable

- (Windows/Linux/macOS) [from Github](https://github.com/wallix/awless/releases/latest)
- (Windows/Linux/macOS) with go `go get github.com/wallix/awless`
- (macOS) using brew: `brew install awless`

# Getting started

If you have previously used `aws` CLI or `aws-shell`, you don't need to do anything! Your credentials will be automatically loaded by `awless` from the `~/.aws/credentials` folder.

Otherwise, get your AWS credentials from [IAM console](https://console.aws.amazon.com/iam/home?#home).
Then, you can either download and store them to `~/.aws/credentials` (Unix) or `%UserProfile%\.aws` (Windows).
Or, export the environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` in your shell session:

	export AWS_ACCESS_KEY_ID=AKIAIOSFODNN7EXAMPLE
	export AWS_SECRET_ACCESS_KEY=wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
	
You can find more information about how to get your AWS credentials on [Amazon Web Services user guide](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-set-up.html).

# Usage

Get main help from the CLI

    $ awless

Get help on a specific command

    $ awless help list

Show config

    $ awless config

Sync your cloud

    $ awless sync

Diff your local cloud to the remote one

    $ awless diff

List various resources

		$ awless list instances
    $ awless list users --format csv # List results as CSV
    $ awless list subnets --local # List the subnets of your local cloud snapshot
		$ awless list roles --sort name,id # List roles sorted by name and id
		...
		
Show details of resources

		$ awless show instance i-abcd1234
    $ awless show user AKIAIOSFODNN7EXAMPLE --local # List the subnets of your local cloud snapshot
		...
		
Show the history of changes to your cloud
 	
		$ awless show revisions # Last 10 revisions (= sync)
		$ awless show revisions -n 3 --properties # Last 3 revisions with their properties
		$ awless show revisions --group-by-week # Group changes by week
		
Connect to instances via SSH

		$ awless ssh i-abcd1234 # Connect to instance by ssh using your keys stored in ~/.awless/keys
		$ awless ssh ubuntu@i-abcd1234 # Use a specific user

Show or delete the history of commands entered in awless

    $ awless history show
    $ awless history delete
		
Display the current AWS user and account

	$ awless whoami

Generate `awless` completion code for bash or zsh

    $ awless completion bash
    $ awless completion zsh

# Bash and Zsh autocomplete

You can easily generate `awless` completion, either for bash or zsh.

## Bash

For Mac OS X, with brew

    $ brew install bash-completion
    $ echo '[ -f /usr/local/etc/bash_completion ] && . /usr/local/etc/bash_completion\n' >> ~/.bashrc
    $ awless completion bash > /usr/local/etc/bash_completion.d/awless

For Ubuntu

    $ sudo apt-get install bash-completion
    $ awless completion bash | sudo tee /etc/bash_completion.d/awless > /dev/null

## Zsh

Test once with

    $ source <(awless completion zsh)

Or add to your ~/.zshrc

    $ echo 'source <(awless completion zsh)\n' >> ~/.zshrc

# Contributing

## Source install

You need first need to [install Go](https://golang.org/doc/install) to build `awless`.
Then, to download awless sources

    go get github.com/wallix/awless

## Test

    $ cd awless
    $ go test -race -cover ./...
		
or (faster) if you have `govendor` (`get -u github.com/kardianos/govendor`)

		$ govendor test +local

## Build and Run

To run from your local pulled repo

    $ go run main.go list instances

To build a local executable

		$ go build .
		
		or
		
		$ GOARCH=amd64 GOOS=linux go build . # Cross compilation to Linux amd64
		
Then

		$ ./awless list instances
		
Or if you want to install to your $GOPATH/bin

		$ go install
		$ awless list instances
