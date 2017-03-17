## v0.0.20 [unreleased]

### Features

- Better error messaging on parsing template errors
- Infra: basic support of RDS: listing, creation and deletion of databases:  `awless list databases`; `awless create/delete database`
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
