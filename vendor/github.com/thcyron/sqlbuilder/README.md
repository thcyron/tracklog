sqlbuilder
==========

[![Travis CI status](https://api.travis-ci.org/thcyron/sqlbuilder.svg)](https://travis-ci.org/thcyron/sqlbuilder)

`sqlbuilder` is a Go library for building SQL queries.

The master branch tracks version 3. The latest stable version is
[2.0.0](https://github.com/thcyron/sqlbuilder/tree/v2.0.0/).

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
query, args, dest := sqlbuilder.Insert().
        Into("customers").
        Set("name", "John").
        Set("phone", "555").
        Build()
res, err := db.Exec(query, args...)
```

**UPDATE**

```go
query, args := sqlbuilder.Update().
        Table("customers").
        Set("name", "John").
        Set("phone", "555").
        Where("id = ?", 1).
        Build()
res, err := db.Exec(query, args...)
```

**DELETE**

```go
query, args := sqlbuilder.Delete().
    From("customers").
    Where("name = ?", "John").
    Build()
res, err := db.Exec(query, args...)
```

Dialects
--------

`sqlbuilder` supports building queries for MySQL, SQLite, and Postgres databases. You
can set the default dialect with:

```go
sqlbuilder.DefaultDialect = sqlbuilder.Postgres
sqlbuilder.Select().From("...")...
```

Or you can specify the dialect explicitly:

```go
sqlbuilder.Select().Dialect(sqlbuilder.Postgres).From("...")...
```

Documentation
-------------

Documentation is available at [GoDoc](https://godoc.org/github.com/thcyron/sqlbuilder).

License
-------

`sqlbuilder` is licensed under the MIT License.
