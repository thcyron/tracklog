package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/thcyron/sqlbuilder"
	"github.com/thcyron/tracklog/pkg/models"
)

type Postgres struct {
	db *sql.DB
}

func (d *Postgres) Open(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *Postgres) UserByID(id int) (*models.User, error) {
	user := new(models.User)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"user"`).
		Map(`"id"`, &user.ID).
		Map(`"username"`, &user.Username).
		Map(`"password"`, &user.Password).
		Map(`"password_version"`, &user.PasswordVersion).
		Where(`"id" = ?`, id).
		Build()

	err := d.db.QueryRow(query, args...).Scan(dest...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *Postgres) UserByUsername(username string) (*models.User, error) {
	user := new(models.User)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"user"`).
		Map(`"id"`, &user.ID).
		Map(`"username"`, &user.Username).
		Map(`"password"`, &user.Password).
		Map(`"password_version"`, &user.PasswordVersion).
		Where(`"username" = ?`, username).
		Build()

	err := d.db.QueryRow(query, args...).Scan(dest...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *Postgres) AddUser(user *models.User) error {
	var id int

	query, args, dest := sqlbuilder.Insert().
		Dialect(sqlbuilder.Postgres).
		Into(`"user"`).
		Set(`"username"`, user.Username).
		Set(`"password"`, user.Password).
		Return(`"id"`, &id).
		Build()

	err := d.db.QueryRow(query, args...).Scan(dest...)
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (d *Postgres) UpdateUser(user *models.User) error {
	query, args := sqlbuilder.Update().
		Dialect(sqlbuilder.Postgres).
		Table(`"user"`).
		Set(`"username"`, user.Username).
		Set(`"password"`, user.Password).
		Set(`"password_version"`, user.PasswordVersion).
		Build()
	_, err := d.db.Exec(query, args...)
	return err
}

func (d *Postgres) DeleteUser(user *models.User) error {
	_, err := d.db.Exec(`DELETE FROM "user" WHERE "id" = $1`, user.ID)
	return err
}

func (d *Postgres) RecentUserLogs(user *models.User, count int) ([]*models.Log, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // read-only transaction

	var (
		log  models.Log
		logs []*models.Log
	)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"log"`).
		Map(`"id"`, &log.ID).
		Map(`"user_id"`, &log.UserID).
		Map(`"name"`, &log.Name).
		Map(`"start"`, &log.Start).
		Map(`"end"`, &log.End).
		Map(`"duration"`, &log.Duration).
		Map(`"distance"`, &log.Distance).
		Where(`"user_id" = ?`, user.ID).
		Order(`"created" DESC`).
		Limit(count).
		Build()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		l := new(models.Log)
		*l = log
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, log := range logs {
		if err := d.getLogTags(tx, log); err != nil {
			return nil, err
		}
	}

	return logs, nil
}

func (d *Postgres) UserLogYears(user *models.User) ([]int, error) {
	var (
		years []int
		year  int
	)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"log"`).
		Map(`DISTINCT EXTRACT(YEAR FROM "start")`, &year).
		Where(`"user_id" = ?`, user.ID).
		Order(`EXTRACT(YEAR FROM "start") ASC`).
		Build()

	rows, err := d.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		years = append(years, year)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return years, nil
}

func (d *Postgres) UserLogByID(user *models.User, id int) (*models.Log, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // read-only transaction

	log := new(models.Log)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"log"`).
		Map(`"id"`, &log.ID).
		Map(`"user_id"`, &log.UserID).
		Map(`"name"`, &log.Name).
		Map(`"start"`, &log.Start).
		Map(`"end"`, &log.End).
		Map(`"duration"`, &log.Duration).
		Map(`"distance"`, &log.Distance).
		Map(`"gpx"`, &log.GPX).
		Where(`"id" = ?`, id).
		Build()

	err = tx.QueryRow(query, args...).Scan(dest...)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := d.getLogTracks(tx, log); err != nil {
		return nil, err
	}
	if err := d.getLogTags(tx, log); err != nil {
		return nil, err
	}

	return log, nil
}

func (d *Postgres) getLogTracks(tx *sql.Tx, log *models.Log) error {
	var (
		track  models.Track
		tracks []*models.Track
	)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"track"`).
		Map(`"id"`, &track.ID).
		Map(`"log_id"`, &track.LogID).
		Map(`COALESCE("name", '')`, &track.Name).
		Map(`"start"`, &track.Start).
		Map(`"end"`, &track.End).
		Map(`"duration"`, &track.Duration).
		Map(`"distance"`, &track.Distance).
		Where(`"log_id" = ?`, log.ID).
		Build()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return err
		}
		t := new(models.Track)
		*t = track
		tracks = append(tracks, t)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	for _, track := range tracks {
		if err := d.getTrackPoints(tx, track); err != nil {
			return err
		}
	}

	log.Tracks = tracks
	return nil
}

