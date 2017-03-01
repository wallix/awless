[![Build Status](https://api.travis-ci.org/wallix/awless.svg?branch=master)](https://travis-ci.org/wallix/awless)
[![Go Report Card](https://goreportcard.com/badge/github.com/wallix/awless)](https://goreportcard.com/report/github.com/wallix/awless)

`awless` is a fast, powerful and easy-to-use command line interface (CLI) to manage Amazon Web Services.

[Twitter](http://twitter.com/awlessCLI) | [Wiki](https://github.com/wallix/awless/wiki)

# Why awless

`awless` will help you

- run frequent actions by using simple commands
- easily explore your infrastructure and cloud resources inter relations via CLI
- ensure smart defaults & security best practices
- manage resources through robust runnable & scriptable templates (see [`awless` templates](https://github.com/wallix/awless/wiki/Templates))
- explore, analyse and query your infrastructure **offline**
- explore, analyse and query your infrastructure **through time**

`awless` brings a new approach to manage AWS infrastructures through CLI.

# Overview

![video of a few `awless` commands](https://raw.githubusercontent.com/wiki/wallix/awless/gif/awless-demo.gif "video of a few `awless` commands")

- Clear and easy listing of multi-region cloud resources (subnets, instances, groups, users, etc.) on AWS EC2, IAM and S3: `awless list`
- Multiple output formats either human (table, trees, ...) or machine readable (csv, json, ...): `--format`
- Explore a resource given only an *id* or name (properties, relations, dependencies, ...): `awless show`
- Creation, update and deletion (CRUD) of cloud resources and complex infrastructure with smart defaults through powerful awless templates: `awless run my-awless-templates/create_my_infra.txt`
- Powerful CRUD CLI one-liner (integrated in the awless templating engine) with: `awless create instance ...`, `awless create vpc ...`, `awless attach policy ...`
- Easy reporting of all the CLI template executions: `awless log`
- Revert of executed templates and resources creation: `awless revert`
- Aliasing of resources through their natural name so you don't have to always use cryptic ids that are impossible to remember
- Inspectors are small CLI utilities to run analysis on your cloud resources graphs: `awless inspect`
- Manual sync mode to fetch & store resources locally. Then query & inspect your cloud offline: `awless sync`
- CLI autocompletion for Unix/Linux's bash and zsh `awless completion`
- [*IN PROGRESS*] A local history and versioning of the changes that occurred in your cloud: `awless history`

# Design concepts

1. [RDF](https://www.w3.org/TR/rdf11-concepts/) is used internally to sync and model cloud resources locally. This permits a good flexibility in modeling while still allowing for DAG (Directed Acyclic Graph) properties and classic graph/tree traversal.
2. Awless templates define a basic DSL (Domain Specific Language) for managing cloud resources. Templates are parsed against a [PEG (parsing expression grammar)](https://en.wikipedia.org/wiki/Parsing_expression_grammar) allowing for robust parsing, AST building/validation and execution of this AST through given official cloud drivers (ex: aws-sdk-go for AWS). More details on awless templates on the [wiki](https://github.com/wallix/awless/wiki/Templates).

# Install

Choose one of the following options:

1. On macOS, use [homebrew](http://brew.sh):  `brew tap wallix/awless; brew install awless`
2. With `curl` (macOS/Linux), run: `curl https://raw.githubusercontent.com/wallix/awless/master/getawless.sh | bash`
3. Download the latest `awless` binaries (Windows/Linux/macOS) [from Github](https://github.com/wallix/awless/releases/latest)
4. If you have Golang already installed, build the source with: `go get github.com/wallix/awless`

# Getting started

## Setup your AWS account with `awless`

You basically need your `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` exported in your environment.

If you have previously used `aws` CLI or `aws-shell`, you don't need to do anything! Your credentials will be automatically loaded by `awless` from the `~/.aws/credentials` folder.

Otherwise, get your AWS credentials from [IAM console](https://console.aws.amazon.com/iam/home?#home).
Then, you can download and store them to `~/.aws/credentials` (Unix) or `%UserProfile%\.aws\credentials` (Windows).

For more options, see [Installation (wiki)](https://github.com/wallix/awless/wiki/Installation#setup-your-aws-account-with-awless).

## Setup shell autocompletion

Awless has commands, subcommands and flag completion. It becomes really useful for CRUD oneliner when managing resources for example.

Read the wiki page for setting autocompletion for [bash](https://github.com/wallix/awless/wiki/Setup-Autocomplete#bash) or [zsh](https://github.com/wallix/awless/wiki/Setup-Autocomplete#zsh).

## Disclaimer

Awless allows for easy resource creation with your cloud provider; We will not be responsible for any cloud costs incurred (even if you create a million instances using awless templates).

## First `awless` commands

`awless` works by performing commands, which query either the AWS services or a local snapshot of the cloud services.

### Listing resources

You can list various resources:

    awless list buckets
    awless list instances --sort launchtime

    # ls is an alias for list
    awless ls users --format csv             
    awless ls roles --sort name,id
    awless ls vpcs --format=json

Listing resources by default performs queries directly to AWS. If you want, you can also query the local snapshot:

    awless list subnets --local

Use `awless list`, `awless list -h` or `awless help list` to see all resources that can be listed.

When dealing with long lists of resources you can filter with the `--filter` flag as such:

    awless list volumes --filter state=in-use --filter volumetype=gp2

    # or with a csv notation
    awless list instances --filter state=running,type=t2.micro 

    # when dealing with name with spaces use
    awless list instances --filter "access key"=my-key
    
For instance, you could list all storage objects in a given bucket using only local data with:

    awless --local ls storageobjects --filter bucketname=pdf-bucket 

Note that filters:

1. ignore case when matching
2. will match when result string contains the search string (ex: `--filter state=Run` will match instances with state `running`)

### Showing resources

`awless show` is quite useful to get a good overview on a resource and to show where its stands in your cloud.

You can either provide the resource _id_, or simpler the resource's _name_. `awless` resolves the id behind the scene (this is the concept of _aliasing_)

    # show instance info: relations to subnets, vpcs, region, ...
    awless show i-34vgbh23jn        

    # show bucket info via its alias: objects it contains, siblings, etc...
    awless show @my-bucket          

    # show user using local data: policy applying to this user, etc...
    # snappy! will not refetch but work with the local graph
    awless show admin-user --local  
    

Basically `awless show` try to maximize the info nicely on your terminal for a given resource

### Creating, Updating and Deleting resources

`awless` provides a powerful template system to interact with cloud infrastructures.

`awless` templates can be used through oneliner shortcut commands:

Using the help:

    awless create                # show what resource can be created
    awless delete -h             # same as above
    awless create instance -h    # show required & extra params for instance creation

Then:

    awless create instance       # will start a prompt for any missing params
    awless delete subnet id=subnet-12345678
    awless attach volume id=vol-12345678 instance=i-12345678

See [Templates (wiki)](https://github.com/wallix/awless/wiki/Templates) for more.

You can also run an `awless` template from a predefined template file with:

    awless run awless-templates/create_instance_ssh.awless

In each case, the CLI guide you through any running of a template (file template or one-liner) so you always have the chance to confirm or quit.

Check out the examples of runnig those commands at [Examples](https://github.com/wallix/awless/wiki/Examples)

Note that you can get inspired with our **in progress** [repo of pre-existing templates](https://github.com/wallix/awless-templates)

### Log & revert executed template commands

To list a detailed account of the last actions you have run on your cloud:

    awless log

Each `awless` command that changes the cloud infrastructure is associated with an unique *id* referencing the (un)successful actions. Using this id you can revert a executed template with:

    awless revert 01B89ZY529E5D7WKDTQHFC0RPA

The CLI guide you through a revert action and you have the chance to confirm or quit.

### Cloud history (in progress)

Using the local auto sync functionality of the cloud resources `awless history` will display in a digested manner the changes that occurred in your infra:

     awless history      # show changes at the resources level
     awless history -p   # show changes including changes in the resources properties
     
*Note*: As model/relations for resources may evolve, if you have any issues with `awless history` between version upgrades, run `rm -Rf ~/.awless/aws/rdf` to start fresh.

### SSH

You can directly ssh to an instance with:

     awless ssh i-abcd1234
     awless ssh ubuntu@i-abcd1234

In the first case, note that `awless` can work out the default ssh user to use given a cloud (ex: `ec2` for AWS)

### Aliasing

When it makes sense we provide the concept of *alias*. Cloud resources ids can be a bit cryptic. An alias is just an already existing name of a resource. Given a alias we resolve the proper resource id. For instance:

        awless ssh my-instance         # ssh to the instance by name. awless resolves its id
        awless delete id=@my-instance  # delete an instance using its name

### Inspectors

Inspectors are small CLI utilities that leverage `awless` graph modeling of cloud resources. Basically an inspector is a program that implements the following interface:

```go
type Inspector interface {
    Inspect(...*graph.Graph) error
    Print(io.Writer)
    Name() string            # name of the inspector
    Services() []string      # name of the services (ec2, iam, s3, ...) the inspector operates on
}
```

Using `awless` cloud resources local synchronisation functionality, you can analyse your data offline (i.e: on your local graphs). There are some builtin inspectors that serve as examples: `pricer`, `bucket_sizer`, etc...

For example, you would run the `bucket_sizer` inspector with:

       $ awless inspect -i bucket_sizer --local
       Bucket           Object count    S3 total storage
       --------         ----------      -----------------
       my-first-bucket     4            0.00358 Gb
       my-other-bucket     1            3.49460 Gb
       third-bucket        422          0.00003 Gb
       fouth-bucket        1000         0.00772 Gb
                                        3.5059 Gb

Note that - as a upcoming feature - using the local infrastructure snapshots (automatically synced), we will be able to run inspectors through time very fast (i.e: all done locally)! For instance, in this case you would see the evolution of your bucket sizing!

# About

`awless` is an open source project created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and François-Xavier Aguessy at Wallix.
`awless` is released under the Apache License and sponsored by [Wallix](https://github.com/wallix).

Contributors are welcome! Please head to [Contributing (wiki)](https://github.com/wallix/awless/wiki/Contributing) to learn more.


