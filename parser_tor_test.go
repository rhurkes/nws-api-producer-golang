package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestBuildTorEvent(t *testing.T) {
	var product Product
	path := "./data/tor-radar-indicated.json"
	json.Unmarshal(helpers.ReadJSONFromFile(path), &product)

	expectedDetails := torDetails{
		Code:               "tor",
		IsTornadoEmergency: false,
		IsPDS:              false,
		IsObserved:         false,
		Issued:             1523664540,
		Name:               "Tornado Warning",
		Source:             "radar indicated rotation",
		Description:        "at 709 pm cdt, a severe thunderstorm capable of producing a tornado was located near mansfield, or 10 miles northeast of ava, moving northeast at 25 mph.",
		Wfo:                "KSGF",
		Polygon: []Coordinates{
			{Lat: 37, Lon: -92.51}, {Lat: 37.06, Lon: -92.61}, {Lat: 37.24, Lon: -92.5}, {Lat: 37.23, Lon: -92.38}, {Lat: 37.12, Lon: -92.26},
		},
		Time:          "0009z",
		Location:      Coordinates{Lat: 37.06, Lon: -92.54},
		MotionDegrees: 217,
		MotionKnots:   22,
	}

	expected := WxEvent{Details: expectedDetails}

	result, err := buildTOREvent(product)
	if err != nil || !helpers.CompareObjects(result, expected) {
		t.Error("TestBuildTorEvent failed")
	}
}

func TestDeriveTORDetailsIsEmergency(t *testing.T) {
	input := "THIS IS A TORNADO EMERGENCY"

	result := deriveTORDetails(input, torDetails{})
	if !result.IsTornadoEmergency {
		t.Error("TestDeriveTORDetailsIsEmergency failed")
	}
}

func TestDeriveTORDetailsIsPDS(t *testing.T) {
	input := "THIS IS A PARTICULARLY DANGEROUS SITUATION."

	result := deriveTORDetails(input, torDetails{})
	if !result.IsPDS {
		t.Error("TestDeriveTORDetailsIsPDS failed")
	}
}

func TestDeriveTORDetailsIsObserved(t *testing.T) {
	input := "BLAH BLAH TORNADO...OBSERVED BLAH BLAH"

	result := deriveTORDetails(input, torDetails{})
	if !result.IsObserved {
		t.Error("TestDeriveTORDetailsIsObserved failed")
	}
}

func TestGetSource(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["Empty string"] = helpers.TestParameters{Input: "", Expected: "unknown"}
	tests["Happy path"] = helpers.TestParameters{Input: "\n\n  source...Weather spotters confirmed tornado. \n\n", Expected: "weather spotters confirmed tornado"}

	for testName, params := range tests {
		result := getSource(params.Input)

		if !helpers.CompareObjects(result, params.Expected) {
			msg := fmt.Sprintf("result: '%v', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetSource - %s failed. %s", testName, msg)
		}
	}
}

func TestGetDescription(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["Empty string"] = helpers.TestParameters{Input: "", Expected: ""}
	tests["Happy path"] = helpers.TestParameters{Input: `\n\* at 709 pm\n\n`, Expected: ""}

	for testName, params := range tests {
		result := getDescription(params.Input)

		if !helpers.CompareObjects(result, params.Expected) {
			msg := fmt.Sprintf("result: '%v', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetDescription - %s failed. %s", testName, msg)
		}
	}
}

func TestGetPolygon(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["Empty string"] = helpers.TestParameters{Input: "", Expected: nil}
	tests["Happy path"] = helpers.TestParameters{Input: "lat...lon 3267 9078 3268 0079 time...", Expected: []Coordinates{{Lat: 32.67, Lon: -90.78}, {Lat: 32.68, Lon: -100.79}}}
	tests["No coords"] = helpers.TestParameters{Input: "lat...lon NO COORDS time...", Expected: nil}

	for testName, params := range tests {
		result := getPolygon(params.Input)

		if !helpers.CompareObjects(result, params.Expected) {
			msg := fmt.Sprintf("result: '%v', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetPolygon - %s failed. %s", testName, msg)
		}
	}
}
