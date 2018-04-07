package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// TestParameters are used for the input value and the expected value in parameterized tests.
type TestParameters struct {
	input    string
	expected interface{}
}

// ReadJSONFromFile reads JSON from a file
func ReadJSONFromFile(filepath string) []byte {
	raw, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err.Error())
	}

	return raw
}

// CompareObjects compares two objects, outputs the difference to stdout
// and returns a bool signifying whether the objects were equal.
func CompareObjects(result interface{}, expected interface{}) bool {
	resultVal, _ := json.Marshal(result)
	expectedVal, _ := json.Marshal(expected)

	if string(expectedVal) != string(resultVal) {
		fmt.Println(fmt.Sprintf("result: %s", string(resultVal)))
		fmt.Println(fmt.Sprintf("expected: %s", string(expectedVal)))

		return false
	}

	return true
}
