package heartrate

import (
	"time"

	"github.com/thcyron/tracklog/pkg/models"
)

type Summary struct {
	Average int
	Rates   []Heartrate
	Zones   map[Zone]int
}

type Heartrate struct {
	BPM  int
	Time time.Time
}

type Zone string

const (
	ZoneRed        Zone = "Red"
	ZoneAnaerobic  Zone = "Anaerobic"
	ZoneAerobic    Zone = "Aerobic"
	ZoneFatBurning Zone = "Fat Burning"
	ZoneEasy       Zone = "Easy"
	ZoneNone       Zone = "None"
)

func (z Zone) String() string {
	return string(z)
}

func (hr Heartrate) Zone() Zone {
	switch {
	case hr.BPM >= 175:
		return ZoneRed
	case hr.BPM >= 164:
		return ZoneAnaerobic
	case hr.BPM >= 153:
		return ZoneAerobic
	case hr.BPM >= 142:
		return ZoneFatBurning
	case hr.BPM >= 131:
		return ZoneEasy
	default:
		return ZoneNone
	}
}

func SummaryForLog(log *models.Log) Summary {
	summary := Summary{
		Zones: make(map[Zone]int),
	}
	sum := uint(0)

	for _, track := range log.Tracks {
		for _, point := range track.Points {
			if point.Heartrate > 0 {
				hr := Heartrate{
					BPM:  int(point.Heartrate),
					Time: point.Time,
				}

				summary.Rates = append(summary.Rates, hr)
				summary.Zones[hr.Zone()]++
				sum += point.Heartrate
			}
		}
	}

	if sum > 0 {
		summary.Average = int(sum / uint(len(summary.Rates)))
	}
	return summary
}
