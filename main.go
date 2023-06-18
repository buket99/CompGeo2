package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"sort"
	"strconv"
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
	//	var berlinArea float64
	//	var bremenArea float64

	bundeslaenderPaths := extractBundesländerPaths(svgString)
	for id, path := range bundeslaenderPaths {
		var bundesland = extractBundesland(path, id)
		// area := calculatePolygonAreaNew(bundesland.coordinates)
		/*	if id == "Berlin" {
				berlinArea = area
			}
			if id == "Bremen" {
				bremenArea = area
			}*/

		// bundesland.area = area
		bundeslaender = append(bundeslaender, bundesland)
		areaGermany = areaGermany + bundesland.area
		fmt.Printf("Fläche für %s: %f\n", id, bundesland.area)
	}
	// Special check for Brandenburg because it will calculate Berlin into the area
	/*for i := range bundeslaender {
		if bundeslaender[i].id == "Brandenburg" {
			// Update the area attribute of the element
			bundeslaender[i].area = bundeslaender[i].area - berlinArea
			break
		}
		if bundeslaender[i].id == "Niedersachsen" {
			// Update the area attribute of the element
			bundeslaender[i].area = bundeslaender[i].area - bremenArea
			break
		}
	}*/
	sort.Slice(bundeslaender, func(i, j int) bool {
		return bundeslaender[i].area > bundeslaender[j].area
	})
	for j := 0; j < len(bundeslaender); j++ {
		fmt.Println("Prozentualer Anteil von  ", bundeslaender[j].id, " ist ", bundeslaender[j].area/areaGermany*100)
	}

	staedtePaths := extractStaedtePaths(svgString)
	for id, path := range staedtePaths {
		var stadt = extractStadt(path, id)
		staedte = append(staedte, stadt)
	}
	for i := 0; i < len(staedte); i++ {
		for j := 0; j < len(bundeslaender); j++ {
			var test = istStadtInBundesland(staedte[i], bundeslaender[j])
			if test == true {
				fmt.Println("City", staedte[i].id, "is in", bundeslaender[j].id)
			}
		}
	}

}

func istStadtInBundesland(stadt City, bundesl Bundesland) bool {
	numPolygons := len(bundesl.coordinates)
	if numPolygons == 0 {
		return false
	}

	intersections := 0
	for _, polygon := range bundesl.coordinates {
		numVertices := len(polygon)
		if numVertices < 3 {
			continue
		}

		for i := 0; i < numVertices; i++ {
			currentVertex := polygon[i]
			nextVertex := polygon[(i+1)%numVertices]

			if (currentVertex.Y > stadt.coordinate.Y) != (nextVertex.Y > stadt.coordinate.Y) &&
				stadt.coordinate.X < (nextVertex.X-currentVertex.X)*(stadt.coordinate.Y-currentVertex.Y)/(nextVertex.Y-currentVertex.Y)+currentVertex.X {
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

	var lastX, lastY float64
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

			prefix = coord[0]

			absX, absY := convertToAbsolute(lastX, lastY, x, y, prefix)
			polygon = append(polygon, Point{absX, absY})

			lastX, lastY = absX, absY
			if prefix == "L" {
				coordinatePoints = append(coordinatePoints, polygon)
				polygon = []Point{}
			}
		}
	}

	var bundesland Bundesland
	bundesland.id = id
	bundesland.area = 0.0

	for _, polygonPoints := range coordinatePoints {
		var points []Point
		for _, point := range polygonPoints {
			points = append(points, Point{X: point.X, Y: point.Y})
		}
		area := calculatePolygonAreaNew(points)
		coordinatePoints = append(coordinatePoints, points)
		bundesland.area += area
	}
	bundesland.coordinates = coordinatePoints

	return bundesland
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
	regexPattern := `([a-zA-Z])([-+]?\d*\.?\d+(?:[eE][-+]?\d+)?),([-+]?\d*\.?\d+(?:[eE][-+]?\d+)?)`
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
		coord1 := match[2]
		coord2 := match[3]
		coordinates[i] = []string{prefix, coord1, coord2}
	}

	return coordinates
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

// pro Bundesland koordinaten zurückgeben
// die Städte einzelne Punkte auslesen
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
