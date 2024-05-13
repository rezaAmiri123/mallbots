package sec

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
)

type (
	StepActionFunc[T any]       func(ctx context.Context, data T) (string, ddd.Command, error)
	StepReplyHandlerFunc[T any] func(ctx context.Context, data T, reply ddd.Reply) error
)
