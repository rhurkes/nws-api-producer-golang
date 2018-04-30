package main

import (
	"regexp"
	"strings"
)

type svrDetails struct {
	// Standard fields
	Code   string
	Issued int64
	Name   string
	Wfo    string

	// Derived fields
	IssuedFor     string
	Polygon       []Coordinates
	Location      Coordinates
	Time          string
	MotionDegrees int
	MotionKnots   int
}

var warningForRegex = regexp.MustCompile(`\n\n\*[\s|\S]+?warning for\.{3}\n([\s|\S]+?)\n\n`)

// Parses products and builds events for Severe Thunderstorm Warnings
func buildSVREvent(product Product) (WxEvent, error) {
	wxEvent := WxEvent{}

	details := svrDetails{
		Code:   strings.ToLower(product.ProductCode),
		Issued: product.IssuanceTime.Unix(),
		Name:   product.ProductName,
		Wfo:    product.IssuingOffice,
	}

	wxEvent.Details = deriveSVRDetails(product.ProductText, details)

	return wxEvent, nil
}

func deriveSVRDetails(text string, details svrDetails) svrDetails {
	lowerCaseText := strings.ToLower(text)
	details.Polygon = getPolygon(lowerCaseText)

	movement := getMovement(lowerCaseText)
	details.Time = movement.Time
	details.Location = movement.Location
	details.MotionDegrees = movement.Degrees
	details.MotionKnots = movement.Knots
	details.IssuedFor = getWarningFor(lowerCaseText)

	return details
}

func getWarningFor(text string) string {
	warningFor := ""

	warningForMatch := warningForRegex.FindStringSubmatch(text)

	if len(warningForMatch) == 2 {
		warningFor = strings.Replace(warningForMatch[1], "...", "", -1)
		warningFor = strings.Replace(warningFor, "  ", "", -1)
		warningFor = strings.Replace(warningFor, "\n", ", ", -1)
	}

	return warningFor
}
