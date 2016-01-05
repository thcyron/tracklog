package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/thcyron/tracklog"
	"github.com/thcyron/tracklog/db"
)

var (
	configFile = flag.String("config", "config.json", "config file")
	username   = flag.String("user", "", "username")
)

func main() {
	flag.Parse()

	if *username == "" {
		die("missing user ID")
	}

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

	user, err := DB.UserByUsername(*username)
	if err != nil {
		die("cannot load user: %s", err)
	}
	if user == nil {
		die("no such user")
	}

	for _, fileName := range flag.Args() {
		if err := importFile(DB, user, fileName); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", fileName, err)
		} else {
			fmt.Println(fileName)
		}
	}
}

func importFile(DB db.DB, user *tracklog.User, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	name := strings.TrimSuffix(path.Base(fileName), ".gpx")
	log, err := tracklog.NewLog(name, data)
	if err != nil {
		return err
	}

	return DB.AddUserLog(user, log)
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "import: "+fmt.Sprintf(format, args...)+"\n")
	os.Exit(1)
}