func (d *Postgres) getTrackPoints(tx *sql.Tx, track *models.Track) error {
	var (
		point  models.Point
		points []*models.Point
	)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"trackpoint"`).
		Map(`"id"`, &point.ID).
		Map(`"track_id"`, &point.TrackID).
		Map(`"point"[0]`, &point.Longitude).
		Map(`"point"[1]`, &point.Latitude).
		Map(`"time"`, &point.Time).
		Map(`"elevation"`, &point.Elevation).
		Map(`"heartrate"`, &point.Heartrate).
		Where(`"track_id" = ?`, track.ID).
		Build()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return err
		}
		p := new(models.Point)
		*p = point
		points = append(points, p)
	}
	if err := rows.Err(); err != nil {
		return err
	}

	track.Points = points
	return nil
}

func (d *Postgres) getLogTags(tx *sql.Tx, log *models.Log) error {
	var (
		tag  string
		tags []string
	)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"log_tag"`).
		Map(`"tag"`, &tag).
		Where(`"log_id" = ?`, log.ID).
		Build()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return err
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	log.Tags = tags
	return nil
}

func (d *Postgres) UserLogsByYear(user *models.User, year int) ([]*models.Log, error) {
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // read-only transaction

	var (
		log  models.Log
		logs []*models.Log
	)

	query, args, dest := sqlbuilder.Select().
		Dialect(sqlbuilder.Postgres).
		From(`"log"`).
		Map(`"id"`, &log.ID).
		Map(`"name"`, &log.Name).
		Map(`"start"`, &log.Start).
		Map(`"end"`, &log.End).
		Map(`"duration"`, &log.Duration).
		Map(`"distance"`, &log.Distance).
		Map(`"gpx"`, &log.GPX).
		Where(`"user_id" = ?`, user.ID).
		Where(`EXTRACT(YEAR FROM "start") = ?`, year).
		Order(`"start" DESC`).
		Build()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(dest...); err != nil {
			rows.Close()
			return nil, err
		}
		l := new(models.Log)
		*l = log
		logs = append(logs, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, log := range logs {
		if err := d.getLogTags(tx, log); err != nil {
			return nil, err
		}
	}

	return logs, nil
}

func (d *Postgres) AddUserLog(user *models.User, log *models.Log) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	query, args, dest := sqlbuilder.Insert().
		Dialect(sqlbuilder.Postgres).
		Into(`"log"`).
		Set(`"user_id"`, user.ID).
		Set(`"start"`, log.Start).
		Set(`"end"`, log.End).
		Set(`"duration"`, log.Duration).
		Set(`"distance"`, log.Distance).
		Set(`"name"`, log.Name).
		Set(`"gpx"`, log.GPX).
		Return(`"id"`, &log.ID).
		Build()

	if err := tx.QueryRow(query, args...).Scan(dest...); err != nil {
		tx.Rollback()
		return err
	}

	for _, track := range log.Tracks {
		if err := d.addLogTrack(tx, log, track); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (d *Postgres) addLogTrack(tx *sql.Tx, log *models.Log, track *models.Track) error {
	var name *string
	if track.Name != "" {
		name = &track.Name
	}
	query, args, dest := sqlbuilder.Insert().
		Dialect(sqlbuilder.Postgres).
		Into(`"track"`).
		Set(`"log_id"`, log.ID).
		Set(`"name"`, name).
		Set(`"start"`, track.Start).
		Set(`"end"`, track.End).
		Set(`"duration"`, track.Duration).
		Set(`"distance"`, track.Distance).
		Return(`"id"`, &track.ID).
		Build()

	if err := tx.QueryRow(query, args...).Scan(dest...); err != nil {
		return err
	}

	for _, point := range track.Points {
		if err := d.addTrackPoint(tx, track, point); err != nil {
			return err
		}
	}

	return nil
}

func (d *Postgres) addTrackPoint(tx *sql.Tx, track *models.Track, point *models.Point) error {
	query, args, dest := sqlbuilder.Insert().
		Dialect(sqlbuilder.Postgres).
		Into(`"trackpoint"`).
		Set(`"track_id"`, track.ID).
		SetSQL(`"point"`, fmt.Sprintf("point(%f,%f)", point.Longitude, point.Latitude)).
		Set(`"time"`, point.Time).
		Set(`"elevation"`, point.Elevation).
		Set(`"heartrate"`, point.Heartrate).
		Return(`"id"`, &point.ID).
		Build()

	if err := tx.QueryRow(query, args...).Scan(dest...); err != nil {
		return err
	}
	return nil
}

func (d *Postgres) UpdateLog(log *models.Log) error {
	tx, err := d.db.Begin()
	if err != nil {
		return nil
	}

	query, args := sqlbuilder.Update().
		Dialect(sqlbuilder.Postgres).
		Table(`"log"`).
		Set(`"name"`, log.Name).
		Where(`"id" = ?`, log.ID).
		Build()

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	if err := d.replaceLogTags(tx, log); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (d *Postgres) replaceLogTags(tx *sql.Tx, log *models.Log) error {
	_, err := tx.Exec(`DELETE FROM "log_tag" WHERE "log_id" = $1`, log.ID)
	if err != nil {
		return err
	}

	for _, tag := range log.Tags {
		query, args, _ := sqlbuilder.Insert().
			Dialect(sqlbuilder.Postgres).
			Into(`"log_tag"`).
			Set(`"log_id"`, log.ID).
			Set(`"tag"`, tag).
			Build()
		_, err = tx.Exec(query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Postgres) DeleteLog(log *models.Log) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM "log" WHERE "id" = $1`, log.ID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
