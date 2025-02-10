package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
)

// Coefs pour l'importance de chaque parametre
const (
	TimeCoef    = 1
	QualityCoef = 1
)

// BenchmarkResult représente une ligne du fichier CSV
// avec les colonnes : Iterations, Goroutines, ExecutionTimeMs, Quality
type BenchmarkResult struct {
	Iterations      int     // Nombre d'itérations
	Goroutines      int     // Nombre de goroutines
	ExecutionTimeMs float64 // Temps d'exécution en millisecondes
	Quality         float64 // Qualité mesurée
}

// LoadCSV charge les résultats de benchmark depuis un fichier CSV
func LoadCSV(filename string) ([]BenchmarkResult, error) {
	// Ouvre le fichier CSV
	file, err := os.Open(filename)
	if err != nil {
		return nil, err // Retourne une erreur si le fichier ne peut pas être ouvert
	}
	defer file.Close() // Assure la fermeture du fichier à la fin

	// Initialise un lecteur CSV
	reader := csv.NewReader(file)
	reader.Comma = ';' // Spécifie le délimiteur CSV comme étant le point-virgule

	// Lis toutes les lignes du fichier CSV
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err // Retourne une erreur si la lecture échoue
	}

	// Vérifie que le fichier contient des données
	if len(records) < 2 {
		return nil, fmt.Errorf("file %s is empty or missing data", filename)
	}

	var results []BenchmarkResult // Tableau pour stocker les résultats

	// Parcourt chaque ligne (sauf l'en-tête)
	for _, record := range records[1:] {
		// Ignore les lignes mal formées
		if len(record) < 4 {
			log.Printf("Skipping malformed row in file %s: %v", filename, record)
			continue
		}
		// Convertit les champs du CSV en types appropriés
		iterations, _ := strconv.Atoi(record[0])
		goroutines, _ := strconv.Atoi(record[1])
		execTime, _ := strconv.ParseFloat(record[2], 64)
		quality, _ := strconv.ParseFloat(record[3], 64)

		// Ignore les lignes avec 1 ou 21 itérations
		if iterations == 1 || iterations == 21 {
			continue
		}

		// Ajoute le résultat à la liste
		results = append(results, BenchmarkResult{
			Iterations:      iterations,
			Goroutines:      goroutines,
			ExecutionTimeMs: execTime,
			Quality:         quality,
		})
	}
	return results, nil // Retourne les résultats
}

// Normalize normalise les valeurs d'un tableau dans l'intervalle [0, 1]
func Normalize(values []float64) []float64 {
	// Trouve les valeurs minimale et maximale
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

	// Initialise le tableau des valeurs normalisées
	normalized := make([]float64, len(values))
	for i, v := range values {
		// Évite la division par zéro si toutes les valeurs sont identiques
		if maxVal != minVal {
			normalized[i] = (v - minVal) / (maxVal - minVal)
		} else {
			normalized[i] = 0
		}
	}
	return normalized
}

func main() {
	// Liste des fichiers CSV à traiter
	files := []string{"new_benchmark_results.csv", "new_new_benchmark_results.csv", "new_new_new_benchmark_results.csv"}
	var allResults []BenchmarkResult // Liste pour stocker tous les résultats

	// Chargement des données depuis chaque fichier CSV
	for _, file := range files {
		results, err := LoadCSV(file)
		if err != nil {
			log.Printf("Error loading CSV file %s: %v", file, err)
			continue
		}
		// Ajout des résultats au tableau global
		allResults = append(allResults, results...)
	}

	// Vérifie que des données valides ont été chargées
	if len(allResults) == 0 {
		log.Fatal("No valid data found in any CSV file.")
	}

	// Prépare les tableaux pour les temps d'exécution et les qualités
	executionTimes := make([]float64, len(allResults))
	qualities := make([]float64, len(allResults))

	// Remplit les tableaux avec les valeurs extraites
	for i, result := range allResults {
		executionTimes[i] = result.ExecutionTimeMs
		qualities[i] = result.Quality
	}

	// Normalise les temps d'exécution et les qualités
	normalizedTimes := Normalize(executionTimes)
	normalizedQualities := Normalize(qualities)

	// Initialisation pour trouver la meilleure configuration
	bestScore := math.MaxFloat64
	var bestResult BenchmarkResult

	// Parcourt toutes les configurations pour trouver la meilleure
	for i, result := range allResults {
		// Calcule le score comme la somme des temps et qualités normalisés
		score := TimeCoef*normalizedTimes[i] + QualityCoef*normalizedQualities[i] // Minimise les deux avec des coefs
		if score < bestScore {
			bestScore = score
			bestResult = result
		}
	}

	// Affiche la meilleure configuration trouvée
	fmt.Printf("Best Configuration:\n")
	fmt.Printf("Iterations: %d, Goroutines: %d\n", bestResult.Iterations, bestResult.Goroutines)
	fmt.Printf("Execution Time (ms): %.2f, Quality: %.2f\n", bestResult.ExecutionTimeMs, bestResult.Quality)
}
