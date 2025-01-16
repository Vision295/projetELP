package main

import (
	"fmt"
	. "mandelbrot/mandelbrot"
	"time"
)

/*
var mandelbrot = NewMandelbrot()

func main() {
	var errOnImagePrint = (&mandelbrot).PrintOnImage(100)
	if errOnImagePrint != nil {
		panic(errOnImagePrint)
	}

	// Create a file to save the image.
	file, errOnImageCreation := os.Create("mandelbrot.png")
	if errOnImageCreation != nil {
		panic(errOnImageCreation)
	}
	defer file.Close()

	// Encode the image to PNG format and save it to the file.
	if err := png.Encode(file, mandelbrot.Image); err != nil {
		panic(err)
	}

	println("Mandelbrot set image saved as mandelbrot.png")
}
*/

func main() {
	// Define image dimensions
	const width, height = 10000, 10000
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
	fileName := fmt.Sprintf("Mandelbrot_image_(%dx%d)_%v_with_%dgoroutines.png", width, height, elapsed, numGoRoutines)
	err = mandelbrot.SaveImage(fileName)
	if err != nil {
		fmt.Println("Error saving image:", err)
	} else {
		fmt.Printf("Mandelbrot_image_(%dx%d)_%v_with_%dgoroutines!\n", width, height, elapsed, numGoRoutines)
	}
}
