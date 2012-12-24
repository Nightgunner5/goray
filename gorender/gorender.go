package gorender

import (
	"../geometry"
	"../kd"
	"container/list"
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"runtime"
	"time"
)

//////////////////
// Utility
//////////////////
func PrintDuration(t time.Duration) {
	if hours := math.Floor(t.Hours()); hours > 0 {
		fmt.Printf("%vh ", int(hours))
	}
	if minutes := math.Mod(math.Floor(t.Minutes()), 60); minutes > 0 {
		fmt.Printf("%1.0vm ", int(minutes))
	}
	if seconds := math.Mod(t.Seconds(), 60); seconds > 0 {
		fmt.Printf("%2.3vs", seconds)
	}
}

func clearLine() {
    fmt.Printf("\r                                                                                                          \r")
}

func ClosestIntersection(shapes *list.List, ray geometry.Ray) (geometry.Shape, float64) {
	var closest geometry.Shape
	bestHit := math.Inf(+1)
	for e := shapes.Front(); e != nil; e = e.Next() {
		shape := e.Value.(geometry.Shape)
		if hit := shape.Intersects(ray); hit > 0 && hit < bestHit {
			bestHit = hit
			closest = shape
		}
	}
	return closest, bestHit
}

type Result struct {
	x, y   int
	colour geometry.Vec3
}

const (
	AIR   = 1.0
	GLASS = 1.5
)

func MonteCarloPixel(results chan Result, scene *geometry.Scene, diffuseMap, causticsMap *kd.KDNode, start, rows int) {
	samples := Config.NumRays
	var px, py, dy, dx float64
	var direction, contribution, delta, colourSamples geometry.Vec3

	for y := start; y < start+rows; y++ {
		py = scene.Height - scene.Height*2*float64(y)/float64(scene.Rows)
		for x := 0; x < scene.Cols; x++ {
			px = -scene.Width + scene.Width*2*float64(x)/float64(scene.Cols)
			colourSamples = geometry.Vec3{0, 0, 0}
			for sample := 0; sample < samples; sample++ {
				dy, dx = rand.Float64()*scene.PixH, rand.Float64()*scene.PixW
				delta = geometry.Vec3{px + dx, py + dy, 0}
				direction = delta.Sub(scene.Camera.Origin).Normalize()

				contribution = Radiance(geometry.Ray{scene.Camera.Origin, direction}, scene, diffuseMap, causticsMap, 0, 1.0)
				colourSamples = colourSamples.Add(contribution.Mult(1.0 / float64(samples)))
			}
			results <- Result{x, y, colourSamples}
		}
	}
}

func CorrectColour(x float64) float64 {
	return math.Pow(x, 1.0/Config.GammaFactor)*255 + 0.5
}

func CorrectColours(v geometry.Vec3) geometry.Vec3 {
	return geometry.Vec3{CorrectColour(v.X),
		CorrectColour(v.Y),
		CorrectColour(v.Z)}
}

func mix(a, b geometry.Vec3, factor float64) geometry.Vec3 {
	return a.Mult(1 - factor).Add(b.Mult(factor))
}

func BloomFilter(img [][]geometry.Vec3, depth int) [][]geometry.Vec3 {
	data := make([][]geometry.Vec3, len(img))
	for i, _ := range data {
		data[i] = make([]geometry.Vec3, len(img[0]))
	}

	const box_width = 2
	factor := 1.0 / math.Pow(2*box_width+1, 2)

	source := img
	for iteration := 0; iteration < depth; iteration++ {
		for y := box_width; y < len(img)-box_width; y++ {
			for x := box_width; x < len(img[0])-box_width; x++ {
				colour := geometry.Vec3{0, 0, 0}
				for dy := -box_width; dy <= box_width; dy++ {
					for dx := -box_width; dx <= box_width; dx++ {
						colour = colour.Add(source[y+dy][x+dx].Mult(factor))
					}
				}
				data[y][x] = colour
			}
		}
        fmt.Printf("\rPost Processing %3.0f%%   \r", 100*float64(iteration)/float64(depth))
		source, data = data, source
	}
	return source
}

