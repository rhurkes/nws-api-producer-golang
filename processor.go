package main

import (
	"encoding/json"
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
	default:
		// TODO how to handle unhandled products
	}

	if parseError != nil {
		logger.Warnf("unable to parse: %v", getNWSProductCode(productType), parseError)
		return
	}

	// TODO why do we have a field? Because the wxevent is valid, we just don't want to process it
	// TODO should negate it so the field defaults to false
	if wxEvent.DoNotPublish {
		return
	}

	wxEvent.Source = config.EventSource
	wxEvent.Ingested = time.Now().UTC()
	payload, _ := json.Marshal(wxEvent)
	topicName := config.Topic
	writeToTopic(payload, &topicName)
}
