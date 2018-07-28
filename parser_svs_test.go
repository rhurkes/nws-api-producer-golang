package main

import (
	"encoding/json"
	"testing"
)

func TestBuildSVSEventNotTornadoEmergency(t *testing.T) {
	var product product
	svsPath := "./test_data/svs-svr-canceled.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)

	_, err := buildSVSEvent(product)
	if err != nil {
		t.Error("TestBuildSVSEventNotTornadoEmergency failed")
	}
}

func TestBuildSVSEventIsTornadoEmergency(t *testing.T) {
	var product product
	svsPath := "./test_data/svs-svr-canceled.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)
	product.ProductText = "THIS IS A TORNADO EMERGENCY"

	expectedDetails := svsDetails{
		Code:               "svs",
		Issued:             1523824740,
		Name:               "Severe Weather Statement",
		Text:               "THIS IS A TORNADO EMERGENCY",
		Wfo:                "KRNK",
		IsTornadoEmergency: true,
	}

	expected := wxEvent{Details: expectedDetails}

	result, err := buildSVSEvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestBuildSVSEventIsTornadoEmergency failed")
	}
}
