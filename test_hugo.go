package main

import (
	// Create and manipulate images

	"fmt"
	"image"
	"image/color"
	"image/png"

	// Complex functions
	"math/cmplx"

	// Handeling file operations (save)
	"os"
)

const NMaxPoints = 100000
const MaxIterations = 1000 // Increase the number of iterations for more precision.
const (
	// Increase the resolution for more precise rendering.
	width, height = 200, 200
	xMin, xMax    = -2.0, 1.0
	yMin, yMax    = -1.5, 1.5
)

var image_size = [2]uint32{width, height}
var rangeX = [2]float64{xMin, xMax}
var rangeY = [2]float64{yMin, yMax}

// mandelbrot determines the color of a point based on the Mandelbrot set calculation.
// It returns a color based on the number of iterations it takes for the sequence to escape.
func mandelbrot_test(c complex128, nbIteration int) color.Color {
	var z complex128
	var x int
	// todo : tester si nbIteration < MaxIterations
	for n := 0; n < nbIteration; n++ {
		if cmplx.Abs(z) > 2 {
			// Map the escape iteration to a color gradient.
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}
		}
		z = z*z + c
	}
	return color.Black // Points in the Mandelbrot set are black.
}

func generate_sample(image_size [2]uint32) [NMaxPoints]complex128 {
	var sample [NMaxPoints]complex128
	// a faire la gestion d'erreurs
	for i := uint32(0); i < image_size[0]; i++ {
		for j := uint32(0); j < image_size[1]; j++ {
			sample[i*image_size[1]+j] = complex(
				float64(i),
				float64(j),
			)
		}
	}
	return sample
}

func printOnImage(image_size [2]uint32, rangeX [2]float64, rangeY [2]float64) *image.RGBA {
	// Create a new blank image with the specified resolution.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Iterate over each pixel in the image.
	/*for py := 0; py < height; py++ {
		y := float64(py)/height*(yMax-yMin) + yMin // Map pixel y-coordinate to complex plane.
		for px := 0; px < width; px++ {
			x := float64(px)/width*(xMax-xMin) + xMin // Map pixel x-coordinate to complex plane.
			c := complex(x, y)                        // Create the complex number for the current pixel.
			img.Set(px, py, mandelbrot_test(c, 100))  // Set the color of the pixel based on the Mandelbrot set.
		}
	} */
	sample := generate_sample(image_size)
	for _, v := range sample {
		fmt.Print(v)
		img.Set(int(real(v)), int(imag(v)), mandelbrot_test(
			complex(
				float64(real(v))/float64(image_size[0])*(rangeX[1]-rangeX[0])+float64(rangeX[0]), // Map pixel y-coordinate to complex plane.
				float64(imag(v))/float64(image_size[1])*(rangeY[1]-rangeY[0])+float64(rangeY[0]), // Map pixel y-coordinate to complex plane.
			), 100),
		)
	}
	return img
}

func main() {
	var img = printOnImage(image_size, rangeX, rangeY)
	// Create a file to save the image.
	file, err := os.Create("mandelbrot.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Encode the image to PNG format and save it to the file.
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	println("Mandelbrot set image saved as mandelbrot.png")
}
