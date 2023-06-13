package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	svgFile := "DeutschlandMitStaedten.svg"
	svgContent, err := ioutil.ReadFile(svgFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	svgString := string(svgContent)

	paths := extractPaths(svgString)

	for id, path := range paths {
		area, err := calculatePathArea(path)
		if err != nil {
			fmt.Printf("Error calculating area for %s: %s\n", id, err)
			continue
		}

		fmt.Printf("Area for %s: %f\n", id, area)
	}
}

func extractPaths(svgString string) map[string]string {
	regexPattern := `<path[^>]+id\s*=\s*["']([^"']+)["'][^>]*\sd\s*=\s*["']([^"']+)["'][^>]*>`
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return nil
	}

	matches := regex.FindAllStringSubmatch(svgString, -1)
	if len(matches) == 0 {
		fmt.Println("No matches found")
		return nil
	}

	paths := make(map[string]string)
	for _, match := range matches {
		if len(match) > 2 {
			id := match[1]
			path := match[2]
			paths[id] = path
		}
	}

	return paths
}

func calculatePathArea(svgPath string) (float64, error) {
	coordinates := extractCoordinates(svgPath)
	if len(coordinates) < 6 {
		return 0, fmt.Errorf("failed to find valid coordinate pairs in SVG path")
	}

	var areas []float64
	var startX, startY, lastX, lastY float64
	var totalArea float64
	var prefix string

	for i := 0; i < len(coordinates)-1; i += 2 {
		x, err := strconv.ParseFloat(coordinates[i], 64)
		if err != nil {
			return 0, err
		}

		y, err := strconv.ParseFloat(coordinates[i+1], 64)
		if err != nil {
			return 0, err
		}

		if i == 0 {
			startX, startY = x, y
			prefix = getPrefix(coordinates[0])
		}

		if i >= 2 && (strings.HasPrefix(coordinates[i], "-") || strings.HasPrefix(coordinates[i+1], "-")) {
			absX, absY := convertToAbsolute(lastX, lastY, x, y, prefix)
			area := calculatePolygonArea(startX, startY, lastX, lastY, absX, absY)
			areas = append(areas, area)
		}

		lastX, lastY = x, y
	}

	for _, area := range areas {
		totalArea += area
	}

	return totalArea, nil
}

func extractCoordinates(svgPath string) []string {
	regexPattern := `[-+]?\d*\.?\d+(?:[eE][-+]?\d+)?`
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil
	}

	matches := regex.FindAllString(svgPath, -1)

	return matches
}

func getPrefix(coord string) string {
	regexPattern := `^[A-Za-z]+`
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return ""
	}

	match := regex.FindString(coord)
	if match != "" {
		return match
	}

	return ""
}

func convertToAbsolute(lastX, lastY, x, y float64, prefix string) (float64, float64) {
	switch prefix {
	case "M", "L":
		return x, y
	case "m", "l":
		return lastX + x, lastY + y
	case "H":
		return x, lastY
	case "h":
		return lastX + x, lastY
	case "V":
		return lastX, y
	case "v":
		return lastX, lastY + y
	}

	return x, y
}
func calculatePolygonArea(x1, y1, x2, y2, x3, y3 float64) float64 {
	area := 0.5 * ((x1*y2 + x2*y3 + x3*y1) - (x2*y1 + x3*y2 + x1*y3))
	return math.Abs(area)
}
