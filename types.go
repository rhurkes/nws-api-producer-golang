package main

import "time"

type nwsProduct int

type productListResponse struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Features []product     `json:"features"`
}

type product struct {
	URI             string    `json:"@id"`
	ID              string    `json:"id"`
	WmoCollectiveID string    `json:"wmoCollectiveId"`
	IssuingOffice   string    `json:"issuingOffice"`
	IssuanceTime    time.Time `json:"issuanceTime"`
	ProductCode     string    `json:"productCode"`
	ProductName     string    `json:"productName"`
	ProductText     string    `json:"productText"`
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
