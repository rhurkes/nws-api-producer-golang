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
	// Impossible to get an error building this request as-is
	req, _ := http.NewRequest("GET", uri, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sigtor.org")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("error fetching uri: " + uri)
	}

	if resp.StatusCode != 200 {
		logger.Errorf("%v %s", resp.StatusCode, uri)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("error reading response body, uri: " + uri)
	}

	return body, nil
}

func buildProductListURI(productType string) string {
	return fmt.Sprintf("%s/products/types/%s", nwsAPIURIBase, productType)
}

func getProductList(productType string) productListResponse {
	uri := buildProductListURI(productType)
	body, err := fetchJSON(client, uri)
	if err != nil {
		logger.Warn(err)
		return productListResponse{}
	}
	var productList productListResponse
	json.Unmarshal(body, &productList)
	return productList
}

func processProductList(productType nwsProduct, features []product) {
	for _, feature := range features {
		// On a fresh launch we capture the first product so we can break early below. This ensures
		// we're only processing new products since the producer started polling.
		if lastSeenProduct[productType] == "" {
			lastSeenProduct[productType] = feature.ID
		}

		// Features are always in the same order, so break the loop if the ID has been seen.
		if feature.ID == lastSeenProduct[productType] {
			break
		}

		productBody, err := fetchJSON(client, feature.URI)
		if err != nil {
			logger.Warn(err)
			continue
		}
		processFeature(productType, productBody)
	}

	// Once feature iteration is done, update last seen with the first feature (if it exists)
	if len(features) > 0 {
		lastSeenProduct[productType] = features[0].ID
	}
}
