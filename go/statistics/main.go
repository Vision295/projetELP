package main

import (
	"fmt"
	. "mandelbrot/mandelbrot"
	"os"
	"strconv"
	"time"
)

func main() {
	// Check if sufficient arguments are provided
	if len(os.Args) < 5 {
		fmt.Println("Error: Insufficient arguments provided.")
		fmt.Println("Usage: go run main.go <width> <height> <numGoRoutines> <nbIteration>")
		os.Exit(1)
	}

	// Parse arguments
	width, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error parsing width: %v\n", err)
		os.Exit(1)
	}

	height, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Error parsing height: %v\n", err)
		os.Exit(1)
	}

	numGoRoutines, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Error parsing numGoRoutines: %v\n", err)
		os.Exit(1)
	}

	nbIteration, err := strconv.Atoi(os.Args[4])
	if err != nil {
		fmt.Printf("Error parsing nbIteration: %v\n", err)
		os.Exit(1)
	}

	// Convert dimensions to uint32 for NewMandelbrot
	mandelbrot := NewMandelbrot(uint32(width), uint32(height))

	start := time.Now()
	err = mandelbrot.PrintOnImage(numGoRoutines, nbIteration)
	elapsed := time.Since(start)

	if err != nil {
		fmt.Println("Error generating Mandelbrot image:", err)
		return
	}

	// Corrected filename generation order
	fileName := fmt.Sprintf("Mandelbrot_%dx%d_%dIterations_%dGoroutines.png", width, height, nbIteration, numGoRoutines)
	err = mandelbrot.SaveImage(fileName)
	if err != nil {
		fmt.Println("Error saving image:", err)
	} else {
		fmt.Printf("Image saved as %s in %v\n", fileName, elapsed)
	}
}
