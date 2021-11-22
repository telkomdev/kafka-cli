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
func NewSaramaPublisher(args *Argument) (*SaramaPublisherImpl, error) {
	// producer config
	config := sarama.NewConfig()
	config.ClientID = "kafka-cli"
	config.Producer.Retry.Max = 10
	config.Producer.Retry.Backoff = 10 * time.Second
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Timeout = 10 * time.Second
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Return.Successes = true
	config.Producer.MaxMessageBytes = 104857599

	if args.Auth {
		config.Metadata.Full = true
		config.Net.SASL.Enable = true
		config.Net.SASL.User = args.Username
		config.Net.SASL.Password = args.Password
		config.Net.SASL.Handshake = true
		config.Net.SASL.Version = sarama.SASLHandshakeV0

		// config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &SaramScramClient{HashGeneratorFcn: SHA512} }
		// config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512

		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}

	// sync producer
	producer, err := sarama.NewSyncProducer(args.Brokers, config)

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
