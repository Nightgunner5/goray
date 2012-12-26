package gorender

import (
	"fmt"
	"github.com/Nightgunner5/goray/geometry"
	"github.com/Nightgunner5/goray/kd"
	"math"
	"math/rand"
)

func EmitterSampling(point, normal geometry.Vec3, shapes []*geometry.Shape, rand *rand.Rand) geometry.Vec3 {
	incommingLight := geometry.Vec3{0, 0, 0}

	for _, shape := range shapes {
		if !shape.Emission.IsZero() {
			// It's a light source
			direction := shape.NormalDir(point).Mult(-1)
			u := direction.Cross(normal).Normalize()
			v := direction.Cross(u).Normalize()

			direction = direction.Add(u.Mult(rand.NormFloat64() * 0.3)).Add(v.Mult(rand.NormFloat64() * 0.3))
			ray := geometry.Ray{point, direction.Normalize()}

			if object, distance := ClosestIntersection(shapes, ray); object == shape {
				incommingLight.AddInPlace(object.Emission.Mult(direction.Dot(normal) / (1 + distance)))
			}
		}
	}
	return incommingLight
}

func Radiance(ray geometry.Ray, scene *geometry.Scene, diffuseMap, causticsMap *kd.KDNode, depth int, alpha float64, rand *rand.Rand) geometry.Vec3 {

	if depth > Config.MinDepth && rand.Float64() > alpha {
		return geometry.Vec3{0, 0, 0}
	}

	if shape, distance := ClosestIntersection(scene.Objects, ray); shape != nil {
		impact := ray.Origin.Add(ray.Direction.Mult(distance))
		normal := shape.NormalDir(impact).Normalize()
		reverse := ray.Direction.Mult(-1)

		contribution := shape.Emission
		outgoing := normal
		if normal.Dot(reverse) < 0 {
			outgoing = normal.Mult(-1)
		}

		if shape.Material == geometry.DIFFUSE {
			var causticLight, directLight geometry.Vec3

			nodes := causticsMap.Neighbors(impact, 0.1)
			for _, e := range nodes {
				photon := causticPhotons[e.Position]
				dist := photon.Location.Distance(impact)
				light := photon.Photon.Mult(outgoing.Dot(photon.Incomming.Mult(-1 / math.Pi * (1 + dist))))
				causticLight.AddInPlace(light)
			}
			if len(nodes) > 0 {
				causticLight = causticLight.Mult(1.0 / float64(len(nodes)))
			}

			directLight = EmitterSampling(impact, normal, scene.Objects, rand)

			u := normal.Cross(reverse).Normalize()
			v := u.Cross(normal).Normalize()

			bounceDirection := u.Mult(rand.NormFloat64() * 0.5).Add(outgoing).Add(v.Mult(rand.NormFloat64() * 0.5))
			bounceRay := geometry.Ray{impact, bounceDirection.Normalize()}
			indirectLight := Radiance(bounceRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.9, rand)
			dot := outgoing.Dot(reverse)
			diffuseLight := geometry.Vec3{
				(shape.Colour.X*(directLight.X+indirectLight.X) + causticLight.X) * dot,
				(shape.Colour.Y*(directLight.Y+indirectLight.Y) + causticLight.Y) * dot,
				(shape.Colour.Z*(directLight.Z+indirectLight.Z) + causticLight.Z) * dot,
			}

			return contribution.Add(diffuseLight)

		}
		if shape.Material == geometry.SPECULAR {
			reflectionDirection := ray.Direction.Sub(normal.Mult(2 * outgoing.Dot(ray.Direction)))
			reflectedRay := geometry.Ray{impact, reflectionDirection.Normalize()}
			incommingLight := Radiance(reflectedRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.99, rand)
			return incommingLight.Mult(outgoing.Dot(reverse))
		}

		if shape.Material == geometry.REFRACTIVE {
			var n1, n2 float64
			if normal.Dot(outgoing) < 0 {
				// Leave the glass
				n1, n2 = GLASS, AIR
			} else {
				n1, n2 = AIR, GLASS
			}

			factor := n1 / n2
			cosTi := normal.Dot(reverse)
			sinTi := math.Sqrt(1 - cosTi*cosTi) // sin² + cos² = 1
			sqrt := math.Sqrt(math.Max(1.0-math.Pow(factor*sinTi, 2), 0))
			// Rs
			top := n1*cosTi - n2*sqrt
			bottom := n1*cosTi + n2*sqrt
			Rs := math.Pow(top/bottom, 2)
			// Rp
			top = n1*sqrt - n2*cosTi
			bottom = n1*sqrt + n2*cosTi
			Rp := math.Pow(top/bottom, 2)

			R := (Rs*Rs + Rp*Rp) / 2.0

			// Approximate:
			R = math.Pow((n1-n2)/(n1+n2), 2)
			// SmallPT formula
			//R = R + (1 - R) * math.Pow(1 - cosTi, 5)
			T := 1.0 - R

			if math.IsNaN(R) {
				fmt.Printf("into: %v, sqrt: %v\n", n2 > n1, sqrt)
				fmt.Printf("cos: %v, sin: %v\n", cosTi, sinTi)
				fmt.Printf("n1: %v, n2: %v\n", n1, n2)
				fmt.Printf("Top: %v, Bottom: %v\n", top, bottom)
				fmt.Printf("Rs: %v, Rp: %v\n", Rs, Rp)
				fmt.Printf("R: %v, T: %v\n", R, T)
				panic("NAN!")
			}

			totalReflection := false
			if n1 > n2 {
				maxAngle := math.Asin(n2 / n1)
				actualAngle := math.Asin(sinTi)

				if actualAngle > maxAngle {
					totalReflection = true
				}
				totalReflection = totalReflection
			}

			if totalReflection {
				reflectionDirection := ray.Direction.Sub(outgoing.Mult(2 * outgoing.Dot(ray.Direction)))
				reflectedRay := geometry.Ray{impact, reflectionDirection.Normalize()}
				return Radiance(reflectedRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.9, rand)
			} else {
				reflectionDirection := ray.Direction.Sub(outgoing.Mult(2 * outgoing.Dot(ray.Direction)))
				reflectedRay := geometry.Ray{impact, reflectionDirection.Normalize()}
				reflectedLight := Radiance(reflectedRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.9, rand).Mult(R)

				nDotI := normal.Dot(ray.Direction)
				trasmittedDirection := ray.Direction.Mult(factor)
				term2 := factor * nDotI
				term3 := math.Sqrt(1 - factor*factor*(1-nDotI*nDotI))

				trasmittedDirection = trasmittedDirection.Add(normal.Mult(term2 - term3))
				transmittedRay := geometry.Ray{impact, trasmittedDirection.Normalize()}
				transmittedLight := Radiance(transmittedRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.9, rand).Mult(T)
				return reflectedLight.Add(transmittedLight).Mult(outgoing.Dot(reverse))
			}
		}
		panic("Material without property encountered!")
	}

	return geometry.Vec3{0, 0, 0}
}
