package cosec

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
	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/edatV2/postgresotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/registry/serdes"
	"github.com/rezaAmiri123/edatV2/sec"
	"github.com/rezaAmiri123/edatV2/stream/jetstream"
	jsserializer "github.com/rezaAmiri123/edatV2/stream/jetstream/serializer"
	"github.com/rezaAmiri123/edatV2/tm"
	"github.com/rezaAmiri123/mallbots/cmd/system"
	"github.com/rezaAmiri123/mallbots/cosec/internal"
	"github.com/rezaAmiri123/mallbots/cosec/internal/constants"
	"github.com/rezaAmiri123/mallbots/cosec/internal/handlers/events"
	"github.com/rezaAmiri123/mallbots/cosec/internal/models"
	"github.com/rezaAmiri123/mallbots/customers/customerspb"
	"github.com/rezaAmiri123/mallbots/internal/config"
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
	r.registry() //
	r.database() //
	r.dispatcher()
	r.inboxStore()
	// r.CustomerRepo() //CustomersRepoKey

	// r.application() //

	r.jsStream() //command

	// Domain Event
	r.domainEventHandler()
	// events.RegisterDomainEventHandlersTx(r.container)

	// Command
	r.commandHandler()
	// events.RegisterCommandHandlersTx(r.container)

	return nil
}

func (r *root) registry() {
	r.container.AddSingleton(constants.RegistryKey, func(c di.Container) (_ any, err error) {
		reg := registry.NewRegistry()
		jsonSerde := serdes.NewJsonSerde(reg)
		// Saga data
		if err = jsonSerde.RegisterKey(internal.CreateOrderSagaName, models.CreateOrderData{}); err != nil {
			return nil, err
		}
		if err = customerspb.RegistrationsWithSerde(jsonSerde); err != nil {
			return nil, err
		}
		if err = orderingpb.RegistrationsWithSerde(jsonSerde); err != nil {
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
		return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
	})
}

// func (r *root) CustomerRepo() {
// 	r.container.AddScoped(constants.CustomersRepoTxKey, func(c di.Container) (any, error) {
// 		return adapters.NewPostgresCustomerRepository(
// 			constants.CustomersTableName,
// 			postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx)),
// 		), nil
// 	})
// 	r.container.AddScoped(constants.CustomersRepoKey, func(c di.Container) (any, error) {
// 		return adapters.NewPostgresCustomerRepository(
// 			constants.CustomersTableName,
// 			postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB)),
// 		), nil
// 	})
// }

// func (r *root) application() {
// 	r.container.AddScoped(constants.ApplicationTxKey, func(c di.Container) (any, error) {
// 		return application.New(
// 			c.Get(constants.CustomersRepoTxKey).(domain.CustomerRepository),
// 			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.AggregateEvent]),
// 		), nil
// 	})
// 	// setup application
// 	r.container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
// 		return application.New(
// 			c.Get(constants.CustomersRepoKey).(domain.CustomerRepository),
// 			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.AggregateEvent]),
// 		), nil
// 	})
// }

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

func (r *root) commandHandler() {
	r.container.AddScoped(constants.ReplyPublisherTxKey, func(c di.Container) (any, error) {
		return am.NewReplyPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			r.amSerializer,
			c.Get(constants.MessagePublisherTxKey).(am.MessagePublisher),
		), nil
	})

	r.container.AddScoped(constants.CommandHandlersKey, func(c di.Container) (any, error) {
		return events.NewCommandHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.ApplicationTxKey).(application.App),
			c.Get(constants.ReplyPublisherTxKey).(am.ReplyPublisher),
			r.amSerializer,
			tm.InboxHandler(c.Get(constants.InboxStoreKey).(tm.InboxStore)),
		), nil
	})
}
func (r *root) inboxStore() {
	r.container.AddScoped(constants.InboxStoreTxKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx))
		return postgres.NewInboxStore(constants.InboxTableName, tx), nil
	})
	r.container.AddSingleton(constants.InboxStoreKey, func(c di.Container) (any, error) {
		db := postgresotel.Trace(r.svc.DB())
		return postgres.NewInboxStore(constants.InboxTableName, db), nil
	})
}

func(r *root)saga(){
	r.container.AddScoped(constants.SagaStoreKey, func(c di.Container) (any, error) {
		reg := c.Get(constants.RegistryKey).(registry.Registry)
		return sec.NewSagaRepository[*models.CreateOrderData](
			reg,
			postgres.NewSagaStore(
				constants.SagasTableName,
				postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx)),
				reg,
			),
		), nil
	})
	r.container.AddSingleton(constants.SagaKey, func(c di.Container) (any, error) {
		return internal.NewCreateOrderSaga(), nil
	})

	// setup application
	r.container.AddScoped(constants.OrchestratorKey, func(c di.Container) (any, error) {
		return sec.NewOrchestrator[*models.CreateOrderData](
			c.Get(constants.SagaKey).(sec.Saga[*models.CreateOrderData]),
			c.Get(constants.SagaStoreKey).(sec.SagaRepository[*models.CreateOrderData]),
			c.Get(constants.CommandPublisherKey).(am.CommandPublisher),
		), nil
	})

}