package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestBuildSELEventSVRWatch(t *testing.T) {
	var product product
	svsPath := "./test_data/sel-svr-watch.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)

	expectedDetails := selDetails{
		IsPDS:       false,
		WatchNumber: 25,
		WatchType:   "severe thunderstorm",
		Status:      "issued",
		IssuedFor:   "southwest arkansas, northwest louisiana, southeast oklahoma, central and northeast texas",
	}

	expected := wxEvent{Data: nwsData{Derived: expectedDetails}}

	result, err := buildSELEvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestBuildSELEventSVRWatch failed")
	}
}

func TestBuildSELEventTORWatch(t *testing.T) {
	var product product
	svsPath := "./test_data/sel-tor-watch.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)

	expectedDetails := selDetails{
		IsPDS:       false,
		WatchNumber: 26,
		WatchType:   "tornado",
		Status:      "issued",
		IssuedFor:   "southern and central indiana, northern kentucky, western and central ohio",
	}

	expected := wxEvent{Data: nwsData{Derived: expectedDetails}}

	result, err := buildSELEvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestBuildSELEventTORWatch failed")
	}
}

func TestGetWatchStats(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: watchStats{}}
	tests["SVR 5"] = TestParameters{Input: "\nsevere thunderstorm watch number 25\n", Expected: watchStats{Type: "severe thunderstorm", Number: 25}}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := getWatchStats(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestGetWatchStats - %s failed. %s", testName, msg)
		}
	}
}

func TestGetStatus(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty string"] = TestParameters{Input: "", Expected: "unknown"}
	tests["Cancelled"] = TestParameters{Input: "the nws storm prediction center has cancelled", Expected: "cancelled"}
	tests["Issued"] = TestParameters{Input: "the nws storm prediction center has issued", Expected: "issued"}
	tests["Not cancelled or issued"] = TestParameters{Input: "the nws storm prediction center has done something weird", Expected: "unknown"}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := getStatus(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestGetStatus - %s failed. %s", testName, msg)
		}
	}
}

func TestGetIssuedFor(t *testing.T) {
	var product product
	svsPath := "./test_data/sel-tor-watch.json"
	json.Unmarshal(ReadJSONFromFile(svsPath), &product)

	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: ""}
	tests["Valid issued for"] = TestParameters{Input: strings.ToLower(product.ProductText), Expected: "southern and central indiana, northern kentucky, western and central ohio"}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := getIssuedFor(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetIssuedFor - %s failed. %s", testName, msg)
		}
	}
}
