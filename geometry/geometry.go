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
type Shape interface {
	Intersects(ray Ray) float64
	Material() int
	Colour() Vec3
	Emission() Vec3
	Position() Vec3
	NormalDir(point Vec3) Vec3
	Size() float64
}

type Sphere struct {
	radius                     float64
	position, emission, colour Vec3
	materialType               int
}

type Plane struct {
	position, emission, colour, normal Vec3
	materialType                       int
}

type Square struct {
	Plane
	width, height float64
}

func (p *Plane) Intersects(r Ray) float64 {
	epsilon := 1e-12

	// Orthogonal
	dot := r.Direction.Dot(p.normal)
	if -epsilon < dot && dot < epsilon {
		return math.Inf(+1)
	}
	return p.position.Sub(r.Origin).Dot(p.normal) / dot
}

func (s *Square) Intersects(r Ray) float64 {
	return 0.0
}

func (s *Sphere) Intersects(ray Ray) float64 {
	difference := s.position.Sub(ray.Origin)
	epsilon := 1e-5
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

func (s *Sphere) Colour() Vec3 {
	return s.colour
}

func (s *Sphere) Position() Vec3 {
	return s.position
}

func (p *Plane) Position() Vec3 {
	return p.position
}

func (p *Plane) Colour() Vec3 {
	return p.colour
}

func (s *Sphere) Material() int {
	return s.materialType
}

func (p *Plane) Material() int {
	return p.materialType
}

func (s *Sphere) Emission() Vec3 {
	return s.emission
}

func (p *Plane) Emission() Vec3 {
	return p.emission
}

func (s *Sphere) NormalDir(point Vec3) Vec3 {
	return point.Sub(s.position)
}

func (p *Plane) NormalDir(point Vec3) Vec3 {
	return p.normal
}

func (s *Square) Size() float64 {
	return s.width * s.height
}

func (s *Sphere) Size() float64 {
	return math.Pi * s.radius * s.radius
}

func (p *Plane) Size() float64 {
	return math.Inf(+1)
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
