package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestBuildSWOEvent_Outlook(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day1-no-severe.json"
	json.Unmarshal(helpers.ReadJSONFromFile(dataPath), &product)

	expected := WxEvent{
		Details: outlookDetails{
			Code:       "swo",
			Issued:     1522524900,
			Name:       "Severe Storm Outlook Narrative (AC)",
			Wfo:        "KWNS",
			SubCode:    "dy1",
			Valid:      "20Z",
			Risk:       "no_severe",
			Summary:    "Thunderstorms are possible from southern Oklahoma across the Ozarks region and over parts of the Florida Peninsula.",
			Forecaster: "Darrow",
		},
	}

	result, err := buildSWOEvent(product)

	if err != nil || !helpers.CompareObjects(result, expected) {
		t.Error("TestBuildSWOEvent_Outlook failed")
	}
}

func TestBuildSWOEvent_MD(t *testing.T) {
	var product Product
	dataPath := "./data/swo-md-continues.json"
	json.Unmarshal(helpers.ReadJSONFromFile(dataPath), &product)

	expected := `{"Source":"","Details":{"Code":"swo","Issued":1522276380,"Name":"Severe Storm Outlook Narrative (AC)","Wfo":"KWNS","SubCode":"mcd","Number":"0190","Affected":"West central through north central Mississippi and adjacent portions of Arkansas/Louisiana","Concerning":"Concerning...Tornado Watch 23...","WatchInfo":"The severe weather threat for Tornado Watch 23 continues.","Valid":"2018-03-28T22:32:00Z","Expires":"2018-03-29T00:30:00Z","WFOs":["meg","jan"],"Summary":"A risk for thunderstorm activity capable of producing damaging wind gusts and a couple of tornadoes will gradually spread across and northeast of the Vicksburg MS area, toward Greenwood and Tupelo, through 7-9 PM CDT.","Forecaster":"Kerr","ImageURI":"http://www.spc.noaa.gov/products/md/2018/mcd0190.gif","Polygon":[{"Lat":33.18,"Lon":-90.84},{"Lat":34.13,"Lon":-90.08},{"Lat":34.49,"Lon":-89.33},{"Lat":34.07,"Lon":-88.56},{"Lat":32.91,"Lon":-89.41},{"Lat":32.2,"Lon":-90.65},{"Lat":31.66,"Lon":-91.55},{"Lat":31.71,"Lon":-91.86},{"Lat":32.45,"Lon":-91.21},{"Lat":33.18,"Lon":-100.84}]},"Ingested":"0001-01-01T00:00:00Z","Summary":""}`

	result, err := buildSWOEvent(product)
	marshalledResult, _ := json.Marshal(result)

	if err != nil || string(marshalledResult) != expected {
		fmt.Println("result: " + string(marshalledResult))
		fmt.Println("expected: " + expected)
		t.Error("Test failed")
	}
}

func TestBuildSWOEvent_Short_Body(t *testing.T) {
	var product Product
	expected := WxEvent{}

	result, err := buildSWOEvent(product)

	if err == nil || !helpers.CompareObjects(result, expected) {
		t.Error("TestBuildSWOEvent_Short_Body failed")
	}
}

func TestParseSWODY_Day1_No_Severe(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day1-no-severe.json"
	json.Unmarshal(helpers.ReadJSONFromFile(dataPath), &product)
	result := parseSWODY(product)

	expected := outlookDetails{
		Code:       "swo",
		Issued:     1522524900,
		Name:       "Severe Storm Outlook Narrative (AC)",
		Wfo:        "KWNS",
		SubCode:    "dy1",
		Valid:      "20Z",
		Risk:       "no_severe",
		Summary:    "Thunderstorms are possible from southern Oklahoma across the Ozarks region and over parts of the Florida Peninsula.",
		Forecaster: "Darrow",
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Day1_No_Severe failed")
	}
}

func TestParseSWODY_Day2_No_Severe(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day2-no-severe.json"
	json.Unmarshal(helpers.ReadJSONFromFile(dataPath), &product)
	result := parseSWODY(product)

	expected := outlookDetails{
		Code:       "swo",
		Issued:     1522516440,
		Name:       "Severe Storm Outlook Narrative (AC)",
		Wfo:        "KWNS",
		SubCode:    "dy2",
		Valid:      "12Z",
		Risk:       "no_severe",
		Summary:    "Isolated thunderstorms may develop across parts of eastern Oklahoma to the Ozark Plateau and lower Tennessee Valley Sunday into Sunday night.  Other storms may develop across the southern and central Florida Peninsula on Sunday.",
		Forecaster: "Darrow",
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Day2_No_Severe failed")
	}
}

