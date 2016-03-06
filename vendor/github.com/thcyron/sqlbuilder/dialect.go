package sqlbuilder

import (
	"github.com/thcyron/sqlbuilder/mysql"
	"github.com/thcyron/sqlbuilder/postgres"
)

// Dialect represents a SQL dialect.
type Dialect interface {
	// Placeholder returns the placeholder binding string for parameter at index idx.
	Placeholder(idx int) string
}

var (
	MySQL    mysql.Dialect    // MySQL
	SQLite   mysql.Dialect    // SQLite (same as MySQL)
	Postgres postgres.Dialect // Postgres
)

var DefaultDialect = MySQL // Default dialect
