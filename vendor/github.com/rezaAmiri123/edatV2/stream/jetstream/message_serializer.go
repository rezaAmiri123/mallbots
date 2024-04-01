package jetstream

import (
	"time"

	"github.com/rezaAmiri123/edatV2/ddd"
)

type (
	MessageSerializer interface {
		Serialize(MessageSerializerData) ([]byte, error)
		Deserialize(data []byte) (MessageSerializerData, error)
	}

	MessageSerializerData struct {
		ID       string       `json:"id"`
		Name     string       `json:"name"`
		Data     []byte       `json:"data"`
		Metadata ddd.Metadata `json:"metadata"`
		SentAt   time.Time    `json:"sent_at"`
	}
)
