package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/flaviostutz/coinbase-vwap/coinbase"
	"github.com/flaviostutz/coinbase-vwap/kafka"
	"github.com/sirupsen/logrus"
)

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	coinbaseWSURL := flag.String("coinbase-ws-url", "wss://ws-feed-public.sandbox.pro.coinbase.com", "Coinbase Websockets API endpoint URL. Defaults to sandbox URL")
	kafkaBrokers := flag.String("kafka-brokers", "", "Kafka broker addresses separated by comma. ex.: kafka1:9092,kafka2:9092")
	productIDs := flag.String("product-ids", "", "Comma separated list of product_id pairs. ex.: BTC-USD,BTC-GBP,BTC-EUR,ETH-BTC")
	flag.Parse()

	switch *logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	if *coinbaseWSURL == "" {
		logrus.Errorf("'coinbase-ws-url' parameter cannot be empty")
		return
	}

	if *productIDs == "" {
		logrus.Errorf("'product-ids' parameter cannot be empty")
		return
	}

	ctx, cancel := context.WithCancel(context.TODO())

	enableKafka := false
	if *kafkaBrokers != "" {
		enableKafka = true
		err := kafka.SetupKafkaProducer(ctx, strings.Split((*kafkaBrokers), ","))
		if err != nil {
			logrus.Errorf("Kafka is enabled but brokers are unreacheable. err=%s", err)
			os.Exit(1)
		}
	}

	logrus.Infof("====Starting coinbase-vwap====")

	logrus.Infof("Connecting to Coinbase Matches Stream...")
	mic := make(chan coinbase.MatchInfo)
	coinbase.SubscribeMatchesChannel(ctx, mic, *coinbaseWSURL, strings.Split(*productIDs, ",")...)

	logrus.Infof("Online VWAP calculations:")
	coinbase.CalculateVWAP(ctx, mic, 200, func(vwap coinbase.VWAPInfo) {
		fmt.Printf("VWAP-200 %s=%s\n\n", vwap.ProductId, vwap.Value.String())
		if enableKafka {
			err := kafka.PublishVWAPToKafka(vwap)
			if err != nil {
				logrus.Warnf("Error sending VWAP to Kafka topic. err=%s", err)
			}
		}
	})

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Kill)
	// for {
	select {
	case <-signals:
		logrus.Infof("Shuting down...")
		cancel()
		return
	}
	// }
}
