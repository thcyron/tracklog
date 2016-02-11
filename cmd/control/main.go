package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/thcyron/tracklog/pkg/config"
	"github.com/thcyron/tracklog/pkg/db"
)

var (
	conf     *config.Config
	database db.DB
)

func init() {
	log.SetFlags(0)
}

func main() {
	flags := flag.NewFlagSet("tracklog-control", flag.ExitOnError)
	flags.Usage = usage
	configFile := flags.String("config", "config.toml", "path to config file")
	flags.Parse(os.Args[1:])
	if flags.NArg() == 0 {
		usage()
	}

	cnf, err := readConfig(*configFile)
	if err != nil {
		log.Fatalln(err)
	}
	conf = cnf

	dbase := db.Driver(conf.DB.Driver)
	if dbase == nil {
		log.Fatalf("unknown database driver %s\n", conf.DB.Driver)
	}
	if err := dbase.Open(conf.DB.DSN); err != nil {
		log.Fatalf("cannot open database: %s\n", err)
	}
	database = dbase

	args := flags.Args()
	switch args[0] {
	case "user":
		userCmd(args[1:])
	case "import":
		importCmd(args[1:])
	default:
		usage()
	}
}

func readConfig(path string) (*config.Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return config.Read(f)
}

func usage() {
	fmt.Fprint(os.Stderr, `usage: tracklog-control [options] <command> [arg...]

Options:

  -config <path>   Path to config file (default: config.toml)

Commands:

  user      Manage user accounts
  import    Import .gpx files
`)

	os.Exit(2)
}
