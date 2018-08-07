package main

import (
	"fmt"
	"strings"
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

func TestGetIssuedFor(t *testing.T) {
	// Strings coming in are normalized to lowercase
	multipleLocationText := strings.ToLower("\n500 \nWGUS53 KGID 020152\nFFWGID\nKSC123-141-020545-\n/O.NEW.KGID.FF.W.0001.180502T0152Z-180502T0545Z/\n/00000.0.ER.000000T0000Z.000000T0000Z.000000T0000Z.OO/\n\nBULLETIN - EAS ACTIVATION REQUESTED\nFlash Flood Warning\nNational Weather Service Hastings NE\n852 PM CDT TUE MAY 1 2018\n\nThe National Weather Service in Hastings has issued a\n\n* Flash Flood Warning for...\n  Mitchell County in north central Kansas...\n  Southeastern Osborne County in north central Kansas...\n\n* Until 1245 AM CDT\n\n* At 844 PM CDT, Doppler radar indicated thunderstorms producing\n  heavy rain across the warned area. Flash flooding is expected to \n  begin shortly. Three to five inches of rain have been estimated to \n  have already fallen for some areas, with potentially another \n  couple of inches of rain before ending Tuesday night.\n\n* Some locations that will experience flooding include...\n  Beloit, Tipton, Asherville, Simpson, Hunter and Victor and along \n  the Solomon River. \n\nLAT...LON 3935 9847 3953 9793 3922 9793 3922 9849\n      3913 9849 3913 9889\n\n$$\n\nHeinlein\n\n")
	multipleLocationExpected := []string{"mitchell county in north central kansas", "southeastern osborne county in north central kansas"}
	singleLocationText := strings.ToLower("\n679 \nWFUS54 KJAN 140024\nTORJAN\nLAC067-140130-\n/O.NEW.KJAN.TO.W.0027.180414T0024Z-180414T0130Z/\n\nBULLETIN - EAS ACTIVATION REQUESTED\nTornado Warning\nNational Weather Service Jackson MS\n724 PM CDT FRI APR 13 2018\n\nThe National Weather Service in Jackson has issued a\n\n* Tornado Warning for...\n  Northwestern Morehouse Parish in northeastern Louisiana...\n\n* Until 830 PM CDT\n         \n* At 724 PM CDT, a confirmed tornado was located 12 miles north of\n  Sterlington, or 13 miles south of Huttig, moving northeast at 30\n  mph.\n\n  HAZARD...Damaging tornado and quarter size hail. \n\n  SOURCE...Radar confirmed tornado. \n\n  IMPACT...Flying debris will be dangerous to those caught without \n           shelter. Mobile homes will be damaged or destroyed. \n           Damage to roofs, windows, and vehicles will occur.  Tree \n           damage is likely. \n\n* This tornadic thunderstorm will remain over mainly rural areas of\n  northwestern Morehouse Parish.\n\nPRECAUTIONARY/PREPAREDNESS ACTIONS...\n\nTo repeat, a tornado is on the ground. TAKE COVER NOW! Move to a\nbasement or an interior room on the lowest floor of a sturdy\nbuilding. Avoid windows. If you are outdoors, in a mobile home, or in\na vehicle, move to the closest substantial shelter and protect\nyourself from flying debris.\n\n&&\n\nLAT...LON 3298 9208 3299 9207 3301 9207 3301 9184\n      3282 9206 3283 9207 3284 9206 3285 9206\n      3287 9208 3293 9207 3296 9208\nTIME...MOT...LOC 0024Z 213DEG 24KT 3286 9212 \n\nTORNADO...OBSERVED\nHAIL...1.00IN\n\n$$\n\n19\n\n")
	singleLocationExpected := []string{"northwestern morehouse parish in northeastern louisiana"}

	tests := map[string]TestParameters{}
	tests["multiple locations"] = TestParameters{Input: multipleLocationText, Expected: multipleLocationExpected}
	tests["single location"] = TestParameters{Input: singleLocationText, Expected: singleLocationExpected}
	tests[""] = TestParameters{Input: "", Expected: []string{}}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := getIssuedFor(str)

		if !CompareObjects(result, params.Expected) {
			msg := fmt.Sprintf("result: '%v', Expected: '%v'", result, params.Expected)
			t.Errorf("TestGetIssuedFor - %s failed. %s", testName, msg)
		}
	}
}
