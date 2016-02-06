package models

import (
	"time"

	"github.com/thcyron/gpx"
	"github.com/thcyron/tracklog/pkg/utils"
)

type Point struct {
	ID        int
	TrackID   int
	Latitude  float64
	Longitude float64
	Time      time.Time
	Elevation float64
	Heartrate uint
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

func (p *Point) DistanceTo(p2 *Point) float64 {
	return utils.Haversine(p.Latitude, p.Longitude, p2.Latitude, p2.Longitude)
}
