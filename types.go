package main

import "time"

type nwsProduct int

type Config struct {
	NWSAPIURIBase  string
	EventSource    string
	RequestDelayMs int
	Topic          string
	UserAgent      string
}

type productListResponse struct {
	Context interface{} `json:"@context"`
	Graph   []product   `json:"@graph"`
}

// Used for both product list and product calls
type product struct {
	URI             string    `json:"@id"`
	ID              string    `json:"id"`
	WmoCollectiveID string    `json:"wmoCollectiveId"`
	IssuingOffice   string    `json:"issuingOffice"`
	IssuanceTime    time.Time `json:"issuanceTime"`
	ProductCode     string    `json:"productCode"`
	ProductName     string    `json:"productName"`

	// ProductText is only populated on product calls
	ProductText string `json:"productText"`
}

type wxEvent struct {
	Source       string
	Details      interface{}
	Ingested     time.Time
	Summary      string
	DoNotPublish bool
}

type coordinates struct {
	Lat float32
	Lon float32
}

type movement struct {
	Time     string
	Location coordinates
	Degrees  int
	Knots    int
}
