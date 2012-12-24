package gorender

import (
	"../geometry"
	"../kd"
	"container/list"
	"fmt"
	"math"
	"math/rand"
)

func EmitterSampling(point, normal geometry.Vec3, shapes *list.List) geometry.Vec3 {
    incommingLight := geometry.Vec3{0, 0, 0}
    
    for e := shapes.Front(); e != nil; e = e.Next() {
		shape := e.Value.(geometry.Shape)
		
		if shape.Emission().Abs() > 0 {
		    // It's a light source
		    direction := shape.NormalDir(point).Mult(-1)
			u := direction.Cross(normal).Normalize()
			v := direction.Cross(u).Normalize()
			
		    direction = direction.Add(u.Mult(rand.NormFloat64() * 0.3)).Add(v.Mult(rand.NormFloat64() * 0.3));
		    ray := geometry.Ray{point, direction.Normalize()};
		    
		    if object, distance := ClosestIntersection(shapes, ray); object == shape {
		        incommingLight = incommingLight.Add(object.Emission().Mult(direction.Dot(normal) / (1 + distance)))
		    }
		}
	}
	return incommingLight
}


func Radiance(ray geometry.Ray, scene *geometry.Scene, diffuseMap, causticsMap *kd.KDNode, depth int, alpha float64) geometry.Vec3 {

	if depth > Config.MinDepth && rand.Float64() > alpha {
		return geometry.Vec3{0, 0, 0}
	}
	
	if shape, distance := ClosestIntersection(scene.Objects, ray); shape != nil {
		impact := ray.Origin.Add(ray.Direction.Mult(distance))
		normal := shape.NormalDir(impact).Normalize()
		reverse := ray.Direction.Mult(-1)

        contribution := shape.Emission();
		outgoing := normal
		if normal.Dot(reverse) < 0 {
			outgoing = normal.Mult(-1)
		}
		
		if shape.Material() == geometry.DIFFUSE {
		    causticLight := geometry.Vec3{0, 0, 0}
		    directLight := geometry.Vec3{0, 0, 0}
		    
		    nodes := causticsMap.Neighbors(impact, 0.1)
            for e := nodes.Front(); e != nil; e = e.Next() {
                photon := e.Value.(*kd.KDNode).Value.(PhotonHit)
                dist := photon.Location.Distance(impact)
                light := photon.Photon.Mult(outgoing.Dot(photon.Incomming.Mult(-1 / math.Pi*(1 + dist))))
                causticLight = causticLight.Add(light)
            }
            if nodes.Len() > 0 {
                causticLight = causticLight.Mult(1.0 / float64(nodes.Len()))
            }
            
            
            directLight = EmitterSampling(impact, normal, scene.Objects)
		    
			u := normal.Cross(reverse).Normalize()
			v := u.Cross(normal).Normalize()
			
			bounceDirection := u.Mult(rand.NormFloat64() * 0.5).Add(outgoing).Add(v.Mult(rand.NormFloat64() * 0.5))
			bounceRay := geometry.Ray{impact, bounceDirection.Normalize()}
            indirectLight := Radiance(bounceRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.9)
            diffuseLight := shape.Colour().MultVec(directLight.Add(indirectLight)).Add(causticLight).Mult(outgoing.Dot(reverse))
            
			return contribution.Add(diffuseLight)

		}
		if shape.Material() == geometry.SPECULAR {
			reflectionDirection := ray.Direction.Sub(normal.Mult(2 * outgoing.Dot(ray.Direction)))
			reflectedRay := geometry.Ray{impact, reflectionDirection.Normalize()}
            incommingLight := Radiance(reflectedRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.99)
			return incommingLight.Mult(outgoing.Dot(reverse))
		}

		if shape.Material() == geometry.REFRACTIVE {
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
			sqrt := math.Sqrt(math.Max(1.0 - math.Pow(factor*sinTi, 2), 0))
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
            R = math.Pow((n1 - n2) / (n1 + n2), 2);
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
                return Radiance(reflectedRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.9)
            } else {
    			reflectionDirection := ray.Direction.Sub(outgoing.Mult(2 * outgoing.Dot(ray.Direction)))
		        reflectedRay := geometry.Ray{impact, reflectionDirection.Normalize()}
                reflectedLight := Radiance(reflectedRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.9).Mult(R)
                
				nDotI := normal.Dot(ray.Direction)
				trasmittedDirection := ray.Direction.Mult(factor)
				term2 := factor * nDotI
				term3 := math.Sqrt(1 - factor*factor*(1-nDotI*nDotI))

				trasmittedDirection = trasmittedDirection.Add(normal.Mult(term2 - term3))
				transmittedRay := geometry.Ray{impact, trasmittedDirection.Normalize()}
				transmittedLight := Radiance(transmittedRay, scene, diffuseMap, causticsMap, depth+1, alpha*0.9).Mult(T)
				return reflectedLight.Add(transmittedLight).Mult(outgoing.Dot(reverse))
			}
		}
		panic("Material without property encountered!")
	}

	return geometry.Vec3{0, 0, 0}
}

/*
func RadianceOLD(ray geometry.Ray, scene *list.List, photonMap *kd.KDNode, depth int, alpha float64) geometry.Vec3 {
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
*/

