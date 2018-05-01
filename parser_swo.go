package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var mdImageURIBase = "http://www.spc.noaa.gov/products/md/"
var numberRegex = regexp.MustCompile(`\n\nMesoscale Discussion (\d{4})\n`)
var validRegex = regexp.MustCompile(`Valid\s((\d{6})Z\s-\s(\d{6})Z)\n\n`)
var concerningRegex = regexp.MustCompile(`\n\n(Concerning\.{3}.*?)\n\n`)
var affectedRegex = regexp.MustCompile(`\n\nAreas affected\.{3}([\S\s]+?)\n\n`)
var watchInfoRegex = regexp.MustCompile(`\n\nValid.+\n\n([\S\s]+?)\n\n`)
var mdSummaryRegex = regexp.MustCompile(`\n\nSUMMARY\.{3}([\S|\s]+?)\n\n`)
var outlookSummaryRegex = regexp.MustCompile(`\n{2}\.{3}SUMMARY\.{3}\n([\S|\s]*?)\n{2}`)
var forecasterRegex = regexp.MustCompile(`\n\n\.{2}(.+?)\.{2}\s\d{2}\/\d{2}\/\d{4}`)
var polygonRegex = regexp.MustCompile(`\d{8}`)
var wfoRegex = regexp.MustCompile(`ATTN\.{3}WFO\.{3}((?:\w{3}\.{3})+)`)

type outlookDetails struct {
	// Standard fields
	Code   string
	Issued int64
	Name   string
	Wfo    string

	// Derived fields
	SubCode    string // dy1, dy2, dy3, d48
	Valid      string // 0600Z, 1300Z, 1630Z, 1730Z, 2000Z, 0100Z
	Risk       string
	Summary    string
	Forecaster string
}

type mdDetails struct {
	// Standard fields
	Code   string
	Issued int64
	Name   string
	Wfo    string

	// Derived fields
	SubCode    string
	Number     string
	Affected   string
	Concerning string
	WatchInfo  string
	Valid      time.Time
	Expires    time.Time
	WFOs       []string
	Summary    string
	Forecaster string
	ImageURI   string
	Polygon    []Coordinates
}

// Parses products and builds events for Severe Storm Outlook Narratives
func buildSWOEvent(product Product) (WxEvent, error) {
	wxEvent := WxEvent{}
	lines := strings.Split(product.ProductText, "\n")

	if len(lines) < 20 {
		msg := fmt.Sprintf("Cannot parse %s: fewer than 20 lines of text", product.ID)
		return wxEvent, errors.New(msg)
	}

	if lines[3] == "SWOMCD" {
		wxEvent.Details = parseSWOMCD(product)
	} else {
		wxEvent.Details = parseSWODY(product)
	}

	return wxEvent, nil
}

func parseSWOMCD(product Product) mdDetails {
	text := product.ProductText
	year := product.IssuanceTime.Year()
	details := mdDetails{
		Code:    strings.ToLower(product.ProductCode),
		SubCode: "mcd",
		Issued:  product.IssuanceTime.Unix(),
		Name:    product.ProductName,
		Wfo:     product.IssuingOffice,
	}

	numberMatch := numberRegex.FindStringSubmatch(text)
	if len(numberMatch) == 2 {
		details.Number = numberMatch[1]
	}

	concerningMatch := concerningRegex.FindStringSubmatch(text)
	if len(concerningMatch) == 2 {
		details.Concerning = concerningMatch[1]
	}

	affectedMatch := affectedRegex.FindStringSubmatch(text)
	if len(affectedMatch) == 2 {
		details.Affected = strings.Replace(affectedMatch[1], "\n", " ", -1)
	}

	watchInfoMatch := watchInfoRegex.FindStringSubmatch(text)
	if len(watchInfoMatch) == 2 {
		details.WatchInfo = strings.Replace(watchInfoMatch[1], "\n", " ", -1)
	}

	summaryMatch := mdSummaryRegex.FindStringSubmatch(text)
	if len(summaryMatch) == 2 {
		details.Summary = strings.Replace(summaryMatch[1], "\n", " ", -1)
	}

	forecasterMatch := forecasterRegex.FindStringSubmatch(text)
	if len(forecasterMatch) == 2 {
		details.Forecaster = forecasterMatch[1]
	}

	wfoMatch := wfoRegex.FindStringSubmatch(text)
	if len(wfoMatch) == 2 {
		var wfos []string
		unfilteredWfos := strings.Split(wfoMatch[1], "...")
		for _, wfo := range unfilteredWfos {
			if len(wfo) > 0 {
				wfos = append(wfos, strings.ToLower(wfo))
			}
		}
		details.WFOs = wfos
	}

	details.ImageURI = fmt.Sprintf("%s%v/mcd%s.gif", mdImageURIBase, year, details.Number)
	polygonMatch := polygonRegex.FindAllString(text, -1)
	details.Polygon = buildPolygon(polygonMatch)

	validRange, err := getValidRange(text, product.IssuanceTime)
	if err == nil {
		details.Valid = validRange[0]
		details.Expires = validRange[1]
	}

	return details
}

