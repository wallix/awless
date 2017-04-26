## v0.0.23 [unreleased]

### Features

- Create and attach role to a user, instance, etc ...
- [#86](https://github.com/wallix/awless/issues/86): SSH using private IP with `--private` flag. Thanks @padilo.
- `awless ssh` now checks the remote host public key before connecting. Check can be disabled with the (insecure) `--disable-strict-host-keychecking` flag.
- [#74](https://github.com/wallix/awless/issues/74): support of encrypted SSH keys for generation `awless create keypair encrypted=true` and in `awless ssh`.
- Create/delete and list Lambda functions
- Create/delete and list autoscaling launch configurations and autoscaling groups
- List EC2 images (AMIs) of which you are the owner: `awless list images`
- Copy an EC2 image from a given region to the current region: `awless copy image name=... source-id=... source-region=...`
- List your IAM access keys: `awless list accesskeys`
- Friendlier (using units: KB, MB, GB) display for storage size (s3objects, volumes, lambda functions)

### Bugfixes

- Update SSH library to fix [CVE-2017-3204](http://www.cve.mitre.org/cgi-bin/cvename.cgi?name=2017-3204).
- Take the file name rather than full path as default name when uploading a s3object
- Correctly create repo on first install on machine with git not installed

## v0.0.22 [2017-04-13]

### Features

- Amazon [**userdata**](http://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html) support. Give the data as local file or remote http file resource. Ex: `awless create instance userdata=/tmp/mydata.sh ...` or `awless create instance userdata=https://gist.github.com/jsmith/5f58272fa5406`.
- Global rename of `storageobject` to `s3object` for shorter typing in CLI.
- awless model/storing is now full RDF ;). Allow exploration of all your infra in RDF tools and ontology editor (Ex: [Protege](http://protege.stanford.edu/))
- Faster, better and simpler RDF & triples management now done through the nifty library [triplestore](https://github.com/wallix/triplestore)
- Ability to use strings with spaces and special characters in template parameters by surrounding them with single or double quotes.
- Loggers are now sent to the stderr file descriptor which makes easier piping and redirecting output.
- Warn when creating an instance without access key.
- ssh: print SSH configuration (`~/.ssh/config`) or the CLI one-liner to connect with SSH using `--print-config` or `--print-cli` flags.
- ssh: better handle when several instances have the same name (e.g., with a running and a terminated instance)
- ssh: more warning; provide help and context on failing connections
- Manage properly secgroups on instances with `awless attach/detach secgroup id=... instance=@my-instance`
- Logging more info when running templates

### Bugfixes

- `awless whoami` now supports displaying info for `root` user and user with org path
- Use `securitygroup` rather than `group` in templates, when appropriate.
- Use `keypair` rather than `key` in templates, when appropriate.
- Fix the fact you could not attach multiple security groups to an instance
- Reverting the creation of a load balancer now waits the deletion of its network interfaces

## v0.0.21 [2017-03-23]

### Features

- `awless whoami` now returns your identity, your attached (i.e. managed), inlined and group policies
- Rudimentary security groups port scanner inspector via `awless inspect -i port_scanner`
- Template: compile time check of undefined or unused references
- Run official remote templates without specifying full url: `awless run repo:create_vpc`
- [#78](https://github.com/wallix/awless/issues/78): Show progress when uploadgin object to storage
- [#81](https://github.com/wallix/awless/issues/81): Global force flag `--force` to bypass confirm prompt

### Bugfixes

- Fix regression: run templates/one-liners failed on `storageobject`, `subscription` entities
- Filtering in `awless list --filter` now works with column types other than string
- Users, groups and policies are now independent of the region
- [#83](https://github.com/wallix/awless/issues/83): Syncing while offline does not clear local cloud infra

## v0.0.20 [2017-03-20]

### Features

- Auto completion of id/name to help fill in easily any missing info before template execution
- Better error messaging on parsing template errors
- Infra: basic support of RDS: listing, creation and deletion of databases and database subnets:  `awless list databases/dbsubnetgroups`; `awless create/delete database/dbsubnetgroup`
- Infra: attach/detach an `instance` to a `targetgroup`
- Infra: delete tag: `awless delete tag`
- Access: create an AWS access key for a user
- DNS: allow to revert creation/deletion of records
- [#80](https://github.com/wallix/awless/issues/80) DNS: return the ChangeInfo id when creating/deleting a record

### Bugfixes

- [#79](https://github.com/wallix/awless/issues/79): `awless list records` do not add new lines between records.
- Better compute table columns width to adjust the number of columns to display exactly to the terminal width.

## v0.0.19 [2017-03-16]

### Features

- [#76](https://github.com/wallix/awless/issues/76): Show private IP and availability zones when listing instances.
- Run remote template when path prefixed with `http`. Ex: `awless run http://github.com/wallix/awless-templates/...`
- Fetch more instances properties when showing instances (ex: network interfaces, public and private DNS, Root device type and name...)
- DNS: listing Route53 zones and records `awless list zones/records`
- DNS: basic creation/deletion of Route53 zones and records `awless create/delete zone/record`
- Infra: detach EBS volumes `awless detach volume`
- Config: enable/disable the syncing of Route53 service `awless config set aws.dns.sync`
- All listing with default format are now Markdown table compatible. 
- Better display of `awless show`. Added `--siblings` flag to display exhaustively all siblings
- Reverse the sorting order when listing instances sorted by "up since"

### Bugfixes

- Fix `awless show` to properly show relations between groups and users

## 0.0.18 [2017-03-13]

### Features

- infra: support the creation/deletion of ELBv2 loadbalancers, listeners and target groups: `awless create loadbalancer/listener/targetgroup`
- infra: add tag `Name` to subnets.
- Format `tsv` supported when listing: `awless list subnets --format tsv`
- Pricer inspector now resolves prices for any regions: `awless inspect -i pricer`

### Bugfixes

- Fix alias, required and extra params parsing in template runs

## 0.0.17 [2017-03-09]

If you have any data or config issues, you can run `rm -Rf ~/.awless/` to start with a fresh install.

### Features

- [#65](https://github.com/wallix/awless/issues/65): `awless ssh`: use existing SSH client if available, otherwise fallback on builtin SSH.
- `awless show` resolves automatically on id, name or arn without any prefixing (previously it was '@')
- [#47](https://github.com/wallix/awless/issues/47): Enable/disable sync per services or resources through config. Ex: `awless config set aws.notification.sync false`, `awless config set  aws.storage.storageobject.sync true`.
- [#55](https://github.com/wallix/awless/issues/55): Dynamically change AWS region/profile with global flags `--aws-region us-west-1` or `--aws-profile myprofile`.
- [#73](https://github.com/wallix/awless/issues/73): `AWS_DEFAULT_REGION` env variable now loaded in `awless`. It takes precedence over `aws.region`.
- [#73](https://github.com/wallix/awless/issues/73): `AWS_DEFAULT_PROFILE` env variable now loaded in `awless`. It takes precedence over `aws.profile`.
- Better output of `awless config list` (doc per variable, etc.).
- Global default menu with clearer one-liner display.
- Simplification of the templating engine using decoupled compile passes.
- Config setters now provide dialogs (ex: `awless config set instance.type` or `awless config set aws.region`).
- [#54](https://github.com/wallix/awless/issues/54): `awless ssh`: specify the keyfile to use with `-i /path/toward/key` flag.
- [#64](https://github.com/wallix/awless/issues/64): `awless ssh`: columns and lines automatically adapt to terminal with/height.

- Attach/detach policy to user/group (see [wiki examples](https://github.com/wallix/awless/wiki/Examples))
- Attach/detach user to group (see [wiki examples](https://github.com/wallix/awless/wiki/Examples))
- List AWS load balancers, target groups and listeners with `awless list loadbalancers/targetgroups/listeners`. Show their relations with, e.g. `awless show LOAD_BALANCER`.

### Bugfixes

- [#12](https://github.com/wallix/awless/issues/12): Support AWS pagination when fetching resources in AWS IAM.
- Template parsing: allow digits in refs; allow regular chars in alias declaration
- Template: all aliases now resolves correctly from file or CLI. Ex: `awless create instance subnet=@my-subnet`

## 0.0.16 [2017-03-01]

### Features

- Allow simple fuzzy search for listing filters. Ex: `awless list instances --filter state=run`
- Revert: waiting instance termination when deleting a vpc/subnet/instance hierarchy.

### Bugfixes

- Fix regression: timeout too low for HTTP requests with AWS.

## 0.0.15 [2017-02-28]

As model/relations for resources may evolve, if you have any issues with models related commands, you can run `rm -Rf ~/.awless/aws/rdf` to start a fresh RDF model.

### Features

- [#6](https://github.com/wallix/awless/issues/6): Create Linux installer shell script: `curl https://raw.githubusercontent.com/wallix/awless/master/getawless.sh | bash`
- [#42](https://github.com/wallix/awless/issues/42), [#60](https://github.com/wallix/awless/issues/60), [#66](https://github.com/wallix/awless/issues/66): Better load AWS credentials (support profile credentials, MFA and crossaccount profile access)
- [#32](https://github.com/wallix/awless/issues/32): Basic support of [SNS](https://aws.amazon.com/sns/) (CRUD for topics and subscriptions)
- [#32](https://github.com/wallix/awless/issues/32): Basic support of [SQS](https://aws.amazon.com/sqs/) (CRUD for queues)
- [#53](https://github.com/wallix/awless/issues/53): Filter results in listings. Ex: `awless ls instances --filter state=running,"Access Key"=my-key` or the equivalent `awless list instances --filter state=running --filter "Access Key"=my-key`
- Better help menus by splitting one-liner template commands from general commands
- Run template: better dialog and remove noisy info
- Template validation: notify on unexpected params; check names unicity against local graph
- Log contextual error instead of hard failure when user has no rights to sync a service

### Bugfixes

- [#57](https://github.com/wallix/awless/issues/57): Properly fetch buckets when they are in the `us-east-1` region.
- [#12](https://github.com/wallix/awless/issues/12): Support AWS pagination when fetching resources in AWS SNS and EC2.

## 0.0.14 [2017-02-21]

As model/relations for resources may evolve, if you have any issues with models related commands, you can run `rm -Rf ~/.awless/aws/rdf` to start a fresh RDF model.

### Features

- [#39](https://github.com/wallix/awless/issues/38), [#38](https://github.com/wallix/awless/issues/33): Remove data collection & sending
- [#33](https://github.com/wallix/awless/issues/33): Ability to set AWS profile using `aws.profile` config key
- Better output for `awless sync`
- `awless ls` now an alias for `awless list`

### Bugfixes

- [#44](https://github.com/wallix/awless/issues/44): Fetch only the S3 buckets and related objects of the current region.
- [#52](https://github.com/wallix/awless/issues/52), [#34](https://github.com/wallix/awless/issues/34): Properly fetch route tables, even if a route contains several destinations.
- [#37](https://github.com/wallix/awless/issues/37): Load the region from database when initializing cloud services rather than `awless` environment.
- [#56](https://github.com/wallix/awless/issues/56): Do not require a VPC as parent of security groups nor route table.
