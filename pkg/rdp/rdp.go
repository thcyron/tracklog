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

	p, q := points[0], points[len(points)-1]
	m := (q.Y - p.Y) / (q.X - p.X)
	t := p.Y - m*p.X

	var (
		dmax  float64
		index int
	)
	for i := 1; i <= len(points)-2; i++ {
		pp := points[i]
		d := distance(pp.X, pp.Y, m, t)
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

func distance(x, y, m, t float64) float64 {
	dt := math.Abs(t - y + m*x)
	phi := 0.5*math.Pi - math.Atan(m)
	return math.Sin(phi) * dt
}
