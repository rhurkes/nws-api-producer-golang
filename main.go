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
	logger          *zap.SugaredLogger
	lastSeenProduct = make(map[nwsProduct]string)
	// TODO better way of handling producer scope
	producer           *kafka.Producer
	kafkaErr           error
	activeProductTypes = []nwsProduct{AreaForecastDiscussion, LocalStormReport, SevereWatch,
		SevereThunderstormWarning, SevereWeatherStatement, StormOutlookNarrative, TornadoWarning}
)

func init() {
	productionLogger, _ := zap.NewProduction()
	defer productionLogger.Sync()
	logger = productionLogger.Sugar()
	logger.Info("Initializing...")
	producer, kafkaErr = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	// Unable to connect to broker is not an error
	if kafkaErr != nil {
		panic(kafkaErr)
	}
}

func processFeature(productType nwsProduct, responseBody []byte) {
	var product product
	var wxEvent wxEvent
	var parseError error

	json.Unmarshal(responseBody, &product)

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
	ticker := time.NewTicker(requestDelayMs * time.Millisecond)
	for range ticker.C {
		for _, productType := range activeProductTypes {
			productCode := getNWSProductCode(productType)
			productList := getProductList(productCode)
			go processProductList(productType, productList.Features)
		}
	}
}
