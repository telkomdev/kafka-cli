package kafka

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Argument struct
type Argument struct {
	Brokers     []string
	Topic       string
	Command     Command
	ShowVersion bool
	Help        func()
	Message     []byte
	Verbose     bool
}

type brokerList []string

func (s *brokerList) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *brokerList) Set(value string) error {
	*s = strings.Split(value, ",")
	return nil
}

// ParseArgument function
func ParseArgument() (*Argument, error) {

	var (
		brokers     brokerList
		topic       string
		message     string
		showVersion bool
		verbose     bool
	)

	argument := &Argument{}

	// sub command
	publishCommand := flag.NewFlagSet("pub", flag.ExitOnError)
	subscribeCommand := flag.NewFlagSet("sub", flag.ExitOnError)

	publishCommand.Var(&brokers, "b", "kafka brokers (you can add multiple brokers using separated comma eg: -b localhost:9091,localhost:9092 ..)")
	publishCommand.Var(&brokers, "broker", "kafka brokers (you can add multiple brokers using separated comma eg: -b localhost:9091,localhost:9092 ..)")
	publishCommand.StringVar(&topic, "t", "", "kafka topic")
	publishCommand.StringVar(&topic, "topic", "", "kafka topic")
	publishCommand.StringVar(&message, "m", "", "message to publish")
	publishCommand.StringVar(&message, "message", "", "message to publish")
	publishCommand.BoolVar(&verbose, "V", false, "verbose mode (if true log will appear otherwise no)")

	subscribeCommand.Var(&brokers, "b", "kafka brokers (you can add multiple brokers using separated comma eg: -b localhost:9091,localhost:9092 ..)")
	subscribeCommand.Var(&brokers, "broker", "you can add multiple brokers using separated comma eg: -b localhost:9091,localhost:9092 ..)")
	subscribeCommand.StringVar(&topic, "t", "", "kafka topic")
	subscribeCommand.StringVar(&topic, "topic", "", "kafka topic")
	subscribeCommand.BoolVar(&verbose, "V", false, "verbose mode (if true log will appear otherwise no)")

	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Usage = func() {
		fmt.Println()
		fmt.Println("usage : kafka-cli pub -broker localhost:9092 -topic my-topic")
		fmt.Println("sub command either pub (publish) or sub (subsriber)")
		fmt.Println("-b | -broker : kafka brokers (you can add multiple brokers using separated comma eg: -b localhost:9091,localhost:9092 ..)")
		fmt.Println("-t | -topic : kafka topic")
		fmt.Println("-h show help")
		fmt.Println("-version show version")
		fmt.Println("-V verbose mode")
	}

	flag.Parse()

	if len(os.Args) < 2 {
		argument.Help = flag.Usage
		return argument, errors.New("require at least one arguments")
	}

	if !strings.Contains(os.Args[1], "version") {
		switch os.Args[1] {
		case "pub":
			publishCommand.Parse(os.Args[2:])
		case "sub":
			subscribeCommand.Parse(os.Args[2:])
		default:
			argument.Help = flag.Usage
			return argument, errors.New("invalid command")
		}
	}

	// publish command parsed
	if publishCommand.Parsed() {

		if len(brokers) <= 0 {
			argument.Help = flag.Usage
			return argument, errors.New("require at least one broker")
		}

		if len(topic) <= 0 {
			argument.Help = flag.Usage
			return argument, errors.New("require a topic name")
		}

		if len(message) <= 0 {
			argument.Help = flag.Usage
			return argument, errors.New("pub command require a message")
		}

		for _, broker := range brokers {
			if broker == "" {
				continue
			}

			argument.Brokers = append(argument.Brokers, broker)
		}

		argument.Topic = topic
		argument.Message = []byte(message)
		argument.Command = CommandFromString(os.Args[1])
	}

	// subscribe command parsed
	if subscribeCommand.Parsed() {

		if len(brokers) <= 0 {
			argument.Help = flag.Usage
			return argument, errors.New("require at least one broker")
		}

		if len(topic) <= 0 {
			argument.Help = flag.Usage
			return argument, errors.New("require a topic name")
		}

		for _, broker := range brokers {
			if broker == "" {
				continue
			}

			argument.Brokers = append(argument.Brokers, broker)
		}

		argument.Topic = topic
		argument.Command = CommandFromString(os.Args[1])
	}

	argument.ShowVersion = showVersion
	argument.Verbose = verbose

	return argument, nil
}
