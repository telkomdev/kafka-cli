package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Shopify/sarama"

	kafka "github.com/telkomdev/kafka-cli"
	// "github.com/telkomdev/kafka-cli/client/kafkagoc"
	"github.com/telkomdev/kafka-cli/client/saramac"
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
	// publisher, err := kafkagoc.NewKafkaGoPublisher(args)
	// if err != nil {
	// 	fmt.Println("error : ", err)

	// 	os.Exit(1)
	// }

	// subscriber, err := kafkagoc.NewKafkaGoSubscriber(args)
	// if err != nil {
	// 	fmt.Println("error : ", err)

	// 	os.Exit(1)
	// }

	publisher, err := saramac.NewSaramaPublisher(args)
	if err != nil {
		fmt.Println("error : ", err)

		os.Exit(1)
	}

	subscriber, err := saramac.NewSaramaSubscriber(args)
	if err != nil {
		fmt.Println("error : ", err)

		os.Exit(1)
	}

	runner := kafka.Runner{
		Publisher:  publisher,
		Subscriber: subscriber,
		Argument:   args,
	}

	if err = runner.Run(ctx); err != nil {
		fmt.Println("error : ", err)

		args.Help()
		os.Exit(1)
	}
}
