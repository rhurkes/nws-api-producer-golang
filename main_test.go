package main

import (
	"fmt"
	"testing"
)

type ProcessFeatureParameter struct {
	productType  nwsProduct
	responseBody string
}

func TestBuildProductListURI(t *testing.T) {
	tests := map[string]TestParameters{}
	tests["Empty String"] = TestParameters{Input: "", Expected: "https://api.weather.gov/products/types/"}
	tests["Product"] = TestParameters{Input: "lsr", Expected: "https://api.weather.gov/products/types/lsr"}

	for testName, params := range tests {
		str, _ := params.Input.(string)
		result := buildProductListURI(str)

		if result != params.Expected {
			msg := fmt.Sprintf("result: '%v', Expected: '%v'", result, params.Expected)
			t.Errorf("TestBuildProductListURI - %s failed. %s", testName, msg)
		}
	}
}

func TestProcessFeature(t *testing.T) {
	tests := map[string]ProcessFeatureParameter{}
	tests["AFD Success"] = ProcessFeatureParameter{productType: AreaForecastDiscussion, responseBody: ""}
	tests["SEL Success"] = ProcessFeatureParameter{productType: SevereWatch, responseBody: ""}
	tests["SVR Success"] = ProcessFeatureParameter{productType: SevereThunderstormWarning, responseBody: ""}
	tests["TOR Success"] = ProcessFeatureParameter{productType: TornadoWarning, responseBody: ""}
	tests["LSR Failure"] = ProcessFeatureParameter{productType: LocalStormReport, responseBody: ""}
	tests["SWO Failure"] = ProcessFeatureParameter{productType: StormOutlookNarrative, responseBody: ""}
	tests["SVS Failure"] = ProcessFeatureParameter{productType: SevereWeatherStatement, responseBody: ""}

	for _, params := range tests {
		processFeature(params.productType, []byte(params.responseBody))
	}
}
