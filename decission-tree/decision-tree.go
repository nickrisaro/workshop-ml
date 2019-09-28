package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/trees"
)

func main() {

	irisData, err := base.ParseCSVToInstances("training.csv", true)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// 0.6 es el prune ratio -> ver que pasa con distintos valores
	tree := trees.NewID3DecisionTree(0.6)
	// tree := trees.NewRandomTree(2)

	tree.Fit(irisData)

	fmt.Println(tree)

	testData, err := base.ParseCSVToInstances("test.csv", true)
	if err != nil {
		log.Fatal(err)
	}

	predictions, err := tree.Predict(testData)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a Confusion Matrix.
	cm, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		log.Fatal(err)
	}

	// Retrieve the accuracy.
	accuracy := evaluation.GetAccuracy(cm)
	fmt.Printf("\nAccuracy: %0.2f\n\n", accuracy)

	_, err = os.Create("predicciones.csv")
	if err != nil {
		log.Fatal(err)
	}

	err = base.SerializeInstancesToCSV(predictions, "predicciones.csv")
	if err != nil {
		log.Fatal(err)
	}
}
