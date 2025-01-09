package main

import (
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
)

// mandelbrot determines the color of a point based on the Mandelbrot set calculation.
func mandelbrot(c complex128) color.Color {
	const maxIterations = 100
	var z complex128
	for n := 0; n < maxIterations; n++ {
		if cmplx.Abs(z) > 2 {
			return color.RGBA{R: uint8(255 - n*5), G: uint8(n * 2), B: uint8(255 - n*10), A: 255}
		}
		z = z*z + c
	}
	return color.Black
}

func main() {
	const (
		width, height = 1024, 1024
		xMin, xMax    = -2.0, 1.0
		yMin, yMax    = -1.5, 1.5
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for py := 0; py < height; py++ {
		y := float64(py)/height*(yMax-yMin) + yMin
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xMax-xMin) + xMin
			c := complex(x, y)
			img.Set(px, py, mandelbrot(c))
		}
	}

	file, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	println("Mandelbrot set image saved as mandelbrot.png")
}
