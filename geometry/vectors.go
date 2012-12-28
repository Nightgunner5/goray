package geometry

import (
	"math"
)

type Float float32

/////////////////////////
// Vectors
/////////////////////////
type Vec3 struct {
	X, Y, Z Float
}

func (v Vec3) Abs() Float {
	return Float(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

func (v Vec3) Truncate() Vec3 {
	const epsilon = 1e-4
	v.X = AdjustEpsilon(epsilon, v.X)
	v.Y = AdjustEpsilon(epsilon, v.Y)
	v.Z = AdjustEpsilon(epsilon, v.Z)
	return v
}

func (v Vec3) Normalize() Vec3 {
	m := v.Abs()
	v.X /= m
	v.Y /= m
	v.Z /= m
	return v
}

func (v Vec3) Add(other Vec3) Vec3 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

func (v *Vec3) AddInPlace(other Vec3) {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
}

func (v Vec3) Sub(other Vec3) Vec3 {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

func (v Vec3) Mult(lambda Float) Vec3 {
	v.X *= lambda
	v.Y *= lambda
	v.Z *= lambda
	return v
}

func (v Vec3) MultVec(other Vec3) Vec3 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

func (v Vec3) Dot(other Vec3) Float {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vec3) SubDot(o, n Vec3) Float {
	return (v.X-o.X)*n.X + (v.Y-o.Y)*n.Y + (v.Z-o.Z)*n.Z
}

func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		v.Y*other.Z - v.Z*other.Y,
		v.Z*other.X - v.X*other.Z,
		v.X*other.Y - v.Y*other.X,
	}
}

func (v Vec3) IsZero() bool {
	return v.X == 0 && v.Y == 0 && v.Z == 0
}

func (v Vec3) Distance(other Vec3) Float {
	return Float(math.Sqrt(float64(v.Distance2(other))))
}

func (v Vec3) Distance2(other Vec3) Float {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}

/////////////////////////
// Ugly util functions
/////////////////////////
func clamp(x, min, max Float) Float {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func (v Vec3) CLAMPF() Vec3 {
	v.X = clamp(v.X, 0, 1)
	v.Y = clamp(v.Y, 0, 1)
	v.Z = clamp(v.Z, 0, 1)
	return v
}

func (v Vec3) CLAMP() Vec3 {
	v.X = clamp(v.X, 0, 255)
	v.Y = clamp(v.Y, 0, 255)
	v.Z = clamp(v.Z, 0, 255)
	return v
}

func (v Vec3) PEAKS(a Float) Vec3 {
	v.X = Float(math.Max(0, float64(v.X-a)))
	v.Y = Float(math.Max(0, float64(v.Y-a)))
	v.Z = Float(math.Max(0, float64(v.Z-a)))
	return v
}
