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
	Width, Height = 10000, 10000 // Width and height of the image
	XMin, XMax    = -2.0, 1.0    // X-axis range for the Mandelbrot set
	YMin, YMax    = -1.5, 1.5    // Y-axis range for the Mandelbrot set
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

// NewMandelbrot initializes a new Mandelbrot set configuration with default values.
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

// ColorConvergence determines the color of a point based on its convergence in the Mandelbrot set.
// It returns a color depending on how quickly the point escapes the set.
func (m *Mandelbrot) ColorConvergence(c complex128, nbIteration int) (color.Color, error) {
	var z complex128
	if nbIteration > m.MaxIterations {
		return nil, fmt.Errorf("tried to iterate %v whereas the maxIterations allowed is %v", nbIteration, MaxIterations)
	}
	// Perform the Mandelbrot set iterations.
	for n := 0; n < nbIteration; n++ {
		if cmplx.Abs(z) > 2 {
			// Return a color based on the number of iterations before escaping.
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}, nil
		}
		z = z*z + c // z = z^2 + c
	}
	return color.Black, nil // If it doesn't escape, return black (point is in the Mandelbrot set)
}

// GenerateSample generates a slice of complex numbers representing the points to be calculated
// for a given range of rows (startRow to endRow).
func (m *Mandelbrot) GenerateSample(startRow, endRow uint32) []complex128 {
	sample := make([]complex128, 0, (endRow-startRow)*uint32(m.Width))
	// Iterate through columns and rows to generate complex numbers
	for i := startRow; i < endRow; i++ {
		for j := uint32(0); j < m.Width; j++ {
			// Each complex number corresponds to a pixel in the image
			sample = append(sample, complex(float64(j), float64(i)))
		}
	}
	return sample
}

// PrintOnImage generates the Mandelbrot image using parallel processing.
// It divides the image into chunks of rows and processes each chunk using a separate goroutine.
func (m *Mandelbrot) PrintOnImage(precision, numGoroutines int) error {
	(*m).Image = image.NewRGBA(image.Rect(0, 0, int(m.Width), int(m.Height))) // Create a new RGBA image

	var wg sync.WaitGroup
	rowsPerGoroutine := int(m.Height) / numGoroutines // Divide the image height by the number of goroutines

	// Launch goroutines to process chunks of rows in parallel.
	for g := 0; g < numGoroutines; g++ {
		// Calculate the range of rows that this goroutine will process
		startRow := uint32(g * rowsPerGoroutine)
		endRow := uint32((g + 1) * rowsPerGoroutine)
		if g == numGoroutines-1 {
			// Ensure the last goroutine processes all remaining rows
			endRow = m.Height
		}

		wg.Add(1) // Increment the wait group counter, indicating a new goroutine is being launched
		go func(start, end uint32) {
			defer wg.Done() // Decrement the wait group counter when the goroutine is done

			// Generate the sample (complex numbers) for the assigned range of rows
			sample := m.GenerateSample(start, end)

			// Process each complex number (pixel) and assign a color based on its convergence
			for _, v := range sample {
				mandelbrotRes, err := m.ColorConvergence(
					complex(
						float64(real(v))/float64(m.Width)*(m.XMax-m.XMin)+m.XMin,
						float64(imag(v))/float64(m.Height)*(m.YMax-m.YMin)+m.YMin,
					), precision)
				if err == nil {
					// Set the pixel color in the image
					(*m).Image.Set(int(real(v)), int(imag(v)), mandelbrotRes)
				}
			}
		}(startRow, endRow) // Pass the row range to the goroutine
	}

	wg.Wait() // Wait for all goroutines to complete
	return nil
}

// SaveImage saves the generated Mandelbrot image as a PNG file at the specified path.
func (m *Mandelbrot) SaveImage(filePath string) error {
	file, err := os.Create(filePath) // Create a new file to save the image
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	err = png.Encode(file, m.Image) // Encode the image to the file as a PNG
	if err != nil {
		return fmt.Errorf("could not encode image to file: %v", err)
	}

	return nil
}

func main() {
	mandelbrot := NewMandelbrot() // Initialize a new Mandelbrot instance

	// Generate the Mandelbrot image using 4 goroutines and the specified precision
	err := mandelbrot.PrintOnImage(MaxIterations, 10)
	if err != nil {
		fmt.Println("Error generating Mandelbrot image:", err)
		return
	}

	// Save the generated image to a PNG file
	err = mandelbrot.SaveImage("mandelbrot.png")
	if err != nil {
		fmt.Println("Error saving image:", err)
	} else {
		fmt.Println("Mandelbrot image saved successfully!")
	}
}