func TestParseSWODY_Day3_No_Severe(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day3-no-severe.json"
	json.Unmarshal(helpers.ReadJSONFromFile(dataPath), &product)
	result := parseSWODY(product)

	expected := outlookDetails{
		Code:       "swo",
		Issued:     1522307280,
		Name:       "Severe Storm Outlook Narrative (AC)",
		Wfo:        "KWNS",
		SubCode:    "dy3",
		Valid:      "",
		Risk:       "no_severe",
		Summary:    "Isolated thunderstorms are possible across portions of the southern Plains and Ozark Plateau as well as across the southern Florida Peninsula.",
		Forecaster: "Mosier",
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Day3_No_Severe failed")
	}
}

func TestParseSWODY_Day48(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day48.json"
	json.Unmarshal(helpers.ReadJSONFromFile(dataPath), &product)
	result := parseSWODY(product)

	expected := outlookDetails{
		Code:       "swo",
		Issued:     1522487160,
		Name:       "Severe Storm Outlook Narrative (AC)",
		Wfo:        "KWNS",
		SubCode:    "d48",
		Valid:      "",
		Risk:       "unknown",
		Summary:    "",
		Forecaster: "Peters",
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Day3_No_Severe failed")
	}
}

func TestParseSWODY_Unknown_Day(t *testing.T) {
	product := Product{WmoCollectiveID: "stuff"}
	result := parseSWODY(product)

	expected := outlookDetails{
		Code:       "",
		Issued:     -62135596800,
		Name:       "",
		Wfo:        "",
		SubCode:    "",
		Valid:      "",
		Risk:       "unknown",
		Summary:    "",
		Forecaster: "",
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Unknown_Day failed")
	}
}

func TestParseSWOMCD(t *testing.T) {
	var product Product
	responsePath := "./data/swo-md-tor-watch-likely.json"
	json.Unmarshal(helpers.ReadJSONFromFile(responsePath), &product)
	expected := `{"Code":"swo","Issued":1522773660,"Name":"Severe Storm Outlook Narrative (AC)","Wfo":"KWNS","SubCode":"mcd","Number":"0205","Affected":"Central and southern Indiana...far northern Kentucky...western Ohio","Concerning":"Concerning...Severe potential...Tornado Watch likely ","WatchInfo":"Probability of Watch Issuance...95 percent","Valid":"2018-04-03T16:41:00Z","Expires":"2018-04-03T19:15:00Z","WFOs":["iln","lmk","iwx","ind","pah","ilx"],"Summary":"Storms are expected to increase in intensity this afternoon with a few tornadoes possible along with large hail. Additional severe storms are likely later this evening.","Forecaster":"Jewell/Hart","ImageURI":"http://www.spc.noaa.gov/products/md/2018/mcd0205.gif","Polygon":[{"Lat":37.82,"Lon":-87.69},{"Lat":38.53,"Lon":-87.76},{"Lat":39.73,"Lon":-87.06},{"Lat":40.62,"Lon":-85.25},{"Lat":40.46,"Lon":-83.56},{"Lat":40.36,"Lon":-83.1},{"Lat":40.12,"Lon":-82.74},{"Lat":39.65,"Lon":-82.75},{"Lat":39.24,"Lon":-83.39},{"Lat":38.8,"Lon":-84.23},{"Lat":38.2,"Lon":-85.03},{"Lat":37.81,"Lon":-85.97},{"Lat":37.73,"Lon":-86.3},{"Lat":37.64,"Lon":-87.21},{"Lat":37.82,"Lon":-87.69}]}`

	result := parseSWOMCD(product)
	marshalledResult, _ := json.Marshal(result)

	if string(marshalledResult) != expected {
		fmt.Println("result: " + string(marshalledResult))
		fmt.Println("expected: " + expected)
		t.Error("Test failed")
	}
}

func TestParseSWOMCD_Continues(t *testing.T) {
	var product Product
	responsePath := "./data/swo-md-continues.json"
	json.Unmarshal(helpers.ReadJSONFromFile(responsePath), &product)
	expected := `{"Code":"swo","Issued":1522276380,"Name":"Severe Storm Outlook Narrative (AC)","Wfo":"KWNS","SubCode":"mcd","Number":"0190","Affected":"West central through north central Mississippi and adjacent portions of Arkansas/Louisiana","Concerning":"Concerning...Tornado Watch 23...","WatchInfo":"The severe weather threat for Tornado Watch 23 continues.","Valid":"2018-03-28T22:32:00Z","Expires":"2018-03-29T00:30:00Z","WFOs":["meg","jan"],"Summary":"A risk for thunderstorm activity capable of producing damaging wind gusts and a couple of tornadoes will gradually spread across and northeast of the Vicksburg MS area, toward Greenwood and Tupelo, through 7-9 PM CDT.","Forecaster":"Kerr","ImageURI":"http://www.spc.noaa.gov/products/md/2018/mcd0190.gif","Polygon":[{"Lat":33.18,"Lon":-90.84},{"Lat":34.13,"Lon":-90.08},{"Lat":34.49,"Lon":-89.33},{"Lat":34.07,"Lon":-88.56},{"Lat":32.91,"Lon":-89.41},{"Lat":32.2,"Lon":-90.65},{"Lat":31.66,"Lon":-91.55},{"Lat":31.71,"Lon":-91.86},{"Lat":32.45,"Lon":-91.21},{"Lat":33.18,"Lon":-100.84}]}`

	result := parseSWOMCD(product)
	marshalledResult, _ := json.Marshal(result)

	if string(marshalledResult) != expected {
		fmt.Println("result: " + string(marshalledResult))
		fmt.Println("expected: " + expected)
		t.Error("Test failed")
	}
}

