package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/mash/gokmeans"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {

	// Descargar los csv de https://www.kaggle.com/shuyangli94/food-com-recipes-and-user-interactions/
	var recetaPasos map[int]int = make(map[int]int)
	var recetaCalificaciones map[int][]int = make(map[int][]int)

	f, err := os.Open("RAW_recipes.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = 12
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Recolectando cantidad de pasos de las recetas")
	for i, record := range trainingData {

		if i == 0 {
			continue
		}
		idReceta, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatal(err)
		}
		cantidadDePasos, err := strconv.Atoi(record[7])
		if err != nil {
			log.Fatal(err)
		}
		recetaPasos[idReceta] = cantidadDePasos
	}

	f, err = os.Open("RAW_interactions.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader = csv.NewReader(f)

	reader.FieldsPerRecord = 5
	trainingData, err = reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nRecolectando calificaciones de las recetas")
	for i, record := range trainingData {

		if i == 0 {
			continue
		}

		idReceta, err := strconv.Atoi(record[1])
		if err != nil {
			log.Fatal(err)
		}
		calificacion, err := strconv.Atoi(record[3])

		if err != nil {
			log.Fatal(err)
		}
		calificaciones := recetaCalificaciones[idReceta]
		recetaCalificaciones[idReceta] = append(calificaciones, calificacion)
	}

	var data []gokmeans.Node

	for idReceta, pasos := range recetaPasos {

		calificacionesReceta, existe := recetaCalificaciones[idReceta]

		if !existe {
			calificacionesReceta = make([]int, 0)
		}

		for _, calificacion := range calificacionesReceta {

			data = append(data, gokmeans.Node{float64(pasos), float64(calificacion)})
		}
	}

	fmt.Printf("\nEntrenando k-means...\n")
	success, centroids := gokmeans.Train(data, 5, 50)
	if !success {
		log.Fatal("No pude generar los clusters")
	}

	fmt.Println("Los centroides son:")
	for _, centroid := range centroids {
		fmt.Println(centroid)
	}

	elementosClusterUno := make(plotter.XYs, 0)
	elementosClusterDos := make(plotter.XYs, 0)
	elementosClusterTres := make(plotter.XYs, 0)
	elementosClusterCuatro := make(plotter.XYs, 0)
	elementosClusterCinco := make(plotter.XYs, 0)

	for _, punto := range data {

		distanciaUno := floats.Distance([]float64{punto[0], punto[1]}, []float64{centroids[0][0], centroids[0][1]}, 2)
		distanciaDos := floats.Distance([]float64{punto[0], punto[1]}, []float64{centroids[1][0], centroids[1][1]}, 2)
		distanciaTres := floats.Distance([]float64{punto[0], punto[1]}, []float64{centroids[2][0], centroids[2][1]}, 2)
		distanciaCuatro := floats.Distance([]float64{punto[0], punto[1]}, []float64{centroids[3][0], centroids[3][1]}, 2)
		distanciaCinco := floats.Distance([]float64{punto[0], punto[1]}, []float64{centroids[4][0], centroids[4][1]}, 2)

		if distanciaUno < distanciaDos && distanciaUno < distanciaTres && distanciaUno < distanciaCuatro && distanciaUno < distanciaCinco {
			elementosClusterUno = append(elementosClusterUno, plotter.XY{X: punto[0], Y: punto[1]})
			continue
		}
		if distanciaDos < distanciaUno && distanciaDos < distanciaTres && distanciaDos < distanciaCuatro && distanciaDos < distanciaCinco {
			elementosClusterDos = append(elementosClusterDos, plotter.XY{X: punto[0], Y: punto[1]})
			continue
		}
		if distanciaTres < distanciaUno && distanciaTres < distanciaDos && distanciaTres < distanciaCuatro && distanciaTres < distanciaCinco {
			elementosClusterTres = append(elementosClusterTres, plotter.XY{X: punto[0], Y: punto[1]})
			continue
		}
		if distanciaCuatro < distanciaUno && distanciaCuatro < distanciaDos && distanciaCuatro < distanciaTres && distanciaCuatro < distanciaCinco {
			elementosClusterCuatro = append(elementosClusterCuatro, plotter.XY{X: punto[0], Y: punto[1]})
			continue
		}
		elementosClusterCinco = append(elementosClusterCinco, plotter.XY{X: punto[0], Y: punto[1]})
	}

	fmt.Println("Generando gráfico...")
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.X.Label.Text = "Pasos"
	p.Y.Label.Text = "Calificación"
	p.Add(plotter.NewGrid())

	sUno, err := plotter.NewScatter(elementosClusterUno)
	if err != nil {
		log.Fatal(err)
	}
	sUno.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	sUno.GlyphStyle.Radius = vg.Points(3)

	sDos, err := plotter.NewScatter(elementosClusterDos)
	if err != nil {
		log.Fatal(err)
	}
	sDos.GlyphStyle.Color = color.RGBA{G: 255, A: 255}
	sDos.GlyphStyle.Radius = vg.Points(3)

	sTres, err := plotter.NewScatter(elementosClusterTres)
	if err != nil {
		log.Fatal(err)
	}
	sTres.GlyphStyle.Color = color.RGBA{B: 255, A: 255}
	sTres.GlyphStyle.Radius = vg.Points(3)

	sCuatro, err := plotter.NewScatter(elementosClusterCuatro)
	if err != nil {
		log.Fatal(err)
	}
	sCuatro.GlyphStyle.Color = color.RGBA{R: 125, B: 125, A: 255}
	sCuatro.GlyphStyle.Radius = vg.Points(3)

	sCinco, err := plotter.NewScatter(elementosClusterCinco)
	if err != nil {
		log.Fatal(err)
	}
	sCinco.GlyphStyle.Color = color.RGBA{G: 125, B: 125, A: 255}
	sCinco.GlyphStyle.Radius = vg.Points(3)

	// Save the plot to a PNG file.
	p.Add(sUno, sDos, sTres, sCuatro, sCinco)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "pasos_calificacion_clusters.png"); err != nil {
		log.Fatal(err)
	}
}
