package kafka

import "errors"

var (
	// ErrorInvalidCommand error
	ErrorInvalidCommand = errors.New("invalid command")

	// ErrorRequiredOneArgument error
	ErrorRequiredOneArgument = errors.New("require at least one arguments")

	// ErrorRequiredOneBroker error
	ErrorRequiredOneBroker = errors.New("require at least one broker")

	// ErrorPubRequredMessage error
	ErrorPubRequredMessage = errors.New("pub command require a message")

	// ErrorRequiredTopicName error
	ErrorRequiredTopicName = errors.New("require a topic name")
)
