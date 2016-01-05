package tracklog

import "time"

type HeartrateSummary struct {
	Average int
	Rates   []Heartrate
	Zones   map[HeartrateZone]int
}

type Heartrate struct {
	BPM  int
	Time time.Time
}

type HeartrateZone string

const (
	HeartrateZoneRed        HeartrateZone = "Red"
	HeartrateZoneAnaerobic  HeartrateZone = "Anaerobic"
	HeartrateZoneAerobic    HeartrateZone = "Aerobic"
	HeartrateZoneFatBurning HeartrateZone = "Fat Burning"
	HeartrateZoneEasy       HeartrateZone = "Easy"
	HeartrateZoneNone       HeartrateZone = "None"
)

func (z HeartrateZone) String() string {
	return string(z)
}

func (hr Heartrate) Zone() HeartrateZone {
	switch {
	case hr.BPM >= 175:
		return HeartrateZoneRed
	case hr.BPM >= 164:
		return HeartrateZoneAnaerobic
	case hr.BPM >= 153:
		return HeartrateZoneAerobic
	case hr.BPM >= 142:
		return HeartrateZoneFatBurning
	case hr.BPM >= 131:
		return HeartrateZoneEasy
	default:
		return HeartrateZoneNone
	}
}

func HeartrateSummaryForLog(log *Log) HeartrateSummary {
	summary := HeartrateSummary{
		Zones: make(map[HeartrateZone]int),
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
