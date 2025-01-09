package main

import (
	// Create and manipulate images
	"image"
	"image/color"
	"image/png"

	// Complex functions
	"math/cmplx"

	// Handeling file operations (save)
	"os"
)

const MaxIterations = 1000 // Increase the number of iterations for more precision.
const (
	// Increase the resolution for more precise rendering.
	width, height = 2000, 2000
	xMin, xMax    = -2.0, 1.0
	yMin, yMax    = -1.5, 1.5
)

// mandelbrot determines the color of a point based on the Mandelbrot set calculation.
// It returns a color based on the number of iterations it takes for the sequence to escape.
func mandelbrot_test(c complex128, nbIteration int) color.Color {
	var z complex128
	// todo : tester si nbIteration < MaxIterations
	for n := 0; n < nbIteration; n++ {
		if cmplx.Abs(z) > 2 {
			// Map the escape iteration to a color gradient.
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}
		}
		z = z*z + c
	}
	return color.Black // Points in the Mandelbrot set are black.
}

func printOnImage() *image.RGBA {
	// Create a new blank image with the specified resolution.
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate over each pixel in the image.
	for py := 0; py < height; py++ {
		y := float64(py)/height*(yMax-yMin) + yMin // Map pixel y-coordinate to complex plane.
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xMax-xMin) + xMin // Map pixel x-coordinate to complex plane.
			c := complex(x, y)                        // Create the complex number for the current pixel.
			img.Set(px, py, mandelbrot_test(c, 100))  // Set the color of the pixel based on the Mandelbrot set.
		}
	}
	return img
}

func main() {
	var img = printOnImage()
	// Create a file to save the image.
	file, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image to PNG format and save it to the file.
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	println("Mandelbrot set image saved as mandelbrot.png")
}
