package mandelbrot

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"sync"
)

/*
func (m *Mandelbrot) PrintOnImage(precision int) error {
	// Create a new blank image with the specified resolution.
	(*m).Image = image.NewRGBA(image.Rect(0, 0, int(m.Width), int(m.Height)))
	// Iterate over each pixel in the image.
	sample := (*m).GenerateSample()
	for _, v := range sample {
		mandelbrotRes, err := (*m).ColorConvergence(
			complex(
				float64(real(v))/float64(m.Width)*((*m).XMax-XMin)+float64(XMin),  // Map pixel y-coordinate to complex plane.
				float64(imag(v))/float64(m.Height)*((*m).YMax-YMin)+float64(YMin), // Map pixel y-coordinate to complex plane.
			), precision)
		if err == nil {
			(*m).Image.Set(int(real(v)), int(imag(v)), mandelbrotRes)
		} else {
			return err
		}
	}
	return nil
}
*/

// PrintOnImage generates the Mandelbrot image using parallel processing.
func (m *Mandelbrot) PrintOnImage(numGoroutines, nbIterations int) error {
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
		go m.ComputeOnSample(&wg, nbIterations, startRow, endRow)

	}

	wg.Wait()
	return nil
}

func (m *Mandelbrot) ComputeOnSample(wg *sync.WaitGroup, nbIterations int, start, end uint32) {
	defer wg.Done()
	for i := start; i < end; i++ {
		for j := uint32(0); j < m.Width; j++ {
			c := complex(
				float64(j)/float64(m.Width)*(m.XMax-m.XMin)+m.XMin,
				float64(i)/float64(m.Height)*(m.YMax-m.YMin)+m.YMin,
			)
			color, _ := m.ColorConvergence(c, nbIterations)
			m.Image.Set(int(j), int(i), color)
		}
	}
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
