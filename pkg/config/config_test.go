package config

import (
	"bytes"
	"testing"
)

var data = []byte(`[server]
development = true
listen_address = ":8080"
csrf_auth_key = "secr"
signing_key = "secret"
mapbox_access_token = "abc"

[db]
driver = "postgres"
dsn = "dbname=tracklog"
`)

func TestRead(t *testing.T) {
	c, err := Read(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if expected := true; c.Server.Development != expected {
		t.Errorf("expected Server.Development = %v; got %v", expected, c.Server.Development)
	}
	if expected := ":8080"; c.Server.ListenAddress != expected {
		t.Errorf("expected Server.ListenAddress = %q; got %q", expected, c.Server.ListenAddress)
	}
	if expected := "secr"; c.Server.CSRFAuthKey != expected {
		t.Errorf("expected Server.CSRFAuthKey = %q; got %q", expected, c.Server.CSRFAuthKey)
	}

	if expected := "postgres"; c.DB.Driver != expected {
		t.Errorf("expected DB.Driver = %q; got %q", expected, c.DB.Driver)
	}
	if expected := "dbname=tracklog"; c.DB.DSN != expected {
		t.Errorf("expected DB.DSN = %q; got %q", expected, c.DB.DSN)
	}
}
