package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/thcyron/tracklog/config"
	"github.com/thcyron/tracklog/db"
	"github.com/thcyron/tracklog/models"
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
	conf, err := config.Read(f)
	if err != nil {
		die("cannot read config file: %s", err)
	}

	DB := db.Driver(conf.DB.Driver)
	if DB == nil {
		die("unknown database driver %s", conf.DB.Driver)
	}
	if err := DB.Open(conf.DB.DSN); err != nil {
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

func importFile(DB db.DB, user *models.User, fileName string) error {
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
	log, err := models.NewLog(name, data)
	if err != nil {
		return err
	}

	return DB.AddUserLog(user, log)
}

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "import: "+fmt.Sprintf(format, args...)+"\n")
	os.Exit(1)
}
