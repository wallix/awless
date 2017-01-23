`awless` is a nice, easy-to-use command line interface (CLI) to manage Amazon Web Services.

# Why awless

`awless` has been created with the idea to run the most frequent actions easily by using simple commands, smart defaults, security best practices and runnable/scriptable templates for resource creations (see `awless` templates).

There is no need to edit manually any line of JSON, deal with policies, etc.
`awless` brings a new approach to manage virtualized infrastructures through CLI.

# Overview

- Clear and easy listing of virtualized resources (subnets, instances, groups, users, etc.): `awless list`
- Multiple output formats either human (tabular, trees, ...) or machine readable (csv, json, ...): `--format`
- Create a local snapshot of the infrastructure (ec2, iam, s3) deployed in the remote cloud: `awless sync`
- Show what has changed on the cloud since the last local snapshot: `awless diff`
- A local history and versioning of the snapshots: `awless show revisions`
- Creation of cloud resources (instances, groups, users, policies) with smart and secure default through powerful awless templates
- CLI autocompletion for Unix/Linux's bash and zsh `awless completion`

# Install

Choose one of the following options:

1. Download the latest `awless` executable (Windows/Linux/macOS) [from Github](https://github.com/wallix/awless/releases/latest)
2. Build the source with Go: Run `go get github.com/wallix/awless` (if `go` is already installed, on Windows/Linux/macOS)
3. On macOS, use [homebrew](http://brew.sh):  `brew install awless`

# Getting started

## Setup your AWS account with `awless`

If you have previously used `aws` CLI or `aws-shell`, you don't need to do anything! Your credentials will be automatically loaded by `awless` from the `~/.aws/credentials` folder.

Otherwise, get your AWS credentials from [IAM console](https://console.aws.amazon.com/iam/home?#home).
Then, you can either download and store them to `~/.aws/credentials` (Unix) or `%UserProfile%\.aws\credentials` (Windows).

For more options, see [Installation (wiki)](https://github.com/wallix/awless/wiki/Installation#setup-your-aws-account-with-awless).

## Setup shell autocompletion

Read the wiki page for setting autocompletion for [bash](https://github.com/wallix/awless/wiki/Setup-Autocomplete#bash) or [zsh](https://github.com/wallix/awless/wiki/Setup-Autocomplete#zsh).

## First `awless` commands

`awless` works by performing commands, which query either the AWS infrastructure or the local snapshot of the infrastructure.

### Listing resources

You can list various resources:

    $ awless list instances
    $ awless list users --format csv
    $ awless list roles --sort name,id

Listing resources by default performs queries directly to AWS.
If you want, you can also query the local snapshot:

    $ awless list subnets --local

See the [manual](https://github.com/wallix/awless/wiki/Commands#awless-list) for a complete reference of `awless list`.

### Updating the local snapshot

<!-- WHY!!! -->

You can synchronize your local snapshot with the current cloud infrastructure on AWS using:

    $ awless sync

Once the sync is done, changes (either to the local model or to the cloud infrastructure) can be tracked easily:

    $ awless diff

### Creating a new resource

<!-- TODO -->

### Much more

		$ awless show revisions --group-by-week
		$ awless ssh ubuntu@i-abcd1234

See [commands (wiki)](https://github.com/wallix/awless/wiki/Commands) for a more complete reference.

# About

`awless` is an open source project created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and François-Xavier Aguessy at Wallix.
`awless` is released under the Apache License and sponsored by [Wallix](https://github.com/wallix).

Contributors are welcome! Please head to [Contributing (wiki)](https://github.com/wallix/awless/wiki/Contributing) to learn more.


