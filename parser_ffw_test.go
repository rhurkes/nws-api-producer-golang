package main

import (
	"encoding/json"
	"testing"
)

func TestParseFFWEvent(t *testing.T) {
	var product product
	svsPath := "./test_data/ffw.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)

	expectedDetails := ffwDetails{
		IsPDS:     false,
		IssuedFor: []string{"mitchell county in north central kansas", "southeastern osborne county in north central kansas"},
		Polygon: []coordinates{
			{Lat: 39.35, Lon: -98.47}, {Lat: 39.53, Lon: -97.93}, {Lat: 39.22, Lon: -97.93}, {Lat: 39.22, Lon: -98.49}, {Lat: 39.13, Lon: -98.49}, {Lat: 39.13, Lon: -98.89},
		},
	}

	expected := wxEvent{Data: nwsData{Derived: expectedDetails}}

	result, err := parseFFWEvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestParseFFWEvent failed")
	}
}

func TestDeriveFFWDetailsIsPDS(t *testing.T) {
	input := "THIS IS A PARTICULARLY DANGEROUS SITUATION."

	result := deriveFFWDetails(input)
	if !result.IsPDS {
		t.Error("TestDeriveFFWDetailsIsPDS failed")
	}
}
