# coinbase-vwap

The goal of this project is to implement a real-time vwap calculator from the Coinbase match stream

## Usage

### Just run it

* git clone this repo to your workspace
* run `docker-compose up --build`
* all unit tests will be run, the program will be compiled and run with default parameters, outputing the latest VWAP to stdout for each ProductId "BTC-USD" and "ETH-BTC"

### Run with Kafka support

* git clone this repo to your workspace
* run `docker-compose -f docker-compose-withkafka.yml up --build`
* program will be compiled, a Kafka broker will be run and the program will start vwap calculations while publishing the results to Kafka
* open http://localhost:19000/topic/vwap-BTC-USD/messages to see published VWAP messages to Kafka Topic
* If you really want to see topic messages being consumed from Kafka and displayed to the terminal, execute the following command in a separate terminal

```sh
docker exec $(docker ps | grep kafka | cut -d' ' -f1) /bin/sh -c '/opt/bitnami/kafka/bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic vwap-BTC-USD --from-beginning'
```

### Parameters

Open docker-compose.yml and adjust Environment parameters:

* `LOG_LEVEL` - Configures log verbosity. Maybe one of info, debug, warning, error. Defaults to info.
* `COINBASE_WS_URL` - Coinbase Websocket endpoint. Defaults to free sandbox URL "wss://ws-feed-public.sandbox.pro.coinbase.com"

## How it works

Basically this program is composed of four components:

* weighted_avg.go: an weighted averager in a running window prepared for financial computations by using 'big' package to avoid floating point errors

* stream_client.go: connects to Coinbase websocket 'match' channel and produces "MatchInfo" messages in an output go channel for being processed by another thread/routine

* vwap_calculator.go: consumes a go channel with "MatchInfo" messages, instantiates one weighted averager for each different "ProductId", calculates the running average and callbacks a method with the newest value

* main.go: orchestrates the program parameters and the components above

## Performance tweaks

A potential bottleneck in our problem is the calculation of long windows in VWAP each time we need its result (in our case we will need to calculate it everytime we receive a new sample).

A regular implementation would have to iterate over all elements in the window twice in order to calculate both terms of the equation WAVG = SUM(V*W)/SUM(W) every time we need to calculate the resulting windows average, resulting in a complexity of O(2N) -> O(N).

In our implementation, we chose to keep partial sums for both SUM(V*W) and SUM(W) by removing "expired" elements from its sum while adding new elements to the sum when adding a new sample to the window. This strategy reduces the "regular" complexity from **O(N) to O(1)**.

## Tests

Execute `go test ./...` for running all tests

