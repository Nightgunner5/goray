package geometry

import (
	"math"
)

/////////////////////////
// Vectors
/////////////////////////
type Vec3 struct {
	X, Y, Z float64
}

func (v Vec3) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vec3) Position() Vec3 {
	return v
}

func (v Vec3) Truncate() Vec3 {
	epsilon := 1e-4
	return Vec3{AdjustEpsilon(epsilon, v.X), AdjustEpsilon(epsilon, v.Y), AdjustEpsilon(epsilon, v.Z)}
}

func (v Vec3) Normalize() Vec3 {
	factor := 1.0 / v.Abs()
	return Vec3{v.X * factor, v.Y * factor, v.Z * factor}
}

func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.X + other.X,
		v.Y + other.Y,
		v.Z + other.Z}
}

func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{v.X - other.X,
		v.Y - other.Y,
		v.Z - other.Z}
}

func (v Vec3) Mult(lambda float64) Vec3 {
	return Vec3{v.X * lambda,
		v.Y * lambda,
		v.Z * lambda}
}

func (v Vec3) MultVec(o Vec3) Vec3 {
	return Vec3{v.X * o.X,
		v.Y * o.Y,
		v.Z * o.Z}
}

func (v Vec3) Dot(other Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X}
}

func (v Vec3) IsZero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

func (me Vec3) Distance(other Vec3) float64 {
	dx, dy, dz := me.X-other.X, me.Y-other.Y, me.Z-other.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func (me Vec3) Distance2(other Vec3) float64 {
	dx, dy, dz := me.X-other.X, me.Y-other.Y, me.Z-other.Z
	return dx*dx + dy*dy + dz*dz
}

/////////////////////////
// Ugly util functions
/////////////////////////
func (v Vec3) CLAMPF() Vec3 {
	return Vec3{math.Max(math.Min(v.X, 1), 0),
		math.Max(math.Min(v.Y, 1), 0),
		math.Max(math.Min(v.Z, 1), 0)}
}

func (v Vec3) CLAMP() Vec3 {
	return Vec3{math.Max(math.Min(v.X, 255), 0),
		math.Max(math.Min(v.Y, 255), 0),
		math.Max(math.Min(v.Z, 255), 0)}
}

func (v Vec3) PEAKS(a float64) Vec3 {
	return Vec3{math.Max(0, v.X-a),
		math.Max(0, v.Y-a),
		math.Max(0, v.Z-a)}
}
