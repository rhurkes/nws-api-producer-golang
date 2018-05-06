package main

import (
	"fmt"
	"strconv"
	"strings"
)

func getNWSProductCode(product nwsProduct) string {
	switch product {
	case 0:
		return "afd"
	case 1:
		return "lsr"
	case 2:
		return "sel"
	case 3:
		return "svr"
	case 4:
		return "svs"
	case 5:
		return "swo"
	case 6:
		return "tor"
	}

	// This case isn't possible but needed for the function signature
	return ""
}

func getLatFromString(input string) float32 {
	if len(input) != 4 {
		logger.Warn(fmt.Sprintf("Unable to parse Lat from '%s'", input))
		return 0
	}

	lat, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", input[0:2], input[2:4]), 32)

	return float32(lat)
}

func getLonFromString(input string) float32 {
	if len(input) != 4 {
		logger.Warn("Unable to parse Lon from '%s'", input)
		return 0
	}

	lonFirstPart := input[0:2]

	if string(lonFirstPart[0]) == "0" {
		lonFirstPart = fmt.Sprintf("%v%s", 1, lonFirstPart)
	}

	lon, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", lonFirstPart, input[2:4]), 32)

	return float32(lon * -1)
}

func getPolygon(text string) []coordinates {
	var polygon []coordinates

	latLonLineMatch := latLonLineRegex.FindStringSubmatch(text)
	if len(latLonLineMatch) != 2 {
		return polygon
	}

	latLonMatches := latLonRegex.FindAllString(latLonLineMatch[0], -1)

	for _, val := range latLonMatches {
		polygon = append(polygon, coordinates{
			Lat: getLatFromString(val[0:4]),
			Lon: getLonFromString(val[5:9]),
		})
	}

	return polygon
}

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
		movement.Location = coordinates{
			Lat: getLatFromString(location[0:4]),
			Lon: getLonFromString(location[5:9]),
		}
	}

	return movement
}

func normalizeString(input string, preserveCase bool) string {
	textWithoutBreaks := strings.Replace(input, "\n", " ", -1)
	trimmedText := strings.TrimSpace(textWithoutBreaks)

	for strings.Contains(trimmedText, "  ") {
		trimmedText = strings.Replace(trimmedText, "  ", " ", -1)
	}

	if preserveCase {
		return trimmedText
	}

	return strings.ToLower(trimmedText)
}

func normalizeFloat(input string) float32 {
	inputString := normalizeString(input, false)
	num, err := strconv.ParseFloat(inputString, 32)
	if err != nil {
		num = 0
	}

	return float32(num)
}

// GetTimezoneOffset takes a three-character timezone string and translates it to an offset.
func GetTimezoneOffset(timezone string) string {
	offset := "0000" // Default to UTC

	switch strings.TrimSpace(strings.ToLower(timezone)) {
	case "hst":
		offset = "1000"
	case "hdt":
		offset = "0900"
	case "akst":
		offset = "0900"
	case "akdt":
		offset = "0800"
	case "pst":
		offset = "0800"
	case "pdt":
		offset = "0700"
	case "mst":
		offset = "0700"
	case "mdt":
		offset = "0600"
	case "cst":
		offset = "0600"
	case "cdt":
		offset = "0500"
	case "est":
		offset = "0500"
	case "edt":
		offset = "0400"
	default:
		logger.Warn("Unrecognized timezone: '%s'", timezone)
	}

	return offset
}
