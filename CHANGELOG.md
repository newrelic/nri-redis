# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 1.9.1  (2022-06-23)

### Changed
 - Bump dependencies
### Added
 - Added support for more distributions:
    RHEL(EL) 9
    Ubuntu 22.04

## 1.9.0  (2022-03-08)
### Added
- Added `redis-log.yml.example` to Linux installers to help setting up log parsing.

## 1.8.2 (2021-10-20)
### Added
Added support for more distributions:
- Debian 11
- Ubuntu 20.10
- Ubuntu 21.04
- SUSE 12.15
- SUSE 15.1
- SUSE 15.2
- SUSE 15.3
- Oracle Linux 7
- Oracle Linux 8

## 1.8.1 (2021-08-27)
### Fixed

Added unit notation to the default interval in the config sample.

Added missing config parameters in the config sample.

## 1.8.0 (2021-08-27)
### Added

Moved default config.sample to [V4](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-newer-configuration-format/), added a dependency for infra-agent version 1.20.0

Please notice that old [V3](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) configuration format is deprecated, but still supported.

## 1.7.1 (2021-08-6)
### Added

A stricter validation of args was introduced without noticing that use_unix_socket was false in defaults, but true in the sample config.

There are users having use_unix_socket=true and then connecting with hostname and port, or use_unix_socket=false and then connecting with the unix socket.

Arg use_unix_socket is not used to define how to connect, but merely the entity name.:

```
Adds the UnixSocketPath value to the entity. If you are monitoring more than one Redis instance on the same host using Unix sockets, then you should set it to true.
```

## 1.7.0 (2021-08-4)
### Added
- Allows Usages of rename-command on Redis Server Installation 
- Support to IPv6 address family as hostname argument 
- Support TLS connections to Redis 

## 1.6.3 (2021-06-7)
### Added
- Added support for ARM and ARM64.

## 1.6.2 (2021-04-22)
### Added
- Upgrade dependency manager to use go mod
- Bumps sdk to v3.6.7 solving multi-instance storage overlapping
- Bumps redigo to v1.8 (redis client library)

## 1.6.1 (2021-03-24)
### Added
- Add arm packages and binaries

## 1.6.0 (2020-10-29)
### Added
- Add print integration version from cli using  `-show_version` flag

## 1.5.1 (2020-09-26)
### Added
- `maxmemoryBytes` metric from the Redis Info is.

## 1.5.0 (2020-08-10)
### Added
- `USE_UNIX_SOCKET` configuration option (default: `false`). Adds the `UnixSocketPath` value to the entity. This helps to uniquely identify your entities when you're monitoring more than one Redis instance on the same host using Unix sockets.
  
## 1.4.0 (2020-01-13)
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
