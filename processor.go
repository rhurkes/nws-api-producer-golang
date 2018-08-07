package main

import (
	"encoding/json"
	"strings"
	"time"
)

func processProducts(activeProductTypes []nwsProduct) {
	for _, productType := range activeProductTypes {
		productCode := getNWSProductCode(productType)
		productList := getProductList(productCode)
		go processProductList(productType, productList.Graph)
	}
}

func processProductList(productType nwsProduct, features []product) {
	var newFeatures = 0
	var productCode = getNWSProductCode(productType)

	for _, feature := range features {
		// On a fresh launch we capture the first product so we can break early below. This ensures
		// we're only processing new products since the producer started polling.
		if lastSeenProduct[productType] == "" {
			logger.Infof("no lastSeenProduct found for %s - setting now", productCode)
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
		newFeatures++
		processFeature(productType, productBody)
	}

	// Once feature iteration is done, update last seen with the first feature (if it exists)
	if len(features) > 0 {
		lastSeenProduct[productType] = features[0].ID
	}

	if newFeatures > 0 {
		logger.Infof("processed %d new features", newFeatures)
	}
}

func processFeature(productType nwsProduct, responseBody []byte) {
	var product product
	var wxEvent wxEvent
	var parseError error

	json.Unmarshal(responseBody, &product)

	switch productType {
	case LocalStormReport:
		wxEvent, parseError = buildLSREvent(product)
	case StormOutlookNarrative:
		wxEvent, parseError = buildSWOEvent(product)
	case AreaForecastDiscussion:
		// Nothing to process
	case TornadoWarning:
		wxEvent, parseError = buildTOREvent(product)
	case SevereWeatherStatement:
		wxEvent, parseError = buildSVSEvent(product)
	case SevereWatch:
		wxEvent, parseError = buildSELEvent(product)
	case SevereThunderstormWarning:
		wxEvent, parseError = buildSVREvent(product)
	case FlashFloodWarning:
		wxEvent, parseError = parseFFWEvent(product)
	default:
		logger.Warn("Unhandled product", product)
	}

	if parseError != nil {
		logger.Warnf("unable to parse: %v", getNWSProductCode(productType), parseError)
		return
	}

	if wxEvent.DoNotPublish {
		logger.Infof("event marked as DoNotPublish")
		return
	}

	wxEvent.Source = conf.EventSource
	wxEvent.Processed = time.Now().UTC()
	wxEvent.Data = nwsData{
		Common: nwsCommonData{
			Code:   strings.ToLower(product.ProductCode),
			ID:     product.ID,
			Issued: product.IssuanceTime,
			Name:   product.ProductName,
			Wfo:    product.IssuingOffice,
			Text:   normalizeString(product.ProductText, false)}}

	payload, _ := json.Marshal(wxEvent)
	topicName := conf.Topic
	writeToTopic(payload, &topicName)
}
