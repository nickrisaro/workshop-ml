package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
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
	default:
		fmt.Println("No sé hacer lo que me pedís")
	}

}

func describirDatos(advertFile *os.File) {

	advertDF := dataframe.ReadCSV(advertFile)
	advertSummary := advertDF.Describe()

	fmt.Println(advertSummary)
}
