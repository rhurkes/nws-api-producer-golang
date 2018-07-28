package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
)

type clientMockSuccess struct{}
type clientMockDoError struct{}
type clientMockNon200 struct{}

func (c *clientMockSuccess) Do(req *http.Request) (*http.Response, error) {
	res := &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewReader([]byte("")))}
	return res, nil
}

func (c *clientMockDoError) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("do error")
}

func (c *clientMockNon200) Do(req *http.Request) (*http.Response, error) {
	res := &http.Response{StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(bytes.NewReader([]byte("")))}
	return res, nil
}

func TestProcessProductList_No_Features(t *testing.T) {
	lastSeenProduct[AreaForecastDiscussion] = ""
	productType := AreaForecastDiscussion
	features := []product{}

	processProductList(productType, features)
	if lastSeenProduct[AreaForecastDiscussion] != "" {
		t.Error("TestProcessProductList_No_Features failure")
	}
}

func TestProcessProductList_Features_Fresh_Load(t *testing.T) {
	lastSeenProduct[AreaForecastDiscussion] = ""
	productType := AreaForecastDiscussion
	firstProduct := product{ID: "id1"}
	secondProduct := product{ID: "id2"}
	features := []product{firstProduct, secondProduct}

	processProductList(productType, features)
	if lastSeenProduct[AreaForecastDiscussion] != "id1" {
		t.Error("TestProcessProductList_No_Features failure")
	}
}

func TestProcessProductList_Features_All_Seen(t *testing.T) {
	oldClient := client
	defer func() { client = oldClient }()
	client = &clientMockSuccess{}

	lastSeenProduct[AreaForecastDiscussion] = "id1"
	productType := AreaForecastDiscussion
	firstProduct := product{ID: "id1"}
	secondProduct := product{ID: "id2"}
	features := []product{firstProduct, secondProduct}

	processProductList(productType, features)
	if lastSeenProduct[AreaForecastDiscussion] != "id1" {
		t.Error("TestProcessProductList_No_Features failure")
	}
}

func TestProcessProductList_Features_None_Seen(t *testing.T) {
	oldClient := client
	defer func() { client = oldClient }()
	client = &clientMockSuccess{}

	lastSeenProduct[AreaForecastDiscussion] = "id3"
	productType := AreaForecastDiscussion
	firstProduct := product{ID: "id1"}
	secondProduct := product{ID: "id2"}
	features := []product{firstProduct, secondProduct}

	processProductList(productType, features)
	if lastSeenProduct[AreaForecastDiscussion] != "id1" {
		t.Error("TestProcessProductList_No_Features failure")
	}
}

func TestGetProductList(t *testing.T) {
	oldClient := client
	defer func() { client = oldClient }()
	client = &clientMockSuccess{}
	bytes, _ := json.Marshal(getProductList("afd"))
	result := string(bytes)

	if result != `{"@context":null,"@graph":null}` {
		t.Error("TestGetProductList failed")
	}
}

func TestFetchJSONDefault(t *testing.T) {
	client = &clientMockSuccess{}
	uri := "http://whatever.whatever"
	result, err := fetchJSON(client, uri)

	if string(result) != "" || err != nil {
		t.Error("TestFetchJSONDefault failed")
	}
}

func TestFetchJSONDoError(t *testing.T) {
	client = &clientMockDoError{}
	_, err := fetchJSON(client, "")

	if err == nil {
		t.Error("TestFetchJSONDoError failed")
	}
}

func TestFetchJSONNon200(t *testing.T) {
	client = &clientMockNon200{}
	_, err := fetchJSON(client, "")

	if err != nil {
		t.Error("TestFetchJSONNon200 failed")
	}
}
