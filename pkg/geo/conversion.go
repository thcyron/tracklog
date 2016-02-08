package geo

import "math"

func ToRad(deg float64) float64 {
	return deg * math.Pi / 180
}
