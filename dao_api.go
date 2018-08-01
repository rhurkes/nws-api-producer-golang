package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var client httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func fetchJSON(client httpClient, uri string) ([]byte, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		var msg = "error building request for uri: " + uri
		logger.Warn(msg)
		return nil, errors.New(msg)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", conf.UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		var msg = "error fetching uri: " + uri
		logger.Warn(msg)
		return nil, errors.New(msg)
	}

	if resp.StatusCode != 200 {
		logger.Warnf("%v %s", resp.StatusCode, uri)
	} else {
		logger.Debugf("%v %s", resp.StatusCode, uri)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		var msg = "error reading response body, uri: " + uri
		logger.Warn(msg)
		return nil, errors.New(msg)
	}

	return body, nil
}

func buildProductListURI(productType string) string {
	return fmt.Sprintf("%s/products/types/%s", conf.NWSAPIURIBase, productType)
}

func getProductList(productType string) productListResponse {
	uri := buildProductListURI(productType)
	body, err := fetchJSON(client, uri)
	if err != nil {
		logger.Warn(err)
		return productListResponse{}
	}
	var productList productListResponse
	err = json.Unmarshal(body, &productList)
	if err != nil {
		logger.Warn(err)
		return productListResponse{}
	}
	return productList
}