var Config struct {
	MinDepth    int
	NumRays     int
	Chunks      int
	GammaFactor float64
	BloomFactor int
}

func Render(scene geometry.Scene) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, scene.Cols, scene.Rows))
	pixels := make(chan Result, 128)

	workload := scene.Rows / Config.Chunks

	startTime := time.Now()
	globals, caustics := GenerateMaps(scene.Objects)
    fmt.Println(" Done!")
	fmt.Printf("Diffuse Map depth: %v Caustics Map depth: %v\n", globals.Depth(), caustics.Depth())
	fmt.Printf("Photon Maps Done. Generation took: ")
	stopTime := time.Now()
	PrintDuration(stopTime.Sub(startTime))
    fmt.Println()

	startTime = time.Now()
	for y := 0; y < scene.Rows; y += workload {
		go MonteCarloPixel(pixels, &scene, globals, caustics, y, workload)
	}

	// Write targets for after effects
	data := make([][]geometry.Vec3, scene.Rows)
	peaks := make([][]geometry.Vec3, scene.Rows)
	for i, _ := range data {
		data[i] = make([]geometry.Vec3, scene.Cols)
		peaks[i] = make([]geometry.Vec3, scene.Cols)
	}

	// Collect results
	var so_far time.Duration
	var highest, lowest geometry.Vec3
	var highValue, lowValue float64
	var memory runtime.MemStats
    numPixels := scene.Rows*scene.Cols
	for i := 0; i < numPixels; i++ {
		// Print progress information every 500 pixels
		if i % 500 == 0 {
            //clearLine()
			fmt.Printf("\rRendering %6.2f%%", 100*float64(i)/float64(scene.Rows*scene.Cols))
			so_far = time.Now().Sub(startTime)
            remaining := time.Duration((so_far.Seconds()/float64(i)) * float64(numPixels-i)) * time.Second
            fmt.Printf(" (Time Remaining: ")
            PrintDuration(remaining)
            fmt.Printf(" at %0.1f pps)                \r", float64(i)/so_far.Seconds())
			runtime.ReadMemStats(&memory)
			//fmt.Printf("M/F/kBs/S/L: %d/%d/%d/%d/%d)", memory.Mallocs, memory.Frees, memory.TotalAlloc/1024, memory.Sys/1024, memory.Lookups)
		}
		pixel := <-pixels

		if low := pixel.colour.Abs(); low < lowValue {
			lowValue = low
			lowest = pixel.colour
		}
		if high := pixel.colour.Abs(); high > highValue {
			highValue = high
			highest = pixel.colour
		}
		data[pixel.y][pixel.x] = pixel.colour.CLAMPF()
		peaks[pixel.y][pixel.x] = pixel.colour.PEAKS(0.8)
	}
	fmt.Println("\rRendering 100.00%")

	bloomed := BloomFilter(peaks, Config.BloomFactor)

	for y := 0; y < len(data); y++ {
		for x := 0; x < len(data[0]); x++ {
			colour := data[y][x].Add(bloomed[y][x])
			colour = CorrectColours(colour).CLAMP()
			img.Set(x, y, color.NRGBA{uint8(colour.X), uint8(colour.Y), uint8(colour.Z), 255})
		}
	}
	stopTime = time.Now()
    clearLine()
	fmt.Println("\rDone!")
	fmt.Printf("Brightest pixel: %v intensity: %v\n", highest, highValue)
	fmt.Printf("Dimmest pixel: %v intensity: %v\n", lowest, lowValue)

	// Print duration
	fmt.Printf("Rendering took ")
	PrintDuration(stopTime.Sub(startTime))
    fmt.Println()

	return img.SubImage(img.Bounds())
}