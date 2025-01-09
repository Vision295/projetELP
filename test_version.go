package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"time"
)

// mandelbrot determines the color of a point based on the Mandelbrot set calculation.
// It returns a color based on the number of iterations it takes for the sequence to escape.
func mandelbrot(c complex128) color.Color {
	const maxIterations = 2000 // Increase the number of iterations for more precision.
	var z complex128
	for n := 0; n < maxIterations; n++ {
		if cmplx.Abs(z) > 2 {
			// Map the escape iteration to a color gradient.
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}
		}
		z = z*z + c
	}
	return color.Black // Points in the Mandelbrot set are black.
}

func main() {
	const (
		// Increase the resolution for more precise rendering.
		width, height = 30000, 30000
		xMin, xMax    = -2.0, 1.0
		yMin, yMax    = -1.5, 1.5
		maxIterations = 2000
	)

	// Start the timer to measure execution time.
	start := time.Now()

	// Create a new blank image with the specified resolution.
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate over each pixel in the image.
	for py := 0; py < height; py++ {
		y := float64(py)/height*(yMax-yMin) + yMin // Map pixel y-coordinate to complex plane.
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xMax-xMin) + xMin // Map pixel x-coordinate to complex plane.
			c := complex(x, y)                        // Create the complex number for the current pixel.
			img.Set(px, py, mandelbrot(c))            // Set the color of the pixel based on the Mandelbrot set.
		}
	}

	// Measure the elapsed time.
	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	fmt.Printf("Time taken to generate Mandelbrot set: %.2f seconds\n", elapsedSeconds)

	// Create a file name with size, iteration information, and elapsed time.
	fileName := fmt.Sprintf("mandelbrot_%dx%d_%diter_%.2fs.png", width, height, maxIterations, elapsedSeconds)
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image to PNG format and save it to the file.
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	fmt.Printf("Mandelbrot set image saved as %s\n", fileName)
}
