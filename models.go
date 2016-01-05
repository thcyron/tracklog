package tracklog

import (
	"bytes"
	"time"

	"github.com/thcyron/gpx"
)

type User struct {
	ID       int
	Username string
	Password string
}

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
}

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

type Point struct {
	ID        int
	TrackID   int
	Latitude  float64
	Longitude float64
	Time      time.Time
	Elevation float64
	Heartrate uint
}

// Speed returns the speed in meters per second.
func (log *Log) Speed() float64 {
	sec := log.Duration
	if sec == 0 {
		return 0
	}
	return float64(log.Distance) / float64(sec)
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

func NewPoint(point gpx.Point) *Point {
	p := &Point{
		Latitude:  point.Latitude,
		Longitude: point.Longitude,
		Time:      point.Time,
		Elevation: point.Elevation,
	}

	ge, err := gpx.ParseGarminTrackPointExtension(point.Extensions)
	if err == nil {
		p.Heartrate = ge.HeartRate
	}

	return p
}
