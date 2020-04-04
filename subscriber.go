package kafka

import (
	"context"
)

//Subscriber interface
type Subscriber interface {
	Subscribe(context.Context, ...string) error
}
