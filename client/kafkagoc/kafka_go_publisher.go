package kafkagoc

import (
	"context"
	"time"

	kafka "github.com/telkomdev/kafka-cli"

	ka "github.com/segmentio/kafka-go"

	// kascram "github.com/segmentio/kafka-go/sasl/scram"
	kaplain "github.com/segmentio/kafka-go/sasl/plain"

	"github.com/segmentio/kafka-go/snappy"
)

//KafkaGoPublisherImpl struct
type KafkaGoPublisherImpl struct {
	writer *ka.Writer
}

//NewKafkaGoPublisherImpl constructor of KafkaGoPublisherImpl
func NewKafkaGoPublisher(args *kafka.Argument) (*KafkaGoPublisherImpl, error) {
	config := ka.WriterConfig{
		Brokers:          args.Brokers,
		Balancer:         &ka.LeastBytes{},
		CompressionCodec: snappy.NewCompressionCodec(),
		BatchTimeout:     5 * time.Millisecond,
		//BatchBytes:       1000000,
	}

	if args.Auth {
		// mechanism, err := kascram.Mechanism(kascram.SHA512, args.Username, args.Password)
		// if err != nil {
		// 	panic(err)
		// }

		mechanism := &kaplain.Mechanism{Username: args.Username, Password: args.Password}

		dialer := &ka.Dialer{
			Timeout:       10 * time.Second,
			DualStack:     true,
			SASLMechanism: mechanism,
		}

		config.Dialer = dialer
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
