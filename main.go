package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

// TODO move to types, NWS Product Type Enums
const (
	AreaForecastDiscussion nwsProduct = iota
	LocalStormReport
	SevereWatch
	SevereThunderstormWarning
	SevereWeatherStatement
	StormOutlookNarrative
	TornadoWarning
)

const configPath = "config.toml"

var (
	config             Config
	logger             *zap.SugaredLogger
	lastSeenProduct    = make(map[nwsProduct]string)
	activeProductTypes = []nwsProduct{AreaForecastDiscussion, LocalStormReport, SevereWatch,
		SevereThunderstormWarning, SevereWeatherStatement, StormOutlookNarrative, TornadoWarning}
)

func init() {
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		fmt.Printf("unable to decode config: %s", configPath)
	}
	productionLogger, _ := zap.NewProduction()
	defer productionLogger.Sync()
	logger = productionLogger.Sugar()
	logger.Info("initializing")
}

func main() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(config.RequestDelayMs))
	processProducts(activeProductTypes)

	for range ticker.C {
		processProducts(activeProductTypes)
	}
}
