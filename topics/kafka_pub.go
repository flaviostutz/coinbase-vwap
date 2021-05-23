package topics

import (
	"context"
	"fmt"

	"github.com/flaviostutz/coinbase-vwap/coinbase"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

var (
	topicWriters map[string]*kafka.Writer
	brokers      []string
)

func init() {
	topicWriters = make(map[string]*kafka.Writer)
}

func SetKafkaAddress(brokers_ []string) {
	brokers = brokers_
}

func PublishVWAPToKafka(ctx context.Context, info coinbase.VWAPInfo) error {
	if len(brokers) == 0 {
		return fmt.Errorf("call SetKafkaAddress before calling Publish")
	}

	topic := "vwap-" + info.ProductId
	// logrus.Debugf("Resolving Kafka topic writer for %s", topic)

	writer, ok := topicWriters[topic]
	if !ok {
		logrus.Infof("Creating Kafka writer for topic %s", topic)
		writer = kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
			//in a real application this should be sent using chan in a separate routine
			//to decouple the receiving of messages from sending to Kafka
			// Async: true,
		})
		topicWriters[topic] = writer
	}

	logrus.Debugf("Sending message vwap %s to topic %s", info.Value.String(), topic)
	err := writer.WriteMessages(ctx, kafka.Message{
		// Key:   []byte(strconv.Itoa(i)),
		Value: []byte(info.Value.String()),
	})
	if err != nil {
		return err
	}
	return nil
}
