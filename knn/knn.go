package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/cluster"
)

var conversorEspecie = map[string]float64{
	"Iris-setosa":     1.0,
	"Iris-versicolor": 2.0,
	"Iris-virginica":  3.0,
}

func main() {

	f, err := os.Open("training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = 5
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	atributos := [][]float64{}
	clasificacion := []float64{}

	for i, record := range trainingData {

		// Salteamos el header
		if i == 0 {
			continue
		}

		// sepal_length,sepal_width,petal_length,petal_width,species
		sepalLength, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		sepalWidth, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		petalLength, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		petalWidth, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		species := record[4]

		atributos = append(atributos, []float64{sepalLength, sepalWidth, petalLength, petalWidth})
		clasificacion = append(clasificacion, conversorEspecie[species])

	}

	model := cluster.NewKNN(3, atributos, clasificacion, base.EuclideanDistance)

	var correctos float64
	var incorrectos float64
	for i, record := range trainingData {

		// Salteamos el header
		if i == 0 {
			continue
		}

		// sepal_length,sepal_width,petal_length,petal_width,species
		sepalLength, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		sepalWidth, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		petalLength, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		petalWidth, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		species := record[4]

		clasificacion := conversorEspecie[species]

		guess, err := model.Predict([]float64{sepalLength, sepalWidth, petalLength, petalWidth})
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if guess[0] == clasificacion {
			correctos++
		} else {
			incorrectos++
		}
	}

	fmt.Printf("Correctos: %v, Incorrectos: %v\n", correctos, incorrectos)
	precision := (correctos / (correctos + incorrectos))
	fmt.Printf("\nPrecisión = %0.2f\n\n", precision)
}