func getValidRange(text string, issued time.Time) ([]time.Time, error) {
	match := validRegex.FindAllStringSubmatch(text, -1)

	if len(match) != 1 || len(match[0]) != 4 {
		return []time.Time{}, errors.New("Unable to parse valid range")
	}

	// Slice indexes and Atoi calls are safe since the regex only returns strings of 6 numbers
	rawStart := match[0][2]
	rawUntil := match[0][3]
	startday, _ := strconv.Atoi(rawStart[0:2])
	starthour, _ := strconv.Atoi(rawStart[2:4])
	startminute, _ := strconv.Atoi(rawStart[4:6])
	untilday, _ := strconv.Atoi(rawUntil[0:2])
	untilhour, _ := strconv.Atoi(rawUntil[2:4])
	untilminute, _ := strconv.Atoi(rawUntil[4:6])

	start := time.Date(issued.Year(), issued.Month(), issued.Day(),
		starthour, startminute, 0, 0, time.UTC)
	if startday != issued.Day() {
		// The valid start is the previous day
		start = start.AddDate(0, 0, -1)
	}

	until := time.Date(issued.Year(), issued.Month(), issued.Day(),
		untilhour, untilminute, 0, 0, time.UTC)
	if untilday != issued.Day() {
		// The valid until is the next day
		until = until.AddDate(0, 0, 1)
	}

	return []time.Time{start, until}, nil
}

func buildPolygon(matches []string) []Coordinates {
	var polygon []Coordinates

	for _, val := range matches {
		lat, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", val[0:2], val[2:4]), 32)

		lonFirstPart := val[4:6]
		if string(lonFirstPart[0]) == "0" {
			lonFirstPart = fmt.Sprintf("%v%s", 1, lonFirstPart)
		}

		lon, _ := strconv.ParseFloat(fmt.Sprintf("%s.%s", lonFirstPart, val[6:8]), 32)

		polygon = append(polygon, Coordinates{
			Lat: float32(lat),
			Lon: float32(lon) * -1,
		})
	}

	return polygon
}

func parseSWODY(product Product) outlookDetails {
	text := product.ProductText
	details := outlookDetails{
		Code:   strings.ToLower(product.ProductCode),
		Issued: product.IssuanceTime.Unix(),
		Name:   product.ProductName,
		Wfo:    product.IssuingOffice,
	}

	switch product.WmoCollectiveID {
	case "ACUS01":
		details.SubCode = "dy1"
	case "ACUS02":
		details.SubCode = "dy2"
	case "ACUS03":
		details.SubCode = "dy3"
	case "ACUS48":
		details.SubCode = "d48"
	default:
		fmt.Println(fmt.Sprintf("Unknown WmoCollectiveID: '%s'", product.WmoCollectiveID))
	}

	if details.SubCode == "dy1" || details.SubCode == "dy2" {
		validMatch := validRegex.FindAllStringSubmatch(text, -1)
		if len(validMatch) == 1 && len(validMatch[0]) == 4 {
			details.Valid = fmt.Sprintf("%sZ", validMatch[0][1][2:4])
		}
	}

	forecasterMatch := forecasterRegex.FindStringSubmatch(text)
	if len(forecasterMatch) == 2 {
		details.Forecaster = forecasterMatch[1]
	}

	details.Risk = getRisk(text)

	if details.SubCode != "d48" {
		summaryMatch := outlookSummaryRegex.FindStringSubmatch(text)
		if len(summaryMatch) == 2 {
			details.Summary = strings.Replace(summaryMatch[1], "\n", " ", -1)
		}
	}

	return details
}

func getRisk(text string) string {
	switch {
	case strings.Contains(text, "...THERE IS A HIGH RISK"):
		return "high"
	case strings.Contains(text, "...THERE IS A MODERATE RISK"):
		return "moderate"
	case strings.Contains(text, "...THERE IS AN ENHANCED RISK"):
		return "enhanced"
	case strings.Contains(text, "...THERE IS A SLIGHT RISK"):
		return "slight"
	case strings.Contains(text, "...THERE IS A MARGINAL RISK"):
		return "marginal"
	case strings.Contains(text, "...NO SEVERE THUNDERSTORM AREAS FORECAST..."):
		return "no_severe"
	default:
		return "unknown"
	}
}
