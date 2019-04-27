package kafka

import (
	"context"
)

// Runner type
type Runner struct {
	Publisher  Publisher
	Subscriber Subscriber
	Argument   *Argument
}

// Run function
func (r *Runner) Run(ctx context.Context) error {
	command := r.Argument.Command

	switch command {
	case PublishCommand:
		return r.Publisher.Publish(r.Argument.Topic, r.Argument.Message)
	case SubscribeCommand:
		return r.Subscriber.Subscribe(ctx, r.Argument.Topic)
	}
	return nil
}
