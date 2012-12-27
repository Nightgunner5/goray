package main

import (
	"flag"
	"fmt"
	"github.com/Nightgunner5/goray/geometry"
	"github.com/Nightgunner5/goray/gorender"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
)

var (
	input    = new(string) //flag.String("in", "default", "The file describing the scene")
	cores    = flag.Int("cores", 2, "The number of cores to use on the machine")
	chunks   = flag.Int("chunks", 8, "The number of chunks to use for parallelism")
	fov      = flag.Int("fov", 90, "The field of view of the rendered image")
	cols     = flag.Int("w", 800, "The width in pixels of the rendered image")
	rows     = flag.Int("h", 600, "The height in pixels of the rendered image")
	seed     = flag.Int64("seed", 1, "The seed for the random number generator")
	output   = flag.String("out", "out.png", "Output file for the rendered scene")
	bloom    = flag.Int("bloom", 10, "The number of iteration to run the bloom filter")
	mindepth = flag.Int("depth", 2, "The minimum recursion depth used for the rays")
	rays     = flag.Int("rays", 10, "The number of rays used to sample each pixel")
	caustics = flag.Int("caustics", 256, "The depth of the caustic photon tracing before the render")
	gamma    = flag.Float64("gamma", 2.2, "The factor to use for gamma correction")
	// Profiling information
	cpuprofile = flag.String("cpuprofile", "", "Write cpu profile informaion to file")
	memprofile = flag.String("memprofile", "", "Write memory profile informaion to file")
)

func main() {
	flag.Parse()

	rand.Seed(*seed)

	gorender.Config.NumRays = *rays
	gorender.Config.Caustics = *caustics
	gorender.Config.BloomFactor = *bloom
	gorender.Config.MinDepth = *mindepth
	gorender.Config.GammaFactor = *gamma

	wantedCPUs := *cores
	if wantedCPUs < 1 {
		wantedCPUs = 1
	}
	fmt.Printf("Running on %v/%v CPU cores\n", wantedCPUs, runtime.NumCPU())
	runtime.GOMAXPROCS(wantedCPUs)

	if wantedCPUs > *chunks {
		*chunks = wantedCPUs * 2
		log.Printf("Warning: changed chunks to %d to accomodate cores setting", *chunks)
	}

	if *rows%*chunks != 0 {
		log.Fatal("The images height needs to be evenly divisible by chunks")
	}

	gorender.Config.Chunks = *chunks

	if *cpuprofile != "" {
		cpupf, err := os.Create(*cpuprofile)
		fmt.Println("Writing CPU profiling information to file:", *cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(cpupf)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		fmt.Println("Writing Memory profiling information to file:", *memprofile)
	} else {
		runtime.MemProfileRate = 0
	}

	file, err := os.OpenFile(*output, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Rendering %vx%v sized image with %v rays per pixel to %v\n", *cols, *rows, *rays, *output)

	// "Real world" frustrum
	height := 2.0
	width := height * (float64(*cols) / float64(*rows)) // Aspect ratio?
	angle := math.Pi * float64(*fov) / 180.0

	scene := geometry.ParseScene(*input, width, height, angle, *cols, *rows)
	img := gorender.Render(scene)

	if err = png.Encode(file, img); err != nil {
		log.Fatal(err)
	}

	if *memprofile != "" {
		mempf, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(mempf)
		defer mempf.Close()
	}
}
