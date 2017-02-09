`awless` is a fast, powerful and easy-to-use command line interface (CLI) to manage Amazon Web Services.

# Why awless

`awless` has been created with the idea to run frequent actions easily by using simple commands, smart defaults, security best practices and runnable/scriptable templates for resource creations (see [`awless` templates](https://github.com/wallix/awless/wiki/Templates)).

There is no need to edit manually any line of JSON, deal with policies, etc.
`awless` brings a new approach to manage virtualized infrastructures through CLI.

# Overview

- Clear and easy listing of cloud resources (subnets, instances, groups, users, etc.) on AWS EC2, IAM and S3: `awless list`
- Multiple output formats either human (table, trees, ...) or machine readable (csv, json, ...): `--format`
- Explore a resource given only an *id* or name (properties, relations, dependencies, ...): `awless show`
- Creation, update and deletion (CRUD) of cloud resources and complex infrastructure with smart defaults through powerful awless templates: `awless run my-awless-templates/create_my_infra.txt`
- Powerful CRUD CLI onliner (integrated in our awless templating engine) with: `awless create instance ...`, `awless create vpc ...`, `awless attach policy ...`
- Easy listing or revert of resources creation: `awless revert`
- A local history and versioning of the changes that occurred in your cloud: `awless history`
- CLI autocompletion for Unix/Linux's bash and zsh `awless completion`

# Design concepts

1. [RDF](https://www.w3.org/TR/rdf11-concepts/) is used internally to sync and model cloud resources locally. This permits a good flexibility in modeling while still allowing for DAG (Directed Acyclic Graph) properties and classic graph/tree traversal.
2. Awless templates define a basic DSL (Domain Specific Language) for managing cloud resources. Templates are parsed against a [PEG (parsing expression grammar)](https://en.wikipedia.org/wiki/Parsing_expression_grammar) allowing for robust parsing, AST building/validation and execution of this AST through given official cloud drivers (ex: aws-sdk-go for AWS). More details on awless templates on the [wiki](https://github.com/wallix/awless/wiki/Templates).

# Install

Choose one of the following options:

1. Download the latest `awless` executable (Windows/Linux/macOS) [from Github](https://github.com/wallix/awless/releases/latest)
2. Build the source with Go: Run `go get github.com/wallix/awless` (if `go` is already installed, on Windows/Linux/macOS)
<!--- 3. On macOS, use [homebrew](http://brew.sh):  `brew install awless` -->

# Getting started

## Setup your AWS account with `awless`

You basically need your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` exported in your environment.

If you have previously used `aws` CLI or `aws-shell`, you don't need to do anything! Your credentials will be automatically loaded by `awless` from the `~/.aws/credentials` folder.

Otherwise, get your AWS credentials from [IAM console](https://console.aws.amazon.com/iam/home?#home).
Then, you can download and store them to `~/.aws/credentials` (Unix) or `%UserProfile%\.aws\credentials` (Windows).

For more options, see [Installation (wiki)](https://github.com/wallix/awless/wiki/Installation#setup-your-aws-account-with-awless).

## Setup shell autocompletion

Awless has commands, subcommands and flag completion. It becomes really useful for CRUD onliner when managing resources for example.

Read the wiki page for setting autocompletion for [bash](https://github.com/wallix/awless/wiki/Setup-Autocomplete#bash) or [zsh](https://github.com/wallix/awless/wiki/Setup-Autocomplete#zsh).

## Disclaimer

Awless allows for easy resource creation with your cloud provider; We will not be responsible for any cloud costs incurred (even if you create a million instances using awless templates).

We also collect a few anonymous data (CLI errors, most frequently used commands and count of resources).

## First `awless` commands

`awless` works by performing commands, which query either the AWS services or a local snapshot of the cloud services.

### Listing resources

You can list various resources:

    $ awless list instances
    $ awless list users --format csv
    $ awless list roles --sort name,id

Listing resources by default performs queries directly to AWS. If you want, you can also query the local snapshot:

    $ awless list subnets --local

See the [manual](https://github.com/wallix/awless/wiki/Commands#awless-list) for a complete reference of `awless list`.

### Creating, Updating and Deleting resources

`awless` provides a powerful template system to interact with cloud infrastructures.

`awless` templates can be used through onliner shortcut commands:

    awless create instance
    awless delete subnet id=subnet-12345678
    awless attach volume id=vol-12345678 instance=i-12345678

See [templates commands (wiki)](https://github.com/wallix/awless/wiki/Templates#Commands) for more commands.

You can also run an `awless` template from a predefined template file with:

    awless run awless-templates/create_instance_ssh.awless

Note that you can get inspired with pre-existing templates from the dedicated git repository: https://github.com/wallix/awless-templates. See [templates (wiki)](https://github.com/wallix/awless/wiki/Templates) for more details about `awless` templates.

### Reverting commands

Each `awless` command that changes the cloud infrastructure is associated with an unique *id* referencing the (un)successful actions.

To list the last actions you have run on your cloud, run

    awless revert -l

Then, you can revert a command with

    awless revert -i 01B89ZY529E5D7WKDTQHFC0RPA # for now, revert only resource creations

### SSH

You can directly ssh to an instance with:

        $ awless ssh i-abcd1234
        $ awless ssh ubuntu@i-abcd1234

In the first case, note that `awless` can work out the default ssh user to use given a cloud (ex: `ec2` for AWS)

### Aliasing

When it makes sense we provide the concept of *alias*. Cloud resources ids can be a bit cryptic. An alias is just an already existing name of a resource. Given a alias we resolve the proper resource id. For instance:

        $ awless ssh my-instance         # ssh to the instance using its name. Behind the scene awless resolve the id
        $ awless delete id=@my-instance  # delete an instance using its name

# About

`awless` is an open source project created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and François-Xavier Aguessy at Wallix.
`awless` is released under the Apache License and sponsored by [Wallix](https://github.com/wallix).

Contributors are welcome! Please head to [Contributing (wiki)](https://github.com/wallix/awless/wiki/Contributing) to learn more.


