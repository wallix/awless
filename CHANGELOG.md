## 0.0.14 [unreleased]

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