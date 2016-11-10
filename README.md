# Awless

A nice, easy-to-use CLI for AWS

## Install

    go install github.com/wallix/awless

## Configure

Export in your shell session `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`

## Usage

Get main help from the CLI

    $ awless

Get help on a specific command

    $ awless help list

Show config

    $ awless config

Discover your infra

    $ awless sync
     Region: eu-west-1, 1 VPC(s)
        1. VPC vpc-00b68c65, 3 subnet(s)
            1. Subnet subnet-0c41ad68, 1 instance(s)
                1. Instance i-ad86f625
            2. Subnet subnet-f5c9dd82, 0 instance(s)
            3. Subnet subnet-267d517f, 0 instance(s)

List various items

    $ awless list users
    $ awless list policies
    $ awless list instances
    $ awless list vpcs
    $ awless list subnets
