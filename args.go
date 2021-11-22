package kafka

import (
	"bufio"
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
	Auth        bool
	Username    string
	Password    string
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
		auth        bool
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
	publishCommand.BoolVar(&auth, "auth", false, "set username and password prompt")

	subscribeCommand.Var(&brokers, "b", "kafka brokers (you can add multiple brokers using separated comma eg: -b localhost:9091,localhost:9092 ..)")
	subscribeCommand.Var(&brokers, "broker", "you can add multiple brokers using separated comma eg: -b localhost:9091,localhost:9092 ..)")
	subscribeCommand.StringVar(&topic, "t", "", "kafka topic")
	subscribeCommand.StringVar(&topic, "topic", "", "kafka topic")
	subscribeCommand.BoolVar(&verbose, "V", false, "verbose mode (if true log will appear otherwise no)")
	subscribeCommand.BoolVar(&auth, "auth", false, "set username and password prompt")

	flag.BoolVar(&showVersion, "version", false, "show version")

	flag.Usage = func() {
		fmt.Println()
		fmt.Printf("---------------------- Kafka CLI (version %s) ----------------------\n", Version)
		fmt.Println("------------------------------------------------------------------------------")
		fmt.Println()
		fmt.Println(`publish usage : kafka-cli pub -broker localhost:9092 -topic my-topic -message "hello world"`)
		fmt.Println()
		fmt.Println(`publish with auth usage : kafka-cli pub -broker localhost:9092 -topic my-topic -message "hello world" -auth`)
		fmt.Println("------------------------------------------------------------------------------")
		fmt.Println(`subscribe usage : kafka-cli sub -broker localhost:9092 -topic my-topic`)
		fmt.Println()
		fmt.Println(`subscribe with auth usage : kafka-cli sub -broker localhost:9092 -topic my-topic -auth`)
		fmt.Println("------------------------------------------------------------------------------")
		fmt.Println("sub command either pub (publish) or sub (subsriber)")
		fmt.Println("-b | -broker : kafka brokers (you can add multiple brokers using separated by comma eg: -b localhost:9091,localhost:9092 ..)")
		fmt.Println("-t | -topic : kafka topic")
		fmt.Println("-h : show help")
		fmt.Println("-version : show version")
		fmt.Println("-V : verbose mode")
		fmt.Println("-auth : prompt sasl auth")
		fmt.Println("------------------------------------------------------------------------------")
	}

	flag.Parse()

	if len(os.Args) < 2 {
		argument.Help = flag.Usage
		return argument, ErrorRequiredOneArgument
	}

	if !strings.Contains(os.Args[1], "version") {
		switch os.Args[1] {
		case "pub":
			publishCommand.Parse(os.Args[2:])
		case "sub":
			subscribeCommand.Parse(os.Args[2:])
		default:
			argument.Help = flag.Usage
			return argument, ErrorInvalidCommand
		}
	}

	// publish command parsed
	if publishCommand.Parsed() {

		if len(brokers) <= 0 {
			argument.Help = flag.Usage
			return argument, ErrorRequiredOneBroker
		}

		if len(topic) <= 0 {
			argument.Help = flag.Usage
			return argument, ErrorRequiredTopicName
		}

		if len(message) <= 0 {
			argument.Help = flag.Usage
			return argument, ErrorPubRequredMessage
		}

		for _, broker := range brokers {
			if broker == "" {
				continue
			}

			argument.Brokers = append(argument.Brokers, strings.Trim(broker, " "))
		}

		argument.Topic = topic
		argument.Message = []byte(message)
		argument.Command = CommandFromString(os.Args[1])
	}

	// subscribe command parsed
	if subscribeCommand.Parsed() {

		if len(brokers) <= 0 {
			argument.Help = flag.Usage
			return argument, ErrorRequiredOneBroker
		}

		if len(topic) <= 0 {
			argument.Help = flag.Usage
			return argument, ErrorRequiredTopicName
		}

		for _, broker := range brokers {
			if broker == "" {
				continue
			}

			argument.Brokers = append(argument.Brokers, strings.Trim(broker, " "))
		}

		argument.Topic = topic
		argument.Command = CommandFromString(os.Args[1])
	}

	if auth {
		var (
			authFields = []string{
				"username: ",
				"password: ",
			}

			responses []string
		)

		authReader := bufio.NewReader(os.Stdin)

		for _, authField := range authFields {

			fmt.Printf("%s", authField)

			response, err := authReader.ReadString('\n')
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			response = strings.TrimSpace(response)

			responses = append(responses, response)
		}

		fmt.Println()

		argument.Username = responses[0]
		argument.Password = responses[1]
	}

	argument.ShowVersion = showVersion
	argument.Verbose = verbose
	argument.Auth = auth

	return argument, nil
}
