package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// Configuration for the parameter sweep
const (
	OutputFile    = "benchmark_results.csv" // File to store the results
	ImageWidth    = 3840                    // Fixed width for the image
	ImageHeight   = 2160                    // Fixed height for the image
	MinIterations = 10                      // Minimum number of iterations
	MaxIterations = 2000                    // Maximum number of iterations
	IterationStep = 40                      // Step size for iterations
	MinGoroutines = 1                       // Minimum number of goroutines
	MaxGoroutines = 2000                    // Maximum number of goroutines
	GoroutineStep = 20                      // Step size for goroutines
)

func main() {
	// Open the CSV file for writing
	file, err := os.Create(OutputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	// Explicitly use a comma delimiter for CSV files (this is default but explicit here)
	writer := csv.NewWriter(file)
	writer.Comma = ';' // Ensure comma delimiter
	defer writer.Flush()

	// Write the header row to the CSV file
	err = writer.Write([]string{"Iterations", "Goroutines", "ExecutionTime(ms)"})
	if err != nil {
		fmt.Printf("Error writing header to CSV: %v\n", err)
		return
	}

	// Sweep through the parameter space
	for iterations := MinIterations; iterations <= MaxIterations; iterations += IterationStep {
		for goroutines := MinGoroutines; goroutines <= MaxGoroutines; goroutines += GoroutineStep {
			// Run the Mandelbrot program with the current parameters
			start := time.Now()
			cmd := exec.Command("go", "run", "main.go", strconv.Itoa(ImageWidth), strconv.Itoa(ImageHeight), strconv.Itoa(iterations), strconv.Itoa(goroutines))
			err := cmd.Run()
			if err != nil {
				fmt.Printf("Error running Mandelbrot program: %v\n", err)
				continue
			}
			executionTime := time.Since(start).Milliseconds()

			// Write the result to the CSV file
			err = writer.Write([]string{
				strconv.Itoa(iterations),
				strconv.Itoa(goroutines),
				strconv.FormatInt(executionTime, 10),
			})
			if err != nil {
				fmt.Printf("Error writing result to CSV: %v\n", err)
				return
			}

			// Print progress to the console
			fmt.Printf("Completed: Iterations=%d, Goroutines=%d, ExecutionTime=%dms\n", iterations, goroutines, executionTime)
		}
	}

	fmt.Println("Benchmarking completed. Results saved to", OutputFile)
}
