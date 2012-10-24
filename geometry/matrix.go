package geometry

import (
	"math"
)

/////////////////////////
// Matrix
/////////////////////////
type Mat4 struct {
	matrix [16]float64
}

func (m Mat4) Mult(v *Vec3) Vec3 {
	return Vec3{m.matrix[0]*v.X + m.matrix[1]*v.Y + m.matrix[2]*v.Z,
		m.matrix[4]*v.X + m.matrix[5]*v.Y + m.matrix[6]*v.Z,
		m.matrix[8]*v.X + m.matrix[9]*v.Y + m.matrix[10]*v.Z}
}

func RotateVector(a float64, axis *Vec3, vec *Vec3) Vec3 {
	sin, cos := math.Sin(a), math.Cos(a)
	x, y, z := axis.X, axis.Y, axis.Z
	m := [16]float64{cos + x*x*(1-cos), x*y*(1-cos) - x*sin, x*z*(1-cos) + y*sin, 0,
		y*x*(1-cos) + z*sin, cos + y*y*(1-cos), y*z*(1-cos) - x*sin, 0,
		z*x*(1-cos) - y*sin, z*y*(1-cos) + x*sin, cos + z*z*(1-cos), 0,
		0, 0, 0, 1}

	return Mat4{m}.Mult(vec).Truncate()
}
