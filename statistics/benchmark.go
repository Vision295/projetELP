package main

import (
	"encoding/csv"
	"fmt"
	"image/png"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Configuration for the parameter sweep
const (
	OutputFile = "benchmark_v3_results.csv" // File to store the results

	ImageWidth  = 3840 // Fixed width for the image
	ImageHeight = 2160 // Fixed height for the image

	MinIterations = 1    // Minimum number of iterations
	MaxIterations = 1000 // Maximum number of iterations
	IterationStep = 20   // Step size for iterations

	MinGoroutines = 30 // Minimum number of goroutines
	MaxGoroutines = 60 // Maximum number of goroutines
	GoroutineStep = 1  // Step size for goroutines
)

func compareImages(width, height, nbIterations, numGoRoutines int, perfectImage string) (float64, error) {
	// Construct the path to the generated image
	path := fmt.Sprintf("Mandelbrot_%dx%d_%dIterations_%dGoroutines.png", width, height, nbIterations, numGoRoutines)

	// Open the generated image
	genFile, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("error opening generated image: %v", err)
	}
	defer genFile.Close()

	// Open the perfect image
	perfFile, err := os.Open(perfectImage)
	if err != nil {
		return 0, fmt.Errorf("error opening perfect image: %v", err)
	}
	defer perfFile.Close()

	// Decode both images
	genImg, err := png.Decode(genFile)
	if err != nil {
		return 0, fmt.Errorf("error decoding generated image: %v", err)
	}
	perfImg, err := png.Decode(perfFile)
	if err != nil {
		return 0, fmt.Errorf("error decoding perfect image: %v", err)
	}

	// Ensure dimensions match
	genBounds := genImg.Bounds()
	perfBounds := perfImg.Bounds()
	if genBounds != perfBounds {
		return 0, fmt.Errorf("image dimensions do not match")
	}

	// Compute the desired metric
	totalDifference := 0.0
	numPixels := float64(genBounds.Dx() * genBounds.Dy())

	for y := genBounds.Min.Y; y < genBounds.Max.Y; y++ {
		for x := genBounds.Min.X; x < genBounds.Max.X; x++ {
			genR, genG, genB, _ := genImg.At(x, y).RGBA()
			perfR, perfG, perfB, _ := perfImg.At(x, y).RGBA()

			// Compute sqrt((R_diff)^2 + (G_diff)^2 + (B_diff)^2)
			rDiff := float64(genR>>8) - float64(perfR>>8)
			gDiff := float64(genG>>8) - float64(perfG>>8)
			bDiff := float64(genB>>8) - float64(perfB>>8)
			pixelDifference := math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff)

			// Accumulate absolute difference
			totalDifference += math.Abs(pixelDifference)
		}
	}

	// Normalize by the square of the number of pixels
	result := totalDifference / numPixels
	return result, nil
}

func main() {
	// Open the CSV file for writing
	file, err := os.Create(OutputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';' // Set the delimiter to semicolon
	defer writer.Flush()

	// Write the header row to the CSV file
	err = writer.Write([]string{"Iterations", "Goroutines", "ExecutionTime(ms)", "Quality"})
	if err != nil {
		fmt.Printf("Error writing header to CSV: %v\n", err)
		return
	}

	// Sweep through the parameter space
	for iterations := MinIterations; iterations <= MaxIterations; iterations += IterationStep {
		for goroutines := MinGoroutines; goroutines <= MaxGoroutines; goroutines += GoroutineStep {
			var totalExecutionTime int64
			var totalScore float64

			// Run the test 10 times and average the results
			for i := 0; i < 10; i++ {
				start := time.Now()
				cmd := exec.Command("go", "run", "main.go",
					strconv.Itoa(ImageWidth),
					strconv.Itoa(ImageHeight),
					strconv.Itoa(goroutines),
					strconv.Itoa(iterations),
				)
				err := cmd.Run()
				if err != nil {
					fmt.Printf("Error running Mandelbrot program: %v\n", err)
					continue
				}
				executionTime := time.Since(start).Milliseconds()
				totalExecutionTime += executionTime

				// Compare the generated image with the "perfect" image and calculate SSIM
				score, err := compareImages(ImageWidth, ImageHeight, iterations, goroutines, "Perfect_Mandelbrot.png")
				// println(ssimScore)
				if err != nil {
					fmt.Println("Error comparing images:", err)
					continue
				}
				totalScore += score
			}

			// Calculate the average execution time and SSIM
			avgExecutionTime := totalExecutionTime / 10
			avgScore := totalScore / 10

			// Write the result to the CSV file
			err = writer.Write([]string{
				strconv.Itoa(iterations),
				strconv.Itoa(goroutines),
				strconv.FormatInt(avgExecutionTime, 10),
				fmt.Sprintf("%.2f", avgScore),
			})
			if err != nil {
				fmt.Printf("Error writing result to CSV: %v\n", err)
				return
			}

			// Print progress to the console
			fmt.Printf("Completed: Iterations=%d, Goroutines=%d, AvgExecutionTime=%dms, AvgSSIM=%.2f\n", iterations, goroutines, avgExecutionTime, avgScore)
		}
	}

	fmt.Println("Benchmarking completed. Results saved to", OutputFile)
}
