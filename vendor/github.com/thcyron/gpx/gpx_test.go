package gpx

import (
	"math"
	"os"
	"testing"
	"time"
)

func TestDecoder(t *testing.T) {
	f, err := os.Open("test/test.gpx")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := NewDecoder(f).Decode()
	if err != nil {
		t.Fatal(err)
	}
	if doc.Version != "1.1" {
		t.Errorf("got wrong version %q", doc.Version)
	}
	if dist := doc.Distance(); math.Abs(dist-1362.370020) > 0.0000001 {
		t.Errorf("got %f distance; expected 1362.370020", dist)
	}
	expectedDuration, err := time.ParseDuration("39m19s")
	if err != nil {
		t.Fatal(err)
	}
	if dur := doc.Duration(); dur != expectedDuration {
		t.Errorf("got %s duration; expected %s", dur, expectedDuration)
	}
	expectedStart := time.Date(2015, 12, 13, 18, 35, 18, 0, time.UTC)
	if start := doc.Start(); !start.Equal(expectedStart) {
		t.Errorf("got %v start; expected %v", start, expectedStart)
	}
	expectedEnd := time.Date(2015, 12, 13, 19, 14, 37, 0, time.UTC)
	if end := doc.End(); !end.Equal(expectedEnd) {
		t.Errorf("got %v end; expected %v", end, expectedEnd)
	}

	metadataTime := time.Date(2015, 12, 13, 18, 35, 18, 0, time.UTC)
	if !doc.Metadata.Time.Equal(metadataTime) {
		t.Errorf("got %v, expected %v", doc.Metadata.Time, metadataTime)
	}

	if l := len(doc.Tracks); l != 1 {
		t.Errorf("got %d track(s); expected 1", l)
	}
	track := doc.Tracks[0]

	if track.Name != "Running" {
		t.Errorf("got %q name; expected %q", track.Name, "Running")
	}
	if l := len(track.Segments); l != 1 {
		t.Errorf("got %d segment(s); expected 1", l)
	}
	seg := track.Segments[0]

	if l := len(seg.Points); l != 9 {
		t.Errorf("got %d points(s); expected 9", l)
	}

	pointTestCases := []struct {
		point Point
		lat   float64
		lon   float64
		ele   float64
		t     time.Time
	}{
		{
			point: seg.Points[0],
			lat:   49.3973693847656250,
			lon:   11.1259574890136719,
			ele:   346.874267578125,
			t:     time.Date(2015, 12, 13, 18, 35, 18, 0, time.UTC),
		},
		{
			point: seg.Points[len(seg.Points)-1],
			lat:   49.3978729248046875,
			lon:   11.1260004043579102,
			ele:   346.11541748046875,
			t:     time.Date(2015, 12, 13, 19, 14, 37, 0, time.UTC),
		},
	}

	for i, testCase := range pointTestCases {
		if testCase.point.Latitude != testCase.lat {
			t.Errorf("point test case %d: got %v latitude; expected %v", i, testCase.point.Latitude, testCase.lat)
		}
		if testCase.point.Longitude != testCase.lon {
			t.Errorf("point test case %d: got %v longitude; expected %v", i, testCase.point.Longitude, testCase.lon)
		}
		if testCase.point.Elevation != testCase.ele {
			t.Errorf("point test case %d: got %v elevation; expected %v", i, testCase.point.Elevation, testCase.ele)
		}
		if !testCase.point.Time.Equal(testCase.t) {
			t.Errorf("point test case %d: got %v time; expected %v", i, testCase.point.Time, testCase.t)
		}
	}
}

func TestDecoderNoGPXTag(t *testing.T) {
	f, err := os.Open("test/no_gpx.gpx")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewDecoder(f).Decode()
	if err != ErrBadRootTag {
		t.Fatal("decoding should fail due to bad root tag")
	}
}

func TestDecoderGPX10(t *testing.T) {
	f, err := os.Open("test/gpx10.gpx")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewDecoder(f).Decode()
	if err != ErrGPX11Only {
		t.Fatal("decoding should fail for GPX 1.0 documents")
	}
}
