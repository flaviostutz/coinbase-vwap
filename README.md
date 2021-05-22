# coinbase-vwap

The goal of this project is to implement a real-time vwap calculator from the Coinbase match stream

## Usage

### Just run it

* git clone this repo to your workspace
* run `docker-compose up --build`
* all unit tests will be run, the program will be compiled and run with default parameters

### Parameters

Open docker-compose.yml and adjust Environment parameters:

* `LOG_LEVEL` - Configures log verbosity. Maybe one of info, debug, warning, error. Defaults to info.
* `COINBASE_WS_URL` - Coinbase Websocket endpoint. Defaults to free sandbox URL "wss://ws-feed-public.sandbox.pro.coinbase.com"

## Tests

Execute `go test ./...` for running all tests

