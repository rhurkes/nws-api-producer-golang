package main

import (
	"encoding/json"
	"testing"
)

func TestParseSEVEventNoWatches(t *testing.T) {
	var product product
	svsPath := "./test_data/sev-no-watches.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)

	expectedWatches := []watchDetails{}
	expectedDetails := sevDetails{Watches: expectedWatches}
	expected := wxEvent{Data: nwsData{Derived: expectedDetails}}

	result, err := parseSEVEvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestParseSEVEvent failed")
	}
}

func TestParseSEVEventTorAndSvrWatch(t *testing.T) {
	var product product
	svsPath := "./test_data/sev-tor-and-svr-watches.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)

	expectedWatches := []watchDetails{}
	expectedDetails := sevDetails{Watches: expectedWatches}
	expected := wxEvent{Data: nwsData{Derived: expectedDetails}}

	result, err := parseSEVEvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestParseSEVEvent failed")
	}
}
