package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/rhurkes/wxNwsProducer/helpers"
)

func TestBuildSELEventSVRWatch(t *testing.T) {
	var product Product
	svsPath := "./data/sel-svr-watch.json"
	json.Unmarshal(helpers.ReadJSONFromFile(svsPath), &product)
	expectedText := "\n727 \nWWUS20 KWNS 031523\nSEL5  \nSPC WW 031523\nARZ000-LAZ000-OKZ000-TXZ000-032300-\n\nURGENT - IMMEDIATE BROADCAST REQUESTED\nSevere Thunderstorm Watch Number 25\nNWS Storm Prediction Center Norman OK\n1025 AM CDT Tue Apr 3 2018\n\nThe NWS Storm Prediction Center has issued a\n\n* Severe Thunderstorm Watch for portions of \n  Southwest Arkansas\n  Northwest Louisiana\n  Southeast Oklahoma\n  Central and Northeast Texas\n\n* Effective this Tuesday morning and evening from 1025 AM until\n  600 PM CDT.\n\n* Primary threats include...\n  Scattered large hail likely with isolated very large hail events\n    to 2.5 inches in diameter possible\n  Scattered damaging wind gusts to 70 mph possible\n\nSUMMARY...Thunderstorms are intensifying over central Texas, and\nwill spread northeastward across the watch area through the\nafternoon.  Other storms will form along an approaching cold front. \nLarge hail and damaging winds will be possible in the strongest\ncells.\n\nThe severe thunderstorm watch area is approximately along and 75\nstatute miles north and south of a line from 50 miles west of Temple\nTX to 40 miles northeast of Shreveport LA. For a complete depiction\nof the watch see the associated watch outline update (WOUS64 KWNS\nWOU5).\n\nPRECAUTIONARY/PREPAREDNESS ACTIONS...\n\nREMEMBER...A Severe Thunderstorm Watch means conditions are\nfavorable for severe thunderstorms in and close to the watch area.\nPersons in these areas should be on the lookout for threatening\nweather conditions and listen for later statements and possible\nwarnings. Severe thunderstorms can and occasionally do produce\ntornadoes.\n\n\u0026\u0026\n\nAVIATION...A few severe thunderstorms with hail surface and aloft to\n2.5 inches. Extreme turbulence and surface wind gusts to 60 knots. A\nfew cumulonimbi with maximum tops to 500. Mean storm motion vector\n26030.\n\n...Hart\n\n"

	expectedDetails := selDetails{
		Code:        "sel",
		Issued:      1522768980,
		Name:        "Severe Local Storm Watch and Watch Cancellation Msg",
		Wfo:         "KWNS",
		IsPDS:       false,
		WatchNumber: 25,
		WatchType:   "severe thunderstorm",
		Status:      "issued",
		Text:        expectedText,
		IssuedFor:   "southwest arkansas, northwest louisiana, southeast oklahoma, central and northeast texas",
	}

	expected := WxEvent{Details: expectedDetails}

	result, err := buildSELEvent(product)
	if err != nil || !helpers.CompareObjects(result, expected) {
		t.Error("TestBuildSELEventSVRWatch failed")
	}
}

