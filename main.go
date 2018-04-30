package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Configuration
const nwsAPIURIBase = "https://api.weather.gov"
const eventSource = "api.weather.gov"
const requestDelayMs = 60 * 1000
const thresholdMinutes float64 = 60

var seenProducts = make(map[string]bool)
var producer, producerErr = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})

func fetchJSON(uri string, logResponse bool) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "sigtor.org")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(fmt.Sprintf("%v %s", resp.StatusCode, uri))

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return body
}

func getProductList(productType string) ProductListResponse {
	var uri string
	var productList ProductListResponse
	var timeFilteredProducts []Product

	uri = fmt.Sprintf("%s/products/types/%s", nwsAPIURIBase, productType)
	body := fetchJSON(uri, true)
	json.Unmarshal(body, &productList)
	now := time.Now().UTC()

	for _, product := range productList.Features {
		if now.Sub(product.IssuanceTime).Minutes() <= thresholdMinutes {
			timeFilteredProducts = append(timeFilteredProducts, product)
		}
	}

	productList.Features = timeFilteredProducts

	return productList
}

func fetchAndProcessData(productType string) {
	productList := getProductList(productType)

	for _, feature := range productList.Features {
		// TODO ugly and needs to be broken out
		_, seen := seenProducts[feature.ID]

		if !seen {
			productBody := fetchJSON(feature.URI, false)
			var product Product
			json.Unmarshal(productBody, &product)

			var wxEvent WxEvent
			var parseError error
			switch productType {
			case "lsr":
				wxEvent, parseError = processLSRProduct(product)
			case "swo":
				wxEvent, parseError = buildSWOEvent(product)
			case "afd":
				wxEvent, parseError = buildAFDEvent(product)
			case "tor":
				wxEvent, parseError = buildTOREvent(product)
			case "svs":
				wxEvent, parseError = buildSVSEvent(product)
			case "sel":
				wxEvent, parseError = buildSELEvent(product)
			case "svr":
				wxEvent, parseError = buildSVREvent(product)
			}

			if parseError != nil {
				fmt.Println(parseError)
				continue
			}

			// Parsing is complete, prepare the payload and deliver it
			wxEvent.Source = eventSource
			wxEvent.Ingested = time.Now().UTC()
			payload, err := json.Marshal(wxEvent)

			if err != nil {
				fmt.Println(fmt.Sprintf("Unable to marshal product: '%s'", product.ID))
				continue
			}

			writeToTopic(payload)
			seenProducts[feature.ID] = true
		}
	}
}

func main() {
	if producerErr != nil {
		panic(producerErr)
	}

	/*ticker := time.NewTicker(requestDelayMs * time.Millisecond)
	for range ticker.C {
		go fetchAndProcessData("afd")
		go fetchAndProcessData("lsr")
		go fetchAndProcessData("sel")
		go fetchAndProcessData("svr")
		go fetchAndProcessData("svs")
		go fetchAndProcessData("swo")
		go fetchAndProcessData("tor")
	}*/

	fetchAndProcessData("swo")
}
