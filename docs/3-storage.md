# Storage

wg-access-server supports 4 storage backends.

| Backend  | Persistent | Supports HA | Use Case                                 |
| -------- | ---------- | ----------- | ---------------------------------------- |
| memory   | ❌         | ❌          | Local development                        |
| sqlite3  | ✔️         | ❌          | Production - single instance deployments |
| postgres | ✔️         | ✔️          | Production - multi instance deployments  |
| mysql    | ✔️         | ❌          | Production - single instance deployments |

## Backends

### Memory

This is the default backend if you're running the binary directly and haven't configured
another storage backend. Data will be lost between restarts. Handy for development.

### SQLite3

This is the default backend if you're running the docker container directly or using docker-compose.

The database file will be written to `/data/db.sqlite3` within the container by default.

Sqlite3 is probably the simplest storage backend to get started with because it doesn't require
any additional setup to be done. It should work out of the box and should be able to support a
large number of users & devices.

Example connection string:

- Relative path: `sqlite3://path/to/db.sqlite3`
- Absolute path: `sqlite3:///absolute/path/to/db.sqlite3`

### PostgreSQL

This backend requires an external Postgres database to be deployed.

Postgres experimentally supports highly-available deployments of wg-access-server
and is the recommended storage backend where possible.
If you have pgbouncer running in front of PostgreSQL, make sure it is running in "Session pooling" mode,
as the more aggressive modes like "Transaction Pooling" break LISTEN/NOTIFY which we rely on.

Example connection string:

- `postgresql://user:password@localhost:5432/database?sslmode=disable`

### MySQL

This backend requires an external Mysql database to be deployed. Mysql flavours should be compatible.
wg-access-server uses [this golang driver](github.com/go-sql-driver/mysql) if you want to check the
compatibility of your favorite flavour.

Example connection string:

- `mysql://user:password@localhost:3306/database?ssl-mode=disabled`

### File (removed)

The `file://` backend was deprecated in 0.3.0 and has been removed in 0.4.0

If you'd like to migrate your `file://` storage to a supported backend you must use
version 0.3.0 and then follow the migration guide below to migrate to a different storage backend.

_Note that the migration tool itself doesn't support the `file://` backend on versions
released after 0.3.0_.

## Migration Between Backends

You can migrate your registered devices between backends using the `wg-access-server migrate <src> <dest>`
command.

The migrate command was added in `v0.3.0` and is provided on a _best effort_ level. As an open source
project any community support here is warmly welcomed.

### Example: `file://` to `sqlite3://`

If you're using the now deprecated `file://` backend you can migrate to `sqlite3://` like this:

```bash
# after upgrading to place1/wg-access-server:v0.3.0
docker exec -it <container-name> wg-access-server migrate file:///data sqlite3:///data/db.sqlite3
```

If you need to do the above within a kubernetes deployment substitute `docker exec` with the equivalent
`kubectl exec` command.

The migrate command is non-destructive but it's always a good idea to take a backup of your data first!

### Example: `sqlite3://` to `postgresql://`

First you'll need to make sure your postgres server is up and that you can connect to it from your
wg-access-server container/pod/vm.

```bash
wg-access-server migrate sqlite3:///data/db.sqlite3 postgresql://user:password@localhost:5432/database?sslmode=disable
```

Remember to update your wg-access-server config to connect to postgres 😀
