package sqlbuilder

import "strconv"

// DBMS represents a DBMS.
type DBMS int

const (
	MySQL    DBMS = iota // MySQL
	Postgres             // Postgres
)

// Placeholder returns the placeholder string for the given index.
func (dbms DBMS) Placeholder(idx int) string {
	switch dbms {
	case MySQL:
		return "?"
	case Postgres:
		return "$" + strconv.Itoa(idx+1)
	default:
		panic("sqlbuilder: unknown DBMS")
	}
}

// Select returns a new SELECT statement.
func (dbms DBMS) Select() SelectStatement {
	return SelectStatement{dbms: dbms}
}

// Insert returns a new INSERT statement.
func (dbms DBMS) Insert() InsertStatement {
	return InsertStatement{dbms: dbms}
}

// Update returns a new UPDATE statement.
func (dbms DBMS) Update() UpdateStatement {
	return UpdateStatement{dbms: dbms}
}

// DefaultDBMS is the DBMS used by the package-level Select, Insert, and Update functions.
var DefaultDBMS = MySQL

// Select returns a new SELECT statement using the default DBMS.
func Select() SelectStatement {
	return DefaultDBMS.Select()
}

// Insert returns a new INSERT statement using the default DBMS.
func Insert() InsertStatement {
	return DefaultDBMS.Insert()
}

// Update returns a new UPDATE statement using the default DBMS.
func Update() UpdateStatement {
	return DefaultDBMS.Update()
}
