package main

import (
	"strings"
)

type ffwDetails struct {
	IsPDS     bool
	IssuedFor []string
	Polygon   []coordinates
}

// Flash Flood Warning Parser
func parseFFWEvent(product product) (wxEvent, error) {
	wxEvent := wxEvent{Data: nwsData{Derived: deriveFFWDetails(product.ProductText)}}

	return wxEvent, nil
}

func deriveFFWDetails(text string) ffwDetails {
	lowerCaseText := strings.ToLower(text)
	details := ffwDetails{}
	details.Polygon = getPolygon(lowerCaseText)
	details.IssuedFor = getIssuedFor(lowerCaseText)
	details.IsPDS = strings.Contains(lowerCaseText, "particularly dangerous situation")

	return details
}
