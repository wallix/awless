## v0.1.11 [2018-06-21]

**Check out our new article** on [Simplified Multi-Factor Authentication](https://medium.com/@awlessCLI/simplified-multi-factor-authentication-for-aws-d703e8d9f332) with `awless`

### Features

- [#71](https://github.com/wallix/awless/issues/71): Add support for Classic load-balancers:

```
    $ awless list classicloadbalancers
    $ awless create classicloadbalancer name=my-loadb subnets=[sub-123,sub-456] listeners=HTTP:80:HTTP:8080 healthcheck-path=/health/ping  securitygroups=sg-54321 tags=Env:Test,Created:Awless
    $ awless update classicloadbalancer name=my-loadb health-interval=10 health-target=HTTP:80/weather/ health-timeout=300 healthy-threshold=10  unhealthy-threshold=5
    $ awless attach classicloadbalancer name=my-loadb instance=@redis-prod-1
    $ awless delete classicloadbalancer name=my-loadb
```

- [#214](https://github.com/wallix/awless/issues/214): `AWS_PROFILE` env variable now loaded in `awless` in addition to the deprecated `AWS_DEFAULT_PROFILE` thanks to @alewando
- Better completion for `attach mfadevice` and `attach user` commands
- [#219](https://github.com/wallix/awless/issues/219): Validate access key and secret key before writing into `~/.aws/credentials` file

### Fixes

- [#220](https://github.com/wallix/awless/issues/220): Add double quotes to CSV output if needed thanks to @lllama
- Fix compilation error in templates with concatenation and reference (c.f. for example in [this template](https://gist.githubusercontent.com/fxaguessy/ef9511bf5ed8f3312904cccb96b818e8/raw/75c0f808220665441055b589be133cf711c64f37/ManageOwnMFA.aws))
- Parse integer beginning with '0' as string (preventing the deletion of the initial '0' for example in `... account.id=0123456789`)

## v0.1.10 [2018-04-13]

### Features

- Much better performance when synchronising all access data (IAM, etc.)
- Create instances now supports distro prompting for CentOS, Amazon Linux 2, CoreOS
   
      $ awless create instance name=myinst distro=amazonlinux:amzn2
      $ awless create instance distro=coreos
      $ awless create instance distro=centos name=myinst

- Avoiding extra throttling: Listing flag `--filter` now passes on the user wanted filtering down to the AWS API when possible so that _less unneeded resources are fetched_, _bandwidth is reduced_ and _some throttling avoided_.
  
  For example:
  
      $ awless ls s3objects --filter bucket=website
      $ awless ls records --filter name=io
      $ awless ls containertasks --filter name=my-task-definition-name

- Support for region embedded in an AWS profile (i.e. shared config files ~/.aws/{credentials,config}). See #181 in Fixes for more details 
      
- [#191](https://github.com/wallix/awless/issues/191) Attach a certificate to a listener with: `awless listener attach id=... certificate=...` (see awless attach listener -h for more)


### Fixes

- [#200](https://github.com/wallix/awless/issues/200): Now paging is supported for s3 objects when listing
- [#196](https://github.com/wallix/awless/issues/196): Regression fix SIGSEV when having AWS config with role assuming
- [#182](https://github.com/wallix/awless/issues/182): Region embedded in profile taken into account and given correct precedence
- [#144](https://github.com/wallix/awless/issues/144): Filtering done on AWS side when listing records for a given zone name
- [#172](https://github.com/wallix/awless/issues/172): Filtering done on AWS side when listing containertasks for a given task definition name

## v0.1.9 [2018-01-16]

**In this release, the local data model has been updated to support multi-account and stale data is removed when upgrading. Local data (ex: used for completion, etc...) will progressively be synced again through your usage of awless. Although, to get all your data now under the new model, you can manually run `'awless sync'`**	

### Features

- Support and seamless sync across multi-account (i.e. multiple profiles) and regions
- Enriched params prompting with optional/skippable but very common params. Can be disabled with `--prompt-only-required` or forced with `--prompt-all` to leverage smart completion for all params
- Automatically complete the username when deleting an access key by its ID, if it is contained in the local graph model:
    * `awless delete accesskey id=ACCESSKEYID`
-  For `awless update stack` param `stackfile` can now slurp yml and json params files. Thanks to @Trane9991 ([#167](https://github.com/wallix/awless/pull/167), [#145](https://github.com/wallix/awless/issues/145))
- Better completion for template parameters independently of their display name
- Aliases can now be resolved to properties other than IDs. For example, they are resolved to ARN in attach/detach/update/delete policy: `awless attach policy arn=@my-policy-name`
- Running only `awless switch` now returns your current region and profile, allowing a quick and short region/profile lookup
- Better completion of slice properties 

### AWS Services

- Listing of Route53 records now contains a new column for aliases [#181](https://github.com/wallix/awless/issues/181)
- Create an image from an existing instance. See `awless create image -h`
    * `awless create image instance=@my-instance-name name=redis-image  description='redis prod image'`
    * `awless create image instance=i-0ee436a45561c04df name=redis-image reboot=true`
    * List your images with `awless ls images --sort created`
    * Delete images with an `awless revert ...` or with `awless delete image id=@redis-image`
- [#169](https://github.com/wallix/awless/issues/169): Start/Stop a RDS database:
    * `awless start database id=my-db-id`
    * `awless stop database id=@my-db-name`
    * `awless restart database id=@my-db-name`
- Restart an EC2 instance
  * `awless restart instance id=id-1234`
  * `awless restart instance ids=@redis-prod-1,@redis-prod-2`
- [#176](https://github.com/wallix/awless/issues/176): Delete a DNS record only by its awless ID (see `awless ls records`) or by its name:
    * `awless delete record id=awls-39ec0618`
    * `awless delete record id=@my.sub.domain.com`

### Fixes

- Fix regression error: errors in dry run showed but where ignored hence user could wrongly confirm to run the template
- Delete a DNS record only by its awless ID

## v0.1.8 [2017-11-29]

### Features

- Better prompting of template parameters
- Overall better logging output of template execution

### AWS Services

- Create a database replica with: `awless create database replica=...`

## v0.1.7 [2017-11-24]

### Features

- Better prompt completion for template parameters
- Create instance/launchconfiguration from community distro names (`awless create instance distro=debian`). In default config value, deprecation of `instance.image` in favor of `instance.distro` (migration should be seamless).
    * `awless create instance distro=redhat:rhel:7.2`
    * `awless create launchconfiguration distro=canonical:ubuntu`
    * `awless create instance distro=debian`
- Quick way to switch to profiles and regions. Ex: `awless switch eu-west-1`, `awless switch mfa us-west-1`
- Create a public subnet in only one command with: `awless create subnet public=true...`
- Save directly your newly created access key in `~/.aws/credentials` with : `awless create accesskey save=true`
- Overall better logging output of template execution

### AWS Services

- Update Cloudfront distribution with: `awless update distribution...`

## v0.1.6 [2017-11-16]

**Overall re-design of AWS commands with full acceptance testing allowing for easier external contribution, greater flexibility and scalability moving forward**

### Features

- [#154](https://github.com/wallix/awless/issues/154): `awless ssh` allow specifying both `--port` and `--through-port`
- [#151](https://github.com/wallix/awless/issues/151): `awless ssh` using ip addresses. Ex: `awless ssh 172.31.68.49 --through 172.31.11.249`
- `awless attach mfadevice` now propose to automatically add the MFA device configuration to `~/.aws/config`
- [#158](https://github.com/wallix/awless/pull/158), [#159](https://github.com/wallix/awless/pull/159): Added bash/zsh completion to regions and profiles. Thanks to @padilo.

## v0.1.5 [2017-10-05]

### Features

- Complete flow to enable MFA for a user, including QRCode generation
- Much better output for `awless log`; default message (or user specified message) stored now in logs
- [#143](https://github.com/wallix/awless/issues/143): Follow CloudFormation stack events: `awless tail stack-events my-stack-name --follow`. Thanks to @Trane9991.
- Support concatenation between `{holes}` and `"quoted strings"` in template with `+` operator: `policy = create policy ... resource="arn:aws:iam::" + {account.id} + ":mfa/${aws:username}"`

### AWS Services

- Manage and listing of MFA devices: `awless create/delete/attach/detach mfadevice`, `awless list mfadevices`
- Support [Network Load Balancers](http://docs.aws.amazon.com/elasticloadbalancing/latest/network/introduction.html): `awless create loadbalancer .... type=network ...`
- Add conditions in policies and support multiple resources `awless create policy ... conditions=\"aws:MultiFactorAuthPresent==true\" resource=arn:aws:iam::0123456789:mfa/test,arn:aws:iam::0123456789:user/test`
- Add conditions in role creation `awless create role name=awless-mfa-role principal-account=0123456789 conditions=\"aws:MultiFactorAuthPresent==true\"`
- List the access keys of all users with `awless list accesskeys` (previously, only current user)
- Fetch role trust policy document: `awless show my-role`

### Fixes

- Exit code is now non zero on template run with KO states

## v0.1.4 [2017-09-21]

### Features

- Local storage of cloud data (RDF store) now done using the NTriples text format instead of a binary format (transition completely transparent for the user). New format allows more friendly git revisioning of data compared to a binary format.
- [#87](https://github.com/wallix/awless/issues/87): Customize columns displayed in `awless list` with `--columns`: `awless ls instances --sort name --columns name,vpc,state,privateip`
- Global `--no-sync` flag to not run any sync on command
- `awless show policy-name/policy-id` now displays the current policy Document (in JSON).

### AWS Services

- Update IAM policies, to add statements with `awless update policy`
- Add ACM certificates in infra:
    - `awless list certificates`
    - `awless create/delete/check certificate domains=my.firstdomain.com,my.seconddomain.com validation-domains=firstdomain.com,seconddomain.com`
- [#123](https://github.com/wallix/awless/issues/123): Listing route tables display the association IDs.

### Fixes
- `awless ssh --through`: no reusing same conn to avoid EOF. Bug: only first user (amazonlinux) was successful (usually ec2-user) !!
- `awless ssh --through`: on new proxy client catching error that where shadowed

## v0.1.3 [2017-09-06]

### Features

- `awless show` command 'not found' error now suggests if resource with same reference exists in other locally synced regions
- `awless` template language now supports lists, for example: `create loadbalancer subnets=[$subnet1, $subnet2]`
- Variables in `awless` template language now support references, holes and lists, for example: `mysecgroups = [$secgroup1, {my.secgroup},sg-123456]`
- `awless` template language now supports *holes* in strings, for example: `create instance name={prefix}database{version}`
- `awless update securitygroup` can now authorize/revoke access from another security group: `update securitygroup id=sg-12345 inbound=authorize portrange=any protocol=tcp securitygroup=sg-23456`
- Template CLI prompt: better TAB completion of resources and their properties
- Man CLI examples for all one liners command. For example, `awless create instance -h` will display relevant CLI examples
- Add `Type` (AWS/Customer managed) and `Attached` (true/false) columns in `awless list policies`
- [#129](https://github.com/wallix/awless/issues/129): flag `--color=always/never` to force enabling/disabling of colored output.

### AWS Services

- List network interfaces with `awless list networkinterfaces`

### Fixes

- Fix regression: listing a resource returned no results when this resource was disabled for sync. Listing should always fetch the resources and display what is on your cloud.
- [#130](https://github.com/wallix/awless/issues/130): Better exit status code in `awless show` command
- Port ranges starting from *0* to *n* are no longer processed as from *n* to *n*.
- `awless ssh --through`: works without an SSH agent running; correct StrictHostkeyChecking; correct display for `--print-config`

## v0.1.2 [2017-08-17]

### Features

- Sync overall speed up and massive reducing in memory consumption
- SSH `--through`: `awless ssh my-priv-inst --through my-pub-inst` allow you to connect to a private instance by going through a public one in ths same VPC. You need to have the same keypair (SSH key) on both instances. 
- Flag `--profile-sync` on `awless sync` to enable live profiling. Will dump `mem` and `cpu` Go profiling files for later inspection
- [#109](https://github.com/wallix/awless/issues/109): Support caching of STS credentials for Multi-Factor Authentication.
- [#126](https://github.com/wallix/awless/issues/126): Flag `--no-alias` in `awless show` force the display of IDs in relations.
- [#126](https://github.com/wallix/awless/issues/126): Reverse sorting when listing resources with flag `--reverse`
- [#120](https://github.com/wallix/awless/issues/120): Profile info is now included in execution logs and appended when suggesting revert action
- [#82](https://github.com/wallix/awless/issues/82): Better template TAB completion (e.g. complete list of parameters)


### AWS Services

- Instance Profiles: List them; attach them to an instance. Ex: `attach instanceprofile name=...`, `awless ls instanceprofiles`
- Replace in one command an InstanceProfile on a given instance with the `replace=true` param. Ex: `attach instanceprofile .... replace=true`
- Update Route53 records with `awless update record`

### Fixes

- [#116](https://github.com/wallix/awless/issues/116) No more sync Out Of Memory

## v0.1.1 [2017-07-06]

### Features

- Detach/Attach rapidly AWS policies to user, group or role with: `attach policy service=ec2 access=readonly group=sysadmin`. More info with `awless attach policy -h`
- Better template TAB completion: suggest on properties, suggest nothing if not relevant
- Create access keys: prompt user to potentially store them locally under a specific profile
- Conveniently prompting and storing locally (~/.aws/credentials) for AWS profile credentials when access keys not found
- `awless ssh`: support SSH agent thanks to @justone
- New `--port` flag for `awless ssh`: specifying non-standard SSH port thanks to @justone
- Use `--no-headers` flag in `awless list` to display the results without headers
- New flag `--values-for` in `awless show` to output machine readable values for resource properties. Ex: `awless show my_instance --values-for name,publicip`
- Sync works on best effort now. Meaning it does not bail out when an error happens (most often it can be an access right issues on some AWS services)
- `awless ls policies` now returns: your managed policies + all policies attached to any users, role or group
- Table display now use full terminal width when possible
- Much friendlier first install

### New AWS Services

- Support of EC2 NAT Gateways: `awless list natgateways` / `awless create/delete natgateway`
- Support [ECR](https://aws.amazon.com/ecr/) repositories and registry: `awless list repositories` / `awless create/delete repository` / `awless authenticate registry`
- Support [ECS](https://aws.amazon.com/ecs/) clusters, services, containerinstances and containers: `awless list containerclusters/containertasks/containerinstances` `awless attach/detach/delete/start/stop containertask`
- Create/Delete [ApplicationAutoScaling](http://docs.aws.amazon.com/ApplicationAutoScaling/latest/APIReference/Welcome.html) scalable target and policies: `awless create/delete appscalingtarget/appscalingpolicy`

### Bugfixes
- Template TAB completion: do not display non relevant id/name listing for each prompt
- Parse successfully template parameters starting with a digit

## v0.1.0 [2017-05-31]

## Features

- Add documentation for all template parameters (`awless create instance -h`, `awless update s3object -h`...)
- Listing with filter invalid keys: return error and help
- `awless whoami` now has flags to return specific account properties only: `--account-only`, `--id-only`, `--name-only`, `--resource-only`, `--type-only`
- Rename template parameters for standardization:
    - `delete keypair id=...` -> `delete keypair name=...`
    - `create listener target=...` -> `create listener targetgroup=...`
    - `delete database skipsnapshot=... snapshotid=...` -> `delete database skip-snapshot=... snapshot=...`
    - `delete dbsubnetgroup id=...` -> `delete dbsubnetgroup name=...`
    - `create queue maxMsgSize=... retentionPeriod=... msgWait=... redrivePolicy=... visibilityTimeout=...` -> `create queue max-msg-size=... retention-period=... msg-wait=... redrive-policy=... visibility-timeout=...`

## v0.0.25 [2017-05-26]

## Features

- [#98](https://github.com/wallix/awless/issues/98): `awless ssh` searches SSH keys in both `~/.awless/keys` and `~/.ssh` folders.
- When `awless ssh` in an instance, you can now specify only `-i keyname`, if the key is stored in `~/.awless/keys` or `~/.ssh`.
- [#99](https://github.com/wallix/awless/issues/99): Suggesting the right command when typing `awless create instance ID` or `awless create ID` rather than `awless create instance id=ID`
- Use a s3 bucket as a public website with `awless update bucket name=my-bucket-name public-website=true`
- Set/update buckets or s3objects predefined ACL (private / public-read / public-read-write / bucket-owner-read...): `awless update s3object acl=public-read`
- List CloudFront distributions: `awless list distributions`
- Create/Update/Check/Delete a CloudFront distribution: `awless create/update/check/delete distribution`
- List CloudFormation stacks: `awless list stacks`
- Create/Update/Delete a CloudFormation stack: `awless create/delete stack`
- `awless log --raw-json` shows the full info stored on template execution (context, fillers used, region, ...). Typically this contextual info can be reused for replay and updates of templates

## v0.0.24 [2017-05-22]

### Features

- Template author is now persisted in awless log using the caller identity
- [#93](https://github.com/wallix/awless/issues/93): Supporting EC2 tags: syncing locally; filtering in `awless list` with --tag, --tag-value, --tag-key
- [#84](https://github.com/wallix/awless/issues/84): Create AMI by importing VM image from S3: `awless import image bucket=my-bucket s3object=my-object`. Add template to create AMI from local VM file (OVA, VMDK ...): `awless run repo:upload_image`.
- Listing pending import image tasks with `awless list importimagetasks`
- Deleting images and optionally its related snapshots `awless delete image delete-snapshots=true`
- Create/Update/Delete login profiles (AWS Console credentials): `awless create/update/delete loginprofile username=...`
- Autowrapping results in tables when too long for `awless list`. No longer truncate results in `--format csv/tsv/json`
- Adjust the width of table columns to the terminal width in `awless show`
- Using local EC2 metadata to set region when installing awless on an EC2 instance
- [#94](https://github.com/wallix/awless/issues/94): Add short flags for `--aws-profile`: `-p` and `--aws-region`: `-r`

### Bugfixes

- Listing in CSV: remove extra spaces; proper listing in TSV (only 1 tab separator)
- Avoid double sync on first install due to pre defined default region value us-east-1
- [#92](https://github.com/wallix/awless/issues/92): Impossible to set a region in config when `aws.region` was empty
- [#89](https://github.com/wallix/awless/issues/89): Fix `awless whoami` when using STS credentials.

## v0.0.23 [2017-05-05]

### Features

- Create and attach role to a user or resource (instance, ...). See an [example](https://github.com/wallix/awless-templates#role-for-resource)
- Get my IP as seen by AWS: `awless whoami --ip-only`. Example: `awless create securitygroup ... cidr=$(awless whoami --ip-only)/32 ...`
- [#86](https://github.com/wallix/awless/issues/86): SSH using private IP with `--private` flag. Thanks @padilo.
- `awless ssh` now checks the remote host public key before connecting. Check can be disabled with the (insecure) `--disable-strict-host-keychecking` flag.
- [#74](https://github.com/wallix/awless/issues/74): support of encrypted SSH keys for generation `awless create keypair encrypted=true` and in `awless ssh`.
- Better documentation of [awless-templates](https://github.com/wallix/awless-templates); listing remote templates in awless with `awless run --list`.
- Friendlier (using units: B, K, M, G) display for storage size (s3objects, volumes, lambda functions)
- Better help for template parameters (ex: `awless create loadbalancer -h`)
- Create/delete and list Lambda functions: `awless list functions` / `awless create/delete function`
- Create/delete/attach/detach and list elastic IPs: `awless list elasticips` / `awless create/delete/attach/detach elasticip`
- Create/delete and list volume snapshots: `awless list snapshots` / `awless create/delete snapshot`
- Create/delete and list autoscaling launch configurations, scaling policies and scaling groups: `awless create/delete launchconfiguration/scalingpolicy/scalinggroup`. See an [example](https://github.com/wallix/awless-templates/#group-of-instances-scaling-with-cpu-consumption)
- Create/delete/start/stop/attach/detach and list cloudwatch alarms. List cloudwatch metrics: `awless list alarms/metrics`
- List EC2 images (AMIs) of which you are the owner: `awless list images`
- Copy an EC2 image from a given region to the current region: `awless copy image name=... source-id=... source-region=...`
- List your IAM access keys: `awless list accesskeys`

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
