package geometry

import (
	"math"
)

type Scene struct {
	Width, Height float64
	Rows, Cols    int
	Objects       []Shape
	Camera        Ray
	Near          float64
	PixW, PixH    float64
}

func ParseScene(filename string, width, height, fov float64, cols, rows int) Scene {

	var shapes []Shape

	shapes = append(shapes, Sphere{2.5, Vec3{2.5, 0.5, -8}, Vec3{0, 0, 0}, Vec3{0.5, 0.5, 0.5}, DIFFUSE})
	shapes = append(shapes, Sphere{5, Vec3{-5, 3, -15}, Vec3{0, 0, 0}, Vec3{1, 1, 1}, SPECULAR})
	shapes = append(shapes, Plane{Vec3{0, -2, -10}, Vec3{0, 0, 0}, Vec3{0, 0, 0.9}, Vec3{0, 10, 1}.Normalize(), DIFFUSE})
	shapes = append(shapes, Plane{Vec3{-10, 0, 0}, Vec3{0, 0, 0}, Vec3{0, 0.9, 0}, Vec3{1, 0, 0}.Normalize(), DIFFUSE})
	shapes = append(shapes, Plane{Vec3{10, 0, -10}, Vec3{0, 0, 0}, Vec3{0.4, 0, 0.4}, Vec3{-1, 0, 0}.Normalize(), DIFFUSE})
	shapes = append(shapes, Plane{Vec3{0, 0, -16}, Vec3{0, 0, 0}, Vec3{0.4, 0.4, 0.4}, Vec3{0, 0, 1}.Normalize(), DIFFUSE})
	shapes = append(shapes, Plane{Vec3{0, 11, -10}, Vec3{0, 0, 0}, Vec3{0, 0.6, 0.6}, Vec3{0, -1, 0}.Normalize(), DIFFUSE})
	shapes = append(shapes, Sphere{1, Vec3{0, 9, -11}, Vec3{2, 2, 2}, Vec3{1, 1, 1}, DIFFUSE})
	shapes = append(shapes, Sphere{1, Vec3{8.2, 0.6, -15}, Vec3{0, 0, 0}, Vec3{1, 1, 1}, DIFFUSE})

	/* SmallPTs set of spheres. Doesn't work with my camera
	shapes = append(shapes, Sphere{1e5, Vec3{1e5+1-50,40.8-52,81.6-295.6},  Vec3{},Vec3{.75,.25,.25}, DIFFUSE})
	shapes = append(shapes, Sphere{1e5, Vec3{50-50,40.8-52, 1e5-295.6},     Vec3{}, Vec3{.75,.75,.75},DIFFUSE})
	shapes = append(shapes, Sphere{1e5, Vec3{50-50,40.8-52,-1e5+170-295.6},  Vec3{},Vec3{}, DIFFUSE})
	shapes = append(shapes, Sphere{1e5, Vec3{50-50, 1e5-52, 81.6-295.6},    Vec3{},Vec3{.75,.75,.75},DIFFUSE})
	shapes = append(shapes, Sphere{1e5, Vec3{50-50,-1e5+81.6-52,81.6-295.6},Vec3{},Vec3{.75,.75,.75},DIFFUSE})
	shapes = append(shapes, Sphere{16.5,Vec3{27-50,16.5-52,47-295.6},Vec3{},Vec3{1,1,1}.Mult(.999), SPECULAR})
	shapes = append(shapes, Sphere{16.5,Vec3{73-50,16.5-52,78-295.6},Vec3{},Vec3{1,1,1}.Mult(.999), REFRACTIVE})
	shapes = append(shapes, Sphere{600, Vec3{50-50,681.6-.27-52,81.6-295.6},Vec3{12,12,12},  Vec3{}, DIFFUSE})
	*/

	/*
	shapes = append(shapes, Plane{Vec3{0, -2, 0}, Vec3{0, 0, 0}, Vec3{0, 0, 0.9}, Vec3{0, 10, 1}.Normalize(), DIFFUSE})
	shapes = append(shapes, Plane{Vec3{-10, 0, 0}, Vec3{0, 0, 0}, Vec3{0, 0.3, 0}, Vec3{1, 0, 0}.Normalize(), DIFFUSE})
	shapes = append(shapes, Plane{Vec3{10, 0, 0}, Vec3{0, 0, 0}, Vec3{0.4, 0, 0.4}, Vec3{-1, 0, 0}.Normalize(), DIFFUSE})
	shapes = append(shapes, Plane{Vec3{0, 0, -15}, Vec3{0, 0, 0}, Vec3{0.4, 0.4, 0.4}, Vec3{0, 0, 1}.Normalize(), DIFFUSE})
	shapes = append(shapes, Plane{Vec3{0, 10, 0}, Vec3{0, 0, 0}, Vec3{0, 0.6, 0.6}, Vec3{0, -1, 0}.Normalize(), DIFFUSE})
	shapes = append(shapes, Sphere{1, Vec3{0, 5, -10}, Vec3{1, 1, 1}, Vec3{1, 1, 1}, DIFFUSE})
	*/
	near := math.Abs(fov / math.Tan(fov/2.0))

	camera := Vec3{0, 0, near}

	return Scene{width, height, rows, cols, shapes,
		Ray{camera, Vec3{0, 0, -1}}, near,
		2 * height / float64(rows),
		2 * width / float64(cols)}
}
