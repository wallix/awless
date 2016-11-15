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

Sync your infra

    $ awless sync

Diff your local infra to the remote one

    $ awless diff
        {
         "RegionMatch": true,
         "ExtraVpcs": null,
         "MissingVpcs": null,
         "ExtraSubnets": [
          "subnet-w8eftweifgw"
         ],
         "MissingSubnets": [
          "subnet-267d517f"
         ],
         "ExtraInstances": [
          "i-uevfwiefbow"
         ],
         "MissingInstances": [
          "i-ad86f625"
         ]
        }

List various items

    $ awless list users
    $ awless list policies
    $ awless list instances
    $ awless list vpcs
    $ awless list subnets
