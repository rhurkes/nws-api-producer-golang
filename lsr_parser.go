package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const reportThresholdMinutes float64 = 60

var timezoneRegex = regexp.MustCompile(`\d{3,4}\s[A|P]M\s(\w{3})\s`)

func getLSRTimezoneOffset(line string) string {
	offset := "0000" // Default offset is UTC
	result := timezoneRegex.FindStringSubmatch(line)

	if len(result) == 2 {
		offset = GetTimezoneOffset(result[1])
	} else {
		fmt.Println(fmt.Sprintf("Unable to parse timezone offset: '%s'", line))
	}

	return offset
}

func processLSRProduct(product Product) (WxEvent, error) {
	wxEvent := WxEvent{}
	details := LSRDetails{}
	lines := strings.Split(product.ProductText, "\n")

	if len(lines) < 16 {
		return wxEvent, errors.New("LSR body missing lines")
	}

	if strings.Contains(lines[5], "SUMMARY") {
		return wxEvent, errors.New("Do not parse summaries")
	}

	remarksLineIndex := -1
	for i, val := range lines {
		if strings.Contains(val, "..REMARKS..") {
			remarksLineIndex = i
			break
		}
	}

	if remarksLineIndex == -1 {
		return wxEvent, errors.New("Remarks section not found, needed for parsing")
	}

	// 2 lines after ..REMARKS.. contains TIME/EVENT/CITY LOCATION/LAT/LON
	currentLine := lines[remarksLineIndex+2]
	rawTime := currentLine[0:7]
	details.Type = normalizeString(currentLine[12:29])
	details.Location = normalizeString(currentLine[29:53])
	details.Lat = normalizeFloat(currentLine[53:58])
	details.Lon = normalizeFloat(currentLine[59:66]) * -1

	// 3 lines after ..REMARKS.. contains DATE/MAG/COUNTY/ST/SOURCE
	currentLine = lines[remarksLineIndex+3]
	rawTime = fmt.Sprintf("%s %s", currentLine[0:10], rawTime)
	// TODO break this out into size, isMeasured?
	details.Mag = normalizeString(currentLine[12:29])
	details.County = normalizeString(currentLine[29:48])
	details.State = normalizeString(currentLine[48:50])
	details.Source = normalizeString(currentLine[50:])

	// 5+ lines after ..REMARKS.. contains actual remarks (if present)
	remarks := ""
	for _, val := range lines[remarksLineIndex+5:] {
		if strings.Contains(val, "&&") || strings.Contains(val, "$$") {
			break
		} else {
			remarks += val
		}
	}

	// Figure out local reported time
	localReportOffset := getLSRTimezoneOffset(lines[7])
	rawTime = fmt.Sprintf("%s-%s", rawTime, localReportOffset)
	reportedTime, err := time.Parse("01/02/2006 1504 PM-0700", rawTime)
	if err == nil {
		details.Reported = reportedTime
		if product.IssuanceTime.Sub(reportedTime).Minutes() > reportThresholdMinutes {
			return wxEvent, fmt.Errorf("Report time (%s) older than threshold (%v)", reportedTime, reportThresholdMinutes)
		}
	} else {
		return wxEvent, fmt.Errorf("Unable to format local reported time: '%s'", rawTime)
	}

	details.Remarks = normalizeString(remarks)
	wxEvent.Details = details

	return wxEvent, nil
}
