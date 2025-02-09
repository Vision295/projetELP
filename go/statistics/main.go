package main

import (
	"fmt"
	. "mandelbrot/mandelbrot"
	"os"
	"strconv"
	"time"
)

func main() {
	// Check if sufficient arguments are provided
	if len(os.Args) < 3 {
		fmt.Println("Error: Insufficient arguments provided.")
		fmt.Println("Usage: go run main.go <numGoRoutines> <nbIteration>")
		os.Exit(1)
	}

	// Parse arguments
	numGoRoutines, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing numGoRoutines: %v\n", err)
		os.Exit(1)
	}

	nbIteration, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error parsing nbIteration: %v\n", err)
		os.Exit(1)
	}

	// Define image dimensions
	const width, height = 3840, 2160

	mandelbrot := NewMandelbrot(width, height)
	/*
	   mandelbrot.XMin = -1
	   mandelbrot.XMax = 0.5
	   mandelbrot.YMin = -0.75
	   mandelbrot.YMax = 0.75
	*/

	start := time.Now()
	fileName := fmt.Sprintf("Mandelbrot_%dx%d_%dIterations_%dGoroutines.png", width, height, nbIteration, numGoRoutines)
	fmt.Printf("Generating Mandelbrot image with %d goroutines and %d iterations...\n", numGoRoutines, nbIteration)
	err = PrintOnImage(mandelbrot, fileName, numGoRoutines, nbIteration)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Printf("Error generating Mandelbrot image: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Mandelbrot_%dx%d_%dIterations_%dGoroutines.png generated in %v\n", width, height, nbIteration, numGoRoutines, elapsed)
}
