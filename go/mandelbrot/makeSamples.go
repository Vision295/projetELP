package mandelbrot

import (
	"image/color"
	"math/cmplx"
)

// ColorConvergence determines the color of a point based on the Mandelbrot set calculation.
// It returns a color based on the number of iterations it takes for the sequence to escape.
// It requires a Mandelbrot to get its MaxIterations so that it doesn't exceed
func (m *Mandelbrot) ColorConvergence(c complex128, nbIteration int) (color.Color, error) {
	var z complex128
	for n := 0; n < nbIteration; n++ {
		if cmplx.Abs(z) > 2 {
			// Map the escape iteration to a color gradient.
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}, nil
		}
		z = z*z + c
	}
	return color.Black, nil // Points in the Mandelbrot set are black.
}

/*
func (m *Mandelbrot) GenerateSample() [NMaxPoints]complex128 {
	var sample [NMaxPoints]complex128
	// a faire la gestion d'erreurs
	for i := uint32(0); i < m.Width; i++ {
		for j := uint32(0); j < m.Height; j++ {
			sample[i*m.Height+j] = complex(float64(i), float64(j))
		}
	}
	return sample
}
*/
