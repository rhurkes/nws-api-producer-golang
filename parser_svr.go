package main

import (
	"regexp"
	"strings"
)

type svrDetails struct {
	IsPDS         bool
	IssuedFor     string
	Polygon       []coordinates
	Location      coordinates
	Time          string
	MotionDegrees int
	MotionKnots   int
}

var warningForRegex = regexp.MustCompile(`\n\n\*[\s|\S]+?warning for\.{3}\n([\s|\S]+?)\n\n`)

// Parses products and builds events for Severe Thunderstorm Warnings
func buildSVREvent(product product) (wxEvent, error) {
	wxEvent := wxEvent{Data: nwsData{Derived: deriveSVRDetails(product.ProductText)}}

	return wxEvent, nil
}

func deriveSVRDetails(text string) svrDetails {
	lowerCaseText := strings.ToLower(text)
	movement := getMovement(lowerCaseText)
	details := svrDetails{}
	details.Polygon = getPolygon(lowerCaseText)
	details.Time = movement.Time
	details.Location = movement.Location
	details.MotionDegrees = movement.Degrees
	details.MotionKnots = movement.Knots
	details.IssuedFor = getWarningFor(lowerCaseText)
	details.IsPDS = strings.Contains(lowerCaseText, "particularly dangerous situation")

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
