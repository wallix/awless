`awless` is a nice, easy-to-use command line interface (CLI) to manage Amazon Web Services.

`awless` is a straightforward replacement for AWS default CLI.

# Why awless

Using Amazon Web Services through the default CLI is often discouraging, namely because of its complexity and API variety.

`awless` has been created with the idea to run the most frequent actions easily, by using simple commands (like `git`) or pre-defined scripts.

There is no need to edit manually any line of JSON, deal with policies, etc.
`awless` brings a new approach to manage virtualized infrastructures through CLI.

# What awless can do

- A simple command to list virtualized resources (subnets, instances, groups, users, etc.): `awless list`
- Several output formats either human (tables, trees) or machine readable (csv, json): `--format`
- The ability to create a local snapshot of the infrastructure deployed in the remote cloud: `awless sync`
- The analysis of what has changed on the cloud since the last local snapshot: `awless diff`
- A git-based versioning of what has been deployed in the cloud: `awless show revisions`
- The simple and secure creation of virtual resources using pre-defined scenarios: `awless create`
- Commands autocomplete for Unix/Linux's bash and zsh `awless completion`

# Install

Choose one of the following options:

1. Download the latest `awless` executable (Windows/Linux/macOS) [from Github](https://github.com/wallix/awless/releases/latest)
2. Build the source with Go: Run `go get github.com/wallix/awless` (if `go` is already installed, on Windows/Linux/macOS)
3. On macOS, use [homebrew](http://brew.sh):  `brew install awless`

# Getting started

## Setup your AWS account with `awless`

If you have previously used `aws` CLI or `aws-shell`, you don't need to do anything! Your credentials will be automatically loaded by `awless` from the `~/.aws/credentials` folder.

Otherwise, get your AWS credentials from [IAM console](https://console.aws.amazon.com/iam/home?#home).
Then, you can either download and store them to `~/.aws/credentials` (Unix) or `%UserProfile%\.aws` (Windows).

For more options, see [Installation](https://github.com/wallix/awless/wiki/Installation#) in the wiki.

## Setup shell autocompletion

Read the wiki page for setting autocompletion for [bash]() or [zsh]().

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

See the [manual](https://github.com/wallix/awless/wiki/awless-Commands#List) for a complete reference of `awless list`.

### Updating the local snapshot

<!-- WHY!!! -->

You can synchronize your local snapshot with the current cloud infrastructure on AWS using:

    $ awless sync

Once the sync is done, changes (either to the local model or to the cloud infrastructure) can be tracked easily:

    $ awless diff

### Creating a new resource

<!-- TODO -->

### Much more

See the [wiki](https://github.com/wallix/awless/wiki/awless-Commands#List) for a more complete reference.

# About

`awless` is an open source project created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and François-Xavier Aguessy at Wallix.
`awless` is released under the Apache License and sponsored by [Wallix](https://github.com/wallix).

Contributors are welcome! Please head to [Contributing (wiki)](https://github.com/wallix/awless/wiki/Contributing) to learn more.


