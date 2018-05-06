package main

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBuildAFDEvent(t *testing.T) {
	var product product
	afdPath := "./data/afd-mpx.json"
	json.Unmarshal(ReadJSONFromFile(afdPath), &product)
	product.ProductText = "afd\ntext"
	productTime, _ := time.Parse(time.RFC3339, "2018-04-14T02:07:00Z")

	expectedDetails := afdDetails{
		Code:   "afd",
		Issued: productTime,
		Name:   "Area Forecast Discussion",
		Text:   "afd text",
		Wfo:    "KMPX",
	}

	expected := wxEvent{Details: expectedDetails}

	result, err := buildAFDEvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestBuildAFDEvent failed")
	}
}
