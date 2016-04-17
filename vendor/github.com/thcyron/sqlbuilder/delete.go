package sqlbuilder

import (
	"strings"
)

// Delete returns a new DELETE statement with the default dialect.
func Delete() DeleteStatement {
	return DeleteStatement{dialect: DefaultDialect}
}

// DeleteStatement represents a DELETE statement.
type DeleteStatement struct {
	dialect Dialect
	table   string
	wheres  []where
	args    []interface{}
}

// Dialect returns a new statement with dialect set to 'dialect'.
func (s DeleteStatement) Dialect(dialect Dialect) DeleteStatement {
	s.dialect = dialect
	return s
}

// From returns a new statement with the table to delete from set to 'table'.
func (s DeleteStatement) From(table string) DeleteStatement {
	s.table = table
	return s
}

// Where returns a new statement with condition 'cond'.
// Multiple Where() are combined with AND.
func (s DeleteStatement) Where(cond string, args ...interface{}) DeleteStatement {
	s.wheres = append(s.wheres, where{cond, args})
	return s
}

// Build builds the SQL query. It returns the query and the argument slice.
func (s DeleteStatement) Build() (query string, args []interface{}) {
	query = "DELETE FROM " + s.table

	if len(s.wheres) > 0 {
		var (
			sqls []string
			idx  int
		)
		for _, w := range s.wheres {
			sql := "(" + w.sql + ")"
			for _, arg := range w.args {
				p := s.dialect.Placeholder(idx)
				idx++
				sql = strings.Replace(sql, "?", p, 1)
				sqls = append(sqls, sql)
				args = append(args, arg)
			}
		}

		query += " WHERE " + strings.Join(sqls, " AND ")
	}

	return
}