func TestBuildSELEventTORWatch(t *testing.T) {
	var product Product
	svsPath := "./data/sel-tor-watch.json"
	json.Unmarshal(helpers.ReadJSONFromFile(svsPath), &product)
	expectedText := "\n281 \nWWUS20 KWNS 031713\nSEL6  \nSPC WW 031713\nINZ000-KYZ000-OHZ000-040000-\n\nURGENT - IMMEDIATE BROADCAST REQUESTED\nTornado Watch Number 26\nNWS Storm Prediction Center Norman OK\n115 PM EDT Tue Apr 3 2018\n\nThe NWS Storm Prediction Center has issued a\n\n* Tornado Watch for portions of \n  Southern and Central Indiana\n  Northern Kentucky\n  Western and Central Ohio\n\n* Effective this Tuesday afternoon and evening from 115 PM until\n  800 PM EDT.\n\n* Primary threats include...\n  A few tornadoes likely with a couple intense tornadoes possible\n  Scattered damaging wind gusts to 70 mph likely\n  Scattered large hail and isolated very large hail events to 2\n    inches in diameter possible\n\nSUMMARY...Thunderstorms are intensifying along the IL/IN border, and\nwill track eastward across the watch area through the afternoon. \nConditions appear favorable for supercell storms capable of large\nhail, damaging winds, and perhaps a strong tornado or two.\n\nThe tornado watch area is approximately along and 70 statute miles\nnorth and south of a line from 40 miles south southwest of Terre\nHaute IN to 20 miles south southeast of Columbus OH. For a complete\ndepiction of the watch see the associated watch outline update\n(WOUS64 KWNS WOU6).\n\nPRECAUTIONARY/PREPAREDNESS ACTIONS...\n\nREMEMBER...A Tornado Watch means conditions are favorable for\ntornadoes and severe thunderstorms in and close to the watch\narea. Persons in these areas should be on the lookout for\nthreatening weather conditions and listen for later statements\nand possible warnings.\n\n\u0026\u0026\n\nOTHER WATCH INFORMATION...CONTINUE...WW 25...\n\nAVIATION...Tornadoes and a few severe thunderstorms with hail\nsurface and aloft to 2 inches. Extreme turbulence and surface wind\ngusts to 60 knots. A few cumulonimbi with maximum tops to 450. Mean\nstorm motion vector 24035.\n\n...Hart\n\n"

	expectedDetails := selDetails{
		Code:        "sel",
		Issued:      1522775580,
		Name:        "Severe Local Storm Watch and Watch Cancellation Msg",
		Wfo:         "KWNS",
		IsPDS:       false,
		WatchNumber: 26,
		WatchType:   "tornado",
		Status:      "issued",
		Text:        expectedText,
		IssuedFor:   "southern and central indiana, northern kentucky, western and central ohio",
	}

	expected := WxEvent{Details: expectedDetails}

	result, err := buildSELEvent(product)
	if err != nil || !helpers.CompareObjects(result, expected) {
		t.Error("TestBuildSELEventTORWatch failed")
	}
}

func TestGetWatchStats(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["Empty String"] = helpers.TestParameters{Input: "", Expected: watchStats{}}
	tests["SVR 5"] = helpers.TestParameters{Input: "\nsevere thunderstorm watch number 25\n", Expected: watchStats{Type: "severe thunderstorm", Number: 25}}

	for testName, params := range tests {
		result := getWatchStats(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestGetWatchStats - %s failed. %s", testName, msg)
		}
	}
}

func TestGetStatus(t *testing.T) {
	tests := map[string]helpers.TestParameters{}
	tests["Empty string"] = helpers.TestParameters{Input: "", Expected: "unknown"}
	tests["Cancelled"] = helpers.TestParameters{Input: "the nws storm prediction center has cancelled", Expected: "cancelled"}
	tests["Issued"] = helpers.TestParameters{Input: "the nws storm prediction center has issued", Expected: "issued"}
	tests["Not cancelled or issued"] = helpers.TestParameters{Input: "the nws storm prediction center has done something weird", Expected: "unknown"}

	for testName, params := range tests {
		result := getStatus(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", params.Input, params.Expected)
			t.Errorf("TestGetStatus - %s failed. %s", testName, msg)
		}
	}
}

func TestGetIssuedFor(t *testing.T) {
	var product Product
	svsPath := "./data/sel-tor-watch.json"
	json.Unmarshal(helpers.ReadJSONFromFile(svsPath), &product)

	tests := map[string]helpers.TestParameters{}
	tests["Empty String"] = helpers.TestParameters{Input: "", Expected: ""}
	tests["Valid issued for"] = helpers.TestParameters{Input: strings.ToLower(product.ProductText), Expected: "southern and central indiana, northern kentucky, western and central ohio"}

	for testName, params := range tests {
		result := getIssuedFor(params.Input)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%s', Expected: '%s'", result, params.Expected)
			t.Errorf("TestGetIssuedFor - %s failed. %s", testName, msg)
		}
	}
}
