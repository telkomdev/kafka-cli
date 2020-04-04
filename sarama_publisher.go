package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
)

//SaramaPublisherImpl struct
type SaramaPublisherImpl struct {
	producer sarama.SyncProducer
}

//NewSaramaPublisher constructor of SaramaPublisherImpl
func NewSaramaPublisher(addresses ...string) (*SaramaPublisherImpl, error) {

	// producer config
	configuration := sarama.NewConfig()
	configuration.ClientID = "kafka-cli"
	configuration.Producer.Retry.Max = 10
	configuration.Producer.Retry.Backoff = 10 * time.Second
	configuration.Producer.RequiredAcks = sarama.WaitForAll
	configuration.Producer.Timeout = 10 * time.Second
	configuration.Producer.Compression = sarama.CompressionSnappy
	configuration.Producer.Return.Successes = true
	configuration.Producer.MaxMessageBytes = 104857599

	// sync producer
	producer, err := sarama.NewSyncProducer(addresses, configuration)

	if err != nil {
		return nil, err
	}

	return &SaramaPublisherImpl{producer}, nil
}

//Publish function
func (publisher *SaramaPublisherImpl) Publish(ctx context.Context, topic string, message []byte) error {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	p, o, err := publisher.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Println("Partition ", p)
	fmt.Println("Offset ", o)

	return nil
}
