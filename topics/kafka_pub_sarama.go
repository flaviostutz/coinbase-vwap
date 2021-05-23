package topics

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/flaviostutz/coinbase-vwap/coinbase"
	"github.com/sirupsen/logrus"
)

var saramaProducer *sarama.AsyncProducer

func SetupKafkaProducer(ctx context.Context, brokerAddress []string) error {
	conf := sarama.NewConfig()
	conf.ClientID = "vwap"
	producer, err := sarama.NewAsyncProducer(brokerAddress, conf)
	if err != nil {
		return err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				producer.AsyncClose()
				return
			case err := <-producer.Errors():
				logrus.Warn("Error publishing message to Kafka. err=%s", err)
			case suc := <-producer.Successes():
				logrus.Debug("Message sent to Kafka with success. topic=%s value=%s", suc.Topic, suc.Value)
			}
		}
	}()
	saramaProducer = &producer
	return nil
}

func PublishVWAPToKafka(info coinbase.VWAPInfo) error {
	if saramaProducer == nil {
		return fmt.Errorf("call SetupKafkaProducer before calling Publish")
	}

	topic := "vwap-" + info.ProductId
	producer := *saramaProducer
	logrus.Debugf("Sending message vwap %s to topic %s", info.Value.String(), topic)
	producer.Input() <- &sarama.ProducerMessage{
		Topic: topic, Key: nil,
		Value: sarama.StringEncoder(info.Value.String()),
	}

	return nil
}
