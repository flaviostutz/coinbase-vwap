FROM golang:1.16.4 AS BUILD

RUN mkdir /app
WORKDIR /app

ADD go.mod .
ADD go.sum .
RUN go mod download

#now build source code
ADD . ./
RUN go test ./...
RUN go build -o /go/bin/coinbase-vwap



FROM golang:1.16.4

ENV LOG_LEVEL 'info'
ENV COINBASE_WS_URL 'wss://ws-feed-public.sandbox.pro.coinbase.com'
ENV PRODUCT_IDS 'ETH-BTC,BTC-EUR'
ENV KAFKA_BROKERS ''

COPY --from=BUILD /go/bin/* /bin/
ADD /startup.sh /
ENTRYPOINT /startup.sh

EXPOSE 4000
