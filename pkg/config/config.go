package config

import (
	"encoding/json"
	"io"
)

type Config struct {
	Server struct {
		Development       bool   `json:"development"`
		ListenAddress     string `json:"listen_address"`
		CSRFAuthKey       string `json:"csrf_auth_key"`
		SigningKey        string `json:"signing_key"`
		MapboxAccessToken string `json:"mapbox_access_token"`
	} `json:"server"`
	DB struct {
		Driver string `json:"driver"`
		DSN    string `json:"dsn"`
	} `json:"db"`
}

func Read(r io.Reader) (*Config, error) {
	config := new(Config)
	if err := json.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}
