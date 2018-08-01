package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestBuildLSREvent_Valid_LSR(t *testing.T) {
	var product product
	lsrHailRemarksPath := "./test_data/lsr-hail-remarks.json"
	json.Unmarshal(ReadJSONFromFile(lsrHailRemarksPath), &product)
	reported, _ := time.Parse(time.RFC3339, "2018-03-26T19:55:00-05:00")

	expectedDetails := lsrDetails{
		Type:        "hail",
		Reported:    reported,
		Location:    "1 e silver",
		Lat:         32.07,
		Lon:         -100.66,
		MagMeasured: false,
		MagValue:    1.25,
		MagUnits:    "inch",
		County:      "coke",
		State:       "tx",
		Source:      "storm chaser",
		Remarks:     "1.25 hail on hwy 208 near silver",
	}

	expected := wxEvent{Data: nwsData{Derived: expectedDetails}}

	result, err := buildLSREvent(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestBuildLSREvent_Valid_LSR failed")
	}
}

func TestBuildLSREvent_Empty_ProductText(t *testing.T) {
	product := product{ProductText: ""}

	_, err := buildLSREvent(product)
	if err == nil {
		t.Error("TestBuildLSREvent_Empty_ProductText failed")
	}
}

func TestBuildLSREvent_Summary(t *testing.T) {
	product := product{ProductText: "0\n1\n2\n3\n4\nSUMMARY\n\n\n\n\n\n\n\n\n\n"}

	_, err := buildLSREvent(product)
	if err == nil {
		t.Error("TestBuildLSREvent_Summary failed")
	}
}

func TestBuildLSREvent_No_Remarks_Flag(t *testing.T) {
	var product product
	lsrHailRemarksPath := "./test_data/lsr-hail-remarks.json"
	json.Unmarshal(ReadJSONFromFile(lsrHailRemarksPath), &product)
	product.ProductText = strings.Replace(product.ProductText, "..REMARKS..", "", 1)

	_, err := buildLSREvent(product)
	if err == nil {
		t.Error("TestBuildLSREvent_No_Remarks_Flag failed")
	}
}

func TestBuildLSREvent_Older_than_Threshold(t *testing.T) {
	var product product
	lsrHailRemarksPath := "./test_data/lsr-hail-remarks.json"
	json.Unmarshal(ReadJSONFromFile(lsrHailRemarksPath), &product)
	product.ProductText = strings.Replace(product.ProductText, "0755 PM", "0655 PM", 1)

	result, err := buildLSREvent(product)
	if !result.DoNotPublish || err != nil {
		t.Error("TestBuildLSREvent_Older_than_Threshold failed")
	}
}

func TestBuildLSREvent_Invalid_Reported_Time(t *testing.T) {
	var product product
	lsrHailRemarksPath := "./test_data/lsr-hail-remarks.json"
	json.Unmarshal(ReadJSONFromFile(lsrHailRemarksPath), &product)
	product.ProductText = strings.Replace(product.ProductText, "0755 PM", "garbage", 1)

	result, err := buildLSREvent(product)
	if !result.DoNotPublish || err != nil {
		t.Error("TestBuildLSREvent_Invalid_Reported_Time failed")
	}
}

func TestGetLSRTimezoneOffset(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: "0000"}
	tests["Valid Timezone"] = TestParameters{Input: "816 PM CDT MON MAR 26 2018", Expected: "0500"}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := getLSRTimezoneOffset(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%v'", result, params.Expected)
			t.Errorf("TestGetLSRTimezoneOffset - %s failed. %s", testName, msg)
		}
	}
}

func TestGetMagnitude(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: magnitude{}}
	tests["Wind Damage"] = TestParameters{Input: "                 ", Expected: magnitude{}}
	tests["Unparsable String"] = TestParameters{Input: "UNPARSABLE", Expected: magnitude{}}
	tests["Estimated Hail"] = TestParameters{Input: "E4.50 INCH", Expected: magnitude{Measured: false, Value: 4.5, Units: "inch"}}
	tests["Measured Hail"] = TestParameters{Input: "M0.75 INCH", Expected: magnitude{Measured: true, Value: .75, Units: "inch"}}
	tests["Estimated Wind"] = TestParameters{Input: "E88 MPH", Expected: magnitude{Measured: false, Value: 88, Units: "mph"}}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := getMagnitude(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%v'", result, params.Expected)
			t.Errorf("TestGetMagnitude - %s failed. %s", testName, msg)
		}
	}
}
