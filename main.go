package main

import (
	"image/png"
	. "mandelbrot/mandelbrot"
	"os"
)

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
