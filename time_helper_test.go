package main

import (
	"fmt"
	"testing"
)

func TestGetTimezoneOffset(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{input: "", expected: "0000"}
	tests["Default"] = TestParameters{input: "default", expected: "0000"}
	tests["HST"] = TestParameters{input: "HST", expected: "1000"}
	tests["HDT"] = TestParameters{input: "HDT", expected: "0900"}
	tests["AKST"] = TestParameters{input: "AKST", expected: "0900"}
	tests["AKDT"] = TestParameters{input: "AKDT", expected: "0800"}
	tests["PST"] = TestParameters{input: "PST", expected: "0800"}
	tests["PDT"] = TestParameters{input: "PDT", expected: "0700"}
	tests["MST"] = TestParameters{input: "MST", expected: "0700"}
	tests["MDT"] = TestParameters{input: "MDT", expected: "0600"}
	tests["CST"] = TestParameters{input: "CST", expected: "0600"}
	tests["CDT"] = TestParameters{input: "CDT", expected: "0500"}
	tests["EST"] = TestParameters{input: "EST", expected: "0500"}
	tests["EDT"] = TestParameters{input: "EDT", expected: "0400"}

	for testName, params := range tests {
		result := GetTimezoneOffset(params.input)

		if result != params.expected {
			msg := fmt.Sprintf("result: '%s', expected: '%s'", params.input, params.expected)
			t.Errorf("TestGetTimezoneOffset - %s failed. %s", testName, msg)
		}
	}
}
