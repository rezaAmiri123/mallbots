package jetstream

import edatlog "github.com/rezaAmiri123/edatV2/log"

type StreamOption func(stream *Stream)

func WithLogger(logger edatlog.Logger)StreamOption{
	return func(stream *Stream) {
		stream.logger = logger
	}
}