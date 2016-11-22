# Awless

A nice, easy-to-use CLI for AWS

## Install

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
      Extras:

      Missings:
        /subnet<subnet-0c41ad68>	"parent_of"@[]	/instance<i-56adc1dd>

List various items

    $ awless list users
    $ awless list policies
    $ awless list instances
    $ awless list vpcs
    $ awless list subnets
