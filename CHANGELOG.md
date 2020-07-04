# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
