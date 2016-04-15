package models

import (
	"bytes"
	"math"
	"time"

	"github.com/thcyron/gpx"
)

type Log struct {
	ID       int
	UserID   int
	Name     string
	Start    time.Time
	End      time.Time
	Duration uint
	Distance uint
	GPX      string
	Tracks   []*Track
	Tags     []string
}

func NewLog(name string, data []byte) (*Log, error) {
	doc, err := gpx.NewDecoder(bytes.NewReader(data)).Decode()
	if err != nil {
		return nil, err
	}

	log := &Log{
		Name:     name,
		Start:    doc.Start(),
		End:      doc.End(),
		Duration: uint(doc.Duration().Seconds()),
		Distance: uint(doc.Distance()),
		GPX:      string(data),
	}

	for _, track := range doc.Tracks {
		for _, segment := range track.Segments {
			t := &Track{
				Name:     track.Name,
				Start:    segment.Start(),
				End:      segment.End(),
				Duration: uint(segment.Duration().Seconds()),
				Distance: uint(segment.Distance()),
			}

			for _, point := range segment.Points {
				t.Points = append(t.Points, NewPoint(point))
			}

			log.Tracks = append(log.Tracks, t)
		}
	}

	return log, nil
}

// Speed returns the speed in meters per second.
func (log *Log) Speed() float64 {
	speed := float64(log.Distance) / float64(log.Duration)
	if math.IsInf(speed, 0) {
		speed = 0
	}
	return speed
}
