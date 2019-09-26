package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/berkmancenter/ridge"
	"github.com/gonum/matrix/mat64"
)

func main() {

	f, err := os.Open("training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.FieldsPerRecord = 4

	rawCSVData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// featureData tiene los valores que usaremos como variables independientes
	// se maneja como un vector pero luego se convertirá a una matriz
	featureData := make([]float64, 4*len(rawCSVData))
	// yData son los valores de nuestra variable dependiente
	yData := make([]float64, len(rawCSVData))

	var featureIndex int
	var yIndex int

	for idx, record := range rawCSVData {

		if idx == 0 {
			continue
		}

		for i, val := range record {

			valParsed, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Fatal("Could not parse float value")
			}

			if i < 3 {

				// Ponemos un 1 para que el algoritmos de ridge calcule un b en la fórmula final
				if i == 0 {
					featureData[featureIndex] = 1
					featureIndex++
				}

				featureData[featureIndex] = valParsed
				featureIndex++
			}

			if i == 3 {

				yData[yIndex] = valParsed
				yIndex++
			}

		}
	}

	features := mat64.NewDense(len(rawCSVData), 4, featureData)
	y := mat64.NewVector(len(rawCSVData), yData)

	// 1.0 es una penalización aplicada a los parámetros
	// Minimiza el overfitting, si es 0 estamos en cuadrados mínimos
	// Si es muy grande puede que nuestro modelo no se adapte a los datos
	r := ridge.New(features, y, 1.0)

	r.Regress()

	c1 := r.Coefficients.At(0, 0)
	c2 := r.Coefficients.At(1, 0)
	c3 := r.Coefficients.At(2, 0)
	c4 := r.Coefficients.At(3, 0)
	fmt.Printf("\nRegression formula:\n")
	fmt.Printf("y = %0.3f + %0.3f TV + %0.3f Radio + %0.3f Newspaper\n\n", c1, c2, c3, c4)

	f, err = os.Open("test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader = csv.NewReader(f)

	reader.FieldsPerRecord = 4
	testData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var mAE float64
	for i, record := range testData {

		if i == 0 {
			continue
		}

		yObserved, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Fatal(err)
		}

		tvVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		radioVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		newspaperVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		yPredicted := predict(tvVal, radioVal, newspaperVal, c1, c2, c3, c4)

		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
	}

	// Output the MAE to standard out.
	fmt.Printf("\nMAE = %0.2f\n\n", mAE)
}

func predict(tv, radio, newspaper float64, c1, c2, c3, c4 float64) float64 {
	return c1 + c2*tv + c3*radio + c4*newspaper
}
