# Tracklog

**Tracklog** is a web application for managing GPX track files written in Go.

![Logs screenshot](doc/logs.jpg)

![Log screenshot](doc/log.jpg)

## Installation

First, make sure you have Go and Node.js installed.

To build the JavaScript and CSS assets, run:

    npm install
    npm run build

Now, fetch dependency packages and build the command line programs:
    
    go get ./...
    (cd cmd/server && go build)
    (cd cmd/import && go build)

Create and initialize a new Postgres database, which will also create a new user
with both username and password set to *admin*:

    createdb tracklog
    psql tracklog < db/postgres.sql

Start the server and point your browser to http://localhost:8080/:

    cmd/server/server -config config.json

You can batch-import your GPX files with:

    cmd/import/import -config config.json -user admin /path/to/your/*.gpx

## License

Tracklog is licensed under the MIT license.
