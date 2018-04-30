package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestBuildAFDEvent(t *testing.T) {
	var product Product
	afdPath := "./data/afd-mpx.json"
	json.Unmarshal(helpers.ReadJSONFromFile(afdPath), &product)
	product.ProductText = "afd\ntext"
	productTime, _ := time.Parse(time.RFC3339, "2018-04-14T02:07:00Z")

	expectedDetails := afdDetails{
		Code:   "afd",
		Issued: productTime,
		Name:   "Area Forecast Discussion",
		Text:   "afd text",
		Wfo:    "KMPX",
	}

	expected := WxEvent{Details: expectedDetails}

	result, err := buildAFDEvent(product)
	if err != nil || !helpers.CompareObjects(result, expected) {
		t.Error("TestBuildAFDEvent failed")
	}
}
