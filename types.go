package main

import "time"

// NWS Product Types
const (
	AreaForecastDiscussion nwsProduct = iota
	LocalStormReport
	SevereWatch
	SevereThunderstormWarning
	SevereWeatherStatement
	StormOutlookNarrative
	TornadoWarning
	FlashFloodWarning
	// TODO SEV
	// TODO Re-order all these
)

type nwsProduct int

type config struct {
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
	Processed    time.Time
	Data         nwsData
	DoNotPublish bool
}

type nwsData struct {
	Common  nwsCommonData
	Derived interface{}
}

type nwsCommonData struct {
	Code   string
	ID     string
	Issued time.Time
	Name   string
	Text   string
	Wfo    string
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
