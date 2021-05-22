package coinbase

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	coinbaseWSURL = "wss://ws-feed-public.sandbox.pro.coinbase.com"
)

func TestStreamClient1(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	//open websocket and keep listening for messages for 10 seconds, according to context
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)

	mic := make(chan MatchInfo)
	SubscribeMatchesChannel(ctx, mic, coinbaseWSURL, "BTC-USD", "ETH-BTC")

	for {
		select {
		case <-ctx.Done():
			//test timeout with no message received yet
			assert.True(t, false)
		default:
			<-mic
			assert.True(t, true)
			//terminate context immediatelly as it was successfull (no need to wait for timeout)
			cancel()
			return
		}
	}
}

func TestStreamClientShouldErr(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	mic := make(chan MatchInfo)
	err := SubscribeMatchesChannel(context.TODO(), mic, coinbaseWSURL, "INVALID_PRODUCT_ID")
	assert.NotNil(t, err)
}
