package sqlbuilder

import (
	"reflect"
	"testing"
)

type customer struct {
	ID    int
	Name  string
	Phone *string
}

func TestSimpleSelect(t *testing.T) {
	c := customer{}

	query, _, dest := Select().
		Dialect(MySQL).
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Map("1+1 AS two", nil).
		Build()

	expectedQuery := "SELECT id, name, phone, 1+1 AS two FROM customers"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedDest := []interface{}{&c.ID, &c.Name, &c.Phone, &nullDest}
	if !reflect.DeepEqual(dest, expectedDest) {
		t.Errorf("bad dest: %v", dest)
	}
}

func TestSimpleSelectWithLimitOffset(t *testing.T) {
	c := customer{}

	query, _, dest := Select().
		Dialect(MySQL).
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Limit(5).
		Offset(10).
		Build()

	expectedQuery := "SELECT id, name, phone FROM customers LIMIT 5 OFFSET 10"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedDest := []interface{}{&c.ID, &c.Name, &c.Phone}
	if !reflect.DeepEqual(dest, expectedDest) {
		t.Errorf("bad dest: %v", dest)
	}
}

func TestSimpleSelectWithJoins(t *testing.T) {
	c := customer{}

	query, _, _ := Select().
		Dialect(MySQL).
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Join("INNER JOIN orders ON orders.customer_id = customers.id").
		Join("LEFT JOIN items ON items.order_id = orders.id").
		Build()

	expectedQuery := "SELECT id, name, phone FROM customers INNER JOIN orders ON orders.customer_id = customers.id LEFT JOIN items ON items.order_id = orders.id"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}
}

func TestSelectWithWhereMySQL(t *testing.T) {
	c := customer{}

	query, args, _ := Select().
		Dialect(MySQL).
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Where("id = ? AND name IS NOT NULL", 9).
		Build()

	expectedQuery := "SELECT id, name, phone FROM customers WHERE (id = ? AND name IS NOT NULL)"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{9}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestSelectWithGroupMySQL(t *testing.T) {
	var count uint
	query, _, _ := Select().Dialect(MySQL).From("customers").Map("COUNT(*)", &count).Group("city").Build()
	expectedQuery := "SELECT COUNT(*) FROM customers GROUP BY city"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}
}

func TestSelectWithWherePostgres(t *testing.T) {
	c := customer{}

	query, args, _ := Select().
		Dialect(Postgres).
		From("customers").
		Map("id", &c.ID).
		Map("name", &c.Name).
		Map("phone", &c.Phone).
		Where("id = ? AND name IS NOT NULL", 9).
		Build()

	expectedQuery := `SELECT id, name, phone FROM customers WHERE (id = $1 AND name IS NOT NULL)`
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{9}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}
