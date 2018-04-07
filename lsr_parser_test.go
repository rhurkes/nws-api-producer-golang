package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestProcessLSRProduct_Valid_LSR(t *testing.T) {
	var product Product
	lsrHailRemarksPath := "./data/lsr-hail-remarks.json"
	json.Unmarshal(ReadJSONFromFile(lsrHailRemarksPath), &product)
	reported, _ := time.Parse(time.RFC3339, "2018-03-26T19:55:00-05:00")

	expectedDetails := LSRDetails{
		Type:     "hail",
		Reported: reported,
		Location: "1 e silver",
		Lat:      32.07,
		Lon:      -100.66,
		Mag:      "e1.25 inch",
		County:   "coke",
		State:    "tx",
		Source:   "storm chaser",
		Remarks:  "1.25 hail on hwy 208 near silver",
	}

	expected := WxEvent{Details: expectedDetails}

	result, err := processLSRProduct(product)
	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestProcessLSRProduct_Valid_LSR failed")
	}
}

func TestProcessLSRProduct_Empty_ProductText(t *testing.T) {
	product := Product{ProductText: ""}

	_, err := processLSRProduct(product)
	if err == nil {
		t.Error("TestProcessLSRProduct_Empty_ProductText failed")
	}
}

func TestProcessLSRProduct_Summary(t *testing.T) {
	product := Product{ProductText: "0\n1\n2\n3\n4\nSUMMARY\n\n\n\n\n\n\n\n\n\n"}

	_, err := processLSRProduct(product)
	if err == nil {
		t.Error("TestProcessLSRProduct_Summary failed")
	}
}

func TestProcessLSRProduct_No_Remarks_Flag(t *testing.T) {
	var product Product
	lsrHailRemarksPath := "./data/lsr-hail-remarks.json"
	json.Unmarshal(ReadJSONFromFile(lsrHailRemarksPath), &product)
	product.ProductText = strings.Replace(product.ProductText, "..REMARKS..", "", 1)

	_, err := processLSRProduct(product)
	if err == nil {
		t.Error("TestProcessLSRProduct_No_Remarks_Flag failed")
	}
}

func TestProcessLSRProduct_Older_than_Threshold(t *testing.T) {
	var product Product
	lsrHailRemarksPath := "./data/lsr-hail-remarks.json"
	json.Unmarshal(ReadJSONFromFile(lsrHailRemarksPath), &product)
	product.ProductText = strings.Replace(product.ProductText, "0755 PM", "0655 PM", 1)

	_, err := processLSRProduct(product)
	if err == nil {
		t.Error("TestProcessLSRProduct_Older_than_Threshold failed")
	}
}

func TestProcessLSRProduct_Invalid_Reported_Time(t *testing.T) {
	var product Product
	lsrHailRemarksPath := "./data/lsr-hail-remarks.json"
	json.Unmarshal(ReadJSONFromFile(lsrHailRemarksPath), &product)
	product.ProductText = strings.Replace(product.ProductText, "0755 PM", "garbage", 1)

	_, err := processLSRProduct(product)
	if err == nil {
		t.Error("TestProcessLSRProduct_Invalid_Reported_Time failed")
	}
}

func TestGetLSRTimezoneOffset(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{input: "", expected: "0000"}
	tests["Valid Timezone"] = TestParameters{input: "816 PM CDT MON MAR 26 2018", expected: "0500"}

	for testName, params := range tests {
		result := getLSRTimezoneOffset(params.input)

		if result != params.expected {
			msg := fmt.Sprintf("result: '%v', expected: '%v'", result, params.expected)
			t.Errorf("TestGetLSRTimezoneOffset - %s failed. %s", testName, msg)
		}
	}
}
