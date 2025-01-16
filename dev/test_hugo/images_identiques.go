package main

import (
	"fmt"
	"image"
	_ "image/jpeg" // pour le décodage JPEG
	_ "image/png"  // pour le décodage PNG
	"os"
)

func main() {
	// Ouvrir les deux fichiers image
	file1, err := os.Open("mandelbrot_30000x30000_2000iter_1576.79s.png")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de l'image 1:", err)
		return
	}
	defer file1.Close()

	file2, err := os.Open("mandelbrot_using_symetry_30000x30000_2000iter_532.24s.png")
	if err != nil {
		fmt.Println("Erreur lors de l'ouverture de l'image 2:", err)
		return
	}
	defer file2.Close()

	// Décoder les images
	img1, _, err := image.Decode(file1)
	if err != nil {
		fmt.Println("Erreur lors du décodage de l'image 1:", err)
		return
	}

	img2, _, err := image.Decode(file2)
	if err != nil {
		fmt.Println("Erreur lors du décodage de l'image 2:", err)
		return
	}

	// Vérifier les dimensions des images
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()

	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		fmt.Println("Les images ont des dimensions différentes.")
		return
	}

	// Comparer pixel par pixel
	for y := bounds1.Min.Y; y < bounds1.Max.Y; y++ {
		for x := bounds1.Min.X; x < bounds1.Max.X; x++ {
			// Obtenir les couleurs des pixels
			color1 := img1.At(x, y)
			color2 := img2.At(x, y)

			// Comparer les pixels
			if color1 != color2 {
				fmt.Printf("Les images diffèrent au pixel (%d, %d)\n", x, y)
				return
			}
		}
	}

	// Si aucune différence n'a été trouvée
	fmt.Println("Les images sont identiques.")
}
