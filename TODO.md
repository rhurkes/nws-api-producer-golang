1. Figure out packages - break out test helper
2. Figure out structs vs classes: https://medium.com/@simplyianm/why-gos-structs-are-superior-to-class-based-inheritance-b661ba897c67
3. Kafka errors on verbose testing
%3|1522520247.123|FAIL|rdkafka#producer-1| [thrd:localhost:9092/bootstrap]: localhost:9092/bootstrap: Connect to ipv6#[::1]:9092 failed: Connection refused
%3|1522520247.123|ERROR|rdkafka#producer-1| [thrd:localhost:9092/bootstrap]: localhost:9092/bootstrap: Connect to ipv6#[::1]:9092 failed: Connection refused
%3|1522520247.123|ERROR|rdkafka#producer-1| [thrd:localhost:9092/bootstrap]: 1/1 brokers are down
4. changing retention time doesn't seem to purge topic
5. Move models somewhere else
6. Logging
7. Mocking

1. Break out LSR mag into value and units and measured/estimated
2. only get active mds
3. only get last outlook of each type


ProbabilisticOutlookPoints: 'pts' as 'pts',
  SevereStormOutlookNarrative: 'swo' as 'swo',
  SevereWeatherStatement: 'svs' as 'svs',
  TornadoWarning: 'tor' as 'tor',
  SevereLocalStormWatch: 'sls' as 'sls',
  WeatherWatchClearanceNotification: 'wcn' as 'wcn',