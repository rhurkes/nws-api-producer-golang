package main

import (
	"encoding/json"
	"io/ioutil"
)

// TestParameters are used for the input value and the expected value in parameterized tests.
type TestParameters struct {
	Input    interface{}
	Expected interface{}
}

// ReadJSONFromFile reads JSON from a file
func ReadJSONFromFile(filepath string) []byte {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Error(err.Error())
	}

	return raw
}

// CompareObjects compares two objects, outputs the difference to stdout
// and returns a bool signifying whether the objects were equal.
func CompareObjects(result interface{}, expected interface{}) bool {
	resultVal, _ := json.Marshal(result)
	expectedVal, _ := json.Marshal(expected)

	if string(expectedVal) != string(resultVal) {
		logger.Infof("result: %s\n", string(resultVal))
		logger.Infof("expected: %s\n", string(expectedVal))

		return false
	}

	return true
}
