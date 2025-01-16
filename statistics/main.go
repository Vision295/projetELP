/*
	On cherche le meilleur temps
	nbIteration = 100
	nbGoroutine = X
	width = 2160
	height = 3840
*/

package main

import (
	"fmt"
	. "mandelbrot/mandelbrot"
	"time"
)

func main() {
	// Define image dimensions
	const width, height = 1000, 1000
	const numGoRoutines = 100
	const nbIteration = 1000

	mandelbrot := NewMandelbrot(width, height)
	/*
		mandelbrot.XMin = -1
		mandelbrot.XMax = 0.5
		mandelbrot.YMin = -0.75
		mandelbrot.YMax = 0.75
	*/

	start := time.Now()
	err := mandelbrot.PrintOnImage(numGoRoutines, nbIteration)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error generating Mandelbrot image:", err)
		return
	}

	// Save the image with a name based on dimensions
	fileName := fmt.Sprintf("Mandelbrot_image_(%dx%d)_%v_with_%dgoroutines.png.png", width, height, elapsed, numGoRoutines)
	err = mandelbrot.SaveImage(fileName)
	if err != nil {
		fmt.Println("Error saving image:", err)
	} else {
		fmt.Printf("Mandelbrot_image_(%dx%d)_%v_with_%dgoroutines!\n", width, height, elapsed, numGoRoutines)
	}
}
