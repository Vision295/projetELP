package mandelbrot

import "image"

// Constants for the range of the Mandelbrot set
const (
	XMin, XMax = -1, 0.5
	YMin, YMax = -0.75, 0.75
)

type Mandelbrot struct {
	Width, Height uint32
	XMin, XMax    float64
	YMin, YMax    float64
	Image         *image.RGBA
}

// NewMandelbrot initializes a new Mandelbrot set configuration with specified dimensions.
func NewMandelbrot(width, height uint32) Mandelbrot {
	return Mandelbrot{
		Width:  width,
		Height: height,
		XMin:   XMin,
		XMax:   XMax,
		YMin:   YMin,
		YMax:   YMax,
		Image:  nil,
	}
}
