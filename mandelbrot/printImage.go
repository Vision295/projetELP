package mandelbrot

import (
	"image"
)

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
