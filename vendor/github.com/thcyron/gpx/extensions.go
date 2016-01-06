package gpx

import (
	"encoding/xml"
	"errors"
)

var (
	ErrNoSuchExtension = errors.New("gpx: no such extension")
)

// GarminTrackPointExtension is Garmin’s TrackPoint extension defined by
// https://www8.garmin.com/xmlschemas/TrackPointExtensionv1.xsd
type GarminTrackPointExtension struct {
	AirTemp   float64 // Air temperature (Celsius)
	WaterTemp float64 // Water temperature (Celsius)
	Depth     float64 // Diving depth (meters)
	HeartRate uint    // Heart rate (beats per minute)
	Cadence   uint    // Cadence (revs per minute)
}

const GarminTrackPointExtensionNS = "http://www.garmin.com/xmlschemas/TrackPointExtension/v1"

// ParseGarminTrackPointExtension tries to parse Garmin’s TrackPoint extension
// from a point’s extensions tokens.
func ParseGarminTrackPointExtension(tokens []xml.Token) (e GarminTrackPointExtension, err error) {
	ts := tokenStream{&sliceTokener{tokens: tokens}}

	if !findExtension(ts, GarminTrackPointExtensionNS, "TrackPointExtension") {
		return e, ErrNoSuchExtension
	}

	for {
		tok, err := ts.Token()
		if err != nil {
			return e, err
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			if se.Name.Space != GarminTrackPointExtensionNS {
				ts.skipTag()
				continue
			}
			switch se.Name.Local {
			case "hr":
				hr, err := ts.consumeInt()
				if err != nil {
					return e, err
				}
				e.HeartRate = uint(hr)
			case "cad":
				cad, err := ts.consumeInt()
				if err != nil {
					return e, err
				}
				e.Cadence = uint(cad)
			case "atemp":
				atemp, err := ts.consumeFloat()
				if err != nil {
					return e, err
				}
				e.AirTemp = atemp
			case "wtemp":
				wtemp, err := ts.consumeFloat()
				if err != nil {
					return e, err
				}
				e.WaterTemp = wtemp
			case "depth":
				depth, err := ts.consumeFloat()
				if err != nil {
					return e, err
				}
				e.Depth = depth
			default:
				ts.skipTag()
			}
		case xml.EndElement:
			return e, nil
		}
	}
}

func findExtension(ts tokenStream, space, local string) bool {
	for {
		tok, err := ts.Token()
		if err != nil {
			return false
		}
		switch tok.(type) {
		case xml.StartElement:
			se := tok.(xml.StartElement)
			if se.Name.Space == space && se.Name.Local == local {
				return true
			}
			ts.skipTag()
		case xml.EndElement:
			return false
		}
	}
}
