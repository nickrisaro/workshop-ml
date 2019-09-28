package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jbrukh/bayesian"
)

const (
	Positivo bayesian.Class = "1"
	Negativo bayesian.Class = "0"
)

func main() {

	fmt.Println("Entrenando con un clasificador sin TF-IDF")
	classifier := bayesian.NewClassifier(Positivo, Negativo)

	entrenar(classifier)
	evaluarFrase(classifier, "This movie is awesome")
	evaluarFrase(classifier, "It's an awful movie")

	fmt.Println("Entrenando con un clasificador con TF-IDF")
	classifier = bayesian.NewClassifierTfIdf(Positivo, Negativo)
	entrenar(classifier)
	classifier.ConvertTermsFreqToTfIdf()

	evaluarFrase(classifier, "This movie is awesome")
	evaluarFrase(classifier, "It's an awful movie")
}

func entrenar(classifier *bayesian.Classifier) {

	f, err := os.Open("imdb_labelled.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.FieldsPerRecord = 2
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, record := range trainingData {

		palabras := record[0]
		categoria := record[1]

		if categoria == "1" {
			classifier.Learn(strings.Split(strings.ToLower(palabras), " "), Positivo)
		} else {
			classifier.Learn(strings.Split(strings.ToLower(palabras), " "), Negativo)
		}
	}
}

func evaluarFrase(classifier *bayesian.Classifier, frase string) {

	scores, _, _ := classifier.ProbScores(strings.Split(strings.ToLower(frase), " "))

	fmt.Println("Probabilidad de que la frase sea positiva", scores[0]*100)
	fmt.Println("Probabilidad de que la frase sea negativa", scores[1]*100)

}
