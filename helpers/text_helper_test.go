package helpers

import (
	"fmt"
	"testing"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestNormalizeString(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: ""}
	tests["Uppercase and Spaces"] = TestParameters{Input: " STUFF ", Expected: "stuff"}

	for testName, params := range tests {
		result := normalizeString(params.Input, false)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestNormalizeString - %s failed. %s", testName, msg)
		}
	}
}

func TestNormalizeStringPreserveCase(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Uppercase"] = TestParameters{Input: "STUFF", Expected: "STUFF"}
	tests["Lowercase"] = TestParameters{Input: "stuff", Expected: "stuff"}

	for testName, params := range tests {
		result := normalizeString(params.Input, true)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestNormalizeString - %s failed. %s", testName, msg)
		}
	}
}

func TestNormalizeFloat(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: float32(0)}
	tests["Invalid String"] = TestParameters{Input: "stuff", Expected: float32(0)}
	tests["0"] = TestParameters{Input: "0", Expected: float32(0)}
	tests["-1"] = TestParameters{Input: "-1", Expected: float32(-1)}
	tests["100"] = TestParameters{Input: "100", Expected: float32(100)}

	for testName, params := range tests {
		result := normalizeFloat(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%v'", result, params.Expected)
			t.Errorf("TestNormalizeFloat - %s failed. %s", testName, msg)
		}
	}
}
