package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type sevDetails struct {
	Watches []watchDetails
}

type watchDetails struct {
	IssuedTs    time.Time
	Polygon     []coordinates
	ExpiresTs   time.Time
	WatchNumber int
	WatchType   string
}

var watchMatchesRegex = regexp.MustCompile(`(:?sevr ([\s|\S])+;)+`)

// TODO name
func parseSEVEvent(product product) (wxEvent, error) {
	wxEvent := wxEvent{Data: nwsData{Derived: deriveSEVDetails(product.ProductText)}}

	return wxEvent, nil
}

// TODO. Should this ever return an error in the parser? What happens if any of these derivations blow up?
// Kind of like the idea of doing as much as I can and still publishing the event
func deriveSEVDetails(text string) sevDetails {
	details := sevDetails{}
	details.Watches = []watchDetails{}
	lowerCaseText := strings.ToLower(text)

	// TODO findallstring vs findallstringsubmatch?
	watchMatches := watchMatchesRegex.FindAllStringSubmatch(lowerCaseText, -1)
	fmt.Printf("\n\nlength: %v", len(watchMatches))
	fmt.Printf("\n\nlength sub: %v", len(watchMatches[0]))
	fmt.Printf("\n\n%v\n\n", watchMatches[0][0])
	fmt.Printf("\n\n%v\n\n", watchMatches[0][1])
	fmt.Printf("\n\n%v\n\n", watchMatches[0][2])

	// for _, watchMatch := range watchMatches {
	// parsedMatch := parseWatchDetails(watchMatch[0])
	// TODO is this really idiomatic? This is annoying not being able to push items
	// details.Watches = append(details.Watches, parsedMatch)
	// }

	return details
}

// TODO how do expires times on the next day work? No examples atm
func parseWatchDetails(text string) watchDetails {
	// fmt.Printf("\n\nMatch: '%v'\n\n", text)

	return watchDetails{}
}
