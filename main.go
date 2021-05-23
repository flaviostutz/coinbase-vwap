package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/flaviostutz/coinbase-vwap/coinbase"
	"github.com/sirupsen/logrus"
)

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	coinbaseWSURL := flag.String("coinbase-ws-url", "wss://ws-feed-public.sandbox.pro.coinbase.com", "Coinbase Websockets API endpoint URL. Defaults to sandbox URL")
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

	logrus.Infof("====Starting coinbase-vwap====")

	logrus.Infof("Connecting to Coinbase Matches Stream...")
	mic := make(chan coinbase.MatchInfo)
	coinbase.SubscribeMatchesChannel(context.TODO(), mic, *coinbaseWSURL, "BTC-USD", "ETH-BTC")

	logrus.Infof("Online VWAP calculations:")
	coinbase.CalculateVWAP(context.TODO(), mic, 200, func(vwap coinbase.VWAPInfo) {
		fmt.Printf("VWAP-200 %s=%s\n", vwap.ProductId, vwap.Value.String())
	})

	for {
		select {
		case <-context.TODO().Done():
			logrus.Debugf("Exiting...")
			os.Exit(0)
		}
	}
}

// opt := handlers.Options{
// 	WFSURL:        *wfsURL,
// 	MongoDBName:   *mongoDBName0,
// 	MongoAddress:  *mongoAddress0,
// 	MongoUsername: *mongoUsername0,
// 	MongoPassword: *mongoPassword0,
// }

// if opt.WFSURL == "" {
// 	logrus.Errorf("'--wfs-url' is required")
// 	os.Exit(1)
// }

// h := handlers.NewHTTPServer(opt)
// err := h.Start()
// if err != nil {
// 	logrus.Errorf("Error starting server. err=%s", err)
// 	os.Exit(1)
// }
