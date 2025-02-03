package mandelbrot

import (
	"image/color"
	"math/cmplx"
)

// ColorConvergence determines the color of a point based on the Mandelbrot set calculation.
// It returns a color based on the number of iterations it takes for the sequence to escape.
func ColorConvergence(c complex128, nbIteration int) (color.RGBA, error) {
	var z complex128
	for n := 0; n < nbIteration; n++ {
		if cmplx.Abs(z) > 2 {
			// Map the escape iteration to a color gradient.
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}, nil
		}
		z = z*z + c
	}
	return color.RGBA{R: 0, G: 0, B: 0, A: 255}, nil // Points in the Mandelbrot set are black.
}
