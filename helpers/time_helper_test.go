package helpers

import (
	"fmt"
	"testing"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestGetTimezoneOffset(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: "0000"}
	tests["Default"] = TestParameters{Input: "default", Expected: "0000"}
	tests["HST"] = TestParameters{Input: "HST", Expected: "1000"}
	tests["HDT"] = TestParameters{Input: "HDT", Expected: "0900"}
	tests["AKST"] = TestParameters{Input: "AKST", Expected: "0900"}
	tests["AKDT"] = TestParameters{Input: "AKDT", Expected: "0800"}
	tests["PST"] = TestParameters{Input: "PST", Expected: "0800"}
	tests["PDT"] = TestParameters{Input: "PDT", Expected: "0700"}
	tests["MST"] = TestParameters{Input: "MST", Expected: "0700"}
	tests["MDT"] = TestParameters{Input: "MDT", Expected: "0600"}
	tests["CST"] = TestParameters{Input: "CST", Expected: "0600"}
	tests["CDT"] = TestParameters{Input: "CDT", Expected: "0500"}
	tests["EST"] = TestParameters{Input: "EST", Expected: "0500"}
	tests["EDT"] = TestParameters{Input: "EDT", Expected: "0400"}

	for testName, params := range tests {
		result := GetTimezoneOffset(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestGetTimezoneOffset - %s failed. %s", testName, msg)
		}
	}
}
