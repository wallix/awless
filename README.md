# Awless

A nice, easy-to-use CLI for AWS

## Install

Until we inline dependencies fetch the following:

    $ go get github.com/aws/aws-sdk-go/...
    $ go get github.com/fatih/color
    $ go get github.com/boltdb/bolt
    $ go get github.com/spf13/viper
    $ go get github.com/spf13/cobra
    $ go get github.com/google/badwolf/...

    $ go get github.com/wallix/awless

or install as a global executable

    $ go install github.com/wallix/awless

## Test

    $ cd awless
    $ go test -race ./...

## Configure

Export in your shell session `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`

## Run

    $ go run main.go list instances

or

    $ go build .
    $ ./awless list instances

## Usage

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

## Awless bash and zsh completion

You can easily generate `awless` completion, either for bash or zsh, thanks to [cobra](https://github.com/spf13/cobra) (bash) and [kubernetes](https://github.com/kubernetes/kubernetes/blob/master/pkg/kubectl/cmd/completion.go) (zsh).  

### Bash

For Mac OS X, with brew

    $ brew install bash-completion
    $ echo '[ -f /usr/local/etc/bash_completion ] && . /usr/local/etc/bash_completion\n' >> ~/.bashrc
    $ awless completion bash > /usr/local/etc/bash_completion.d/awless

For Ubuntu

    $ sudo apt-get install bash-completion
    $ sudo awless completion bash > /etc/bash_completion.d/awless

### Zsh

    $ source <(awless completion zsh)
