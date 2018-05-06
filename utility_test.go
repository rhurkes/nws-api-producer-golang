package main

import (
	"fmt"
	"testing"
)

func TestGetNWSProductCode(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["AFD"] = TestParameters{Input: AreaForecastDiscussion, Expected: "afd"}
	tests["LSR"] = TestParameters{Input: LocalStormReport, Expected: "lsr"}
	tests["SEL"] = TestParameters{Input: SevereWatch, Expected: "sel"}
	tests["SVR"] = TestParameters{Input: SevereThunderstormWarning, Expected: "svr"}
	tests["SVS"] = TestParameters{Input: SevereWeatherStatement, Expected: "svs"}
	tests["SWO"] = TestParameters{Input: StormOutlookNarrative, Expected: "swo"}
	tests["TOR"] = TestParameters{Input: TornadoWarning, Expected: "tor"}

	for testName, params := range tests {
		nwsProduct, _ := params.Input.(nwsProduct)
		result := getNWSProductCode(nwsProduct)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%v'", result, params.Expected)
			t.Errorf("TestGetNWSProductCode - %s failed. %s", testName, msg)
		}
	}
}

func TestNormalizeString(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: ""}
	tests["Uppercase and Spaces"] = TestParameters{Input: " STUFF ", Expected: "stuff"}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := normalizeString(str, false)

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
		str, _ := params.Input.(string)
		result := normalizeString(str, true)

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
		str, _ := params.Input.(string)
		result := normalizeFloat(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%v'", result, params.Expected)
			t.Errorf("TestNormalizeFloat - %s failed. %s", testName, msg)
		}
	}
}

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
		str, _ := params.Input.(string)
		result := GetTimezoneOffset(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestGetTimezoneOffset - %s failed. %s", testName, msg)
		}
	}
}

func TestGetLatFromString(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["3267"] = TestParameters{Input: "3267", Expected: float32(32.67)}
	tests["0967"] = TestParameters{Input: "0967", Expected: float32(9.67)}
	tests[""] = TestParameters{Input: "", Expected: float32(0)}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := getLatFromString(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetLatFromString - %s failed. %s", testName, msg)
		}
	}
}

func TestGetLonFromString(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["9967"] = TestParameters{Input: "9967", Expected: float32(-99.67)}
	tests["0067"] = TestParameters{Input: "0067", Expected: float32(-100.67)}
	tests[""] = TestParameters{Input: "", Expected: float32(0)}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := getLonFromString(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetLonFromString - %s failed. %s", testName, msg)
		}
	}
}
