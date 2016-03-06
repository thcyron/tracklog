package sqlbuilder

import "strings"

// Insert returns a new INSERT statement with the default dialect.
func Insert() InsertStatement {
	return InsertStatement{dialect: DefaultDialect}
}

type insertSet struct {
	col string
	arg interface{}
	raw bool
}

type insertRet struct {
	sql  string
	dest interface{}
}

// InsertStatement represents an INSERT statement.
type InsertStatement struct {
	dialect Dialect
	table   string
	sets    []insertSet
	rets    []insertRet
}

// Dialect returns a new statement with dialect set to 'dialect'.
func (s InsertStatement) Dialect(dialect Dialect) InsertStatement {
	s.dialect = dialect
	return s
}

// Into returns a new statement with the table to insert into set to 'table'.
func (s InsertStatement) Into(table string) InsertStatement {
	s.table = table
	return s
}

// Set returns a new statement with column 'col' set to value 'val'.
func (s InsertStatement) Set(col string, val interface{}) InsertStatement {
	s.sets = append(s.sets, insertSet{col, val, false})
	return s
}

// SetSQL returns a new statement with column 'col' set to the raw SQL expression 'sql'.
func (s InsertStatement) SetSQL(col, sql string) InsertStatement {
	s.sets = append(s.sets, insertSet{col, sql, true})
	return s
}

// Return returns a new statement with a RETURNING clause.
func (s InsertStatement) Return(col string, dest interface{}) InsertStatement {
	s.rets = append(s.rets, insertRet{sql: col, dest: dest})
	return s
}

// Build builds the SQL query. It returns the SQL query and the argument slice.
func (s InsertStatement) Build() (query string, args []interface{}, dest []interface{}) {
	var cols, vals []string
	idx := 0

	for _, set := range s.sets {
		cols = append(cols, set.col)

		if set.raw {
			vals = append(vals, set.arg.(string))
		} else {
			args = append(args, set.arg)
			vals = append(vals, s.dialect.Placeholder(idx))
			idx++
		}
	}

	query = "INSERT INTO " + s.table + " (" + strings.Join(cols, ", ") + ") VALUES (" + strings.Join(vals, ", ") + ")"

	if len(s.rets) > 0 {
		var args []string
		for _, ret := range s.rets {
			args = append(args, ret.sql)
			dest = append(dest, ret.dest)
		}
		query += " RETURNING " + strings.Join(args, ", ")
	}

	return
}
