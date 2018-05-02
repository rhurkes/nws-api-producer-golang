package main

import (
	"time"
)

type afdDetails struct {
	Code   string
	Issued time.Time
	Name   string
	Text   string
	Wfo    string
}

// Parses products and builds events for Area Forecast Discussions
func buildAFDEvent(product Product) (WxEvent, error) {
	wxEvent := WxEvent{}

	wxEvent.Details = &afdDetails{
		Code:   "afd",
		Issued: product.IssuanceTime,
		Name:   product.ProductName,
		Wfo:    product.IssuingOffice,
		Text:   normalizeString(product.ProductText, false),
	}

	return wxEvent, nil
}
