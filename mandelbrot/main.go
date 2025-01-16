package mandelbrot

import "image"

// Constants for the range of the Mandelbrot set
const (
	/*
		(xMax - xMin) / (yMax - yMin) = 16:9
	*/
	XMin, XMax = -1.5 * 2.0, 1.5 * 1.0
	YMin, YMax = -1.5 * 0.84375, 1.5 * 0.84375
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
