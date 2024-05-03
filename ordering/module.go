package ordering

import (
	"context"
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/rezaAmiri123/edatV2/am"
	amserializer "github.com/rezaAmiri123/edatV2/am/serializer"
	"github.com/rezaAmiri123/edatV2/amotel"
	"github.com/rezaAmiri123/edatV2/amprom"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/es"
	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/edatV2/postgresotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/registry/serdes"
	"github.com/rezaAmiri123/edatV2/stream/jetstream"
	jsserializer "github.com/rezaAmiri123/edatV2/stream/jetstream/serializer"
	"github.com/rezaAmiri123/edatV2/tm"
	"github.com/rezaAmiri123/mallbots/cmd/system"
	"github.com/rezaAmiri123/mallbots/internal/config"
	"github.com/rezaAmiri123/mallbots/ordering/internal/application"
	"github.com/rezaAmiri123/mallbots/ordering/internal/constants"
	"github.com/rezaAmiri123/mallbots/ordering/internal/domain"
	"github.com/rezaAmiri123/mallbots/ordering/internal/handlers/events"
	"github.com/rezaAmiri123/mallbots/ordering/internal/handlers/grpcserver"
	"github.com/rezaAmiri123/mallbots/ordering/orderingpb"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

type Service interface {
	Config() config.Config
	DB() *sql.DB
	JS() nats.JetStreamContext
	Mux() *chi.Mux
	RPC() *grpc.Server
	Logger() zerolog.Logger
	// Waiter() waiter.Waiter
}

type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	r := &root{
		container:    di.New(),
		ctx:          ctx,
		svc:          mono,
		amSerializer: amserializer.NewJsonSerializer(),
	}
	return r.Startup()
}

type root struct {
	container    di.Container
	ctx          context.Context
	svc          Service
	amSerializer am.MessageSerializer
}

func (r *root) Startup() (err error) {
	r.registry()
	r.database()
	r.dispatcher()
	r.esAggregateStore()
	r.esOrderRepo()

	r.application()

	if err = grpcserver.RegisterServerTx(r.container, r.svc.RPC()); err != nil {
		return err
	}

	r.jsStream()
	r.domainEventHandler()
	events.RegisterDomainEventHandlersTx(r.container)

	return nil
}

func (r *root) registry() {
	r.container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.NewRegistry()
		JsonSerde := serdes.NewJsonSerde(reg)
		if err := domain.RegistrationsWithSerde(JsonSerde); err != nil {
			return nil, err
		}
		if err := orderingpb.RegistrationsWithSerde(JsonSerde); err != nil {
			return nil, err
		}
		return reg, nil
	})
}

func (r *root) database() {
	r.container.AddScoped(constants.DatabaseTxKey, func(c di.Container) (any, error) {
		return r.svc.DB().Begin()
	})
	r.container.AddScoped(constants.DatabaseKey, func(c di.Container) (any, error) {
		return r.svc.DB(), nil
	})
}

func (r *root) dispatcher() {
	r.container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.Event](), nil
	})
}

func (r *root) esAggregateStore() {
	r.container.AddScoped(constants.AggregateStoreTxKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.AggregateStoreWithMiddleware(
			postgres.NewEventStore(constants.EventsTableName, tx, reg),
			postgres.NewSnapshotStore(constants.SnapshotsTableName, tx, reg),
		), nil
	})
	r.container.AddScoped(constants.AggregateStoreKey, func(c di.Container) (any, error) {
		db := postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB))
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return es.AggregateStoreWithMiddleware(
			postgres.NewEventStore(constants.EventsTableName, db, reg),
			postgres.NewSnapshotStore(constants.SnapshotsTableName, db, reg),
		), nil
	})
}

func (r *root) esOrderRepo() {
	r.container.AddScoped(constants.OrdersRepoTxKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Order](
			domain.OrderAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreTxKey).(es.AggregateStore),
		), nil
	})
	r.container.AddScoped(constants.OrdersRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Order](
			domain.OrderAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})
}

func (r *root) application() {
	r.container.AddScoped(constants.ApplicationTxKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.OrdersRepoTxKey).(domain.OrderRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
	// setup application
	r.container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.OrdersRepoKey).(domain.OrderRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
}

func (r *root) jsStream() {
	stream := jetstream.NewStream(r.svc.Config().Nats.Stream, r.svc.JS(), jsserializer.NewJsonSerializer())

	sentCounter := amprom.SentMessagesCounter(constants.ServiceName)
	r.container.AddScoped(constants.MessagePublisherTxKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx))
		outboxStore := postgres.NewOutboxStore(constants.OutboxTableName, tx)
		return am.NewMessagePublisher(
			stream,
			amotel.OtelMessageContextInjector(),
			sentCounter,
			tm.OutboxPublisher(outboxStore),
		), nil
	})

	r.container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessageCounter(constants.ServiceName),
		), nil
	})
	r.container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		return am.NewEventPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherTxKey).(am.MessagePublisher),
			r.amSerializer,
		), nil
	})

	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		postgres.NewOutboxStore(constants.OutboxTableName, r.svc.DB()),
	)

	go func() {
		logger := r.svc.Logger()
		err := outboxProcessor.Start(r.ctx)
		if err != nil {
			logger.Error().Err(err).Msg("stores outbox processor encountered an error")
		}
	}()
}

func (r *root) domainEventHandler() {
	r.container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return events.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
}
