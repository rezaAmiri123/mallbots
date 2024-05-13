package sec

import (
	"github.com/rezaAmiri123/edatV2/ddd"
)

type stepResult[T any] struct {
	ctx         *SagaContext[T]
	destination string
	cmd         ddd.Command
	err         error
}
