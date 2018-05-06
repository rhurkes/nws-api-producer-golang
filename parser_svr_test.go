package main

import (
	"encoding/json"
	"testing"
)

func TestBuildSVR(t *testing.T) {
	var product product
	svsPath := "./data/svr.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)

	expectedDetails := svrDetails{
		Code:      "svr",
		Issued:    1523658960,
		Name:      "Severe Thunderstorm Warning",
		Wfo:       "KDMX",
		IsPDS:     false,
		IssuedFor: "western greene county in west central iowa, eastern carroll county in west central iowa",
		Polygon: []coordinates{
			{Lat: 42.21, Lon: -94.75}, {Lat: 42.21, Lon: -94.34}, {Lat: 41.91, Lon: -94.52}, {Lat: 41.91, Lon: -94.75},
		},
		Location:      coordinates{Lat: 41.98, Lon: -94.62},
		Time:          "2236z",
		MotionDegrees: 206,
		MotionKnots:   24,
	}

	expected := wxEvent{Details: expectedDetails}

	result, err := buildSVREvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestBuildSVR failed")
	}
}

func TestDeriveSVRDetailsIsPDS(t *testing.T) {
	input := "THIS IS A PARTICULARLY DANGEROUS SITUATION."

	result := deriveSVRDetails(input, svrDetails{})
	if !result.IsPDS {
		t.Error("TestDeriveSVRDetailsIsPDS failed")
	}
}
