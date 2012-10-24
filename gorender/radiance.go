package gorender

import (
    "../geometry"
    "../kd"
    "container/list"
    "math/rand"
    "fmt"
)

func Radiance(ray geometry.Ray, scene *list.List, photonMap *kd.KDNode, depth int, alpha float64) geometry.Vec3 {
	// Russian roulette to stop tracing
	if depth > Config.MinDepth && rand.Float64() > alpha {
		return geometry.Vec3{0, 0, 0}
	}

	if shape, distance := ClosestIntersection(scene, ray); shape != nil {
		impact := ray.Origin.Add(ray.Direction.Mult(distance))
		normal := shape.NormalDir(impact).Normalize()
		reverse := ray.Direction.Mult(-1)

		contribution := geometry.Vec3{0, 0, 0}
		outgoing := normal
		if normal.Dot(reverse) < 0 {
			outgoing = normal.Mult(-1)
		}

		// Look up information in photon map
		nodes := photonMap.Neighbors(impact, 1)
		directLight := geometry.Vec3{0, 0, 0}
		for e := nodes.Front(); e != nil; e = e.Next() {
			photon := e.Value.(*kd.KDNode).Value.(PhotonHit)
			delta := impact.Distance(photon.Location)
			//light := photon.Photon.Mult(outgoing.Dot(photon.Incomming.Mult(-1)) / (10 * delta))
			light := photon.Photon.Mult(1.0 / (10 * delta))
			directLight = directLight.Add(light)
		}
		directLight = directLight.Mult(1.0 / (1 + float64(nodes.Len())))
		contribution = contribution.Add(directLight.MultVec(shape.Colour()))
		return contribution

		switch shape.Material() {
		case geometry.DIFFUSE:
			contribution = contribution.Add(shape.Colour().MultVec(EmitterSampling(impact, shape, scene)))

			u := normal.Cross(reverse).Normalize()
			v := u.Cross(normal).Normalize()
			//r1, r2 := math.Pi*rand.Float64(), 2*math.Pi*rand.Float64()
			//bounce := u.Mult(math.Cos(r2)).Add(normal).Add(v.Mult(math.Sin(r1)))
			bounce := u.Mult(rand.NormFloat64() * 0.5).Add(outgoing).Add(v.Mult(rand.NormFloat64() * 0.5))
			secondary := Radiance(geometry.Ray{impact, bounce.Normalize()}, scene, photonMap, depth+1, alpha*0.66)
			//return contribution.Add(shape.Colour().MultVec(secondary).Mult(outgoing.Dot(bounce) / alpha))
			return contribution.Add(shape.Colour().MultVec(secondary).Mult(outgoing.Dot(bounce) / alpha))

		case geometry.SPECULAR:
			bounce := ray.Direction.Sub(normal.Mult(2 * outgoing.Dot(ray.Direction)))
			// TODO: Solve problem with lights on specular surface
			secondary := Radiance(geometry.Ray{impact, bounce.Normalize()}, scene, photonMap, depth+1, alpha*0.9)
			//return contribution.Add(secondary.Mult(0.99).Mult(outgoing.Dot(bounce) / alpha))
			return contribution.Add(secondary.Mult(0.99))
		case geometry.REFRACTIVE:
			fmt.Println("HAHA!")
		}
	}
	return geometry.Vec3{0, 0, 0}
}
