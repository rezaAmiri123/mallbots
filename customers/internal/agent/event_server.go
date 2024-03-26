package agent

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/stream/jetstream"
	amserializer "github.com/rezaAmiri123/edatV2/stream/jetstream/serializer"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
)

func (a *Agent) setupEventServer() (err error) {
	var stream am.MessageStream
	switch a.config.StreamType {
	case "nats":
		stream, err = a.getNatsStream()
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("event server typeis unknown")
	}

	a.container.AddSingleton(constants.StreamKey, func(c di.Container) (any, error) {
		return stream, nil
	})

	return nil
}

func (a *Agent) getNatsStream() (am.MessageStream, error) {
	nc, err := nats.Connect(a.config.Nats.URL)
	if err != nil {
		return nil, err
	}
	// defer nc.Close()
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	_, err = js.AddStream(&nats.StreamConfig{
		Name:     a.config.Nats.Stream,
		Subjects: []string{fmt.Sprintf("%s.>", a.config.Nats.Stream)},
	})
	if err != nil {
		return nil, err
	}

	var serializer jetstream.MessageSerializer

	switch a.config.SerdeType {
	default:
		serializer = amserializer.NewJsonSerializer()
	}

	stream := jetstream.NewStream(a.config.Nats.Stream, js, serializer)

	return stream, nil
}
