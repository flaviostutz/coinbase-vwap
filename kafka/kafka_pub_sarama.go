package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/flaviostutz/coinbase-vwap/coinbase"
	"github.com/sirupsen/logrus"
)

var (
	saramaProducer      sarama.AsyncProducer
	producerInitialized bool
)

func SetupKafkaProducer(ctx context.Context, brokerAddresses []string) error {
	logrus.Infof("Setting up Kafka producer for brokers %v...", brokerAddresses)
	conf := sarama.NewConfig()
	conf.ClientID = "vwap"
	producer, err := sarama.NewAsyncProducer(brokerAddresses, conf)
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
				logrus.Warnf("Error publishing message to Kafka. err=%s", err)
			case suc := <-producer.Successes():
				logrus.Debugf("Message sent to Kafka with success. topic=%s value=%s", suc.Topic, suc.Value)
			}
		}
	}()
	saramaProducer = producer
	producerInitialized = true
	return nil
}

func PublishVWAPToKafka(info coinbase.VWAPInfo) error {
	if !producerInitialized {
		return fmt.Errorf("call SetupKafkaProducer before calling Publish")
	}

	topic := "vwap-" + info.ProductId
	logrus.Debugf("Sending message %s to kafka topic %s", info.Value.String(), topic)
	saramaProducer.Input() <- &sarama.ProducerMessage{
		Topic: topic, Key: nil,
		Value: sarama.StringEncoder(info.Value.String()),
	}

	return nil
}
