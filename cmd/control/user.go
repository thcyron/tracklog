package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"unicode/utf8"

	"github.com/thcyron/tracklog/pkg/models"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func userUsage() {
	fmt.Fprint(os.Stderr, `usage: tracklog-control user <command> [arg...]

Commands:

  add         Add a new user
  delete      Delete a user
  password    Change a user’s password`)
	os.Exit(2)
}

func userCmd(args []string) {
	if len(args) == 0 {
		userUsage()
	}

	switch args[0] {
	case "add":
		addUserCmd(args[1:])
	case "delete":
		deleteUserCmd(args[1:])
	case "password":
		passwordUserCmd(args[1:])
	default:
		userUsage()
	}
}

func addUserUsage() {
	fmt.Fprint(os.Stderr, `usage: tracklog-control user add <username>`)
	os.Exit(2)
}

func addUserCmd(args []string) {
	if len(args) != 1 {
		addUserUsage()
	}

	username := args[0]
	if username == "" {
		addUserUsage()
	}

	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	if utf8.RuneCount(password) == 0 {
		log.Fatalln("zero-length password not allowed")
	}

	pwhash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	fmt.Print("\n")
	if err != nil {
		log.Fatalln(err)
	}

	user := &models.User{
		Username: username,
		Password: string(pwhash),
	}

	if err := database.AddUser(user); err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("user %s created\n", username)
}

func deleteUserCmd(args []string) {
	if len(args) != 1 {
		addUserUsage()
	}

	username := args[0]
	if username == "" {
		addUserUsage()
	}

	user, err := database.UserByUsername(username)
	if err != nil {
		log.Fatalln(err)
	}
	if user == nil {
		log.Fatalf("no such user: %s\n", username)
	}

	if err := database.DeleteUser(user); err != nil {
		log.Fatalln(err)
	}

	log.Printf("user %s deleted\n", username)
}

func passwordUserCmd(args []string) {
	if len(args) != 1 {
		addUserUsage()
	}

	username := args[0]
	if username == "" {
		addUserUsage()
	}

	user, err := database.UserByUsername(username)
	if err != nil {
		log.Fatalln(err)
	}
	if user == nil {
		log.Fatalf("no such user: %s\n", username)
	}

	fmt.Print("Password: ")
	password, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	if utf8.RuneCount(password) == 0 {
		log.Fatalln("zero-length password not allowed")
	}

	pwhash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	fmt.Print("\n")
	if err != nil {
		log.Fatalln(err)
	}
	user.Password = string(pwhash)

	if err := database.UpdateUser(user); err != nil {
		log.Fatalln(err)
	}

	log.Printf("password for %s changed\n", username)
}
