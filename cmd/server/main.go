package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/thcyron/tracklog"
	"github.com/thcyron/tracklog/db"
	"github.com/thcyron/tracklog/server"
)

var (
	configFile = flag.String("config", "config.json", "config file")
)

func main() {
	flag.Parse()

	f, err := os.Open(*configFile)
	if err != nil {
		die("cannot open config file: %s", err)
	}
	defer f.Close()
	config, err := tracklog.ReadConfig(f)
	if err != nil {
		die("cannot read config file: %s", err)
	}

	DB := db.Driver(config.DB.Driver)
	if DB == nil {
		die("unknown database driver %s", config.DB.Driver)
	}
	if err := DB.Open(config.DB.DSN); err != nil {
		die("cannot open database: %s", err)
	}

	s, err := server.New(config, DB)
	if err != nil {
		die("%s", err)
	}

	log.Fatalln(http.ListenAndServe(config.Server.ListenAddress, s))
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "server: "+fmt.Sprintf(format, args...)+"\n")
	os.Exit(1)
}
