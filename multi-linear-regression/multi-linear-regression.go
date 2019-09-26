package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/sajari/regression"
)

func main() {

	f, err := os.Open("training.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = 4
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Vamos a modelar las ventas en base a la inversión de TV y radio
	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "TV")
	r.SetVar(1, "Radio")

	for i, record := range trainingData {

		// Salteamos el header.
		if i == 0 {
			continue
		}

		yVal, err := strconv.ParseFloat(record[3], 64)
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

		r.Train(regression.DataPoint(yVal, []float64{tvVal, radioVal}))
	}

	r.Run()

	fmt.Printf("\nRegression Formula:\n%v\n\n", r.Formula)

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

		// Salteamos el header.
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

		// Predecimos y con nuestro modelo
		yPredicted, err := r.Predict([]float64{tvVal, radioVal})

		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
	}

	fmt.Printf("MAE = %0.2f\n\n", mAE)
}
