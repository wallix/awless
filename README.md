[![Build Status](https://api.travis-ci.org/wallix/awless.svg?branch=master)](https://travis-ci.org/wallix/awless)
[![Go Report Card](https://goreportcard.com/badge/github.com/wallix/awless)](https://goreportcard.com/report/github.com/wallix/awless)

<img src="https://user-images.githubusercontent.com/808274/33351381-5b9a0d00-d458-11e7-91ed-cf7ada7237c1.png" alt="terminal icon" width="48">
`awless` is a powerful, innovative and small surface command line interface (CLI) to manage Amazon Web Services.

[Twitter](http://twitter.com/awlessCLI) | [Wiki](https://github.com/wallix/awless/wiki) | [Changelog](https://github.com/wallix/awless/blob/master/CHANGELOG.md#readme)

# Why awless

`awless` stands out by providing the following features:

- small and hierarchical set of commands
- create and revert fully-fledged infrastructures through a new simple and powerful templating language (see [`awless` templates (wiki)](https://github.com/wallix/awless/wiki/Templates))
- local log of all your cloud modifications done through `awless`
- exploration of your cloud infrastructure and resources relations, **even offline** using a local graph storage
- greater output's readability with numerous machine and human friendly formats
- ensure smart defaults & security best practices
- connect easily through smart SSH to your private & public instances


# Install

Choose one of the following options:

1. On macOS, use [homebrew](http://brew.sh):  `brew tap wallix/awless; brew install awless`
2. With `curl` (macOS/Linux), run: `curl https://raw.githubusercontent.com/wallix/awless/master/getawless.sh | bash`
3. Download the latest `awless` binaries (Windows/Linux/macOS) [from Github](https://github.com/wallix/awless/releases/latest)
4. If you have Golang already installed, install from the source with: `go get -u github.com/wallix/awless`

# Main Features

<p align="center">
  <a href="https://raw.githubusercontent.com/wiki/wallix/awless/apng/awless-demo.png"><img src="https://raw.githubusercontent.com/wiki/wallix/awless/apng/awless-demo.png" alt="video of a few awless commands"></a>
<br/>
<em>Note that the video above is in <a href="https://en.wikipedia.org/wiki/APNG">APNG</a> and requires a recent browser.</em>
</p>

- Clear and easy listing of multi-region cloud resources (subnets, instances, users, buckets, records, etc.) on AWS EC2, IAM, S3, RDS, AutoScaling, SNS, SQS, Route53, CloudWatch, CloudFormation, Lambda, etc.: `awless list`
- Output formats either human (Markdown-compatible tables) or machine readable (csv, tsv, json, ...): `--format`
- Listing filters via *resources properties* or *resources tags*: `--filter property=val`, `--tag Key=Value`
- Explore a resource given only a *name* (or id/arn) showing its properties, relations, dependencies, etc.: `awless show`
- Creation, update and deletion of complex infrastructures with smart defaults and sound autocomplete through awless templates: `awless run my-awless-templates/create_my_infra.txt`
- Powerful CRUD CLI one-liner (integrated in the awless templating engine) with: `awless create instance ...`, `awless create vpc ...`, `awless attach policy ...`
- Leveraging AWS `userdata` to provision instance on creation given remote (i.e http) or local scripts: `awless create instance ... userdata=http://...` 
- Easy reporting of all the CLI template executions: `awless log`
- Revert of executed templates and resources creation: `awless revert`
- Clean and simple SSH to public & private instances using only a name: `awless ssh my-production-instance`, `awless ssh redis-prod --through jump-server`
- Resolve public images dynamically (i.e. independant of the region specific AMI id): `awless search images canonical:ubuntu:xenial --latest-id`
- Aliasing of resources through their natural name so you don't have to always use cryptic ids that are impossible to remember
- Inspectors are small CLI utilities to run analysis on your cloud resources graphs: `awless inspect`
- Manual sync mode to fetch & store resources locally. Then query & inspect your cloud offline: `awless sync`
- CLI autocompletion for Unix/Linux's bash and zsh `awless completion`

# Getting started

Take the tour at [Getting Started (wiki)](https://github.com/wallix/awless/wiki/Getting-Started).

Or read the [introductory blog post about awless](https://medium.com/@hbbio/awless-io-a-mighty-cli-for-aws-a0d48bdb59a4).

More articles:

   - [Simplified user management for AWS](https://medium.com/@awlessCLI/simplified-user-management-for-aws-6f828ccab387)
   - [InfoWorld: Production-grade deployment of WordPress](https://www.infoworld.com/article/3230547/cloud-computing/awless-tutorial-try-a-smarter-cli-for-aws.html)

# About

`awless` is an open source project created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and François-Xavier Aguessy at WALLIX.
`awless` is released under the Apache License and sponsored by [Wallix](https://github.com/wallix).

    Disclaimer: Awless allows for easy resource creation with your cloud provider;
    we will not be responsible for any cloud costs incurred (even if you create a 
    million instances using awless templates).

Contributors are welcome! Please head to [Contributing (wiki)](https://github.com/wallix/awless/wiki/Contributing) to learn more.
Note that `awless` uses [triplestore](https://github.com/wallix/triplestore) another project developped at WALLIX.
