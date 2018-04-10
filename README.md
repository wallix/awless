

[![Build Status](https://api.travis-ci.org/wallix/awless.svg?branch=master)](https://travis-ci.org/wallix/awless)
[![Go Report Card](https://goreportcard.com/badge/github.com/wallix/awless)](https://goreportcard.com/report/github.com/wallix/awless)

<img src="https://user-images.githubusercontent.com/808274/33351381-5b9a0d00-d458-11e7-91ed-cf7ada7237c1.png" alt="terminal icon" width="48"> `awless` is a powerful, innovative and small surface command line interface (CLI) to manage Amazon Web Services.

[Twitter](http://twitter.com/awlessCLI) | [Wiki](https://github.com/wallix/awless/wiki) | [Changelog](https://github.com/wallix/awless/blob/master/CHANGELOG.md#readme)

# Why awless

`awless` stands out by providing the following features:

- small and hierarchical set of commands
- create and revert fully-fledged infrastructures through a new simple and powerful templating language (see [`awless` templates (wiki)](https://github.com/wallix/awless/wiki/Templates))
- local log of all your cloud modifications done through `awless`
- exploration of your cloud infrastructure and resources relations, **even offline** using a local graph storage
- greater output's readability with numerous machine and human friendly formats
- ensure smart defaults & security best practices
- connect easily using awless' **smart SSH** to your private & public instances

# Install

Choose one of the following options:

1. On macOS, use [homebrew](http://brew.sh):  `brew tap wallix/awless; brew install awless`
2. With `curl` (macOS/Linux), run: `curl https://raw.githubusercontent.com/wallix/awless/master/getawless.sh | bash`
3. Download the latest `awless` binaries (Windows/Linux/macOS) [from Github](https://github.com/wallix/awless/releases/latest)
4. If you have Golang already installed, install from the source with: `go get -u github.com/wallix/awless`

# Main features

<p align="center">
  <a href="https://raw.githubusercontent.com/wiki/wallix/awless/apng/awless-demo.png"><img src="https://raw.githubusercontent.com/wiki/wallix/awless/apng/awless-demo.png" alt="video of a few awless commands"></a>
<br/>
<em>Note that the video above is in <a href="https://en.wikipedia.org/wiki/APNG">APNG</a> and requires a recent browser.</em>
</p>

- **Aliasing of resources through their natural name** so you don't have to always use cryptic ids that are impossible to remember
- `awless show` : Explore a resource given only a *name* (or id/arn) showing its properties, relations, dependencies, etc.
- `awless run` : Creation, update and deletion of complex infrastructures with smart defaults and sound autocomplete through awless templates

      $ awless run my-awless-templates/create_my_infra.aws

- **Hundreds of powerful CRUD CLI one-liners** integrated in the awless templating engine:

      $ awless create instance -h
      $ awless create vpc -h
      $ awless attach policy -h
      etc.

- `awless log` : Easy reporting of all the CLI template executions
- `awless revert` : Revert of executed templates and resources creation
- Create instances straight from a distro name. No need to know the region or AMI ;) (_free tier community bare distro only_, see `awless create instance -h`)

      $ awless create instance distro=debian
      $ awless create instance distro=coreos
      $ awless create instance distro=redhat::7.2 type=t2.micro
      $ awless create instance distro=debian:debian:jessie lock=true
      $ awless create instance distro=amazonlinux:amzn2
      etc.

- Leveraging AWS `userdata` to provision instance on creation from remote (i.e http) or local scripts: `awless create instance ... userdata=/home/john/...` 
- `awless ssh` : Clean and simple SSH to public & private instances using only a name

      $ awless ssh my-production-instance
      $ awless ssh redis-prod --through jump-server
      $ awless ssh 34.215.29.221
      $ awless ssh db-private --private
      $ awless ssh 172.31.77.151 --port 2222 --through my-proxy --through-port 23
      etc.

- `awless switch` : Switch easily between AWS accounts (i.e. profile) and regions

       $ awless switch admin eu-west-2
       $ awless switch us-west-1
       $ awless switch mfa
       etc.

- `awless inspect` : Leverage _experimental_ inspectors which are small CLI utilities to run analysis on your cloud resources graphs
- `awless list` : Clear and easy listing of multi-region cloud resources (subnets, instances, users, buckets, records, etc.) on AWS EC2, IAM, S3, RDS, AutoScaling, SNS, SQS, Route53, CloudWatch, CloudFormation, Lambda, etc.
- `awless sync` : Explicitly and manually sync to fetch & store resources locally. Then query & inspect your cloud offline
- Output formats either human (Markdown-compatible tables) or machine readable (csv, tsv, json, ...): `--format`
- Listing filters via *resources properties* or *resources tags*: `--filter property=val`, `--tag Env=Production`, `--tag-value Purchased`, `--tag-key Dept,Internal`
- `awless completion` : CLI autocompletion for Unix/Linux's bash and zsh 

# Getting started

Take the tour at [Getting Started (wiki)](https://github.com/wallix/awless/wiki/Getting-Started).

Or read the [introductory blog post about awless](https://medium.com/@hbbio/awless-io-a-mighty-cli-for-aws-a0d48bdb59a4).

More articles:

   - [Simplified user management for AWS](https://medium.com/@awlessCLI/simplified-user-management-for-aws-6f828ccab387)
   - [InfoWorld: Production-grade deployment of WordPress](https://www.infoworld.com/article/3230547/cloud-computing/awless-tutorial-try-a-smarter-cli-for-aws.html)
   - [Easy create & tear down of a multi-AZ CockroachDB cluster](https://github.com/wallix/awless-templates/tree/master/cockroachdb)

# Awards

- [Top 50 Developer Tools of 2017](https://stackshare.io/posts/top-developer-tools-2017)
- [InfoWorld Bossie Awards 2017](https://www.infoworld.com/article/3227920/cloud-computing/bossie-awards-2017-the-best-cloud-computing-software.html#slide12)

# About

`awless` is an open source project created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and François-Xavier Aguessy at WALLIX.
`awless` is released under the Apache License and sponsored by [Wallix](https://github.com/wallix).

    Disclaimer: Awless allows for easy resource creation with your cloud provider;
    we will not be responsible for any cloud costs incurred (even if you create a 
    million instances using awless templates).

Contributors are welcome! Please head to [Contributing (wiki)](https://github.com/wallix/awless/wiki/Contributing) to learn more.
Note that `awless` uses [triplestore](https://github.com/wallix/triplestore) another project developped at WALLIX.
