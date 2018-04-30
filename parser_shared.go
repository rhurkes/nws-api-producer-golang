package main

import (
	"fmt"
	"strconv"
)

type movement struct {
	Time     string
	Location Coordinates
	Degrees  int
	Knots    int
}

func getLatFromString(input string) float32 {
	if len(input) != 4 {
		fmt.Println(fmt.Sprintf("Unable to parse Lat from '%s'", input))
		return 0
	}

	lat, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", input[0:2], input[2:4]), 32)

	return float32(lat)
}

func getLonFromString(input string) float32 {
	if len(input) != 4 {
		fmt.Println(fmt.Sprintf("Unable to parse Lon from '%s'", input))
		return 0
	}

	lonFirstPart := input[0:2]

	if string(lonFirstPart[0]) == "0" {
		lonFirstPart = fmt.Sprintf("%v%s", 1, lonFirstPart)
	}

	lon, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", lonFirstPart, input[2:4]), 32)

	return float32(lon * -1)
}

// TODO share with SVR
func getPolygon(text string) []Coordinates {
	var polygon []Coordinates

	latLonLineMatch := latLonLineRegex.FindStringSubmatch(text)
	if len(latLonLineMatch) != 2 {
		return polygon
	}

	latLonMatches := latLonRegex.FindAllString(latLonLineMatch[0], -1)

	for _, val := range latLonMatches {
		polygon = append(polygon, Coordinates{
			Lat: getLatFromString(val[0:4]),
			Lon: getLonFromString(val[5:9]),
		})
	}

	return polygon
}

// TODO share with SVR
func getMovement(text string) movement {
	movement := movement{}
	movementMatch := movementRegex.FindStringSubmatch(text)

	if len(movementMatch) == 5 {
		movement.Time = movementMatch[1]
		degrees, err := strconv.Atoi(movementMatch[2])

		if err == nil {
			movement.Degrees = degrees
		}

		knots, err := strconv.Atoi(movementMatch[3])
		if err == nil {
			movement.Knots = knots
		}

		location := movementMatch[4]
		movement.Location = Coordinates{
			Lat: getLatFromString(location[0:4]),
			Lon: getLonFromString(location[5:9]),
		}
	}

	return movement
}
