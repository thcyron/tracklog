package utils

import "testing"

func TestDuration(t *testing.T) {
	testCases := []struct {
		seconds  uint
		duration string
	}{
		{
			seconds:  1,
			duration: "0:00:01",
		},
		{
			seconds:  60,
			duration: "0:01:00",
		},
		{
			seconds:  61,
			duration: "0:01:01",
		},
		{
			seconds:  3600,
			duration: "1:00:00",
		},
		{
			seconds:  36001,
			duration: "10:00:01",
		},
		{
			seconds:  24 * 3600,
			duration: "1d 0:00:00",
		},
		{
			seconds:  24*3600 + 61,
			duration: "1d 0:01:01",
		},
		{
			seconds:  25*3600 + 61,
			duration: "1d 1:01:01",
		},
	}

	for i, testCase := range testCases {
		s := Duration(testCase.seconds).String()
		if s != testCase.duration {
			t.Errorf("test case %d: expected %q; got %q", i, testCase.duration, s)
		}
	}
}

func TestDistance(t *testing.T) {
	testCases := []struct {
		meters   uint
		distance string
	}{
		{
			meters:   1,
			distance: "0.00 km",
		},
		{
			meters:   10,
			distance: "0.01 km",
		},
		{
			meters:   100,
			distance: "0.10 km",
		},
		{
			meters:   1000,
			distance: "1.00 km",
		},
		{
			meters:   1011,
			distance: "1.01 km",
		},
	}

	for i, testCase := range testCases {
		s := Distance(testCase.meters).String()
		if s != testCase.distance {
			t.Errorf("test case %d: expected %q; got %q", i, testCase.distance, s)
		}
	}
}

func TestSpeed(t *testing.T) {
	testCases := []struct {
		mps   float64 // meters per second
		speed string
	}{
		{
			mps:   0.2777,
			speed: "1.00 km/h",
		},
		{
			mps:   2.7777,
			speed: "10.00 km/h",
		},
		{
			mps:   18.46145,
			speed: "66.46 km/h",
		},
	}

	for i, testCase := range testCases {
		s := Speed(testCase.mps).String()
		if s != testCase.speed {
			t.Errorf("test case %d: expected %q; got %q", i, testCase.speed, s)
		}
	}
}

func TestPace(t *testing.T) {
	testCases := []struct {
		spm  float64 // seconds per meter
		pace string
	}{
		{
			spm:  0.060,
			pace: "1:00 min/km",
		},
		{
			spm:  0.601,
			pace: "10:01 min/km",
		},
	}

	for i, testCase := range testCases {
		s := Pace(testCase.spm).String()
		if s != testCase.pace {
			t.Errorf("test case %d: expected %q; got %q", i, testCase.pace, s)
		}
	}
}
