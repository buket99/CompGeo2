package main

import (
	"fmt"
	"io/ioutil"
)

/*
Aufgabe:
	- Lese SVG-Datei
	- Ermittle die Flächen der einzelnen Bundesländer
	- Bzgl. der Koordinaten von Städten, finde heraus in welchem Bundesland diese liegen
*/

func main() {
	// Read the SVG file
	data, err := ioutil.ReadFile("DeutschlandMitStaedten.svg")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Convert the byte slice to a string
	svgContent := string(data)

	// Process the SVG content as needed
	fmt.Println(svgContent)
}

func flaecheninhalt() float64 {
	return 0
}

func findeBundesland() float64 {
	return 0
}
