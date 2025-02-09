package main

import (
	"fmt"
	. "mandelbrot/mandelbrot"
	"time"
)

func main() {
	// Define image dimensions
	const width, height = 3840, 2160
	const numGoRoutines = 10
	const nbIteration = 10000

	mandelbrot := NewMandelbrot(width, height)
	/*
		mandelbrot.XMin = -1
		mandelbrot.XMax = 0.5
		mandelbrot.YMin = -0.75
		mandelbrot.YMax = 0.75
	*/

	start := time.Now()
	fileName := fmt.Sprintf("Mandelbrot_image_(%dx%d)_with_%dgoroutines.png.png", width, height, numGoRoutines)
	err := PrintOnImage(mandelbrot, fileName, numGoRoutines, nbIteration)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error generating Mandelbrot image:", err)
		return
	}

	// Save the image with a name based on dimensions
	if err != nil {
		fmt.Println("Error saving image:", err)
	} else {
		fmt.Printf("Mandelbrot_image_(%dx%d)_%v_with_%dgoroutines!\n", width, height, elapsed, numGoRoutines)
	}
}
