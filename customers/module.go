package customers

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rezaAmiri123/edatV2/am"
	amserializer "github.com/rezaAmiri123/edatV2/am/serializer"
	"github.com/rezaAmiri123/edatV2/amotel"
	"github.com/rezaAmiri123/edatV2/amprom"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/edatV2/postgresotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/stream/jetstream"
	jsserializer "github.com/rezaAmiri123/edatV2/stream/jetstream/serializer"
	"github.com/rezaAmiri123/edatV2/tm"
	"github.com/rezaAmiri123/mallbots/cmd/system"
	"github.com/rezaAmiri123/mallbots/customers/customerspb"
	"github.com/rezaAmiri123/mallbots/customers/internal/adapters"
	"github.com/rezaAmiri123/mallbots/customers/internal/application"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
	"github.com/rezaAmiri123/mallbots/customers/internal/domain"
	"github.com/rezaAmiri123/mallbots/customers/internal/handlers/events"
	"github.com/rezaAmiri123/mallbots/customers/internal/handlers/grpcserver"
	"github.com/rs/zerolog"
)

type monoService interface {
	system.System
}
type Module struct{}

func (m Module) Startup(ctx context.Context, mono system.Service) (err error) {
	return Root(ctx, mono)
}

func Root(ctx context.Context, svc system.Service) (err error) {
	fmt.Println("****************"," customer root is running")
	cfg := svc.Config()
	amSerializer := amserializer.NewJsonSerializer()
	container := di.New()
	// setup Driven adapters
	container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.NewRegistry()
		if err := customerspb.Registrations(reg); err != nil {
			return nil, err
		}
		return reg, nil
	})
	
	stream := jetstream.NewStream(cfg.Nats.Stream, svc.JS(),jsserializer.NewJsonSerializer())
	container.AddSingleton(constants.DomainDispatcherKey, func(c di.Container) (any, error) {
		return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
	})
	container.AddScoped(constants.DatabaseTransactionKey, func(c di.Container) (any, error) {
		return svc.DB().Begin()
	})
	container.AddScoped(constants.CustomersRepoKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresCustomerRepository(
			constants.CustomersTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
		), nil
	})
	sentCounter := amprom.SentMessagesCounter(constants.ServiceName)
	container.AddScoped(constants.MessagePublisherKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		outboxStore := postgres.NewOutboxStore(constants.OutboxTableName, tx)
		return am.NewMessagePublisher(
			stream,
			amotel.OtelMessageContextInjector(),
			sentCounter,
			tm.OutboxPublisher(outboxStore),
		), nil
	})
	container.AddSingleton(constants.MessageSubscriberKey, func(c di.Container) (any, error) {
		return am.NewMessageSubscriber(
			stream,
			amotel.OtelMessageContextExtractor(),
			amprom.ReceivedMessageCounter(constants.ServiceName),
		), nil
	})
	container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		return am.NewEventPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
			amSerializer,
		), nil
	})
	container.AddScoped(constants.ReplyPublisherKey, func(c di.Container) (any, error) {
		return am.NewReplyPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			amSerializer,
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
		), nil
	})
	container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return postgres.NewInboxStore(constants.InboxTableName, tx), nil
	})
	// Prometheus counters
	customersRegistered := promauto.NewCounter(prometheus.CounterOpts{
		Name: constants.CustomersRegisteredCount,
	})

	// setup application
	container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.NewInstrumentedApp(application.New(
			c.Get(constants.CustomersRepoKey).(domain.CustomerRepository),
			c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent]),
		), customersRegistered), nil
	})
	container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return events.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		return events.NewCommandHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationKey).(application.App),
			c.Get(constants.ReplyPublisherKey).(am.ReplyPublisher),
			amSerializer,
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})
	outboxProcessor := tm.NewOutboxProcessor(
		stream,
		postgres.NewOutboxStore(constants.OutboxTableName, svc.DB()),
	)

	// setup Driver adapters
	if err = grpcserver.RegisterServerTx(container, svc.RPC()); err != nil {
		return err
	}
	// if err = rest.RegisterGateway(ctx, svc.Mux(), svc.Config().Rpc.Address()); err != nil {
	// 	return err
	// }
	// if err = rest.RegisterSwagger(svc.Mux()); err != nil {
	// 	return err
	// }
	events.RegisterDomainEventHandlersTx(container)
	if err = events.RegisterCommandHandlersTx(container); err != nil {
		return err
	}
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil

}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, logger zerolog.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			logger.Error().Err(err).Msg("customers outbox processor encountered an error")
		}
	}()
}
