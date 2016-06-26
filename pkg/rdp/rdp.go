package rdp

import "math"

type Point struct {
	X, Y float64
	Data interface{}
}

// Reduce reduces points using the Ramer–Douglas–Peucker algorithm.
func Reduce(points []Point, epsilon float64) []Point {
	if len(points) <= 2 || epsilon <= 0 {
		return points
	}

	var (
		p, q  = points[0], points[len(points)-1]
		dmax  float64
		index int
	)
	for i := 1; i <= len(points)-2; i++ {
		pp := points[i]
		d := distance(pp, p, q)
		if d > dmax {
			dmax = d
			index = i
		}
	}
	if index > 0 && dmax > epsilon {
		res1 := Reduce(points[0:index], epsilon)
		res2 := Reduce(points[index:], epsilon)
		return append(res1, res2...)
	}

	return []Point{p, q}
}

func distance(p0, p1, p2 Point) float64 {
	return math.Abs((p2.Y-p1.Y)*p0.X-(p2.X-p1.X)*p0.Y+p2.X*p1.Y-p2.Y*p1.X) /
		math.Sqrt(math.Pow(p2.Y-p1.Y, 2)+math.Pow(p2.X-p1.X, 2))
}
