package geometry

import (
	"math"
)

/////////////////////////
// Utility
/////////////////////////
func AdjustEpsilon(e float64, x float64) float64 {
	if -e < x && x < e {
		return 0
	}
	return x
}

/////////////////////////
// Geometry
/////////////////////////
type Shape struct {
	intersects func(*Shape, Ray) float64
	Material   int
	Colour     Vec3
	Emission   Vec3
	Position   Vec3
	normalDir  func(*Shape, Vec3) Vec3
	Size       float64

	normal Vec3
	radius float64
}

func (s *Shape) Intersects(ray Ray) float64 {
	return s.intersects(s, ray)
}

func (s *Shape) NormalDir(point Vec3) Vec3 {
	return s.normalDir(s, point)
}

func Sphere(radius float64, position, emission, colour Vec3, materialType int) *Shape {
	return &Shape{
		intersects: sphereIntersects,
		Material:   materialType,
		Colour:     colour,
		Emission:   emission,
		Position:   position,
		normalDir:  sphereNormal,
		Size:       math.Pi * radius * radius,

		radius: radius,
	}
}

func Plane(position, emission, colour, normal Vec3, materialType int) *Shape {
	return &Shape{
		intersects: planeIntersects,
		Material:   materialType,
		Colour:     colour,
		Emission:   emission,
		Position:   position,
		normalDir:  planeNormal,
		Size:       math.Inf(+1),

		normal: normal,
	}
}

func planeIntersects(s *Shape, r Ray) float64 {
	const epsilon = 1e-12

	// Orthogonal
	dot := r.Direction.Dot(s.normal)
	if -epsilon < dot && dot < epsilon {
		return math.Inf(+1)
	}
	return s.Position.Sub(r.Origin).Dot(s.normal) / dot
}

func sphereIntersects(s *Shape, ray Ray) float64 {
	difference := s.Position.Sub(ray.Origin)
	const epsilon = 1e-5
	dot := difference.Dot(ray.Direction)
	hypotenuse := dot*dot - difference.Dot(difference) + s.radius*s.radius

	if hypotenuse < 0 {
		return math.Inf(+1)
	}

	hypotenuse = math.Sqrt(hypotenuse)
	if diff := dot - hypotenuse; diff > epsilon {
		return diff
	}
	if diff := dot + hypotenuse; diff > epsilon {
		return diff
	}
	return math.Inf(+1)
}

func sphereNormal(s *Shape, point Vec3) Vec3 {
	return point.Sub(s.Position)
}

func planeNormal(s *Shape, point Vec3) Vec3 {
	return s.normal
}

/////////////////////////
// Rays
/////////////////////////
type Ray struct {
	Origin, Direction Vec3
}

/////////////////////////
// CONSTANTS
/////////////////////////
const (
	DIFFUSE    = 1
	SPECULAR   = 2
	REFRACTIVE = 3
)
