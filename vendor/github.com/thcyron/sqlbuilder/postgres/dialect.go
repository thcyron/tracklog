package postgres

import "strconv"

// Dialect is a dialect for the PostgreSQL DBMS.
type Dialect struct{}

// Placeholder implements the sqlbuilder.Dialect interface.
func (p Dialect) Placeholder(idx int) string { return "$" + strconv.Itoa(idx+1) }
