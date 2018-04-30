package main

import (
	"encoding/json"
	"testing"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestBuildSVSEventNotTornadoEmergency(t *testing.T) {
	var product Product
	svsPath := "./data/svs-svr-canceled.json"
	json.Unmarshal(helpers.ReadJSONFromFile(svsPath), &product)

	_, err := buildSVSEvent(product)
	if err == nil {
		t.Error("TestBuildSVSEventNotTornadoEmergency failed")
	}
}

func TestBuildSVSEventIsTornadoEmergency(t *testing.T) {
	var product Product
	svsPath := "./data/svs-svr-canceled.json"
	json.Unmarshal(helpers.ReadJSONFromFile(svsPath), &product)
	product.ProductText = "THIS IS A TORNADO EMERGENCY"

	expectedDetails := svsDetails{
		Code:               "svs",
		Issued:             1523824740,
		Name:               "Severe Weather Statement",
		Text:               "THIS IS A TORNADO EMERGENCY",
		Wfo:                "KRNK",
		IsTornadoEmergency: true,
	}

	expected := WxEvent{Details: expectedDetails}

	result, err := buildSVSEvent(product)
	if err != nil || !helpers.CompareObjects(result, expected) {
		t.Error("TestBuildSVSEventIsTornadoEmergency failed")
	}
}
