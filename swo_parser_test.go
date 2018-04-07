package main

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestBuildSWOEvent_Outlook(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day1-no-severe.json"
	json.Unmarshal(ReadJSONFromFile(dataPath), &product)

	expected := WxEvent{
		Details: outlookDetails{
			ProductCode: "swo",
			ProductType: "dy1",
			Valid:       "20Z",
			Risk:        "no_severe",
			Summary:     "Thunderstorms are possible from southern Oklahoma across the Ozarks region and over parts of the Florida Peninsula.",
			Forecaster:  "Darrow",
		},
	}

	result, err := buildSWOEvent(product)

	if err != nil || !CompareObjects(result, expected) {
		t.Error("TestBuildSWOEvent_Outlook failed")
	}
}

func TestBuildSWOEvent_MD(t *testing.T) {
	var product Product
	dataPath := "./data/swo-md.json"
	json.Unmarshal(ReadJSONFromFile(dataPath), &product)

	expected := `{"Source":"","Details":{"ProductCode":"swo","ProductType":"mcd","Number":"0190","Affected":"West central through north central Mississippi and adjacent portions of Arkansas/Louisiana","Concerning":"Concerning...Tornado Watch 23...","WatchInfo":"The severe weather threat for Tornado Watch 23 continues.","Valid":"2018-03-28T22:32:00Z","Expires":"2018-03-29T00:30:00Z","WFOs":["meg","jan"],"Summary":"A risk for thunderstorm activity capable of producing damaging wind gusts and a couple of tornadoes will gradually spread across and northeast of the Vicksburg MS area, toward Greenwood and Tupelo, through 7-9 PM CDT.","Forecaster":"Kerr","ImageURI":"http://www.spc.noaa.gov/products/md/2018/mcd0190.gif","Polygon":[{"Lat":33.18,"Lon":-90.84},{"Lat":34.13,"Lon":-90.08},{"Lat":34.49,"Lon":-89.33},{"Lat":34.07,"Lon":-88.56},{"Lat":32.91,"Lon":-89.41},{"Lat":32.2,"Lon":-90.65},{"Lat":31.66,"Lon":-91.55},{"Lat":31.71,"Lon":-91.86},{"Lat":32.45,"Lon":-91.21},{"Lat":33.18,"Lon":-100.84}]},"Ingested":"0001-01-01T00:00:00Z","Description":""}`

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

	if err == nil || !CompareObjects(result, expected) {
		t.Error("TestBuildSWOEvent_Short_Body failed")
	}
}

func TestParseSWODY_Day1_No_Severe(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day1-no-severe.json"
	json.Unmarshal(ReadJSONFromFile(dataPath), &product)
	result := parseSWODY(product)

	expected := outlookDetails{
		ProductCode: "swo",
		ProductType: "dy1",
		Valid:       "20Z",
		Risk:        "no_severe",
		Summary:     "Thunderstorms are possible from southern Oklahoma across the Ozarks region and over parts of the Florida Peninsula.",
		Forecaster:  "Darrow",
	}

	if !CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Day1_No_Severe failed")
	}
}

func TestParseSWODY_Day2_No_Severe(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day2-no-severe.json"
	json.Unmarshal(ReadJSONFromFile(dataPath), &product)
	result := parseSWODY(product)

	expected := outlookDetails{
		ProductCode: "swo",
		ProductType: "dy2",
		Valid:       "12Z",
		Risk:        "no_severe",
		Summary:     "Isolated thunderstorms may develop across parts of eastern Oklahoma to the Ozark Plateau and lower Tennessee Valley Sunday into Sunday night.  Other storms may develop across the southern and central Florida Peninsula on Sunday.",
		Forecaster:  "Darrow",
	}

	if !CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Day2_No_Severe failed")
	}
}

func TestParseSWODY_Day3_No_Severe(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day3-no-severe.json"
	json.Unmarshal(ReadJSONFromFile(dataPath), &product)
	result := parseSWODY(product)

	expected := outlookDetails{
		ProductCode: "swo",
		ProductType: "dy3",
		Valid:       "",
		Risk:        "no_severe",
		Summary:     "Isolated thunderstorms are possible across portions of the southern Plains and Ozark Plateau as well as across the southern Florida Peninsula.",
		Forecaster:  "Mosier",
	}

	if !CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Day3_No_Severe failed")
	}
}

