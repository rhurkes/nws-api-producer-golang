package main

import (
	"fmt"
	"strings"
)

// GetTimezoneOffset takes a three-character timezone string and translates it to an offset.
func GetTimezoneOffset(timezone string) string {
	offset := "0000" // Default to UTC

	switch strings.TrimSpace(strings.ToLower(timezone)) {
	case "hst":
		offset = "1000"
	case "hdt":
		offset = "0900"
	case "akst":
		offset = "0900"
	case "akdt":
		offset = "0800"
	case "pst":
		offset = "0800"
	case "pdt":
		offset = "0700"
	case "mst":
		offset = "0700"
	case "mdt":
		offset = "0600"
	case "cst":
		offset = "0600"
	case "cdt":
		offset = "0500"
	case "est":
		offset = "0500"
	case "edt":
		offset = "0400"
	default:
		fmt.Println(fmt.Sprintf("Unrecognized timezone: '%s'", timezone))
	}

	return offset
}
