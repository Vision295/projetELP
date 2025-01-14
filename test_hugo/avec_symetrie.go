package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
)

// mandelbrot determines the color of a point based on the Mandelbrot set calculation.
// It returns a color based on the number of iterations it takes for the sequence to escape.
func mandelbrot(c complex128) color.Color {
	const maxIterations = 3000 // Increase the number of iterations for more precision.
	var z complex128
	for n := 0; n < maxIterations; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 { // Avoid square root for performance
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
		width, height = 50000, 50000
		xMin, xMax    = -2.0, 1.0
		yMin, yMax    = -1.5, 1.5
		maxIterations = 3000
	)

	// Start the timer to measure execution time.
	start := time.Now()

	// Create a new blank image with the specified resolution.
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Iterate over the top half of the image.
	for py := 0; py < height/2; py++ {
		y := float64(py)/height*(yMax-yMin) + yMin // Map pixel y-coordinate to complex plane.
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xMax-xMin) + xMin // Map pixel x-coordinate to complex plane.
			c := complex(x, y)                        // Create the complex number for the current pixel.
			color := mandelbrot(c)                    // Get the color for the pixel.
			img.Set(px, py, color)                    // Set the color for the top half.
			img.Set(px, height-py-1, color)           // Mirror the color for the bottom half.
		}
	}

	// Measure the elapsed time.
	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	fmt.Printf("Time taken to generate Mandelbrot set: %.2f seconds\n", elapsedSeconds)

	// Create a file name with size, iteration information, and elapsed time.
	fileName := fmt.Sprintf("mandelbrot_using_symetry_%dx%d_%diter_%.2fs.png", width, height, maxIterations, elapsedSeconds)
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
