package sqlbuilder

import (
	"strconv"
	"strings"
)

var nullDest interface{}

// Select returns a new SELECT statement with the default dialect.
func Select() SelectStatement {
	return SelectStatement{dialect: DefaultDialect}
}

// SelectStatement represents a SELECT statement.
type SelectStatement struct {
	dialect Dialect
	table   string
	selects []sel
	joins   []join
	wheres  []where
	lock    bool
	limit   *int
	offset  *int
	order   string
	group   string
	having  string
}

type sel struct {
	col  string
	dest interface{}
}

type join struct {
	sql  string
	args []interface{}
}

// Dialect returns a new statement with dialect set to 'dialect'.
func (s SelectStatement) Dialect(dialect Dialect) SelectStatement {
	s.dialect = dialect
	return s
}

// From returns a new statement with the table to select from set to 'table'.
func (s SelectStatement) From(table string) SelectStatement {
	s.table = table
	return s
}

// Map returns a new statement with column 'col' selected and scanned
// into 'dest'. 'dest' may be nil if the value should not be scanned.
func (s SelectStatement) Map(col string, dest interface{}) SelectStatement {
	if dest == nil {
		dest = nullDest
	}
	s.selects = append(s.selects, sel{col, dest})
	return s
}

// Join returns a new statement with JOIN expression 'sql'.
func (s SelectStatement) Join(sql string, args ...interface{}) SelectStatement {
	s.joins = append(s.joins, join{sql, args})
	return s
}

// Where returns a new statement with condition 'cond'. Multiple conditions
// are combined with AND.
func (s SelectStatement) Where(cond string, args ...interface{}) SelectStatement {
	s.wheres = append(s.wheres, where{cond, args})
	return s
}

// Limit returns a new statement with the limit set to 'limit'.
func (s SelectStatement) Limit(limit int) SelectStatement {
	s.limit = &limit
	return s
}

// Offset returns a new statement with the offset set to 'offset'.
func (s SelectStatement) Offset(offset int) SelectStatement {
	s.offset = &offset
	return s
}

// Order returns a new statement with ordering 'order'.
// Only the last Order() is used.
func (s SelectStatement) Order(order string) SelectStatement {
	s.order = order
	return s
}

// Group returns a new statement with grouping 'group'.
// Only the last Group() is used.
func (s SelectStatement) Group(group string) SelectStatement {
	s.group = group
	return s
}

// Having returns a new statement with HAVING condition 'having'.
// Only the last Having() is used.
func (s SelectStatement) Having(having string) SelectStatement {
	s.having = having
	return s
}

// Lock returns a new statement with FOR UPDATE locking.
func (s SelectStatement) Lock() SelectStatement {
	s.lock = true
	return s
}

// Build builds the SQL query. It returns the query, the argument slice,
// and the destination slice.
func (s SelectStatement) Build() (query string, args []interface{}, dest []interface{}) {
	var cols []string
	idx := 0

	if len(s.selects) > 0 {
		for _, sel := range s.selects {
			cols = append(cols, sel.col)
			if sel.dest == nil {
				dest = append(dest, &nullDest)
			} else {
				dest = append(dest, sel.dest)
			}
		}
	} else {
		cols = append(cols, "1")
		dest = append(dest, &nullDest)
	}
	query = "SELECT " + strings.Join(cols, ", ") + " FROM " + s.table

	for _, join := range s.joins {
		sql := join.sql
		for _, arg := range join.args {
			sql = strings.Replace(sql, "?", s.dialect.Placeholder(idx), 1)
			idx++
			args = append(args, arg)
		}
		query += " " + sql
	}

	if len(s.wheres) > 0 {
		var sqls []string
		for _, where := range s.wheres {
			sql := "(" + where.sql + ")"
			for _, arg := range where.args {
				sql = strings.Replace(sql, "?", s.dialect.Placeholder(idx), 1)
				idx++
				args = append(args, arg)
			}
			sqls = append(sqls, sql)
		}
		query += " WHERE " + strings.Join(sqls, " AND ")
	}

	if s.order != "" {
		query += " ORDER BY " + s.order
	}

	if s.group != "" {
		query += " GROUP BY " + s.group
	}

	if s.having != "" {
		query += " HAVING " + s.having
	}

	if s.limit != nil {
		query += " LIMIT " + strconv.Itoa(*s.limit)
	}

	if s.offset != nil {
		query += " OFFSET " + strconv.Itoa(*s.offset)
	}

	if s.lock {
		query += " FOR UPDATE"
	}

	return
}
