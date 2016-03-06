package mysql

// Dialect is a dialect for the MySQL DBMS.
type Dialect struct{}

// Placeholder implements the sqlbuilder.Dialect interface.
func (d Dialect) Placeholder(idx int) string { return "?" }
