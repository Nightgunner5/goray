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
	Material int
	Colour   Vec3
	Emission Vec3
	Position Vec3
	Size     float64

	normal Vec3
	radius float64
	kind   int
}

const (
	kindSphere = iota
	kindPlane
	kindCube
)

func (s *Shape) Intersects(ray *Ray) float64 {
	switch s.kind {
	case kindSphere:
		return sphereIntersects(s, ray)
	case kindPlane:
		return planeIntersects(s, ray)
	case kindCube:
		return cubeIntersects(s, ray)
	}
	panic("unreachable")
}

func (s *Shape) NormalDir(point Vec3) Vec3 {
	switch s.kind {
	case kindSphere:
		return sphereNormal(s, point)
	case kindPlane:
		return planeNormal(s, point)
	case kindCube:
		return cubeNormal(s, point)
	}
	panic("unreachable")
}

func Plane(position, emission, colour, normal Vec3, materialType int) *Shape {
	return &Shape{
		Material: materialType,
		Colour:   colour,
		Emission: emission,
		Position: position,
		Size:     math.Inf(+1),

		normal: normal,
		kind:   kindPlane,
	}
}

func Sphere(radius float64, position, emission, colour Vec3, materialType int) *Shape {
	return &Shape{
		Material: materialType,
		Colour:   colour,
		Emission: emission,
		Position: position,
		Size:     math.Pi * radius * radius,

		radius: radius,
		kind:   kindSphere,
	}
}

func Cube(radius float64, position, emission, colour Vec3, materialType int) *Shape {
	return &Shape{
		Material: materialType,
		Colour:   colour,
		Emission: emission,
		Position: position,
		Size:     radius * radius * radius * 8,

		radius: radius,
		kind:   kindCube,
	}
}

var positiveInfinity = math.Inf(+1)

func intersectPlane(origin, normal Vec3, r *Ray) float64 {
	const epsilon = 1e-12

	// Orthogonal
	dot := r.Direction.Dot(normal)
	if -epsilon < dot && dot < epsilon {
		return positiveInfinity
	}
	return origin.SubDot(r.Origin, normal) / dot
}

func planeIntersects(s *Shape, r *Ray) float64 {
	return intersectPlane(s.Position, s.normal, r)
}

func sphereIntersects(s *Shape, ray *Ray) float64 {
	difference := s.Position.Sub(ray.Origin)
	const epsilon = 1e-5
	dot := difference.Dot(ray.Direction)
	hypotenuse := dot*dot - difference.Dot(difference) + s.radius*s.radius

	if hypotenuse < 0 {
		return positiveInfinity
	}

	hypotenuse = math.Sqrt(hypotenuse)
	if diff := dot - hypotenuse; diff > epsilon {
		return diff
	}
	if diff := dot + hypotenuse; diff > epsilon {
		return diff
	}
	return positiveInfinity
}

func cubeIntersects(s *Shape, r *Ray) float64 {
	// TODO: optimize this heavily
	min := positiveInfinity
	for i := 0; i < 6; i++ {
		var normal Vec3
		switch i {
		case 0:
			normal.X = -s.radius
		case 1:
			normal.X = s.radius
		case 2:
			normal.Y = -s.radius
		case 3:
			normal.Y = s.radius
		case 4:
			normal.Z = -s.radius
		case 5:
			normal.Z = s.radius
		}
		dist := intersectPlane(s.Position.Add(normal), normal, r)
		if dist > 0 && dist < min {
			diff := r.Origin.Add(r.Direction.Mult(dist)).Sub(s.Position)
			if -s.radius <= diff.X && diff.X <= s.radius &&
				-s.radius <= diff.Y && diff.Y <= s.radius &&
				-s.radius <= diff.Z && diff.Z <= s.radius {
				min = dist
			}
		}
	}

	return min
}

func planeNormal(s *Shape, point Vec3) Vec3 {
	return s.normal
}

func sphereNormal(s *Shape, point Vec3) Vec3 {
	return point.Sub(s.Position)
}

func cubeNormal(s *Shape, point Vec3) Vec3 {
	// TODO: optimize this heavily
	max := 0.0
	var bestNormal Vec3
	diff := point.Sub(s.Position)
	for i := 0; i < 6; i++ {
		var normal Vec3
		switch i {
		case 0:
			normal.X = -s.radius
		case 1:
			normal.X = s.radius
		case 2:
			normal.Y = -s.radius
		case 3:
			normal.Y = s.radius
		case 4:
			normal.Z = -s.radius
		case 5:
			normal.Z = s.radius
		}
		dot := normal.Dot(diff)
		if dot > max {
			max = dot
			bestNormal = normal
		}
	}

	return bestNormal
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
