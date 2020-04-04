package kafka

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
	ka "github.com/segmentio/kafka-go"
	_ "github.com/segmentio/kafka-go/snappy"
)

//KafkaGoSubscriberImpl struct
type KafkaGoSubscriberImpl struct {
	reader *ka.Reader
}

//NewKafkaGoSubscriber constructor of KafkaGoSubscriberImpl
func NewKafkaGoSubscriber(topic string, addresses ...string) (*KafkaGoSubscriberImpl, error) {
	config := kafka.ReaderConfig{
		Brokers:  addresses,
		Topic:    topic,
		GroupID:  "kafka-cli-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		// CommitInterval: time.Second, // flushes commits to Kafka every second
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
			message, err := s.reader.ReadMessage(ctx)
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
			//s.reader.CommitMessages(ctx, message)
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
