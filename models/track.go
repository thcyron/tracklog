package models

import "time"

type Track struct {
	ID       int
	LogID    int
	Name     string
	Start    time.Time
	End      time.Time
	Duration uint
	Distance uint
	Points   []*Point
}
