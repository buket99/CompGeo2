package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func extractCoordinates(svgPath string) ([]float64, error) {
	// Define the regular expression pattern to match the "d" parameter
	regexPattern := `[MmLlHhVvZz]+((?:-?\d+(?:\.\d+)?(?:e-?\d+)?[, ]*)+)`

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
	var allCoordinates []float64
	for _, match := range matches {
		coordinates := match[1]

		// Split the coordinates by spaces, commas, or other separators
		coordinatesArr := strings.FieldsFunc(coordinates, func(r rune) bool {
			return r == ',' || r == ' '
		})

		// Process the coordinates to obtain absolute values
		var absoluteCoordinates []float64
		var lastX, lastY float64
		for i := 0; i < len(coordinatesArr); i += 2 {
			cmd := coordinatesArr[i]

			if len(coordinatesArr)-i < 2 {
				break // Handle edge case of insufficient coordinates
			}

			if strings.ToUpper(cmd) == "Z" {
				continue // Skip Z command
			}

			x, _ := strconv.ParseFloat(coordinatesArr[i+0], 64)
			y, _ := strconv.ParseFloat(coordinatesArr[i+1], 64)

			// Handle SVG commands that modify the last position
			switch strings.ToUpper(cmd) {
			case "M", "L":
				// Move or Line commands
				x += lastX
				y += lastY
				absoluteCoordinates = append(absoluteCoordinates, x, y)
				lastX, lastY = x, y
			case "H":
				// Horizontal Line command
				x += lastX
				absoluteCoordinates = append(absoluteCoordinates, x, lastY)
				lastX = x
			case "V":
				// Vertical Line command
				y += lastY
				absoluteCoordinates = append(absoluteCoordinates, lastX, y)
				lastY = y
			}
		}

		allCoordinates = append(allCoordinates, absoluteCoordinates...)
	}

	return allCoordinates, nil
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