func TestGetValidRange_Until_On_Next_Day(t *testing.T) {
	text := "\n\nValid 282232Z - 290030Z\n\n"
	issued, _ := time.Parse(time.RFC3339, "2018-03-28T22:34:00-00:00")
	expectedStart, _ := time.Parse(time.RFC3339, "2018-03-28T22:32:00Z")
	expectedUntil, _ := time.Parse(time.RFC3339, "2018-03-29T00:30:00Z")
	expected := [2]time.Time{expectedStart, expectedUntil}

	result, err := getValidRange(text, issued)
	if err != nil {
		t.Error("TestGetValidRange_Until_On_Next_Day failed")
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestGetValidRange_Until_On_Next_Day failed")
	}
}

func TestGetValidRange_Same_Day(t *testing.T) {
	text := "\n\nValid 290032Z - 290130Z\n\n"
	issued, _ := time.Parse(time.RFC3339, "2018-03-29T00:34:00-00:00")
	expectedStart, _ := time.Parse(time.RFC3339, "2018-03-29T00:32:00Z")
	expectedUntil, _ := time.Parse(time.RFC3339, "2018-03-29T01:30:00Z")
	expected := [2]time.Time{expectedStart, expectedUntil}

	result, err := getValidRange(text, issued)
	if err != nil {
		t.Error("TestGetValidRange_Same_Day failed")
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestGetValidRange_Same_Day failed")
	}
}

func TestGetValidRange_Start_On_Previous_Day(t *testing.T) {
	text := "\n\nValid 282359Z - 290030Z\n\n"
	issued, _ := time.Parse(time.RFC3339, "2018-03-29T00:01:00-00:00")
	expectedStart, _ := time.Parse(time.RFC3339, "2018-03-28T23:59:00Z")
	expectedUntil, _ := time.Parse(time.RFC3339, "2018-03-29T00:30:00Z")
	expected := [2]time.Time{expectedStart, expectedUntil}

	result, err := getValidRange(text, issued)
	if err != nil {
		t.Error("TestGetValidRange_Start_On_Previous_Day failed")
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestGetValidRange_Start_On_Previous_Day failed")
	}
}

func TestGetValidRange_Start_No_Regex_Match(t *testing.T) {
	text := "\n\nValid 82359Z - 90030Z\n\n"
	issued, _ := time.Parse(time.RFC3339, "2018-03-29T00:01:00-00:00")
	expected := []time.Time{}

	result, err := getValidRange(text, issued)
	if err == nil {
		t.Error("TestGetValidRange_Start_No_Regex_Match failed")
	}

	if !helpers.CompareObjects(result, expected) {
		t.Error("TestGetValidRange_Start_No_Regex_Match failed")
	}
}

func TestGetRisk(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["Empty String"] = helpers.TestParameters{Input: "", Expected: "unknown"}
	tests["Invalid Text"] = helpers.TestParameters{Input: "invalid text", Expected: "unknown"}
	tests["High Risk"] = helpers.TestParameters{Input: "...THERE IS A HIGH RISK ...THERE IS A MODERATE RISK", Expected: "high"}
	tests["Moderate Risk"] = helpers.TestParameters{Input: "...THERE IS A MODERATE RISK ...THERE IS AN ENHANCED RISK", Expected: "moderate"}
	tests["Enhanced Risk"] = helpers.TestParameters{Input: "...THERE IS AN ENHANCED RISK ...THERE IS A SLIGHT RISK", Expected: "enhanced"}
	tests["Slight Risk"] = helpers.TestParameters{Input: "...THERE IS A SLIGHT RISK ...THERE IS A MARGINAL RISK", Expected: "slight"}
	tests["Marginal Risk"] = helpers.TestParameters{Input: "...THERE IS A MARGINAL RISK ...NO SEVERE THUNDERSTORM AREAS FORECAST...", Expected: "marginal"}
	tests["No Severe"] = helpers.TestParameters{Input: "...NO SEVERE THUNDERSTORM AREAS FORECAST...", Expected: "no_severe"}

	for testName, params := range tests {
		result := getRisk(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestGetRisk - %s failed. %s", testName, msg)
		}
	}

}
