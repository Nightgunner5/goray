package gorender

import (
	"../geometry"
	"../kd"
	"container/list"
	"fmt"
	"math"
	"math/rand"
)

////////////////////
// Photon Mapping
////////////////////
type PhotonHit struct {
	Location, Photon, Incomming geometry.Vec3
    Depth uint8
}

func (p PhotonHit) Position() geometry.Vec3 {
	return p.Location
}

type RayFunc func(*list.List, geometry.Shape, geometry.Ray, geometry.Vec3, chan<- PhotonHit, float64, int)

func CausticPhoton(scene *list.List, emitter geometry.Shape, ray geometry.Ray, colour geometry.Vec3, result chan<- PhotonHit, alpha float64, depth int) {
	if rand.Float64() > alpha {
		return
	}
	if shape, distance := ClosestIntersection(scene, ray); shape != nil {
		impact := ray.Origin.Add(ray.Direction.Mult(distance))
		if  emitter == shape {
			// Leave the emitter first
			nextRay := geometry.Ray{impact, ray.Direction}
			CausticPhoton(scene, emitter, nextRay, colour, result, alpha, depth)
		} else {
			normal := shape.NormalDir(impact).Normalize()
			reverse := ray.Direction.Mult(-1)
			outgoing := normal
			if normal.Dot(reverse) < 0 {
				outgoing = normal.Mult(-1)
			}
			outgoing = outgoing
			//fmt.Println("Hit something else!")
			if depth > 0 {
				strength := colour.Mult(1.0 / (alpha + distance))
				result <- PhotonHit{impact, strength, ray.Direction, uint8(depth)}
			}

			// Specular objects makes reflections
            /*
			if shape.Material() == geometry.SPECULAR {
				reflection := ray.Direction.Sub(normal.Mult(2 * outgoing.Dot(ray.Direction)))
				reflectedRay := geometry.Ray{impact, reflection.Normalize()}
				CausticPhoton(scene, shape, reflectedRay, colour, result, alpha*0.9, depth+1)
			}
            */
			// Refracting objects makes refractions
            if shape.Material() == geometry.REFRACTIVE {
                
                var n1, n2 float64
                var into bool
                if normal.Dot(outgoing) < 0 {
                    // Leave the glass
                    n1, n2 = GLASS, AIR
                    into = false
                } else {
                    n1, n2 = AIR, GLASS
                    into = true
                }
                
                factor := n1 / n2
                cosTi  := normal.Dot(reverse)
                sinTi  := math.Sqrt(1 - cosTi*cosTi) // sin² + cos² = 1
                sqrt   := math.Sqrt(math.Max(1.0 - math.Pow(factor*sinTi, 2), 0))
                // Rs
                top    := n1*cosTi - n2*sqrt
                bottom := n1*cosTi + n2*sqrt
                Rs     := math.Pow(top/bottom, 2)
                // Rp
                top     = n1*sqrt - n2*cosTi
                bottom  = n1*sqrt + n2*cosTi
                Rp     := math.Pow(top/bottom, 2)

                R := (Rs*Rs + Rp*Rp) / 2.0
                T := 1.0 - R
                
                if math.IsNaN(R) {
                    fmt.Printf("into: %v, sqrt: %v\n", into, sqrt)
                    fmt.Printf("cos: %v, sin: %v\n", cosTi, sinTi)
                    fmt.Printf("n1: %v, n2: %v\n", n1, n2)
                    fmt.Printf("Top: %v, Bottom: %v\n", top, bottom)
                    fmt.Printf("Rs: %v, Rp: %v\n", Rs, Rp)
                    fmt.Printf("R: %v, T: %v\n", R, T)
                    return
                }

                totalReflection := false
                if n1 > n2 {
                    maxAngle    := math.Asin(n2 / n1)
                    actualAngle := math.Asin(sinTi)

                    if actualAngle > maxAngle {
                        totalReflection = true
                    }
                    totalReflection = totalReflection
                }
                
                reflectionDirection := ray.Direction.Sub(normal.Mult(2 * normal.Dot(ray.Direction)))
                reflectedRay := geometry.Ray{impact, reflectionDirection.Normalize()}
                reflectedRay = reflectedRay
                //CausticPhoton(scene, emitter, reflectedRay, colour, result, alpha*0.9, depth+1)
                    
                if true {
                    nDotI := normal.Dot(ray.Direction);
                    trasmittedDirection := ray.Direction.Mult(factor)
                    term2 := factor * nDotI
                    term3 := math.Sqrt(1 - factor*factor*(1 - nDotI*nDotI))
                    
                    trasmittedDirection = trasmittedDirection.Add(normal.Mult(term2 - term3))
                    
                    //transmittedRay := geometry.Ray{impact, ray.Direction}
                    transmittedRay := geometry.Ray{impact, trasmittedDirection.Normalize()}
                    CausticPhoton(scene, emitter, transmittedRay, colour, result, alpha*0.9, depth+1)
                }
            }
		}
	}
}

