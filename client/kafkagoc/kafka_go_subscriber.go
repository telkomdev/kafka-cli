package kafkagoc

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	kafka "github.com/telkomdev/kafka-cli"

	ka "github.com/segmentio/kafka-go"

	// kascram "github.com/segmentio/kafka-go/sasl/scram"
	kaplain "github.com/segmentio/kafka-go/sasl/plain"
	_ "github.com/segmentio/kafka-go/snappy"
)

//KafkaGoSubscriberImpl struct
type KafkaGoSubscriberImpl struct {
	reader *ka.Reader
}

//NewKafkaGoSubscriber constructor of KafkaGoSubscriberImpl
func NewKafkaGoSubscriber(args *kafka.Argument) (*KafkaGoSubscriberImpl, error) {
	config := ka.ReaderConfig{
		Brokers:        args.Brokers,
		Topic:          args.Topic,
		GroupID:        "kafka-cli-group",
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // flushes commits to Kafka every second
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

	reader := ka.NewReader(config)
	return &KafkaGoSubscriberImpl{reader}, nil
}

//Subscribe function
func (s *KafkaGoSubscriberImpl) Subscribe(ctx context.Context, topics ...string) error {
	signals := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	// wait until get signal from OS
	signal.Notify(signals, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	go waitOSNotify(signals, done)

	go func() {
		for {
			message, err := s.reader.FetchMessage(ctx)
			if err != nil {
				if err == io.EOF {
					fmt.Println("read message end ")
					break
				} else {
					fmt.Println("read message error ", err)
				}
			}

			fmt.Println()
			fmt.Printf("\nMessage = %s\nTopic = %s", string(message.Value), message.Topic)
			s.reader.CommitMessages(ctx, message)
		}
	}()

	<-done

	// close subscriber after exit
	err := s.close()
	if err != nil {
		return err
	}

	return nil
}

func waitOSNotify(kill chan os.Signal, done chan bool) {
	select {
	case <-kill:
		fmt.Println("process interrupted")
		done <- true
		return
	}

}

func (s *KafkaGoSubscriberImpl) close() error {
	fmt.Println("closing subscriber")
	return s.reader.Close()
}
