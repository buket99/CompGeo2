package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

// 1. Lesen Sie die SVG-Datei 'DeutschlandMitStaedten.svg'
// 2. Ermitteln Sie die Flächen der einzelnen Bundesländer (bezüglich der in der Datei verwendeten Skala)
// 3. Am Ende der Datei befinden sich Koordinaten von Städten,
//    Versuchen Sie herauszufinden (bzw. lassen Sie das Ihren Rechner machen ;-), in welchem Bundesland diese
//    jeweils liegen.

type SVG struct {
	XMLName xml.Name `xml:"svg"`
	Width   string   `xml:"width,attr"`
	Height  string   `xml:"height,attr"`
}

func main() {
	file, err := os.Open("DeutschlandMitStaedten.svg")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var svg SVG
	err = xml.NewDecoder(file).Decode(&svg)
	if err != nil {
		fmt.Println("Error decoding XML:", err)
		return
	}

	fmt.Printf("SVG width: %s, height: %s", svg.Width, svg.Height)

}
