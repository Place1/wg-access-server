# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [master]

### Added

- The VPN DNS proxy feature can now be disabled using config: `dns.enabled = false`
  - When disabled the `DNS` wireguard config value will be omitted from client wg config files
  - When disabled the DNS proxy will not be started server-side (i.e. port 53 won't be used)

### Changes

- The admin UI will now show the device owner's name or email if available.
- The admin UI will now show the auth provider for a given device if more than 1 auth provider is in use.
- Bug fix: upstream dns now correctly configured using resolvconf if not set in config file, flag or envvar.

## [0.1.1]

### Changes

- Helm chart bug fixes and improvements

## [0.1.0]

### Added

- Added support for an admin account. An admin can see all devices registered
  with the server.
- Added support for networking isolation modes. You can now allow/deny VPN LAN,
  Server LAN and internet traffic. Selective network CIDRs can be white listed.
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
