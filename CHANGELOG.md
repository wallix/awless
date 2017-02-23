## 0.0.15 [unreleased]

As model/relations for resources may evolve, if you have any issues with models related commands, you can run `rm -Rf ~/.awless/aws/rdf` to start a fresh RDF model.

### Features

- [#6](https://github.com/wallix/awless/issues/6): Create Linux installer shell script: `curl https://raw.githubusercontent.com/wallix/awless/master/getawless.sh | bash`
- Better help menus by splitting one-liner template commands from general commands
- [#32](https://github.com/wallix/awless/issues/32): Basic support of [SNS](https://aws.amazon.com/sns/) (CRUD for topics and subscriptions)
- Run template: better dialog and remove noisy info
- Template validation: notify on unexpected params; check names unicity against local graph
- Log contextual error instead of hard failure when user has no rights to sync a service

### Bugfixes

- [#57](https://github.com/wallix/awless/issues/57): Properly fetch buckets when they are in the `us-east-1` region.

## 0.0.14

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