package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/thcyron/tracklog/pkg/config"
	"github.com/thcyron/tracklog/pkg/db"
	"github.com/thcyron/tracklog/pkg/server"
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
	conf, err := config.Read(f)
	if err != nil {
		die("cannot read config file: %s", err)
	}
	if err := config.Check(conf); err != nil {
		die("invalid config file: %s", err)
	}

	DB := db.Driver(conf.DB.Driver)
	if DB == nil {
		die("unknown database driver %s", conf.DB.Driver)
	}
	if err := DB.Open(conf.DB.DSN); err != nil {
		die("cannot open database: %s", err)
	}

	s, err := server.New(conf, DB)
	if err != nil {
		die("%s", err)
	}

	log.Fatalln(http.ListenAndServe(conf.Server.ListenAddress, s))
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "server: "+fmt.Sprintf(format, args...)+"\n")
	os.Exit(1)
}
