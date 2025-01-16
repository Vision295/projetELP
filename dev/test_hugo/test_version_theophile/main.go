package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
)

const NMaxPoints = 100000
const MaxIterations = 1000
const (
	Width, Height = 10000, 10000
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

// ColorConvergence determines the color of a point based on the Mandelbrot set calculation.
func (m *Mandelbrot) ColorConvergence(c complex128, nbIteration int) (color.Color, error) {
	var z complex128
	if nbIteration > m.MaxIterations {
		return nil, fmt.Errorf("tried to iterate %v whereas the maxIterations allowed is %v", nbIteration, MaxIterations)
	}
	for n := 0; n < nbIteration; n++ {
		if cmplx.Abs(z) > 2 {
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}, nil
		}
		z = z*z + c
	}
	return color.Black, nil
}

func (m *Mandelbrot) GenerateSample(startRow, endRow uint32) []complex128 {
	sample := make([]complex128, 0, (endRow-startRow)*uint32(m.Width))
	for i := startRow; i < endRow; i++ {
		for j := uint32(0); j < m.Width; j++ {
			sample = append(sample, complex(float64(j), float64(i)))
		}
	}
	return sample
}

// PrintOnImage generates the Mandelbrot image using parallel processing.
func (m *Mandelbrot) PrintOnImage(precision, numGoroutines int) error {
	(*m).Image = image.NewRGBA(image.Rect(0, 0, int(m.Width), int(m.Height)))

	var wg sync.WaitGroup
	rowsPerGoroutine := int(m.Height) / numGoroutines

	// Launch goroutines for processing chunks of rows.
	for g := 0; g < numGoroutines; g++ {
		startRow := uint32(g * rowsPerGoroutine)
		endRow := uint32((g + 1) * rowsPerGoroutine)
		if g == numGoroutines-1 {
			endRow = m.Height
		}

		wg.Add(1)
		go func(start, end uint32) {
			defer wg.Done()
			sample := m.GenerateSample(start, end)
			for _, v := range sample {
				mandelbrotRes, err := m.ColorConvergence(
					complex(
						float64(real(v))/float64(m.Width)*(m.XMax-m.XMin)+m.XMin,
						float64(imag(v))/float64(m.Height)*(m.YMax-m.YMin)+m.YMin,
					), precision)
				if err == nil {
					(*m).Image.Set(int(real(v)), int(imag(v)), mandelbrotRes)
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()
	return nil
}

// SaveImage saves the generated Mandelbrot image to a PNG file.
func (m *Mandelbrot) SaveImage(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, m.Image)
	if err != nil {
		return fmt.Errorf("could not encode image to file: %v", err)
	}

	return nil
}

func main() {
	mandelbrot := NewMandelbrot()
	err := mandelbrot.PrintOnImage(MaxIterations, 4) // Using 4 goroutines
	if err != nil {
		fmt.Println("Error generating Mandelbrot image:", err)
		return
	}

	err = mandelbrot.SaveImage("mandelbrot.png")
	if err != nil {
		fmt.Println("Error saving image:", err)
	} else {
		fmt.Println("Mandelbrot image saved successfully!")
	}
}
