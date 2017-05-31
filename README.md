[![Build Status](https://api.travis-ci.org/wallix/awless.svg?branch=master)](https://travis-ci.org/wallix/awless)
[![Go Report Card](https://goreportcard.com/badge/github.com/wallix/awless)](https://goreportcard.com/report/github.com/wallix/awless)

`awless` is a fast, powerful and easy-to-use command line interface (CLI) to manage Amazon Web Services.

[Twitter](http://twitter.com/awlessCLI) | [Wiki](https://github.com/wallix/awless/wiki) | [Changelog](https://github.com/wallix/awless/blob/master/CHANGELOG.md#readme)

# Why awless

`awless` will help you achieve your goals without leaving your terminal:

- run frequent actions by using simple commands
- get nice and readable output (for humans) that machine know how to parse too
- explore and query your infrastructure and cloud resources, even **offline**
- ensure smart defaults & security best practices
- write and run powerful templates (see [`awless` templates](https://github.com/wallix/awless/wiki/Templates))
- connect to your instances easily

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
<em>Note that this video is in <a href="https://en.wikipedia.org/wiki/APNG">APNG</a>. On Chrome, you need <a href="https://chrome.google.com/webstore/detail/apng/ehkepjiconegkhpodgoaeamnpckdbblp">an extension</a> to view it.</em>
</p>

- Clear and easy listing of multi-region cloud resources (subnets, instances, groups, users, etc.) on AWS EC2, IAM and S3: `awless list`
- Output formats either human (Markdown-compatible tables, trees) or machine readable (csv, tsv, json, ...): `--format`
- Listing can be filtered via *resource properties* or *resources tags* 
- Explore a resource given only an *id*, name or arn (properties, relations, dependencies, ...): `awless show`
- Creation, update and deletion (CRUD) of cloud resources and complex infrastructure with smart defaults and sound autocomplete through powerful awless templates: `awless run my-awless-templates/create_my_infra.txt`
- Powerful CRUD CLI one-liner (integrated in the awless templating engine) with: `awless create instance ...`, `awless create vpc ...`, `awless attach policy ...`
- Leveraging AWS `userdata` to provision instance on creation from remote (i.e http) or local scripts: `awless create instance userdata=http://...` 
- Easy reporting of all the CLI template executions: `awless log`
- Revert of executed templates and resources creation: `awless revert`
- Clean and simple ssh to instances: `awless ssh`
- Resolve public images dynamically (i.e. independant of the region specific AMI id): `awless search images canonical:ubuntu:xenial --id-only`
- Aliasing of resources through their natural name so you don't have to always use cryptic ids that are impossible to remember
- Inspectors are small CLI utilities to run analysis on your cloud resources graphs: `awless inspect`
- Manual sync mode to fetch & store resources locally. Then query & inspect your cloud offline: `awless sync`
- CLI autocompletion for Unix/Linux's bash and zsh `awless completion`

# Getting started

Read more at [awless wiki: Getting Started](https://github.com/wallix/awless/wiki/Getting-Started)

# About

`awless` is an open source project created by Henri Binsztok, Quentin Bourgerie, Simon Caplette and François-Xavier Aguessy at Wallix.
`awless` is released under the Apache License and sponsored by [Wallix](https://github.com/wallix).

    Disclaimer: Awless allows for easy resource creation with your cloud provider; we will not be responsible for any cloud costs incurred (even if you create a million instances using awless templates).

Contributors are welcome! Please head to [Contributing (wiki)](https://github.com/wallix/awless/wiki/Contributing) to learn more.