func TestParseSWODY_Day48(t *testing.T) {
	var product Product
	dataPath := "./data/swo-day48.json"
	json.Unmarshal(ReadJSONFromFile(dataPath), &product)
	result := parseSWODY(product)

	expected := outlookDetails{
		ProductCode: "swo",
		ProductType: "d48",
		Valid:       "",
		Risk:        "unknown",
		Summary:     "",
		Forecaster:  "Peters",
	}

	if !CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Day3_No_Severe failed")
	}
}

func TestParseSWODY_Unknown_Day(t *testing.T) {
	product := Product{WmoCollectiveID: "stuff"}
	result := parseSWODY(product)

	expected := outlookDetails{
		ProductCode: "swo",
		ProductType: "",
		Valid:       "",
		Risk:        "unknown",
		Summary:     "",
		Forecaster:  "",
	}

	if !CompareObjects(result, expected) {
		t.Error("TestParseSWODY_Unknown_Day failed")
	}
}

func TestParseSWOMCD(t *testing.T) {
	var product Product
	responsePath := "./data/swo-md.json"
	json.Unmarshal(ReadJSONFromFile(responsePath), &product)
	expected := `{"ProductCode":"swo","ProductType":"mcd","Number":"0190","Affected":"West central through north central Mississippi and adjacent portions of Arkansas/Louisiana","Concerning":"Concerning...Tornado Watch 23...","WatchInfo":"The severe weather threat for Tornado Watch 23 continues.","Valid":"2018-03-28T22:32:00Z","Expires":"2018-03-29T00:30:00Z","WFOs":["meg","jan"],"Summary":"A risk for thunderstorm activity capable of producing damaging wind gusts and a couple of tornadoes will gradually spread across and northeast of the Vicksburg MS area, toward Greenwood and Tupelo, through 7-9 PM CDT.","Forecaster":"Kerr","ImageURI":"http://www.spc.noaa.gov/products/md/2018/mcd0190.gif","Polygon":[{"Lat":33.18,"Lon":-90.84},{"Lat":34.13,"Lon":-90.08},{"Lat":34.49,"Lon":-89.33},{"Lat":34.07,"Lon":-88.56},{"Lat":32.91,"Lon":-89.41},{"Lat":32.2,"Lon":-90.65},{"Lat":31.66,"Lon":-91.55},{"Lat":31.71,"Lon":-91.86},{"Lat":32.45,"Lon":-91.21},{"Lat":33.18,"Lon":-100.84}]}`

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

	if !CompareObjects(result, expected) {
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

	if !CompareObjects(result, expected) {
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

	if !CompareObjects(result, expected) {
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

	if !CompareObjects(result, expected) {
		t.Error("TestGetValidRange_Start_No_Regex_Match failed")
	}
}

func TestGetRisk(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{input: "", expected: "unknown"}
	tests["Invalid Text"] = TestParameters{input: "invalid text", expected: "unknown"}
	tests["High Risk"] = TestParameters{input: "...THERE IS A HIGH RISK ...THERE IS A MODERATE RISK", expected: "high"}
	tests["Moderate Risk"] = TestParameters{input: "...THERE IS A MODERATE RISK ...THERE IS AN ENHANCED RISK", expected: "moderate"}
	tests["Enhanced Risk"] = TestParameters{input: "...THERE IS AN ENHANCED RISK ...THERE IS A SLIGHT RISK", expected: "enhanced"}
	tests["Slight Risk"] = TestParameters{input: "...THERE IS A SLIGHT RISK ...THERE IS A MARGINAL RISK", expected: "slight"}
	tests["Marginal Risk"] = TestParameters{input: "...THERE IS A MARGINAL RISK ...NO SEVERE THUNDERSTORM AREAS FORECAST...", expected: "marginal"}
	tests["No Severe"] = TestParameters{input: "...NO SEVERE THUNDERSTORM AREAS FORECAST...", expected: "no_severe"}

	for testName, params := range tests {
		result := getRisk(params.input)

		if result != params.expected {
			msg := fmt.Sprintf("result: '%s', expected: '%s'", params.input, params.expected)
			t.Errorf("TestGetRisk - %s failed. %s", testName, msg)
		}
	}

}
