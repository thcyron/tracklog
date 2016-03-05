package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestInsertMySQL(t *testing.T) {
	query, args, _ := MySQL.Insert().
		Into("customers").
		Set("name", "John").
		Set("phone", "555").
		SetSQL("created_at", "NOW()").
		Build()

	expectedQuery := "INSERT INTO customers (name, phone, created_at) VALUES (?, ?, NOW())"
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{"John", "555"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestInsertPostgres(t *testing.T) {
	query, args, _ := Postgres.Insert().
		Into("customers").
		Set("name", "John").
		Set("phone", "555").
		SetSQL("created_at", "NOW()").
		Build()

	expectedQuery := `INSERT INTO customers (name, phone, created_at) VALUES ($1, $2, NOW())`
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{"John", "555"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}
}

func TestInsertReturningPostgres(t *testing.T) {
	var id, one uint

	query, args, dest := Postgres.Insert().
		Into("customers").
		Set("name", "John").
		Set("phone", "555").
		SetSQL("created_at", "NOW()").
		Return("id", &id).
		Return("1", &one).
		Build()

	expectedQuery := `INSERT INTO customers (name, phone, created_at) VALUES ($1, $2, NOW()) RETURNING id, 1`
	if query != expectedQuery {
		t.Errorf("bad query: %s", query)
	}

	expectedArgs := []interface{}{"John", "555"}
	if !reflect.DeepEqual(args, expectedArgs) {
		t.Errorf("bad args: %v", args)
	}

	expectedDest := []interface{}{&id, &one}
	if !reflect.DeepEqual(dest, expectedDest) {
		t.Errorf("bad dest: %v", dest)
	}
}
