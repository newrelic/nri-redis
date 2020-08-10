# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 1.5.0 (2020-08-10)
### Added
- `USE_UNIX_SOCKET` configuration option (default: `false`). Adds the `UnixSocketPath` value to the entity. This helps to uniquely identify your entities when you're monitoring more than one Redis instance on the same host using Unix sockets.
  
## 1.5.0 (2020-01-13)
### Added
- `CONFIG_INVENTORY` configuration option (default: true). Set it to `false` to avoid invoking the Redis
  `CONFIG` command when querying for inventory data. This option is useful in environments where the Redis
  `CONFIG` command is prohibited (e.g. AWS ElastiCache).

### Changed
- Avoid invoking the `CONFIG` command if the Inventory data is skipped.

## 1.3.0 (2019-11-18)
### Changed
- Renamed the integration executable from nr-redis to nri-redis in order to be consistent with the package naming. **Important Note:** if you have any security module rules (eg. SELinux), alerts or automation that depends on the name of this binary, these will have to be updated.
## 1.2.1 (2019-08-05)
## Fixed
* Omitted `masterauth` inventory entry. It is now submitted as `(omitted entry)`.

## 1.2.0 (2019-04-29)
### Added
- Upgraded to SDK v3.1.5. This version implements [the aget/integrations
  protocol v3](https://github.com/newrelic/infra-integrations-sdk/blob/cb45adacda1cd5ff01544a9d2dad3b0fedf13bf1/docs/protocol-v3.md),
  which enables [name local address replacement](https://github.com/newrelic/infra-integrations-sdk/blob/cb45adacda1cd5ff01544a9d2dad3b0fedf13bf1/docs/protocol-v3.md#name-local-address-replacement).
  and could change your entity names and alarms. For more information, refer
  to:

  - https://docs.newrelic.com/docs/integrations/integrations-sdk/file-specifications/integration-executable-file-specifications#h2-loopback-address-replacement-on-entity-names
  - https://docs.newrelic.com/docs/remote-monitoring-host-integration://docs.newrelic.com/docs/remote-monitoring-host-integrations

## 1.1.0 (2019-04-08)
### Added
- Upgraded to SDKv3.
- Remote monitoring option. It enables monitoring multiple instances,
  more information can be found at the [official documentation page](https://docs.newrelic.com/docs/remote-monitoring-host-integrations).

## 1.0.1 (2018-09-07)
### Changed
- Update Makefile

## 1.0.0 (2018-08-02)
### Added
- Initial version: Includes non-keyspace and keyspace Metrics and Inventory data
