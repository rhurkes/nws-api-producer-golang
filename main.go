package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"go.uber.org/zap"
)

const configPath = "config.toml"

var (
	conf               config
	logger             *zap.SugaredLogger
	lastSeenProduct    = make(map[nwsProduct]string)
	activeProductTypes = []nwsProduct{AreaForecastDiscussion, LocalStormReport, SevereWatch,
		SevereThunderstormWarning, SevereWeatherStatement, StormOutlookNarrative, TornadoWarning}
)

func init() {
	if _, err := toml.DecodeFile(configPath, &conf); err != nil {
		fmt.Printf("unable to decode config: %s", configPath)
	}
	productionLogger, _ := zap.NewDevelopment() // NewProduction()
	defer productionLogger.Sync()
	logger = productionLogger.Sugar()
	logger.Info("initializing")
}

func main() {
	ticker := time.NewTicker(time.Millisecond * time.Duration(conf.RequestDelayMs))
	processProducts(activeProductTypes)

	for range ticker.C {
		processProducts(activeProductTypes)
	}
}
