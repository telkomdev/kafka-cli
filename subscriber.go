package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
)

//Subscriber interface
type Subscriber interface {
	Subscribe(context.Context, string) error
}

//SubscriberImpl struct
type SubscriberImpl struct {
	c       sarama.ConsumerGroup
	handler *SubscriberHandler
}

// SubscriberHandler struct will implement ConsumerGroupHandler
type SubscriberHandler struct {
	wait chan bool
}

//NewSubscriber constructor of SubscriberImpl
func NewSubscriber(addresses ...string) (*SubscriberImpl, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumerGroup(addresses, "kafka-cli-group", config)
	if err != nil {
		panic(err)
	}

	subscriberHandler := &SubscriberHandler{}

	return &SubscriberImpl{c: consumer, handler: subscriberHandler}, nil
}

//Subscribe function
func (s *SubscriberImpl) Subscribe(ctx context.Context, topics ...string) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for {
			s.handler.wait = make(chan bool, 0)
			err := s.c.Consume(ctx, topics, s.handler)
			if err != nil {
				return
			}
		}
	}()
	<-s.handler.wait

	<-signals

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
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}

	return nil
}
