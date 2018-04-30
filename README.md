# Build & Run
TODO

# Kafka
`kafkacat -b localhost -t wx.nws.api`

# Implemented Products - https://en.wikipedia.org/wiki/Specific_Area_Message_Encoding
- `LSR` DONE Local Storm Report. Only grabs the last 30 minutes on a fresh start as there can be 1000+ on the product list response.
- `SWO` DONE Severe Storm Outlook Narrative. Includes the 1/2/3/4-8 day outlooks (ACUS01/02/03/48) and Mesoscale Discussions (ACUS11).
        MDs contain their own coordinates and do not have a corresponding PTS. I think I only care about the SWODY1 for MVP.
- `SEL` DONE Severe Local Storm Watch and Watch Cancellation Msg. Issued when watches are issued. Has the watch text.
- `TOR` DONE Tornado Warning. Can contain Tornado Emergency status (https://en.wikipedia.org/wiki/Tornado_emergency)
- `SVR` Severe Thunderstorm Warning.
- `SVS` DONE Severe Weather Statement. Can contain Tornado Emergency status (https://en.wikipedia.org/wiki/Tornado_emergency)
- `AFD` DONE Area Forecast Discussion.

# Notes
Any watch or warning can have PDS terminology

# Potential Products
- `WOU` Watch Outline Update. Has the counties listed out for each watch.
- `WCN` Weather Watch Clearance Notification. Shows issued watches, but several minutes later. Generated by WFOs. Don't use this for anything.
- `SPS` Special Weather Statement. Why would I want this?
- `FFA` Flash Flood Watch
- `FLA` Flood Watch
- `FFW` Flash Flood Warning
- `FLW` Flood Warning

# Only Playable when/if mapping is added
- `PTS` Probabilistic Outlook Points. Contains coordinates for SWO outlooks (WUUS01/02/03/48).
- `SEV` Shows coordinates for all active watches.

# LSRs
## Helpful URLs
https://mesonet.agron.iastate.edu/request/gis/lsrs.phtml
## Event Types
FLOOD, HAIL, TSTM WND DMG, SNOW, HEAVY RAIN, NON-TSTM WND GST
## Regex
MAG and REMARKS may be empty

0700 PM     TSTM WND DMG     1 N CRAFTON             33.36N  97.90W
(\d{4} [A|P]M)\s{5}([A-Z|\s]+)\s{2,}(.+)\s+(\d{2}\.\d{2})N (\d{3}\.\d{2})W

03/26/2018  E1.25 INCH       COKE               TX   STORM CHASER

# TODO
IMMEDIATE:
2. implement SEL
4. SVR
5. redo and standardize all product schemas
6. watch probabilities

2. Figure out structs vs classes: https://medium.com/@simplyianm/why-gos-structs-are-superior-to-class-based-inheritance-b661ba897c67
3. Kafka errors on verbose testing
%3|1522520247.123|FAIL|rdkafka#producer-1| [thrd:localhost:9092/bootstrap]: localhost:9092/bootstrap: Connect to ipv6#[::1]:9092 failed: Connection refused
%3|1522520247.123|ERROR|rdkafka#producer-1| [thrd:localhost:9092/bootstrap]: localhost:9092/bootstrap: Connect to ipv6#[::1]:9092 failed: Connection refused
%3|1522520247.123|ERROR|rdkafka#producer-1| [thrd:localhost:9092/bootstrap]: 1/1 brokers are down
4. changing retention time doesn't seem to purge topic
6. Logging
7. Mocking

1. Break out LSR mag into value and units and measured/estimated