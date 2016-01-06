sqlbuilder
==========

[![Travis CI status](https://api.travis-ci.org/thcyron/sqlbuilder.svg)](https://travis-ci.org/thcyron/sqlbuilder)

`sqlbuilder` is a Go library for building SQL queries.

The current version is [2.0.0](https://github.com/thcyron/sqlbuilder/tree/v2.0.0/).
`sqlbuilder` follows [Semantic Versioning](http://semver.org/).

Installation
------------

    go get github.com/thcyron/sqlbuilder

Examples
--------

**SELECT**

```go
query, args, dest := sqlbuilder.Select().
        From("customers").
        Map("id", &customer.ID).
        Map("name", &customer.Name).
        Map("phone", &customer.Phone).
        Order("id DESC").
        Limit(1).
        Build()

err := db.QueryRow(query, args...).Scan(dest...)
```

**INSERT**

```go
query, args := sqlbuilder.Insert().
        Into("customers").
        Set("name", "John").
        Set("phone", "555").
        Build()
err := db.Exec(query, args...)
```

**UPDATE**

```go
query, args := sqlbuilder.Update().
        Table("customers").
        Set("name", "John").
        Set("phone", "555").
        Where("id = ?", 1).
        Build()
err := db.Exec(query, args...)
```

Supported DBMS
--------------

`sqlbuilder` supports building queries for MySQL and Postgres databases. You
can set the default DBMS used by the package-level Select, Update and Insert
functions with:

```go
sqlbuilder.DefaultDBMS = sqlbuilder.Postgres
sqlbuilder.Select().From("...")...
```

or you can specify the DBMS explicitly:

```go
sqlbuilder.Postgres.Select().From("...")...
```

You typically set the default DBMS in `init()`:

```go
func init() {
        sqlbuilder.DefaultDBMS = sqlbuilder.Postgres
}
```

Documentation
-------------

Documentation is available at [Godoc](http://godoc.org/github.com/thcyron/sqlbuilder).

License
-------

`sqlbuilder` is licensed under the MIT license.
