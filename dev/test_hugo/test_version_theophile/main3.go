package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"sync"
	"time"
)

// Constants for the range of the Mandelbrot set
const (
	XMin, XMax = -2.0, 1.0
	YMin, YMax = -1.5, 1.5
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

// ColorConvergence determines the color of a point based on its convergence in the Mandelbrot set.
func (m *Mandelbrot) ColorConvergence(c complex128) color.Color {
	var z complex128
	const maxIterations = 1000 // Fixed maximum iterations
	for n := 0; n < maxIterations; n++ {
		if cmplx.Abs(z) > 2 {
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}
		}
		z = z*z + c
	}
	return color.Black
}

// PrintOnImage generates the Mandelbrot image using parallel processing.
func (m *Mandelbrot) PrintOnImage(numGoroutines int) error {
	m.Image = image.NewRGBA(image.Rect(0, 0, int(m.Width), int(m.Height)))

	var wg sync.WaitGroup
	rowsPerGoroutine := int(m.Height) / numGoroutines

	for g := 0; g < numGoroutines; g++ {
		startRow := uint32(g * rowsPerGoroutine)
		endRow := uint32((g + 1) * rowsPerGoroutine)
		if g == numGoroutines-1 {
			endRow = m.Height
		}

		wg.Add(1)
		go func(start, end uint32) {
			defer wg.Done()
			for i := start; i < end; i++ {
				for j := uint32(0); j < m.Width; j++ {
					c := complex(
						float64(j)/float64(m.Width)*(m.XMax-m.XMin)+m.XMin,
						float64(i)/float64(m.Height)*(m.YMax-m.YMin)+m.YMin,
					)
					color := m.ColorConvergence(c)
					m.Image.Set(int(j), int(i), color)
				}
			}
		}(startRow, endRow)
	}

	wg.Wait()
	return nil
}

// SaveImage saves the generated Mandelbrot image as a PNG file.
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
	// Define image dimensions
	const width, height = 10000, 10000
	const numGoRoutines = 1000

	mandelbrot := NewMandelbrot(width, height)

	start := time.Now()
	err := mandelbrot.PrintOnImage(numGoRoutines)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error generating Mandelbrot image:", err)
		return
	}

	// Save the image with a name based on dimensions
	fileName := fmt.Sprintf("Mandelbrot_image_(%dx%d)_%v_with_%dgoroutines!", width, height, elapsed, numGoRoutines)
	err = mandelbrot.SaveImage(fileName)
	if err != nil {
		fmt.Println("Error saving image:", err)
	} else {
		fmt.Printf("Mandelbrot_image_(%dx%d)_%v_with_%dgoroutines!\n", width, height, elapsed, numGoRoutines)
	}
}
