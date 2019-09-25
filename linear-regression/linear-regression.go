package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/kniren/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Decime lo que querés hacer")
		return
	}

	advertFile, err := os.Open("Advertising.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer advertFile.Close()

	advertDF := dataframe.ReadCSV(advertFile)

	switch os.Args[1] {

	case "--desc":
		describirDatos(advertDF)
	case "--hist":
		generarHistograma(advertDF)
	case "--disp":
		generarDiagramaDeDispersion(advertDF)
	case "--sets":
		generarDatasets(advertDF)
	case "--train":
		entrenar()
	default:
		fmt.Println("No sé hacer lo que me pedís")
	}

}

func describirDatos(advertDF dataframe.DataFrame) {

	advertSummary := advertDF.Describe()

	fmt.Println(advertSummary)
}

func generarHistograma(advertDF dataframe.DataFrame) {

	// Crea un histograma para cada columna del archivo
	for _, colName := range advertDF.Names() {

		plotVals := make(plotter.Values, advertDF.Nrow())
		for i, floatVal := range advertDF.Col(colName).Float() {
			plotVals[i] = floatVal
		}

		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.Title.Text = fmt.Sprintf("Histograma de %s", colName)

		h, err := plotter.NewHist(plotVals, 16)
		if err != nil {
			log.Fatal(err)
		}

		h.Normalize(1)

		p.Add(h)

		guardarGrafico(p, colName+"_hist.png")

	}
}

func generarDiagramaDeDispersion(advertDF dataframe.DataFrame) {

	yVals := advertDF.Col("Sales").Float()

	// Creamos un diagrama de dispersión entre las ventas y cada una de las otras columnas
	for _, colName := range advertDF.Names() {

		pts := make(plotter.XYs, advertDF.Nrow())

		for i, floatVal := range advertDF.Col(colName).Float() {
			pts[i].X = floatVal
			pts[i].Y = yVals[i]
		}

		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.X.Label.Text = colName
		p.Y.Label.Text = "Sales"
		p.Add(plotter.NewGrid())

		s, err := plotter.NewScatter(pts)
		if err != nil {
			log.Fatal(err)
		}
		s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 255}
		s.GlyphStyle.Radius = vg.Points(3)

		p.Add(s)

		guardarGrafico(p, colName+"_scatter.png")

	}
}

func guardarGrafico(p *plot.Plot, nombre string) {

	if err := p.Save(4*vg.Inch, 4*vg.Inch, nombre); err != nil {
		log.Fatal(err)
	}
}

func generarDatasets(advertDF dataframe.DataFrame) {

	trainingNum := (4 * advertDF.Nrow()) / 5
	testNum := advertDF.Nrow() / 5
	if trainingNum+testNum < advertDF.Nrow() {
		trainingNum++
	}

	trainingIdx := make([]int, trainingNum)
	testIdx := make([]int, testNum)

	for i := 0; i < trainingNum; i++ {
		trainingIdx[i] = i
	}

	for i := 0; i < testNum; i++ {
		testIdx[i] = trainingNum + i
	}

	trainingDF := advertDF.Subset(trainingIdx)
	testDF := advertDF.Subset(testIdx)

	guardarDataSet(trainingDF, "training.csv")
	guardarDataSet(testDF, "test.csv")
}

func guardarDataSet(dataSet dataframe.DataFrame, nombreArchivo string) {

	f, err := os.Create(nombreArchivo)
	if err != nil {
		log.Fatal(err)
	}

	w := bufio.NewWriter(f)

	if err := dataSet.WriteCSV(w); err != nil {
		log.Fatal(err)
	}
}

func entrenar() {

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

	// Vamos a generar un modelo con las ventas como variable dependiente
	// y lo invertido en TV como variable dependiente
	var r regression.Regression
	r.SetObserved("Sales")
	r.SetVar(0, "TV")

	for i, record := range trainingData {

		// Salteamos el header
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

		r.Train(regression.DataPoint(yVal, []float64{tvVal}))
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

		// Predecimos y con nuestro modelo
		yPredicted, err := r.Predict([]float64{tvVal})

		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
	}

	fmt.Printf("MAE = %0.2f\n\n", mAE)
}
