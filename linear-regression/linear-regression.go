package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
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

	switch os.Args[1] {

	case "--desc":
		describirDatos(advertFile)
	case "--hist":
		generarHistograma(advertFile)
	default:
		fmt.Println("No sé hacer lo que me pedís")
	}

}

func describirDatos(advertFile *os.File) {

	advertDF := dataframe.ReadCSV(advertFile)
	advertSummary := advertDF.Describe()

	fmt.Println(advertSummary)
}

func generarHistograma(advertFile *os.File) {
	advertDF := dataframe.ReadCSV(advertFile)

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

		if err := p.Save(4*vg.Inch, 4*vg.Inch, colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}
	}
}
