package main

import (
	"strings"
)

type svsDetails struct {
	// Standard fields
	Code   string
	Issued int64
	Name   string
	Text   string
	Wfo    string

	// Derived fields
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
	wxEvent.Details = &svsDetails{
		Code:               strings.ToLower(product.ProductCode),
		Issued:             product.IssuanceTime.Unix(),
		Name:               product.ProductName,
		Wfo:                product.IssuingOffice,
		Text:               normalizeString(product.ProductText, true),
		IsTornadoEmergency: true,
	}

	return wxEvent, nil
}
