package geometry

import (
	"container/list"
	"math"
)

type Scene struct {
	Width, Height float64
	Rows, Cols    int
	Objects       *list.List
	Camera        Ray
	Near          float64
	PixW, PixH    float64
}

func ParseScene(filename string, width, height, fov float64, cols, rows int) Scene {

	shapes := list.New()
    
	shapes.PushBack(Sphere{3, Vec3{5, 1, -7}, Vec3{0, 0, 0}, Vec3{0, 0, 0}, REFRACTIVE})
	shapes.PushBack(Sphere{5, Vec3{-5, 3, -15}, Vec3{0, 0, 0}, Vec3{1, 1, 1}, SPECULAR})
	shapes.PushBack(Plane{Vec3{0, -2, -10}, Vec3{0, 0, 0}, Vec3{0, 0, 0.9}, Vec3{0, 10, 1}.Normalize(), DIFFUSE})
	shapes.PushBack(Plane{Vec3{-10, 0, 0}, Vec3{0, 0, 0}, Vec3{0, 0.9, 0}, Vec3{1, 0, 0}.Normalize(), DIFFUSE})
	shapes.PushBack(Plane{Vec3{10, 0, -10}, Vec3{0, 0, 0}, Vec3{0.4, 0, 0.4}, Vec3{-1, 0, 0}.Normalize(), DIFFUSE})
	shapes.PushBack(Plane{Vec3{0, 0, -16}, Vec3{0, 0, 0}, Vec3{0.4, 0.4, 0.4}, Vec3{0, 0, 1}.Normalize(), DIFFUSE})
	shapes.PushBack(Plane{Vec3{0, 11, -10}, Vec3{0, 0, 0}, Vec3{0, 0.6, 0.6}, Vec3{0, -1, 0}.Normalize(), DIFFUSE})
	shapes.PushBack(Sphere{1, Vec3{0, 9, -11}, Vec3{1, 1, 1}, Vec3{1, 1, 1}, DIFFUSE})
	shapes.PushBack(Sphere{1, Vec3{8.2, 0.6, -15}, Vec3{0, 0, 0}, Vec3{1, 1, 1}, DIFFUSE})
    
    /*
	shapes.PushBack(Plane{Vec3{0, -2, 0}, Vec3{0, 0, 0}, Vec3{0, 0, 0.9}, Vec3{0, 10, 1}.Normalize(), DIFFUSE})
	shapes.PushBack(Plane{Vec3{-10, 0, 0}, Vec3{0, 0, 0}, Vec3{0, 0.3, 0}, Vec3{1, 0, 0}.Normalize(), DIFFUSE})
	shapes.PushBack(Plane{Vec3{10, 0, 0}, Vec3{0, 0, 0}, Vec3{0.4, 0, 0.4}, Vec3{-1, 0, 0}.Normalize(), DIFFUSE})
	shapes.PushBack(Plane{Vec3{0, 0, -15}, Vec3{0, 0, 0}, Vec3{0.4, 0.4, 0.4}, Vec3{0, 0, 1}.Normalize(), DIFFUSE})
	shapes.PushBack(Plane{Vec3{0, 10, 0}, Vec3{0, 0, 0}, Vec3{0, 0.6, 0.6}, Vec3{0, -1, 0}.Normalize(), DIFFUSE})
	shapes.PushBack(Sphere{1, Vec3{0, 5, -10}, Vec3{1, 1, 1}, Vec3{1, 1, 1}, DIFFUSE})
    */
	near := math.Abs(fov / math.Tan(fov/2.0))

	return Scene{width, height, rows, cols, shapes,
		Ray{Vec3{0, 0, near}, Vec3{0, 0, -1}}, near,
		2 * height / float64(rows),
		2 * width / float64(cols)}
}
