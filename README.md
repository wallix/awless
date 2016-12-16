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
- [Concepts](#concepts)
	- [Security](#security)
	- [Infrastructures versioning](#infrastructures-versioning)
	- [RDF](#rdf)
- [Contributing](#contributing)
	- [Source install](#source-install)
	- [Test](#test)
	- [Run](#run)

<!-- /TOC -->

# Overview

Using Amazon Web Services through CLI is often discouraging, namely because of its complexity and API variety.
`awless` has been created with the idea to make this much easier, by using simple and intuitive commands and without editing any line of JSON or dealing with policies.
It brings a totally new approach to manage virtualized infrastructures through CLI.

Awless provides:

- A simple command to list virtualized resources (subnets, instances, groups, users, etc.): `awless list`
- Several output formats either human (list, trees,...) or machine readable (json): `--format`
- The ability to download a local snapshot of the infrastructure deployed in the remote cloud: `awless sync`
- The analysis of what has changed on the cloud since the last local snapshot: `awless diff`
- A git-based versioning of what has been deployed in the cloud: `awless ...`
- The simple and secure creation of virtual resources using pre-defined scenarios: `awless create`
- Commands and flags autocomplete for Unix/Linux's bash and zsh `awless completion`

# Install

Download the latest awless executable

- (Windows/Linux/macOS) [from Github](https://github.com/wallix/awless/releases/latest)
- (Windows/Linux/macOS) with go `go install github.com/wallix/awless`
- (macOS) using brew: `brew install awless`

# Getting started

Export in your shell session `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`


# Usage

Get main help from the CLI

    $ awless

Get help on a specific command

    $ awless help list

Show config

    $ awless config

Sync your infra

    $ awless sync

Diff your local infra to the remote one

    $ awless diff

List various items

    $ awless list users
    $ awless list policies
    $ awless list instances
    $ awless list vpcs
    $ awless list subnets

Show or delete the history of commands entered in awless

    $ awless history show
    $ awless history delete

Generate awless CLI completion code for bash or zsh

    $ awless completion bash
    $ awless completion zsh

# Bash and Zsh autocomplete

You can easily generate `awless` completion, either for bash or zsh, thanks to [cobra](https://github.com/spf13/cobra) (bash) and [kubernetes](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubectl/cmd/completion.go) (zsh).  

## Bash

For Mac OS X, with brew

    $ brew install bash-completion
    $ echo '[ -f /usr/local/etc/bash_completion ] && . /usr/local/etc/bash_completion\n' >> ~/.bashrc
    $ awless completion bash > /usr/local/etc/bash_completion.d/awless

For Ubuntu

    $ sudo apt-get install bash-completion
    $ sudo awless completion bash > /etc/bash_completion.d/awless

## Zsh

Test once with

    $ source <(awless completion zsh)

Or add to your ~/.zshrc

    $ echo 'source <(awless completion zsh)\n' >> ~/.zshrc

# Concepts

## Security

A major drawback of complexity is often the lack of security.
Managing virtualized infrastructures requires to well understand many network and security concepts, which may be time consuming.
As a result, these infrastructures are often configured with non-secure default parameters.
`awless` assists you in the creation of resources, ensuring that their security is properly configured.

## Infrastructures versioning

One main innovative feature of `awless` is the infrastructure versioning.
As soon as you sync your infrastructure, `awless` will keep a local history of what is deployed in your cloud.
This allows henceforth to track the changes in the cloud, keep an history of changes, show diff between versions, etc.

## RDF

Under the hood, `awless` processes and stores the cloud infrastructure in [RDF](https://www.w3.org/RDF/), using Google's [Badwolf](https://github.com/google/badwolf/) go library.
It allows to use advanced algorithms while being agnostic of the handled data.


# Contributing

## Source install

Until we inline dependencies fetch the following:

    $ go get github.com/aws/aws-sdk-go/...
    $ go get github.com/fatih/color
    $ go get github.com/boltdb/bolt
    $ go get github.com/spf13/viper
    $ go get github.com/spf13/cobra
    $ go get github.com/google/badwolf/...
		$ go get github.com/olekukonko/tablewriter

    $ go get github.com/wallix/awless

or install as a global executable

    $ go install github.com/wallix/awless

## Test

    $ cd awless
    $ go test -race ./...

## Run

    $ go run main.go list instances

or

    $ go build .
    $ ./awless list instances
