package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/thcyron/tracklog/pkg/models"
)

func importUsage() {
	fmt.Fprint(os.Stderr, `usage: tracklog-control import [options] <file...>

Options:

  -user <username>    User to import files for`)
	os.Exit(2)
}

func importCmd(args []string) {
	flags := flag.NewFlagSet("tracklog-control import", flag.ExitOnError)
	flags.Usage = usage
	username := flags.String("user", "", "username")
	flags.Parse(args)
	if *username == "" || flags.NArg() == 0 {
		usage()
	}

	user, err := database.UserByUsername(*username)
	if err != nil {
		log.Fatalf("cannot load user: %s\n", err)
	}
	if user == nil {
		log.Fatalf("no such user: %s\n", *username)
	}

	for _, fileName := range flags.Args() {
		if err := importFile(user, fileName); err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", fileName, err)
		} else {
			fmt.Println(fileName)
		}
	}
}

func importFile(user *models.User, fileName string) error {
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

	return database.AddUserLog(user, log)
}
