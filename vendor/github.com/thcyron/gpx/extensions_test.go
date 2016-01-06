package gpx

import (
	"os"
	"reflect"
	"testing"
)

func TestGarminTrackPointExtension(t *testing.T) {
	f, err := os.Open("test/test.gpx")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := NewDecoder(f).Decode()
	if err != nil {
		t.Fatal(err)
	}

	point := doc.Tracks[0].Segments[0].Points[0]
	ext, err := ParseGarminTrackPointExtension(point.Extensions)
	if err != nil {
		t.Fatal(err)
	}

	expectedExt := GarminTrackPointExtension{
		HeartRate: 126,
		Cadence:   81,
		AirTemp:   23,
		WaterTemp: 19,
		Depth:     9,
	}
	if !reflect.DeepEqual(ext, expectedExt) {
		t.Errorf("got %#v extension; expected %#v", ext, expectedExt)
	}
}
