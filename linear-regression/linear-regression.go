package main

import (
	"fmt"
	"image/color"
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
	case "--disp":
		generarDiagramaDeDispersion(advertFile)
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

		guardarGrafico(p, colName+"_hist.png")

	}
}

func generarDiagramaDeDispersion(advertFile *os.File) {

	advertDF := dataframe.ReadCSV(advertFile)

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
