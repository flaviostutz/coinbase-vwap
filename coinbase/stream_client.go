package coinbase

import (
	"context"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

type StreamResponse struct {
	// "type": "last_match",
	Type    string `json:"type"`
	Message string `json:"message"`
	Reason  string `json:"reason"`
	// "trade_id": 3168503,
	// "maker_order_id": "989df982-552b-47db-a247-2cb34386d25d",
	// "taker_order_id": "6c5e5552-130a-4460-9fb2-9084291814c7",
	// "side": "sell",
	// "size": "0.01",
	Size *big.Float `json:"size"`
	// "price": "0.06961",
	Price *big.Float `json:"price"`
	// "product_id": "ETH-BTC",
	ProductId string `json:"product_id"`
	// "sequence": 25855546,
	// "time": "2021-05-21T05:04:10.251496Z"
}

type MatchInfo struct {
	ProductId string
	Price     *big.Float
	Size      *big.Float
}

type OnMatchInfo func(MatchInfo)

//SubscribeMatchesChannel Connects to Coinbase websocket matches channel and publishes each "match" message to the matchInfoChan channel
func SubscribeMatchesChannel(ctx context.Context, matchInfoChan chan<- MatchInfo, websocketURL string, product_ids ...string) error {

	//prepare subscribe message
	if len(product_ids) == 0 {
		return fmt.Errorf("'product_ids' cannot be empty")
	}
	pids := ""
	for _, s := range product_ids {
		if len(pids) > 0 {
			pids += ","
		}
		pids += fmt.Sprintf("\"%s\"", s)
	}

	subscribeMsg :=
		fmt.Sprintf(
			`{
		"type": "subscribe",
		"product_ids": [
			%s
		],
		"channels": [
			{
				"name": "matches",
				"product_ids": [
					%s
				]
			}
		]
	}`, pids, pids)

	//connect to coinbase ws
	logrus.Debugf("Connecting to coinbase websocket...")
	conn, err := websocket.Dial(websocketURL, "", "http://localhost")
	if err != nil {
		return err
	}

	//send subscribe message
	logrus.Debugf("Sending subscription config...")
	err = websocket.Message.Send(conn, subscribeMsg)
	if err != nil {
		return err
	}

	//wait for subscription confirmation essage
	resp := StreamResponse{}
	err = websocket.JSON.Receive(conn, &resp)
	if err != nil {
		return fmt.Errorf("Error while reading coinbase websocket. err=%s", err.Error())
	}
	if resp.Type == "error" {
		return fmt.Errorf("Subscription error. err=%s %s", resp.Message, resp.Reason)
	}
	if resp.Type != "subscriptions" {
		return fmt.Errorf("Subscription confirmation message was not received")
	}

	logrus.Infof("Subscription confirmed")

	//receive messages on a separate go routine
	go func() {
		logrus.Debugf("Starting to receive messages from match stream")
		defer conn.Close()

		for {
			select {
			case <-ctx.Done():
				if err != nil {
					logrus.Warnf("Error closing WS socket. err=%s", err.Error())
				}
				return
			default:
				resp := StreamResponse{}
				err = websocket.JSON.Receive(conn, &resp)
				if err != nil {
					logrus.Errorf("Error while reading coinbase websocket. err=%s", err.Error())
					return
				}
				logrus.Debugf("Match stream msg: %v", resp)
				if resp.Type == "match" {
					//notify received message
					matchInfoChan <- MatchInfo{
						ProductId: resp.ProductId,
						Price:     resp.Price,
						Size:      resp.Size,
					}
					continue
				}
			}
		}
	}()

	return nil
}
