package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

func extractCoordinates(svgPath string) ([]float64, error) {
	// Define the regular expression pattern to match the "d" parameter
	regexPattern := `([MmZzLlHhVvCcSsQqTtAa])((?:-?\d+(?:\.\d+)?(?:e-?\d+)?[, ]*)+)`

	// Compile the regular expression
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}

	// Find all matches of the pattern in the SVG path
	matches := regex.FindAllStringSubmatch(svgPath, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("failed to find 'd' parameter in SVG path")
	}

	// Extract the coordinates from the matched substrings
	var absoluteCoordinates []float64
	var lastX, lastY float64
	for _, match := range matches {
		command := match[1]
		coordinates := match[2]

		// Split the coordinates by spaces, commas, or other separators
		coordinatesArr := strings.FieldsFunc(coordinates, func(r rune) bool {
			return r == ',' || r == ' '
		})

		// Process the coordinates to obtain absolute values
		for i := 0; i < len(coordinatesArr); i += 2 {
			x, _ := strconv.ParseFloat(coordinatesArr[i], 64)
			y, _ := strconv.ParseFloat(coordinatesArr[i+1], 64)

			// Handle absolute and relative coordinates
			if unicode.IsUpper(rune(command[0])) {
				// Absolute coordinate
				lastX, lastY = x, y
			} else {
				// Relative coordinate
				x += lastX
				y += lastY
				lastX, lastY = x, y
			}

			absoluteCoordinates = append(absoluteCoordinates, x, y)
		}
	}

	return absoluteCoordinates, nil
}

func calculatePolygonArea(coordinates []float64) float64 {
	var area float64
	numCoordinates := len(coordinates)

	// Calculate the area using the shoelace formula
	for i := 0; i < numCoordinates; i += 2 {
		x1 := coordinates[i]
		y1 := coordinates[i+1]
		x2 := coordinates[(i+2)%numCoordinates]
		y2 := coordinates[(i+3)%numCoordinates]

		area += (x1 * y2) - (x2 * y1)
	}

	area = math.Abs(area / 2)

	return area
}

func main() {
	// Read SVG file
	filePath := "DeutschlandMitStaedten.svg"
	svgData, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	svgContent := string(svgData)

	// Define the regular expression pattern to match SVG paths
	regexPathPattern := `<path[^>]*d="([^"]+)"[^>]*>`

	// Compile the regular expression for SVG paths
	regexPath, err := regexp.Compile(regexPathPattern)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Find all matches of SVG paths in the SVG content
	pathMatches := regexPath.FindAllStringSubmatch(svgContent, -1)

	// Iterate over each path and calculate the area
	for i, match := range pathMatches {
		svgPath := match[1]

		coordinates, err := extractCoordinates(svgPath)
		if err != nil {
			fmt.Printf("Error extracting coordinates for path %d: %s\n", i+1, err)
			continue
		}

		area := calculatePolygonArea(coordinates)
		fmt.Printf("Area of path %d: %.2f\n", i+1, area)
	}
}
