package mandelbrot

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
)

// PrintOnImage generates the Mandelbrot image using parallel processing.
func PrintOnImage(m Mandelbrot, filePath string, numGoroutines, nbIterations int) error {
	/*
		prints mandelbrot onto an image
			using numGoroutines goroutines slicing the image into vertical slices
			with a precision of nbIteration iterations
	*/
	// initial values
	var wg sync.WaitGroup
	rowsPerGoroutine := m.Height / numGoroutines

	// creates a list of rows to store the result of computations in each goroutine
	rowList := make(chan [][]color.RGBA, numGoroutines)
	rowOrders := make(chan int, numGoroutines)

	for routineStep := 0; routineStep < numGoroutines; routineStep++ {
		startRow := int(routineStep * rowsPerGoroutine)
		endRow := int((routineStep + 1) * rowsPerGoroutine)

		if routineStep == numGoroutines-1 {
			endRow = m.Height
		}

		wg.Add(1)

		// starts a go routine to compute points from startRow to endRow
		// it will compute the image in numGoroutine vertical sections
		go ComputeOnSample(rowList, rowOrders, m, &wg, nbIterations, routineStep, startRow, endRow)

	}

	// waits termination of all goroutines
	wg.Wait()

	SaveImage(rowList, rowOrders, m, numGoroutines, filePath)
	return nil
}

func ComputeOnSample(rowList chan [][]color.RGBA, rowOrder chan int, m Mandelbrot, wg *sync.WaitGroup, nbIterations, routineStep, start, end int) error {
	/*
		computes mandelbrot on a sample from where the x and y coordinate varies like this :
			x : from start to end
			y : from 0 to width
			lr : list of the rows of the image
	*/
	// ensures the waitgroup gets a return value after execution of this method
	defer wg.Done()
	colors := make([][]color.RGBA, end-start)

	// Initialize each row in the 2D slice
	for i := 0; i < end-start; i++ {
		colors[i] = make([]color.RGBA, m.Width)
	}

	for i := 0; i < end-start; i++ {
		for j := 0; j < m.Width; j++ {
			c := complex(
				float64(j)/float64(m.Width)*(m.XMax-m.XMin)+m.XMin,
				float64(i+start)/float64(m.Height)*(m.YMax-m.YMin)+m.YMin,
			)
			err := error(nil)
			colors[i][j], err = ColorConvergence(c, nbIterations)
			// sets the pixel of coordinate (i, j) to color : color.
			if err != nil {
				return fmt.Errorf("tried to apply a color to a pixel out of image \n coordinate : (%v, %v) ", i, j)
			}
		}
	}
	// send colors to channel
	rowList <- colors
	// sends index to channel
	rowOrder <- routineStep
	return nil
}

// SaveImage saves the generated Mandelbrot image as a PNG file.
func SaveImage(rowList chan [][]color.RGBA, rowOrder chan int, m Mandelbrot, numGoroutines int, filePath string) error {
	/*
		Prints the output of whats computed (which is inside m.Image) onto a
			file in filePath path
	*/
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}

	imageList := make([][]color.RGBA, m.Height)
	for i := 0; i < m.Height; i++ {
		imageList[i] = make([]color.RGBA, m.Width)
	}

	close(rowList)
	close(rowOrder)
	// need to recreate the image here
	for val := range rowList {
		index := <-rowOrder
		for j := 0; j < len(val); j++ {
			imageList[index*len(val)+j] = val[j]
		}
	}

	image := image.NewRGBA(image.Rect(0, 0, m.Width, m.Height))

	for i := 0; i < len(imageList); i++ {
		for j := 0; j < len(imageList[i]); j++ {
			image.Set(j, i, imageList[i][j])
		}
	}

	// ensure the closure of the file before putting the pixels into it
	defer file.Close()
	err = png.Encode(file, image)

	if err != nil {
		return fmt.Errorf("could not encode image to file: %v", err)
	}

	return nil
}
