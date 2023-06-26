package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"regexp"
	"sort"
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
	var bundeslaender []Bundesland
	var staedte []City
	var areaGermany float64
	svgString := string(svgContent)
	bundeslaenderPaths := extractBundesländerPaths(svgString)
	for id, path := range bundeslaenderPaths {
		var bundesland = extractBundesland(path, id)
		bundeslaender = append(bundeslaender, bundesland)
		areaGermany = areaGermany + bundesland.area
	}
	sort.Slice(bundeslaender, func(i, j int) bool {
		return bundeslaender[i].area > bundeslaender[j].area
	})
	for j := 0; j < len(bundeslaender); j++ {
		fmt.Printf("Fläche für %s: %f\n", bundeslaender[j].id, bundeslaender[j].area)
		fmt.Println("Prozentualer Anteil von  ", bundeslaender[j].id, " ist ", bundeslaender[j].area/areaGermany*100)
	}
	fmt.Println("The area of Germany is: ", areaGermany)
	staedtePaths := extractStaedtePaths(svgString)
	for id, path := range staedtePaths {
		var stadt = extractStadt(path, id)
		staedte = append(staedte, stadt)
	}
	for i := 0; i < len(staedte); i++ {
		for j := 0; j < len(bundeslaender); j++ {
			var test = isCoordinateInBundesland(staedte[i].coordinate, bundeslaender[j].coordinates)
			if test == true {
				fmt.Println("Die Stadt", staedte[i].id, "ist in", bundeslaender[j].id)
			}
		}
	}
}

func isCoordinateInBundesland(coordinates Point, polygon [][]Point) bool {
	numPolygons := len(polygon)
	if numPolygons == 0 {
		return false
	}

	intersections := 0
	for _, polygon := range polygon {
		numVertices := len(polygon)
		if numVertices < 3 {
			continue
		}

		for i := 0; i < numVertices; i++ {
			currentVertex := polygon[i]
			nextVertex := polygon[(i+1)%numVertices]

			if (currentVertex.Y > coordinates.Y) != (nextVertex.Y > coordinates.Y) &&
				coordinates.X < (nextVertex.X-currentVertex.X)*(coordinates.Y-currentVertex.Y)/(nextVertex.Y-currentVertex.Y)+currentVertex.X {
				intersections++
			}
		}
	}

	return intersections%2 != 0
}

func extractStaedtePaths(svgString string) map[string]string {
	regexPattern := `<path[^>]+id\s*=\s*["']([^"']+)["'][^>]*\s+sodipodi:cx\s*=\s*["']([^"']+)["'][^>]*\s+sodipodi:cy\s*=\s*["']([^"']+)["'][^>]*>`
	regex := regexp.MustCompile(regexPattern)

	matches := regex.FindAllStringSubmatch(svgString, -1)
	pathMap := make(map[string]string)

	for _, match := range matches {
		if len(match) >= 4 {
			path := match[1]
			id := match[0]
			pathMap[path] = id
		}
	}

	return pathMap
}

func extractStadt(path string, id string) City {
	var stadt City
	regexPattern := `sodipodi:cx\s*=\s*"([^"]+)"\s+sodipodi:cy\s*=\s*"([^"]+)"`
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return stadt
	}

	match := regex.FindStringSubmatch(path)
	if len(match) < 3 {
		return stadt
	}

	stadt.id = id
	stadt.coordinate.X, err = strconv.ParseFloat(match[1], 64)
	stadt.coordinate.Y, err = strconv.ParseFloat(match[2], 64)

	return stadt
}

