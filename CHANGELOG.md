# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v0.4.4]

### Changed

- The default AllowedIPs setting was changed from "0.0.0.0/1, 128.0.0.0/1" to "0.0.0.0/0".

## [v0.4.3]

### Changed

- The device list on the website now updates a little less frequently.
- The device list now always shows the "last seen" field to hopefully
  better reflect what the "connected" status means.
- The metadata scraping loop has been updated to be more efficient when
  there are many disconnected peers compared to connected peers.
- The metadata scraping algorithm is now more friendly for HA deployments.

## [v0.4.2]

### Bug Fixes

- The vpn Allowed IPs setting is now correctly enforced.

## [v0.4.1]

### Bug Fixes

- Fixed a bug that caused devices to get disconnected intermittently
- The helm template now respects the "replicas" value

## [v0.4.0]

### Added

- High availability (HA) is now supported when using the `postgresql://` storage backend.
  You can now deploy multiple replicas of wg-access-server pointing to the same Postgres DB.
- The wireguard service can now be disabled via the config file. Helpful for developing
  on Mac and Windows.

### Removed

- The `file://` storage backend was deprecated in v0.3.0 and has now been removed.
  See the v0.3.0 changelog entry for more information about migrating your data.

## [v0.3.0]

### Added

- arm64 and arm/v7 docker image support + github actions thanks to [@timtorChen](https://github.com/Place1/wg-access-server/pull/73)

### Changed

- the wireguard private key is now required when the storage backend is persistent (i.e. not `memory://`)
- configuration flags, environment variables and file properties have been refactored for consistency
  * all configuration file properties (excluding auth providers) can now be set via flags and environment variables
  * all environment variables are prefixed with `WG_` to avoid collisions in hosted environments like Kubernetes
  * all flags & environment variables are named consistently
  * **breaking:** no functionality has been removed but you'll need to update any flags/envvars that you're using

### Deprecations

- deprecated support for having no admin account
  * a config error will be thrown in v0.4.0 if an admin account is not configured
  * see the README.md for examples on setting the admin account
- deprecated `file://` storage in favour of `sqlite3://`
  * will be removed in v0.4.0
  * there is now a storage `migrate` command that you can use to move your data to a different storage backend
  * see the docs for migrating your data: https://place1.github.io/wg-access-server/3-storage/#example-file-to-sqlite3

## [0.2.5]

### Added

- Admin users can now delete devices from the "all devices" page (issue [#57](https://github.com/Place1/wg-access-server/issues/57))

### Bug Fixes

- Fixes website routing to solve 404s (issue [#56](https://github.com/Place1/wg-access-server/issues/56))

## [0.2.4]

### Bug Fixes

- Improved config validation and error reporting (issue [#58](https://github.com/Place1/wg-access-server/issues/58) [#61](https://github.com/Place1/wg-access-server/issues/61))

## [0.2.3]

### Added

- Helm chart now supports configuring a LoadBalancer service for the web ui ([@nqngo](https://github.com/Place1/wg-access-server/pull/60))

## [0.2.2]

### Changed

- Changed the default "AllowedIPs" to `0.0.0.0/0`

## [0.2.1]

### Changed

- The "is connected" now shows devices as connected if they've been active within the last 3 minutes
- Improved handling of oidc/gitlab authentication with domain verification when a user hasn't set their email

## [0.2.0]

### Added

- New SQL storage backend supporting SQLite, MySQL and PostgreSQL ([@halkeye](https://github.com/Place1/wg-access-server/pull/37))
- Support for mapping claims from an OIDC auth backend to wg-access-server claims using a simple rule
  syntax ([@halkeye](https://github.com/Place1/wg-access-server/pull/39)). You can use this feature
  to decide which user has the 'admin' claim based on your own OIDC claims.
- The VPN DNS proxy feature can now be disabled using config: `dns.enabled = false`
  - When disabled the `DNS` wireguard config value will be omitted from client wg config files
  - When disabled the DNSasd proxy will not be started server-side (i.e. port 53 won't be used)
- Config options to change the web, wireguard and dns ports.
- Better instructions for connecting a linux device ([@nfg](https://github.com/Place1/wg-access-server/pull/38))
- More helm chart flexibility ([@halkeye](https://github.com/Place1/wg-access-server/pull/33))

### Changes

- The admin UI will now show the device owner's name or email if available.
- The admin UI will now show the auth provider for a given device if more than 1 auth provider is in use.
- Bug fix: upstream dns now correctly configured using resolvconf if not set in config file, flag or envvar.

### Removed

- dns port configuration was removed because wireguard client's only support port 53 for dns

### How to upgrade

- If you've been using the `storage.directory="/some/path"` config value then
  you'll need to update it to `storage=file:///some/path`
- If you've been using the `--storage-directory=/some/path` cli flag then
  you'll need to update it to `--storage="file:///some/path"`
- If you've been using the `STORAGE_DIRECTORY=/some/path` environment variable then
  you'll need to update it to `STORAGE="file:///some/path"`

## [0.1.1]

### Changes

- Helm chart bug fixes and improvements

## [0.1.0]

### Added

- Added support for an admin account. An admin can see all devices registered
  with the server.
- Added support for configuring "AllowedIPs"
- New docker compose example ([@antoniebou13](https://github.com/Place1/wg-access-server/pull/13))
- Added a helm chart
- Added a basic kubernetes quickstart.yaml manifest (based on helm template)
- Added a documentation site based on [mkdocs](https://www.mkdocs.org/). Hosted
  on github pages (still a wip!)

## [0.0.9]

### Changed

- Some UI/UX improvements

## [0.0.8]

### Added

- Added an embedded DNS proxy

### Changed

- Completely re-implemented the auth subsystem to avoid trying to integrate
  with Dex. OIDC, Gitlab and Basic auth are supported.

## [0.0.0] -> [0.0.7]

MVP :)
