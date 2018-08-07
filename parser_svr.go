package main

import (
	"regexp"
	"strings"
)

type svrDetails struct {
	IsPDS         bool
	IssuedFor     []string
	Polygon       []coordinates
	Location      coordinates
	Time          string
	MotionDegrees int
	MotionKnots   int
}

var svrWarningForRegex = regexp.MustCompile(`\n\n\*[\s|\S]+?warning for\.{3}\n([\s|\S]+?)\n\n`)

// Severe Thunderstorm Warnings parser
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
	details.IssuedFor = getIssuedFor(lowerCaseText)
	details.IsPDS = strings.Contains(lowerCaseText, "particularly dangerous situation")

	return details
}
