package main

import (
	"strings"
)

type svsDetails struct {
	IsTornadoEmergency bool
}

// Parses products and builds events for Severe Weather Statements - currently only produces
// messages for Tornado Emergencies.
func buildSVSEvent(product product) (wxEvent, error) {
	wxEvent := wxEvent{DoNotPublish: true}
	lowerCaseText := strings.ToLower(product.ProductText)

	if !strings.Contains(lowerCaseText, "tornado emergency") {
		return wxEvent, nil
	}

	wxEvent.DoNotPublish = false
	wxEvent.Data = nwsData{Derived: svsDetails{IsTornadoEmergency: true}}

	return wxEvent, nil
}
