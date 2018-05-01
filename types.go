package main

import "time"

// ProductListResponse type
type ProductListResponse struct {
	Context  []interface{} `json:"@context"`
	Type     string        `json:"type"`
	Features []Product     `json:"features"`
}

// Product type
type Product struct {
	URI             string    `json:"@id"`
	ID              string    `json:"id"`
	WmoCollectiveID string    `json:"wmoCollectiveId"`
	IssuingOffice   string    `json:"issuingOffice"`
	IssuanceTime    time.Time `json:"issuanceTime"`
	ProductCode     string    `json:"productCode"`
	ProductName     string    `json:"productName"`
	ProductText     string    `json:"productText"`
}

// WxEvent type
type WxEvent struct {
	Source   string
	Details  interface{}
	Ingested time.Time
	Summary  string
}

// Coordinates type
type Coordinates struct {
	Lat float32
	Lon float32
}
