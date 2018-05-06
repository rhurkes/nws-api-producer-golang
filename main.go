package main

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

// NWS Product Type Enums
const (
	AreaForecastDiscussion nwsProduct = iota
	LocalStormReport
	SevereWatch
	SevereThunderstormWarning
	SevereWeatherStatement
	StormOutlookNarrative
	TornadoWarning
)

const (
	nwsAPIURIBase  = "https://api.weather.gov"
	eventSource    = "api.weather.gov"
	requestDelayMs = 60 * 1000
	topic          = "queue.wx.events"
)

var (
	logger                *zap.SugaredLogger
	lastSeenProduct       = make(map[nwsProduct]string)
	producer, producerErr = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	activeProductTypes    = []nwsProduct{AreaForecastDiscussion, LocalStormReport, SevereWatch,
		SevereThunderstormWarning, SevereWeatherStatement, StormOutlookNarrative, TornadoWarning}
)

func init() {
	productionLogger, _ := zap.NewProduction()
	defer productionLogger.Sync()
	logger = productionLogger.Sugar()
	logger.Info("Initializing...")
}

func processFeature(productType nwsProduct, responseBody []byte) {
	var product product
	json.Unmarshal(responseBody, &product)

	var wxEvent wxEvent
	var parseError error
	switch productType {
	case LocalStormReport:
		wxEvent, parseError = processLSRProduct(product)
	case StormOutlookNarrative:
		wxEvent, parseError = buildSWOEvent(product)
	case AreaForecastDiscussion:
		wxEvent, parseError = buildAFDEvent(product)
	case TornadoWarning:
		wxEvent, parseError = buildTOREvent(product)
	case SevereWeatherStatement:
		wxEvent, parseError = buildSVSEvent(product)
	case SevereWatch:
		wxEvent, parseError = buildSELEvent(product)
	case SevereThunderstormWarning:
		wxEvent, parseError = buildSVREvent(product)
	}

	if parseError != nil {
		logger.Warnf("unable to parse: %v", getNWSProductCode(productType), parseError)
		return
	}

	if wxEvent.DoNotPublish {
		return
	}

	wxEvent.Source = eventSource
	wxEvent.Ingested = time.Now().UTC()
	payload, _ := json.Marshal(wxEvent)
	topicName := topic
	writeToTopic(payload, &topicName)
}

func main() {
	if producerErr != nil {
		logger.Fatal("Unable to start producer", producerErr)
		panic(producerErr)
	}

	ticker := time.NewTicker(requestDelayMs * time.Millisecond)
	for range ticker.C {
		for _, productType := range activeProductTypes {
			productCode := getNWSProductCode(productType)
			productList := getProductList(productCode)
			go processProductList(productType, productList.Features)
		}
	}
}
