package kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

//SaramaSubscriberImpl struct
type SaramaSubscriberImpl struct {
	c       sarama.ConsumerGroup
	handler *SubscriberHandler
}

// SubscriberHandler struct will implement ConsumerGroupHandler
type SubscriberHandler struct {
	wait chan bool
}

//NewSaramaSubscriber constructor of SaramaSubscriberImpl
func NewSaramaSubscriber(args *Argument) (*SaramaSubscriberImpl, error) {
	config := sarama.NewConfig()

	kafkaVersion, _ := sarama.ParseKafkaVersion("2.1.1")
	config.Version = kafkaVersion
	config.Consumer.Return.Errors = true

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

	consumer, err := sarama.NewConsumerGroup(args.Brokers, "kafka-cli-group", config)
	if err != nil {
		return nil, err
	}

	subscriberHandler := &SubscriberHandler{}

	return &SaramaSubscriberImpl{c: consumer, handler: subscriberHandler}, nil
}

//Subscribe function
func (s *SaramaSubscriberImpl) Subscribe(ctx context.Context, topics ...string) error {
	go func() {
		for {
			s.handler.wait = make(chan bool, 0)
			err := s.c.Consume(ctx, topics, s.handler)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()
	<-s.handler.wait

	// wait until get signal from OS
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	// close consumer after exit
	err := s.c.Close()
	if err != nil {
		return err
	}

	return nil
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (handler *SubscriberHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(handler.wait)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (handler *SubscriberHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (handler *SubscriberHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		fmt.Println()
		log.Printf("\nMessage = %s\nTopic = %s", string(message.Value), message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
