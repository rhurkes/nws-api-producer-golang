package main

import (
	"regexp"
	"strings"
)

type torDetails struct {
	IsTornadoEmergency bool
	IsPDS              bool
	IsObserved         bool
	Source             string
	Description        string
	Polygon            []coordinates
	Location           coordinates
	Time               string
	MotionDegrees      int
	MotionKnots        int
}

var sourceRegex = regexp.MustCompile(`\n{2}\s{2}source...(.+)\.\s?\n{2}`)
var descriptionRegex = regexp.MustCompile(`\n\*\s(at\s[\S|\s]+?)\n\n`)
var movementRegex = regexp.MustCompile(`\ntime...mot...loc\s(\d{4}z)\s(\d+)\D{3}\s(\d+)kt\s(\d{4}\s\d{4})`)
var latLonLineRegex = regexp.MustCompile(`lat...lon\s([\s|\S]+)time\.{3}`)
var latLonRegex = regexp.MustCompile(`(\d{4}\s\d{4})`)

// Parses products and builds events for Tornado Warnings
func buildTOREvent(product product) (wxEvent, error) {
	wxEvent := wxEvent{Data: nwsData{Derived: deriveTORDetails(product.ProductText)}}

	return wxEvent, nil
}

func deriveTORDetails(text string) torDetails {
	lowerCaseText := strings.ToLower(text)
	details := torDetails{}
	details.IsTornadoEmergency = strings.Contains(lowerCaseText, "tornado emergency")
	details.IsPDS = strings.Contains(lowerCaseText, "particularly dangerous situation")
	details.IsObserved = strings.Contains(lowerCaseText, "tornado...observed")
	details.Source = getSource(lowerCaseText)
	details.Description = getDescription(lowerCaseText)
	details.Polygon = getPolygon(lowerCaseText)
	movement := getMovement(lowerCaseText)
	details.Time = movement.Time
	details.Location = movement.Location
	details.MotionDegrees = movement.Degrees
	details.MotionKnots = movement.Knots

	return details
}

func getSource(text string) string {
	source := "unknown"
	sourceMatch := sourceRegex.FindStringSubmatch(text)
	if len(sourceMatch) == 2 {
		source = normalizeString(sourceMatch[1], false)
	}

	return source
}

func getDescription(text string) string {
	description := ""
	descriptionMatch := descriptionRegex.FindStringSubmatch(text)
	if len(descriptionMatch) == 2 {
		description = normalizeString(descriptionMatch[1], false)
	}

	return description
}
