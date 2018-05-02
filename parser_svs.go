package main

import (
	"errors"
	"strings"

	"github.com/rhurkes/wxNwsProducer/helpers"
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
func buildSVSEvent(product Product) (WxEvent, error) {
	wxEvent := WxEvent{}
	lowerCaseText := strings.ToLower(product.ProductText)

	if !strings.Contains(lowerCaseText, "tornado emergency") {
		return wxEvent, errors.New("Ignoring since not tornado emergency: svs " + product.ID)
	}

	wxEvent.Details = &svsDetails{
		Code:               strings.ToLower(product.ProductCode),
		Issued:             product.IssuanceTime.Unix(),
		Name:               product.ProductName,
		Wfo:                product.IssuingOffice,
		Text:               helpers.NormalizeString(product.ProductText, true),
		IsTornadoEmergency: true,
	}

	return wxEvent, nil
}
