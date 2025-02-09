package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

// Représente une ligne du fichier csv avec les 4 colonnes
type BenchmarkResult struct {
	Iterations      int
	Goroutines      int
	ExecutionTimeMs float64
	Quality         float64
}

// Load un fichier csv
func LoadCSV(filename string) ([]BenchmarkResult, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("file %s is empty or missing data", filename)
	}

	var results []BenchmarkResult

	// Parcourt chaque ligne (sauf l'en-tête)
	for _, record := range records[1:] {
		if len(record) < 4 {
			log.Printf("Skipping malformed row in file %s: %v", filename, record)
			continue
		}

		iterations, _ := strconv.Atoi(record[0])
		goroutines, _ := strconv.Atoi(record[1])
		execTime, _ := strconv.ParseFloat(record[2], 64)
		quality, _ := strconv.ParseFloat(record[3], 64)

		// Ignore les lignes avec 1 ou 21 itérations car sinon qualité est nulle comme mesure
		if iterations == 1 || iterations == 21 {
			continue
		}

		results = append(results, BenchmarkResult{
			Iterations:      iterations,
			Goroutines:      goroutines,
			ExecutionTimeMs: execTime,
			Quality:         quality,
		})
	}
	return results, nil
}

// Normalise les valeurs d'un tableau dans l'intervalle [0, 1]
func Normalize(values []float64) []float64 {

	minVal := math.MaxFloat64
	maxVal := -math.MaxFloat64

	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	normalized := make([]float64, len(values))
	for i, v := range values {
		if maxVal != minVal {
			normalized[i] = (v - minVal) / (maxVal - minVal)
		} else {
			normalized[i] = 0
		}
	}
	return normalized
}

func main() {
	files := []string{"new_benchmark_results.csv", "new_new_benchmark_results.csv", "new_new_new_benchmark_results.csv"}
	var allResults []BenchmarkResult

	for _, file := range files {
		results, err := LoadCSV(file)
		if err != nil {
			log.Printf("Error loading CSV file %s: %v", file, err)
			continue
		}
		allResults = append(allResults, results...)
	}

	if len(allResults) == 0 {
		log.Fatal("No valid data found in any CSV file.")
	}

	executionTimes := make([]float64, len(allResults))
	qualities := make([]float64, len(allResults))

	for i, result := range allResults {
		executionTimes[i] = result.ExecutionTimeMs
		qualities[i] = result.Quality
	}

	normalizedTimes := Normalize(executionTimes)
	normalizedQualities := Normalize(qualities)

	bestScore := math.MaxFloat64
	var bestResult BenchmarkResult

	for i, result := range allResults {
		score := normalizedTimes[i] + normalizedQualities[i] // Minimise les deux, ajouter des coefs en fonction de l'importance
		if score < bestScore {
			bestScore = score
			bestResult = result
		}
	}

	fmt.Printf("Best Configuration:\n")
	fmt.Printf("Iterations: %d, Goroutines: %d\n", bestResult.Iterations, bestResult.Goroutines)
	fmt.Printf("Execution Time (ms): %.2f, Quality: %.2f\n", bestResult.ExecutionTimeMs, bestResult.Quality)
}
