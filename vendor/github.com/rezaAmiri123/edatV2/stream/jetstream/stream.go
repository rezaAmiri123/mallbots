package jetstream

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rezaAmiri123/edatV2/am"
	edatlog "github.com/rezaAmiri123/edatV2/log"
)

const maxRetries = 5

type Stream struct {
	streamName string
	js         nats.JetStreamContext
	mu         sync.Mutex
	subs       []*nats.Subscription
	serializer MessageSerializer
	logger     edatlog.Logger
}

var _ am.MessageStream = (*Stream)(nil)

func NewStream(streamName string, js nats.JetStreamContext, serializer MessageSerializer, options ...StreamOption) *Stream {
	stream := &Stream{
		streamName: streamName,
		js:         js,
		serializer: serializer,
		logger:     edatlog.DefaultLogger,
	}

	for _, option := range options {
		option(stream)
	}

	return stream
}
func (s *Stream) Publish(ctx context.Context, topicName string, msg am.Message) (err error) {
	data, err := s.serializer.Serialize(MessageSerializerData{
		ID:       msg.ID(),
		Name:     msg.MessageName(),
		Data:     msg.Data(),
		Metadata: msg.Metadata(),
		SentAt:   msg.SentAt(),
	})
	if err != nil {
		return err
	}

	natsMsg := &nats.Msg{
		Subject: msg.Subject(),
		Data:    data,
	}

	var p nats.PubAckFuture
	p, err = s.js.PublishMsgAsync(natsMsg, nats.MsgId(msg.ID()))
	if err != nil {
		return err
	}

	// retry a handful of times to publish the messages
	go func(future nats.PubAckFuture, tries int) {
		var err error

		for {
			select {
			case <-future.Ok(): // publish acknowledged
				return
			case <-future.Err(): // error ignored; try again
				// TODO add some variable delay between tries
				tries = tries - 1
				if tries <= 0 {
					s.logger.Error(fmt.Sprintf("unable to publish message after %d tries", maxRetries), edatlog.Error(err))
					return
				}
				future, err = s.js.PublishMsgAsync(future.Msg())
				if err != nil {
					// TODO do more than give up
					s.logger.Error("failed to publish a message", edatlog.Error(err))
					return
				}
			}
		}
	}(p, maxRetries)

	return
}

func (s *Stream) Subscribe(topicName string, handler am.MessageHandler, options ...am.SubscriberOption) (am.Subscription, error) {
	var err error

	s.mu.Lock()
	defer s.mu.Unlock()

	subCfg := am.NewSubscriberConfig(options)

	opts := []nats.SubOpt{
		nats.MaxDeliver(subCfg.MaxRedeliver()),
	}
	cfg := &nats.ConsumerConfig{
		MaxDeliver:     subCfg.MaxRedeliver(),
		DeliverSubject: topicName,
		FilterSubject:  topicName,
	}

	if groupName := subCfg.GroupName(); groupName != "" {
		cfg.DeliverSubject = groupName
		cfg.DeliverGroup = groupName
		cfg.Durable = groupName

		opts = append(opts,
			nats.Bind(s.streamName, groupName),
			nats.Durable(groupName),
		)
	}

	if ackType := subCfg.AckType(); ackType != am.AckTypeAuto {
		ackWait := subCfg.AckWait()

		cfg.AckPolicy = nats.AckExplicitPolicy
		cfg.AckWait = ackWait

		opts = append(opts,
			nats.AckExplicit(),
			nats.AckWait(ackWait),
		)
	} else {
		cfg.AckPolicy = nats.AckNonePolicy
		opts = append(opts, nats.AckNone())
	}

	_, err = s.js.AddConsumer(s.streamName, cfg)
	if err != nil {
		return nil, err
	}

	var sub *nats.Subscription
	if groupName := subCfg.GroupName(); groupName == "" {
		sub, err = s.js.Subscribe(topicName, s.handleMsg(subCfg, handler), opts...)
	} else {
		sub, err = s.js.QueueSubscribe(topicName, groupName, s.handleMsg(subCfg, handler), opts...)
	}

	s.subs = append(s.subs, sub)

	return subscription{sub}, nil

}

func (s *Stream) Unsubscribe() error {
	for _, sub := range s.subs {
		if !sub.IsValid() {
			continue
		}
		err := sub.Drain()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Stream) handleMsg(cfg am.SubscriberConfig, handler am.MessageHandler) nats.MsgHandler {
	var filters map[string]struct{}
	if len(cfg.MessageFilters()) > 0 {
		filters = make(map[string]struct{})
		for _, value := range cfg.MessageFilters() {
			filters[value] = struct{}{}
		}
	}

	return func(natsMsg *nats.Msg) {
		var err error

		m, err := s.serializer.Deserialize(natsMsg.Data)
		if err != nil {
			s.logger.Warn("failed to unmarshal the *nats.Msg", edatlog.Error(err))
			return
		}

		if filters != nil {
			if _, exists := filters[m.Name]; !exists {
				err = natsMsg.Ack()
				if err != nil {
					s.logger.Warn("failed to Ack a filtered message", edatlog.Error(err))
				}
				return
			}
		}

		msg := &rawMessage{
			id:         m.ID,
			name:       m.Name,
			subject:    natsMsg.Subject,
			data:       m.Data,
			metadata:   m.Metadata,
			sentAt:     m.SentAt,
			receivedAt: time.Now(),
			acked:      false,
			ackFu:      func() error { return natsMsg.Ack() },
			nackFn:     func() error { return natsMsg.Nak() },
			extendFn:   func() error { return natsMsg.InProgress() },
			killFn:     func() error { return natsMsg.Term() },
		}

		wCtx, cancel := context.WithTimeout(context.Background(), cfg.AckWait())
		defer cancel()

		errc := make(chan error)
		go func() {
			errc <- handler.HandleMessage(wCtx, msg)
		}()

		if cfg.AckType() == am.AckTypeAuto {
			err = msg.Ack()
			if err != nil {
				s.logger.Warn("failed to auto-Ack a message", edatlog.Error(err))
			}
		}

		select {
		case err = <-errc:
			if err == nil {
				if ackErr := msg.Ack(); ackErr != nil {
					s.logger.Warn("failed to auto-Ack a message", edatlog.Error(err))
				}
				return
			}
			s.logger.Error("error while handling message", edatlog.Error(err))
			if nakErr := msg.NAck(); nakErr != nil {
				s.logger.Warn("failed to Nack a message", edatlog.Error(err))
			}
		case <-wCtx.Done():
			// TODO logging?
			return
		}
	}
}
