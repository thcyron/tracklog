package db

import "github.com/thcyron/tracklog/models"

type DB interface {
	Open(dsn string) error

	UserByID(id int) (*models.User, error)
	UserByUsername(username string) (*models.User, error)

	RecentUserLogs(user *models.User, count int) ([]*models.Log, error)
	UserLogYears(user *models.User) ([]int, error)
	UserLogByID(user *models.User, id int) (*models.Log, error)
	UserLogsByYear(user *models.User, year int) ([]*models.Log, error)
	AddUserLog(user *models.User, log *models.Log) error
	UpdateLog(log *models.Log) error
	DeleteLog(log *models.Log) error
}

func Driver(name string) DB {
	switch name {
	case "postgres":
		return new(Postgres)
	default:
		return nil
	}
}
