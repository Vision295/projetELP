package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
	"time"
)

// mandelbrot determines the color of a point based on the Mandelbrot set calculation.
// It returns a color based on the number of iterations it takes for the sequence to escape.
func mandelbrot(c complex128, maxIterations int) color.Color {
	var z complex128
	for n := 0; n < maxIterations; n++ {
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}
		}
		z = z*z + c
	}
	return color.Black
}

func renderZone(img *image.RGBA, startY, endY, width, height, maxIterations int, xMin, xMax, yMin, yMax float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for y := startY; y < endY; y++ {
		yCoord := float64(y)/float64(height)*(yMax-yMin) + yMin
		for x := 0; x < width; x++ {
			xCoord := float64(x)/float64(width)*(xMax-xMin) + xMin
			c := complex(xCoord, yCoord)
			color := mandelbrot(c, maxIterations)
			img.Set(x, y, color)
		}
	}
}

func main() {
	const (
		width, height = 50000, 50000
		xMin, xMax    = -2.0, 1.0
		yMin, yMax    = -1.5, 1.5
		numGoroutines = 8
		maxIterations = 2000
	)

	start := time.Now()

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	var wg sync.WaitGroup
	rowsPerGoroutine := (height / 2) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		startY := i * rowsPerGoroutine
		endY := startY + rowsPerGoroutine
		if i == numGoroutines-1 {
			endY = height / 2
		}
		wg.Add(1)
		go renderZone(img, startY, endY, width, height, maxIterations, xMin, xMax, yMin, yMax, &wg)
	}

	wg.Wait()

	for y := 0; y < height/2; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, height-y-1, img.At(x, y))
		}
	}

	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	fmt.Printf("Time taken to generate Mandelbrot set: %.2f seconds\n", elapsedSeconds)

	fileName := fmt.Sprintf("mandelbrot_with_%dgoroutines_%dx%d_%diter_%.2fs.png", numGoroutines, width, height, maxIterations, elapsedSeconds)
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	fmt.Printf("Mandelbrot set image saved as %s\n", fileName)
}
