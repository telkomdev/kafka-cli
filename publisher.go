package kafka

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

//Publisher interface
type Publisher interface {
	Publish(string, []byte) error
}

//PublisherImpl struct
type PublisherImpl struct {
	producer sarama.SyncProducer
}

//NewPublisher constructor of PublisherImpl
func NewPublisher(addresses ...string) (*PublisherImpl, error) {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

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

	return &PublisherImpl{producer}, nil
}

//Publish function
func (publiher *PublisherImpl) Publish(topic string, message []byte) error {
	// publish sync
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	p, o, err := publiher.producer.SendMessage(msg)
	if err != nil {
		return err
	}

	fmt.Println("Partition ", p)
	fmt.Println("Offset ", o)

	return nil
}