package db

import "github.com/thcyron/tracklog"

type DB interface {
	Open(dsn string) error

	UserByID(id int) (*tracklog.User, error)
	UserByUsername(username string) (*tracklog.User, error)

	RecentUserLogs(user *tracklog.User, count int) ([]*tracklog.Log, error)
	UserLogYears(user *tracklog.User) ([]int, error)
	UserLogByID(user *tracklog.User, id int) (*tracklog.Log, error)
	UserLogsByYear(user *tracklog.User, year int) ([]*tracklog.Log, error)
	AddUserLog(user *tracklog.User, log *tracklog.Log) error
	UpdateLog(log *tracklog.Log) error
	DeleteLog(log *tracklog.Log) error
}

func Driver(name string) DB {
	switch name {
	case "postgres":
		return new(Postgres)
	default:
		return nil
	}
}
