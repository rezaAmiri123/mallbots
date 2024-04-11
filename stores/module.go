package stores

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
	"github.com/rezaAmiri123/mallbots/stores/internal/adapters"
	"github.com/rezaAmiri123/mallbots/stores/internal/application"
	"github.com/rezaAmiri123/mallbots/stores/internal/constants"
	"github.com/rezaAmiri123/mallbots/stores/internal/domain"
	"github.com/rezaAmiri123/mallbots/stores/internal/handlers/events"
	"github.com/rezaAmiri123/mallbots/stores/internal/handlers/grpcserver"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
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
	// return Root(ctx, mono)
}

type root struct {
	container    di.Container
	ctx          context.Context
	svc          Service
	amSerializer am.MessageSerializer
}

func (r *root) registry() {
	r.container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.NewRegistry()
		serde := serdes.NewJsonSerde(reg)
		if err := domain.Registrations(serde); err != nil {
			return nil, err
		}
		if err := storespb.RegistrationsWithSerde(serde); err != nil {
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

func (r *root) esStoreRepo() {
	r.container.AddScoped(constants.StoresRepoTxKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Store](
			domain.StoreAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreTxKey).(es.AggregateStore),
		), nil
	})
	r.container.AddScoped(constants.StoresRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Store](
			domain.StoreAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})
}
func (r *root) esProductRepo() {
	r.container.AddScoped(constants.ProductsRepoTxKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Product](
			domain.ProductAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreTxKey).(es.AggregateStore),
		), nil
	})
	r.container.AddScoped(constants.ProductsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Product](
			domain.ProductAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})
}

func (r *root) mallRepo() {
	r.container.AddScoped(constants.MallRepoTxKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresMallRepository(
			constants.MallTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx)),
		), nil
	})
	r.container.AddScoped(constants.MallRepoKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresMallRepository(
			constants.MallTableName,
			postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB)),
		), nil
	})
}

func (r *root) catalogRepo() {
	r.container.AddScoped(constants.CatalogRepoTxKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresCatalogRepository(
			constants.CatalogTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx)),
		), nil
	})
	r.container.AddScoped(constants.CatalogRepoKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresCatalogRepository(
			constants.CatalogTableName,
			postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB)),
		), nil
	})
}

func (r *root) application() {
	// setup application
	r.container.AddScoped(constants.ApplicationTxKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.StoresRepoTxKey).(domain.StoreRepository),
			c.Get(constants.MallRepoTxKey).(domain.MallRepository),
			c.Get(constants.ProductsRepoTxKey).(domain.ProductRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
	// setup application
	r.container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.StoresRepoKey).(domain.StoreRepository),
			c.Get(constants.MallRepoKey).(domain.MallRepository),
			c.Get(constants.ProductsRepoKey).(domain.ProductRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), nil
	})
}
func (r *root) domainEventHandler() {
	r.container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return events.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})

}

func (r *root) mallHandler() {
	r.container.AddScoped(constants.MallHandlersTxKey, func(c di.Container) (any, error) {
		return events.NewMallHandlers(c.Get(constants.MallRepoTxKey).(domain.MallRepository)), nil
	})

	r.container.AddScoped(constants.MallHandlersKey, func(c di.Container) (any, error) {
		return events.NewMallHandlers(c.Get(constants.MallRepoKey).(domain.MallRepository)), nil
	})
}

func (r *root) catalogHandler() {
	r.container.AddScoped(constants.CatalogHandlersTxKey, func(c di.Container) (any, error) {
		return events.NewCatalagHandlers(c.Get(constants.CatalogRepoTxKey).(domain.CatalogRepository)), nil
	})

	r.container.AddScoped(constants.CatalogHandlersKey, func(c di.Container) (any, error) {
		return events.NewCatalagHandlers(c.Get(constants.CatalogRepoKey).(domain.CatalogRepository)), nil
	})
}

func (r *root) Startup() (err error) {
	r.registry()
	r.database()
	r.jsStream()

	r.dispatcher()
	r.esAggregateStore()
	r.esStoreRepo()
	r.esProductRepo()
	r.mallRepo()
	r.catalogRepo()
	r.application()

	r.domainEventHandler()
	r.mallHandler()
	r.catalogHandler()

	// }
	// func Root(ctx context.Context, svc Service) (err error) {

	// container.AddScoped(constants.InboxStoreTxKey, func(c di.Container) (any, error) {
	// 	tx := postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx))
	// 	return postgres.NewInboxStore(constants.InboxTableName, tx), nil
	// })

	// setup Driver adapters
	if err = grpcserver.RegisterServerTx(r.container, r.svc.RPC()); err != nil {
		return err
	}

	// if err = storespb.RegisterAsyncAPI(svc.Mux()); err != nil {
	// 	return err
	// }

	events.RegisterDomainEventHandlersTx(r.container)
	events.RegisterMallHandlersTx(r.container)
	events.RegisterCatalagHandlersTx(r.container)

	return nil
}

// func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
// 	go func() {
// 		err := outboxProcessor.Start(ctx)
// 		if err != nil {
// 			logger.Error().Err(err).Msg("stores outbox processor encountered an error")
// 		}
// 	}()
// }
