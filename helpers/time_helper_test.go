package helpers

import (
	"fmt"
	"testing"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestGetTimezoneOffset(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["Empty String"] = helpers.TestParameters{Input: "", Expected: "0000"}
	tests["Default"] = helpers.TestParameters{Input: "default", Expected: "0000"}
	tests["HST"] = helpers.TestParameters{Input: "HST", Expected: "1000"}
	tests["HDT"] = helpers.TestParameters{Input: "HDT", Expected: "0900"}
	tests["AKST"] = helpers.TestParameters{Input: "AKST", Expected: "0900"}
	tests["AKDT"] = helpers.TestParameters{Input: "AKDT", Expected: "0800"}
	tests["PST"] = helpers.TestParameters{Input: "PST", Expected: "0800"}
	tests["PDT"] = helpers.TestParameters{Input: "PDT", Expected: "0700"}
	tests["MST"] = helpers.TestParameters{Input: "MST", Expected: "0700"}
	tests["MDT"] = helpers.TestParameters{Input: "MDT", Expected: "0600"}
	tests["CST"] = helpers.TestParameters{Input: "CST", Expected: "0600"}
	tests["CDT"] = helpers.TestParameters{Input: "CDT", Expected: "0500"}
	tests["EST"] = helpers.TestParameters{Input: "EST", Expected: "0500"}
	tests["EDT"] = helpers.TestParameters{Input: "EDT", Expected: "0400"}

	for testName, params := range tests {
		result := GetTimezoneOffset(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestGetTimezoneOffset - %s failed. %s", testName, msg)
		}
	}
}
