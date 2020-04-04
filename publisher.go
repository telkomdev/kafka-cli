package kafka

import (
	"context"
)

//Publisher interface
type Publisher interface {
	Publish(context.Context, string, []byte) error
}
