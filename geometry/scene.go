package geometry

import (
	"math"
)

type Scene struct {
	Width, Height Float
	Rows, Cols    int
	Objects       []*Shape
	Camera        Ray
	Near          Float
	PixW, PixH    Float
}

func ParseScene(filename string, width, height, fov Float, cols, rows int) Scene {
	var shapes []*Shape

	// light source
	shapes = append(shapes, Sphere(
		1,                // radius
		Vec3{-4, 0, -10}, // position
		Vec3{2, 2, 2},    // emission
		Vec3{1, 1, 1},    // colour
		DIFFUSE,          // material
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

	// left metal sphere
	shapes = append(shapes, Sphere(
		2.5,                 // radius
		Vec3{-3.5, 0.5, -6}, // position
		Vec3{0, 0, 0},       // emission
		Vec3{1, 1, 1},       // colour
		SPECULAR,            // material
	))

	// left plastic sphere
	shapes = append(shapes, Sphere(
		0.75,                // radius
		Vec3{-2, -1.25, -2}, // position
		Vec3{0, 0, 0},       // emission
		Vec3{0.8, 0.2, 0.4}, // colour
		DIFFUSE,             // material
	))

	// left glass sphere
	shapes = append(shapes, Sphere(
		0.9,                  // radius
		Vec3{-0.5, -1.1, -1}, // position
		Vec3{0, 0, 0},        // emission
		Vec3{1, 1, 1},        // colour
		REFRACTIVE,           // material
	))

	// right plastic sphere
	shapes = append(shapes, Sphere(
		2.5,                 // radius
		Vec3{4, 0.5, -9.5},  // position
		Vec3{0, 0, 0},       // emission
		Vec3{0.5, 0.5, 0.5}, // colour
		DIFFUSE,             // material
	))

	// right glass sphere
	shapes = append(shapes, Sphere(
		1.5,                 // radius
		Vec3{4.5, -0.5, -7}, // position
		Vec3{0, 0, 0},       // emission
		Vec3{1, 1, 1},       // colour
		REFRACTIVE,          // material
	))

	near := Float(math.Abs(float64(fov) / math.Tan(float64(fov/2.0))))

	camera := Vec3{0, 0, near}

	return Scene{width, height, rows, cols, shapes,
		Ray{camera, Vec3{0, 0, -1}}, near,
		2 * height / Float(rows),
		2 * width / Float(cols)}
}
