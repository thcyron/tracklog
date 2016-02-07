package rdp

import (
	"reflect"
	"testing"
)

func TestReduce(t *testing.T) {
	testCases := []struct {
		points  []Point
		reduced []Point
		epsilon float64
	}{
		{
			points: []Point{
				Point{X: 0, Y: 0},
				Point{X: 1, Y: 1},
			},
			reduced: []Point{
				Point{X: 0, Y: 0},
				Point{X: 1, Y: 1},
			},
			epsilon: 1,
		},
		{
			points: []Point{
				Point{X: 0, Y: 0},
				Point{X: 1, Y: 1},
				Point{X: 2, Y: 2},
			},
			reduced: []Point{
				Point{X: 0, Y: 0},
				Point{X: 2, Y: 2},
			},
			epsilon: 1,
		},
		{
			points: []Point{
				Point{X: 0, Y: 0},
				Point{X: 1, Y: 1},
				Point{X: 2, Y: 2},
			},
			reduced: []Point{
				Point{X: 0, Y: 0},
				Point{X: 1, Y: 1},
				Point{X: 2, Y: 2},
			},
			epsilon: 0,
		},
		{
			points: []Point{
				Point{X: 0, Y: 0},
				Point{X: 1, Y: 3},
				Point{X: 2, Y: 2},
			},
			reduced: []Point{
				Point{X: 0, Y: 0},
				Point{X: 1, Y: 3},
				Point{X: 2, Y: 2},
			},
			epsilon: 1,
		},
	}

	for i, testCase := range testCases {
		r := Reduce(testCase.points, testCase.epsilon)
		if !reflect.DeepEqual(r, testCase.reduced) {
			t.Errorf("test case %d: reduced points do not match: %v <> %v", i, r, testCase.reduced)
		}
	}
}
