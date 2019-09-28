package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/ensemble"
	"github.com/sjwhitworth/golearn/evaluation"
)

func main() {

	irisData, err := base.ParseCSVToInstances("iris.csv", true)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// Creamos un bosque con 10 árboles y 2 atributos por árbol
	// En general se usa sqrt(atributos) como default para la
	// cantidad de atributos por árbol
	rf := ensemble.NewRandomForest(10, 2)

	// Entrena y evalúa el modelo en 5 subconjuntos del dataset
	cv, err := evaluation.GenerateCrossFoldValidationConfusionMatrices(irisData, rf, 5)
	if err != nil {
		log.Fatal(err)
	}

	// Calculamos la precisión en base a los datos de las matrices de confusión
	mean, variance := evaluation.GetCrossValidatedMetric(cv, evaluation.GetAccuracy)
	stdev := math.Sqrt(variance)

	fmt.Printf("\nPrecisión\n%.2f (+/- %.2f)\n\n", mean, stdev*2)
}
