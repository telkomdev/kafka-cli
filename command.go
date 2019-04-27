package kafka

// Command tyoe
type Command int

const (
	// PublishCommand command
	PublishCommand Command = iota

	//SubscribeCommand command
	SubscribeCommand
)

// String function
func (c Command) String() string {
	switch c {
	case PublishCommand:
		return "publish"
	case SubscribeCommand:
		return "subscribe"
	default:
		panic("command not found")
	}
}

// CommandFromString function
func CommandFromString(c string) Command {
	switch c {
	case "pub":
		return PublishCommand
	case "sub":
		return SubscribeCommand
	default:
		panic("command not found")
	}
}
