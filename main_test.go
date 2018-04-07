package main

import (
	"fmt"
	"testing"
)

func TestNormalizeString(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{input: "", expected: ""}
	tests["Uppercase and Spaces"] = TestParameters{input: " STUFF ", expected: "stuff"}

	for testName, params := range tests {
		result := normalizeString(params.input)

		if result != params.expected {
			msg := fmt.Sprintf("result: '%s', expected: '%s'", params.input, params.expected)
			t.Errorf("TestNormalizeString - %s failed. %s", testName, msg)
		}
	}
}

func TestNormalizeFloat(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{input: "", expected: float32(0)}
	tests["Invalid String"] = TestParameters{input: "stuff", expected: float32(0)}
	tests["0"] = TestParameters{input: "0", expected: float32(0)}
	tests["-1"] = TestParameters{input: "-1", expected: float32(-1)}
	tests["100"] = TestParameters{input: "100", expected: float32(100)}

	for testName, params := range tests {
		result := normalizeFloat(params.input)

		if result != params.expected {
			msg := fmt.Sprintf("result: '%v', expected: '%v'", result, params.expected)
			t.Errorf("TestNormalizeFloat - %s failed. %s", testName, msg)
		}
	}
}
