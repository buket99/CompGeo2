package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strings"

	"github.com/tdewolff/parse"
)

type Point struct {
	X float64
	Y float64
}

type Path struct {
	Points []Point
}

func parseSVGPath(svgPath string) (*Path, error) {
	path := &Path{}
	r := strings.NewReader(svgPath)

	for {
		t, _, _ := parse.Next(r)
		if t == parse.ErrorToken {
			break
		}

		if t == parse.CommandToken {
			c, _ := parse.ReadCommand(r)

			switch c {
			case 'm', 'M':
				for {
					p, err := parse.ReadFloats(r, 2)
					if err != nil {
						break
					}

					path.Points = append(path.Points, Point{X: p[0], Y: p[1]})
				}
			case 'l', 'L':
				for {
					p, err := parse.ReadFloats(r, 2)
					if err != nil {
						break
					}

					path.Points = append(path.Points, Point{X: p[0], Y: p[1]})
				}
			case 'z', 'Z':
				// Ignore close path commands
			default:
				return nil, fmt.Errorf("Unsupported SVG path command: %c", c)
			}
		}
	}

	return path, nil
}

func calculatePathArea(path *Path) float64 {
	area := 0.0
	n := len(path.Points)

	if n < 3 {
		return area
	}

	for i := 0; i < n-1; i++ {
		area += (path.Points[i].X * path.Points[i+1].Y) - (path.Points[i+1].X * path.Points[i].Y)
	}

	area += (path.Points[n-1].X * path.Points[0].Y) - (path.Points[0].X * path.Points[n-1].Y)

	return math.Abs(area) / 2.0
}

func main() {
	// Read the SVG data from a file
	svgData, err := ioutil.ReadFile("path/to/svg/file.svg")
	if err != nil {
		fmt.Println("Failed to read SVG file:", err)
		return
	}

	// Convert the SVG data to a string
	svgString := string(svgData)

	// Define the regular expression pattern to match the path elements
	pathRegex := regexp.MustCompile(`<path\s+[^>]*d="([^"]+)"[^>]*>`)

	// Find all matches of the pattern in the SVG data
	matches := pathRegex.FindAllStringSubmatch(svgString, -1)

	// Iterate over the matches and extract the path coordinates
	for _, match := range matches {
		if len(match) >= 2 {
			svgPath := match[1]

			// Parse the SVG path to get the coordinates
			path, err := parseSVGPath(svgPath)
			if err != nil {
				fmt.Println("Failed to parse SVG path:", err)
				continue
			}

			// Calculate the area surrounded by the path
			area := calculatePathArea(path)
			fmt.Println("Path Area:", area)
		}
	}
}
