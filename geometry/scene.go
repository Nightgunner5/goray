package geometry

import (
	"math"
)

type Scene struct {
	Width, Height float64
	Rows, Cols    int
	Objects       []*Shape
	Camera        Ray
	Near          float64
	PixW, PixH    float64
}

func ParseScene(filename string, width, height, fov float64, cols, rows int) Scene {
	var shapes []*Shape

	// high centered light source
	shapes = append(shapes, Sphere(
		1,               // radius
		Vec3{0, 4, -10}, // position
		Vec3{2, 2, 2},   // emission
		Vec3{1, 1, 1},   // colour
		DIFFUSE,         // material
	))

	// rear wall
	shapes = append(shapes, Plane(
		Vec3{0, 0, -12},     // position
		Vec3{0, 0, 0},       // emission
		Vec3{0.6, 0.6, 0.6}, // colour
		Vec3{0, 0, 1},       // normal
		DIFFUSE,             // material
	))

	// floor
	shapes = append(shapes, Plane(
		Vec3{0, -2, 0},    // position
		Vec3{0, 0, 0},     // emission
		Vec3{0, 0.2, 0.4}, // colour
		Vec3{0, 1, 0},     // normal
		DIFFUSE,           // material
	))

	// ceiling
	shapes = append(shapes, Plane(
		Vec3{0, 6, 0},       // position
		Vec3{0, 0, 0},       // emission
		Vec3{0.6, 0.4, 0.2}, // colour
		Vec3{0, -1, 0},      // normal
		DIFFUSE,             // material
	))

	// left wall
	shapes = append(shapes, Plane(
		Vec3{-6, 0, 0},      // position
		Vec3{0, 0, 0},       // emission
		Vec3{0.2, 0.6, 0.2}, // colour
		Vec3{1, 0, 0},       // normal
		DIFFUSE,             // material
	))

	// right wall
	shapes = append(shapes, Plane(
		Vec3{6, 0, 0},       // position
		Vec3{0, 0, 0},       // emission
		Vec3{0.6, 0.2, 0.4}, // colour
		Vec3{-1, 0, 0},      // normal
		DIFFUSE,             // material
	))

	// left metalic sphere
	shapes = append(shapes, Sphere(
		2.5,                 // radius
		Vec3{-3.5, 0.5, -6}, // position
		Vec3{0, 0, 0},       // emission
		Vec3{1, 1, 1},       // colour
		SPECULAR,            // material
	))

	/*// central glass cube
	shapes = append(shapes, Cube(
		0.9,                  // radius
		Vec3{-0.5, -1.1, -2}, // position
		Vec3{0, 0, 0},        // emission
		Vec3{1, 1, 1},        // colour
		REFRACTIVE,           // material
	))*/

	// right rear plastic sphere
	shapes = append(shapes, Sphere(
		2,                   // radius
		Vec3{4, 0, -11},     // position
		Vec3{0, 0, 0},       // emission
		Vec3{0.5, 0.5, 0.5}, // colour
		DIFFUSE,             // material
	))

	near := math.Abs(fov / math.Tan(fov/2.0))

	camera := Vec3{0, 0, near}

	return Scene{width, height, rows, cols, shapes,
		Ray{camera, Vec3{0, 0, -1}}, near,
		2 * height / float64(rows),
		2 * width / float64(cols)}
}