func extractBundesland(path string, id string) Bundesland {
	var coordinatePoints [][]Point
	coordinates := extractCoordinates(path)
	var lastX, lastY, startX, startY float64
	var prefix string
	var polygon []Point

	for _, coord := range coordinates {
		if len(coord) > 1 {
			x, err := strconv.ParseFloat(coord[1], 64)
			if err != nil {
				// Handle error
			}

			y, err := strconv.ParseFloat(coord[2], 64)
			if err != nil {
				// Handle error
			}

			if lastX == 0 && lastY == 0 {
				lastX, lastY = x, y
			}

			if startX == 0 && startY == 0 {
				startX, startY = x, y
			}

			prefix = coord[0]
			if prefix == "z" {
				coordinatePoints = append(coordinatePoints, polygon)
				polygon = []Point{}
				lastX, lastY = 0, 0
				startX, startY = 0, 0
				continue
			}
			absX, absY := convertToAbsolute(lastX, lastY, x, y, prefix)
			polygon = append(polygon, Point{absX, absY})

			lastX, lastY = absX, absY

		}
	}

	var bundesland Bundesland
	bundesland.id = id
	bundesland.area = 0.0

	// check if island or hole and then decide if + area oder - area
	for _, polygonPoints := range coordinatePoints {
		area := calculatePolygonAreaNew(polygonPoints)
		if checkIfIsland(polygonPoints, coordinatePoints) {
			bundesland.area += area
		} else {
			bundesland.area -= area
		}
	}

	bundesland.coordinates = coordinatePoints

	return bundesland
}

func checkIfIsland(polygon []Point, allPolygons [][]Point) bool {
	var randomPoint = getRandomPointInsidePolygon(polygon)
	for _, polygonPoints := range allPolygons {
		var coordinates [][]Point
		if areArraysEqual(polygon, polygonPoints) {
			continue
		}
		coordinates = append(coordinates, polygonPoints)
		if isCoordinateInBundesland(randomPoint, coordinates) {
			if calculatePolygonAreaNew(polygon) < calculatePolygonAreaNew(polygonPoints) {
				return false
			}
		}
	}
	return true
}

func areArraysEqual(arr1, arr2 []Point) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for i := 0; i < len(arr1); i++ {
		if arr1[i] != arr2[i] {
			return false
		}
	}

	return true
}
func getRandomPointInsidePolygon(polygon []Point) Point {
	var minX, maxX, minY, maxY float64
	var coordinates [][]Point
	for _, point := range polygon {
		if point.X < minX || minX == 0 {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY || minY == 0 {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	for {
		randomX := minX + rand.Float64()*(maxX-minX)
		randomY := minY + rand.Float64()*(maxY-minY)
		randomPoint := Point{randomX, randomY}
		coordinates = append(coordinates, polygon)

		if isCoordinateInBundesland(randomPoint, coordinates) {
			return randomPoint
		}
	}
}

func extractBundesländerPaths(svgString string) map[string]string {
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

func calculatePolygonAreaNew(coordinates []Point) float64 {
	n := len(coordinates)
	area := 0.0

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += (coordinates[i].X + coordinates[j].X) * (coordinates[i].Y - coordinates[j].Y)
	}

	return math.Abs(area / 2.0)
}
func extractCoordinates(svgPath string) [][]string {
	regexPattern := `([a-zA-Z])([-+]?\d*\.?\d+(?:[eE][-+]?\d+)?(?:,[-+]?\d*\.?\d+(?:[eE][-+]?\d+)?)?)?`
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil
	}

	matches := regex.FindAllStringSubmatch(svgPath, -1)
	if matches == nil {
		return nil
	}

	coordinates := make([][]string, len(matches))
	for i, match := range matches {
		prefix := match[1]
		coords := match[2]

		// Handle case when coordinates are empty
		if coords == "" {
			if prefix == "z" {
				coords = "0,0" // Set default coordinates for "z" command
			} else {
				continue // Skip if no coordinates are available
			}
		}

		coordValues := strings.Split(coords, ",")
		coord1 := coordValues[0]
		coord2 := ""
		if len(coordValues) > 1 {
			coord2 = coordValues[1]
		}

		coordinates[i] = []string{prefix, coord1, coord2}
	}

	return coordinates
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

type Point struct {
	X, Y float64
}

type Bundesland struct {
	id          string
	coordinates [][]Point
	area        float64
}

type City struct {
	id         string
	coordinate Point
}
