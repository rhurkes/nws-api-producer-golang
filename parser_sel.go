package main

import (
	"regexp"
	"strconv"
	"strings"
)

type watchStats struct {
	Type   string
	Number int
}

type selDetails struct {
	// Standard fields
	Code   string
	Issued int64
	Name   string
	Text   string
	Wfo    string

	// Derived fields
	IsPDS       bool
	WatchNumber int
	WatchType   string
	Status      string // issued, cancelled, unknown
	IssuedFor   string
}

var watchTypeAndNumberRegex = regexp.MustCompile(`\n(.+)\swatch number\s(\d{1,3})\n`)
var issuedForRegex = regexp.MustCompile(`\n\nthe nws storm prediction center has issued a\n\n\*\s[\s|\S]+?watch for portions of\s\n([\s|\S]+?)\n\n`)

// Parses products and builds events for Severe Local Storm Watch and Watch Cancellation Msg.
// Issued when watches are issued. Has the watch text.
func buildSELEvent(product product) (wxEvent, error) {
	wxEvent := wxEvent{}

	details := selDetails{
		Code:   strings.ToLower(product.ProductCode),
		Issued: product.IssuanceTime.Unix(),
		Name:   product.ProductName,
		Wfo:    product.IssuingOffice,
		Text:   product.ProductText,
	}

	wxEvent.Details = deriveSELDetails(product.ProductText, details)

	return wxEvent, nil
}

func deriveSELDetails(text string, details selDetails) selDetails {
	lowerCaseText := strings.ToLower(text)
	details.IsPDS = strings.Contains(lowerCaseText, "this is a particularly dangerous situation")
	stats := getWatchStats(lowerCaseText)
	details.WatchType = stats.Type
	details.WatchNumber = stats.Number
	details.Status = getStatus(lowerCaseText)
	details.IssuedFor = getIssuedFor(lowerCaseText)

	return details
}

func getWatchStats(text string) watchStats {
	stats := watchStats{}
	match := watchTypeAndNumberRegex.FindStringSubmatch(text)

	if len(match) == 3 {
		stats.Type = match[1]
		stats.Number, _ = strconv.Atoi(match[2])
	}

	return stats
}

func getIssuedFor(text string) string {
	issuedFor := ""
	issuedForMatch := issuedForRegex.FindStringSubmatch(text)

	if len(issuedForMatch) != 2 {
		return issuedFor
	}

	issuedFor = strings.Replace(issuedForMatch[1], "\n", ", ", -1)
	issuedFor = strings.Replace(issuedFor, "   ", " ", -1)
	issuedFor = strings.Trim(issuedFor, " ")

	return issuedFor
}

func getStatus(text string) string {
	status := "unknown"

	if strings.Contains(text, "the nws storm prediction center has cancelled") {
		status = "cancelled"
	} else if strings.Contains(text, "the nws storm prediction center has issued") {
		status = "issued"
	}

	return status
}
