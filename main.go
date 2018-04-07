package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ProductType is one of the supported NWS product types
type ProductType string

const lsrProductType ProductType = "lsr"
const swoProductType ProductType = "swo"
const nwsAPIURIBase = "https://api.weather.gov"
const eventSource = "api.weather.gov"
const requestDelayMs = 60 * 1000
const thresholdMinutes float64 = 60

var seenProducts = make(map[string]bool)
var producer, producerErr = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})

func writeToTopic(data []byte) {
	// Delivery report handler for produced messages
	go func() {
		for e := range producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic asynchronously
	topic := "wx.nws.api" // TODO why can't this be a const?
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)

	// Wait for message deliveries
	producer.Flush(15 * 1000)
}

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

func getProductList(productType ProductType) ProductListResponse {
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

// TODO ugly and needs to be broken out
func fetchAndProcessData(productType ProductType) {
	productList := getProductList(productType)

	for _, feature := range productList.Features {
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

// TODO test newline replace
func normalizeString(input string) string {
	return strings.ToLower(strings.TrimSpace(strings.Replace(input, "\n", " ", -1)))
}

func normalizeFloat(input string) float32 {
	inputString := normalizeString(input)
	num, err := strconv.ParseFloat(inputString, 32)
	if err != nil {
		num = 0
	}

	return float32(num)
}

func main() {
	if producerErr != nil {
		panic(producerErr)
	}

	/*ticker := time.NewTicker(requestDelayMs * time.Millisecond)
	for range ticker.C {
		go fetchAndProcessData(lsrProductType)
		go fetchAndProcessData(swoProductType)
	}*/

	//go fetchAndProcessData(lsrProductType)
	fetchAndProcessData(swoProductType)
}

// ProductListResponse model
type ProductListResponse struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Features []Product     `json:"features"`
}

// Product model
type Product struct {
	URI             string    `json:"@id"`
	ID              string    `json:"id"`
	WmoCollectiveID string    `json:"wmoCollectiveId"`
	IssuingOffice   string    `json:"issuingOffice"`
	IssuanceTime    time.Time `json:"issuanceTime"`
	ProductCode     string    `json:"productCode"`
	ProductName     string    `json:"productName"`
	ProductText     string    `json:"productText"`
	Details         interface{}
}

// LSRDetails model
type LSRDetails struct {
	Type     string
	Datetime time.Time
	Reported time.Time
	Mag      string
	Lat      float32
	Lon      float32
	Location string
	County   string
	State    string
	Source   string
	Remarks  string
}

// WxEvent model
type WxEvent struct {
	Source      string
	Details     interface{}
	Ingested    time.Time
	Description string
}

type Coordinates struct {
	Lat float32
	Lon float32
}
