package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/cmplx"
	"os"
	"time"
)

func mandelbrot(c complex128, nbIterations int) color.Color {
	/*
		Determine la couleur d'un point de l'ensemble de mandelbrot
		L'ensemble de mandelbrot est l'ensemble des points bornés c du plan complexe tq
			{ z_0 = 0
			{ z_(n+1) = (z_n)^2 + c

		Prend en entrée un point complexe c et le nombre d'iérations voulues n
	*/
	var z complex128
	for n := 0; n < nbIterations; n++ {
		//Si |z| > 2 on sort en mappant le gradient de couleur
		if cmplx.Abs(z) > 2 {
			return color.RGBA{R: uint8(255 - n%32*8), G: uint8(n % 64 * 4), B: uint8(255 - n%16*16), A: 255}
		}
		z = z*z + c
	}
	return color.Black
}

func main() {

	//On définie les constantes qui vont nous servir
	// Modifier les valeurs pour changer la résolution de l'image
	const (
		hauteur, largeur = 10000, 10000
		xMin, xMax       = -2.0, 1.0
		yMin, yMax       = -1.5, 1.5
		maxIterations    = 2000
	)

	//On commence le timmer pour mesurer le temps d'execution
	start := time.Now()

	//On crée l'image vierge avec la résolution voulue
	img := image.NewRGBA((image.Rect(0, 0, hauteur, largeur)))

	//On applique mandelbrot sur chaque pixel de l'image
	for py := 0; py < hauteur; py++ {
		y := float64(py)/hauteur*(yMax-yMin) + yMin
		for px := 0; px < largeur; px++ {
			x := float64(px)/largeur*(xMax-xMin) + xMin
			c := complex(x, y)
			img.Set(px, py, mandelbrot(c))
		}
	}

	//On mesure le temps qui s'est écoulé et on l'affiche
	elapsed := time.Since(start)
	elapsedSeconds := elapsed.Seconds()
	fmt.Printf("Temps pour générer l'image %.2f secondes \n", elapsedSeconds)

	//Crée un fichier avec un nom adéquat pour stocker l'image
	nom := fmt.Sprintf("mandelbrot_%dx%d_%diter_%.2fs", hauteur, largeur, maxIterations, elapsedSeconds)
	file, err := os.Create(nom)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}

	// Affiche le nom du fichier, si le programme s'est bien exécuté
	fmt.Printf("Image de l'ensemble de Mandelbrot sauvgarder en : %s", nom)
}