func DiffusePhoton(scene *list.List, emitter geometry.Shape, ray geometry.Ray, colour geometry.Vec3, result chan<- PhotonHit, alpha float64, depth int) {
	if rand.Float64() > alpha {
		return
	}
	if shape, distance := ClosestIntersection(scene, ray); shape != nil {
		impact := ray.Origin.Add(ray.Direction.Mult(distance))
        
		if depth == 0 && emitter == shape {
			// Leave the emitter first
			nextRay := geometry.Ray{impact, ray.Direction}
			DiffusePhoton(scene, emitter, nextRay, colour, result, alpha, depth)
		} else {
			normal := shape.NormalDir(impact).Normalize()
			reverse := ray.Direction.Mult(-1)
			outgoing := normal
			if normal.Dot(reverse) < 0 {
				outgoing = normal.Mult(-1)
			}
			outgoing = outgoing
			//fmt.Println("Hit something else!")
            strength := colour.Mult(alpha / (1 + distance))
            result <- PhotonHit{impact, strength, ray.Direction, uint8(depth)}

			if shape.Material() == geometry.DIFFUSE {
				// Random bounce for color bleeding
				u := normal.Cross(reverse).Normalize()
				v := u.Cross(normal).Normalize()
				bounce := u.Mult(rand.NormFloat64() * 0.5).Add(outgoing).Add(v.Mult(rand.NormFloat64() * 0.5))
				bounceRay := geometry.Ray{impact, bounce.Normalize()}
                bleedColour := colour.MultVec(shape.Colour()).Mult(alpha / (1 + distance))
				DiffusePhoton(scene, shape, bounceRay, bleedColour, result, alpha*0.66, depth+1)
			}
			// Store Shadow Photons
			shadowRay := geometry.Ray{impact, ray.Direction}
			DiffusePhoton(scene, shape, shadowRay, geometry.Vec3{0, 0, 0}, result, alpha*0.66, depth+1)
		}
	}
}

func PhotonChunk(scene *list.List, traceFunc RayFunc, shape geometry.Shape, factor, start, chunksize int, result chan<- PhotonHit, done chan<- bool) {

	for i := 0; i < chunksize; i++ {
		longitude := (start*chunksize + i) / factor
		latitude := (start*chunksize + i) % factor

		//fmt.Println("Lo La:", longitude, latitude)

		sign := -2.0*float64(longitude%2.0) + 1.0
		phi := 2.0 * math.Pi * float64(longitude) / float64(factor)
		theta := math.Pi * float64(latitude) / float64(factor)

		//fmt.Println("S, T, P:", sign, theta, phi)

		x, y, z := math.Sin(theta)*math.Cos(phi),
			sign*math.Cos(theta),
			math.Sin(theta)*math.Sin(phi)

		direction := geometry.Vec3{x, y, z}
		ray := geometry.Ray{shape.Position(), direction.Normalize()}
		traceFunc(scene, shape, ray, shape.Emission(), result, 1.0, 0)
	}
	done <- true
}

func PhotonMapping(scene *list.List, factor int, rayFunc RayFunc) *list.List {

	result := list.New()
	photons := factor * factor * 2
	chunks := 8
	chunksize := photons / chunks

	for e := scene.Front(); e != nil; e = e.Next() {
		shape := e.Value.(geometry.Shape)
		hits := make(chan PhotonHit)
		done := make(chan bool)
		if shape.Emission().Abs() > 0 {
			for start := 0; start < chunks; start++ {
				go PhotonChunk(scene, rayFunc, shape, factor, start, chunksize, hits, done)
			}

			go func() {
				for start := 0; start < chunks; start++ {
					<-done
				}
				close(hits)
			}()

			count := 0
			const tick = 10000
			fmt.Printf("Tracing %v photons through the scene ", photons)
			for photon := range hits {
				result.PushBack(photon)
				count++
				if count%tick == 0 {
					fmt.Printf(".")
					if count%(10*tick) == 0 {
						clearLine()
						fmt.Printf("Tracing %v photons through the scene ", photons)
					}
				}
			}
			fmt.Printf("\rTraced %v photons to %v intersections in the scene.          \n", photons, count)
		}
	}
	return result
}

func GenerateMaps(scene *list.List) (*kd.KDNode, *kd.KDNode) {
	caustics := PhotonMapping(scene, 128, CausticPhoton)
	globals := PhotonMapping(scene, 16, DiffusePhoton)
	fmt.Printf("Building KD-trees ...")

	globalsChannel := kd.AsyncNew(globals, 3)
	causticsChannel := kd.AsyncNew(caustics, 3)
	return <-globalsChannel, <-causticsChannel
}
