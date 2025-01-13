package mandelbrot

import "image"

const NMaxPoints = 100000
const MaxIterations = 1000
const (
	Width, Height = 200, 200
	XMin, XMax    = -2.0, 1.0
	YMin, YMax    = -1.5, 1.5
)

var ImageSizeval = [2]uint32{Width, Height}
var RangeXval = [2]float64{XMin, XMax}
var RangeYval = [2]float64{YMin, YMax}

type Mandelbrot struct {
	NMaxPoints    int
	MaxIterations int
	Width, Height uint32
	XMin, XMax    float64
	YMin, YMax    float64
	Image         *image.RGBA
}

func NewMandelbrot() Mandelbrot {
	return Mandelbrot{
		NMaxPoints:    NMaxPoints,
		MaxIterations: MaxIterations,
		Width:         Width,
		Height:        Height,
		XMin:          XMin,
		XMax:          XMax,
		YMin:          YMin,
		YMax:          YMax,
		Image:         nil,
	}
}
