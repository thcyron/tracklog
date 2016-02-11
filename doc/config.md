# Config

Tracklog is configured via a TOML file.

## `server` Section

The `server` section configures the server.

| Key                   | Required | Description                                            |
| --------------------- | -------- | ------------------------------------------------------ |
| `development`         | no       | Enables development mode when set to `true`            |
| `listen_address`      | yes      | Listen address of the HTTP server                      |
| `csrf_auth_key`       | yes      | Secret key used for CSRF protection (must be 32 bytes) |
| `signing_key`         | yes      | Secret key used for signing tokens                     |
| `mapbox_access_token` | no       | Mapbox access token for Mapbox maps                    |

## `db` section

The `db` section configures the database.

| Key                   | Required | Description                                            |
| --------------------- | -------- | ------------------------------------------------------ |
| `driver`              | yes      | Driver                                                 |
| `dsn`                 | yes      | Data source name                                       |
