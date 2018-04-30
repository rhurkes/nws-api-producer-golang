package main

import (
	"fmt"
	"testing"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestGetLatFromString(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["3267"] = helpers.TestParameters{Input: "3267", Expected: float32(32.67)}
	tests["0967"] = helpers.TestParameters{Input: "0967", Expected: float32(9.67)}
	tests[""] = helpers.TestParameters{Input: "", Expected: float32(0)}

	for testName, params := range tests {
		result := getLatFromString(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetLatFromString - %s failed. %s", testName, msg)
		}
	}
}

func TestGetLonFromString(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["9967"] = helpers.TestParameters{Input: "9967", Expected: float32(-99.67)}
	tests["0067"] = helpers.TestParameters{Input: "0067", Expected: float32(-100.67)}
	tests[""] = helpers.TestParameters{Input: "", Expected: float32(0)}

	for testName, params := range tests {
		result := getLonFromString(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetLonFromString - %s failed. %s", testName, msg)
		}
	}
}
