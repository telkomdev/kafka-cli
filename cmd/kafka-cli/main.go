package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"

	kafka "github.com/musobarlab/kafka-cli"
)

func main() {
	args, err := kafka.ParseArgument()
	if err != nil {
		fmt.Println("error : ", err)

		args.Help()
		os.Exit(1)
	}

	if args.ShowVersion {
		fmt.Printf("%s version %s\n", os.Args[0], kafka.Version)
		os.Exit(0)
	}

	if args.Verbose {
		sarama.Logger = log.New(os.Stdout, "kafka-cli-", log.Ltime)
	}

	ctx := context.Background()
	publisher, err := kafka.NewKafkaGoPublisher(args.Topic, args.Brokers...)
	if err != nil {
		fmt.Println("error : ", err)

		os.Exit(1)
	}

	subsriber, err := kafka.NewKafkaGoSubscriber(args.Topic, args.Brokers...)
	if err != nil {
		fmt.Println("error : ", err)

		os.Exit(1)
	}

	runner := kafka.Runner{
		Publisher:  publisher,
		Subscriber: subsriber,
		Argument:   args,
	}

	if err = runner.Run(ctx); err != nil {
		fmt.Println("error : ", err)

		args.Help()
		os.Exit(1)
	}
}
