package geo

import (
	"math"
)

func EquirectangularProjection(λ, φ, φ1 float64) (x, y float64) {
	return λ * math.Cos(φ1), φ
}
