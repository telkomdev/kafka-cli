package kafka

import (
	"context"
	"time"

	ka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/snappy"
)

//KafkaGoPublisherImpl struct
type KafkaGoPublisherImpl struct {
	writer *ka.Writer
}

//NewKafkaGoPublisherImpl constructor of KafkaGoPublisherImpl
func NewKafkaGoPublisher(topic string, addresses ...string) (*KafkaGoPublisherImpl, error) {
	config := ka.WriterConfig{
		Brokers:          addresses,
		Balancer:         &ka.LeastBytes{},
		CompressionCodec: snappy.NewCompressionCodec(),
		BatchTimeout:     5 * time.Millisecond,
		//BatchBytes:       1000000,
	}

	writer := ka.NewWriter(config)
	return &KafkaGoPublisherImpl{writer}, nil
}

//Publish function
func (publisher *KafkaGoPublisherImpl) Publish(ctx context.Context, topic string, message []byte) error {
	defer func() { publisher.writer.Close() }()

	msg := ka.Message{
		Topic: topic,
		Value: message,
	}

	err := publisher.writer.WriteMessages(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}
