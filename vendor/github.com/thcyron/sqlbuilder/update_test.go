package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestUpdateMySQL(t *testing.T) {
	query, args := MySQL.Update().
		Table("customers").
		Set("name", "John").
		Set("phone", "555").
		Build()

	expectedQuery := "UPDATE `customers` SET `name` = ?, `phone` = ?"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{"John", "555"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestUpdatePostgres(t *testing.T) {
	query, args := Postgres.Update().
		Table("customers").
		Set("name", "John").
		Set("phone", "555").
		Build()

	expectedQuery := `UPDATE "customers" SET "name" = $1, "phone" = $2`
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{"John", "555"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestUpdateWithWhereMySQL(t *testing.T) {
	query, args := MySQL.Update().
		Table("customers").
		Set("name", "John").
		Set("phone", "555").
		Where("id = ?", 9).
		Build()

	expectedQuery := "UPDATE `customers` SET `name` = ?, `phone` = ? WHERE (id = ?)"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{"John", "555", 9}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestUpdateWithWherePostgres(t *testing.T) {
	query, args := Postgres.Update().
		Table("customers").
		Set("name", "John").
		Set("phone", "555").
		Where("id = ?", 9).
		Build()

	expectedQuery := `UPDATE "customers" SET "name" = $1, "phone" = $2 WHERE (id = $3)`
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{"John", "555", 9}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestUpdateReuse(t *testing.T) {
	baseStatement := MySQL.Update().Table("customers").Set("name", "John")

	query, args := baseStatement.Set("phone", "555").Build()
	expectedQuery := "UPDATE `customers` SET `name` = ?, `phone` = ?"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}
	expectedArgs := []interface{}{"John", "555"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}

	query, args = baseStatement.Set("city", "Berlin").Build()
	expectedQuery = "UPDATE `customers` SET `name` = ?, `city` = ?"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}
	expectedArgs = []interface{}{"John", "Berlin"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}
