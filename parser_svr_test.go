package main

import (
	"encoding/json"
	"testing"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestBuildSVR(t *testing.T) {
	var product Product
	svsPath := "./data/svr.json"
	json.Unmarshal(helpers.ReadJSONFromFile(svsPath), &product)

	expectedDetails := svrDetails{
		Code:      "svr",
		Issued:    1523658960,
		Name:      "Severe Thunderstorm Warning",
		Wfo:       "KDMX",
		IssuedFor: "western greene county in west central iowa, eastern carroll county in west central iowa",
		Polygon: []Coordinates{
			{Lat: 42.21, Lon: -94.75}, {Lat: 42.21, Lon: -94.34}, {Lat: 41.91, Lon: -94.52}, {Lat: 41.91, Lon: -94.75},
		},
		Location:      Coordinates{Lat: 41.98, Lon: -94.62},
		Time:          "2236z",
		MotionDegrees: 206,
		MotionKnots:   24,
	}

	expected := WxEvent{Details: expectedDetails}

	result, err := buildSVREvent(product)
	if err != nil || !helpers.CompareObjects(result, expected) {
		t.Error("TestBuildSVR failed")
	}
}
