package baskets

import (
	"context"
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/go-chi/chi/v5"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/rezaAmiri123/edatV2/am"
	amserializer "github.com/rezaAmiri123/edatV2/am/serializer"
	"github.com/rezaAmiri123/edatV2/amotel"
	"github.com/rezaAmiri123/edatV2/amprom"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/edatV2/es"
	edatgrpc "github.com/rezaAmiri123/edatV2/grpc"
	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/edatV2/postgresotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/registry/serdes"
	"github.com/rezaAmiri123/edatV2/stream/jetstream"
	jsserializer "github.com/rezaAmiri123/edatV2/stream/jetstream/serializer"
	"github.com/rezaAmiri123/edatV2/tm"
	"github.com/rezaAmiri123/mallbots/baskets/basketspb"
	"github.com/rezaAmiri123/mallbots/baskets/internal/adapters"
	"github.com/rezaAmiri123/mallbots/baskets/internal/application"
	"github.com/rezaAmiri123/mallbots/baskets/internal/constants"
	"github.com/rezaAmiri123/mallbots/baskets/internal/domain"
	"github.com/rezaAmiri123/mallbots/baskets/internal/handlers/events"
	"github.com/rezaAmiri123/mallbots/baskets/internal/handlers/grpcserver"
	"github.com/rezaAmiri123/mallbots/cmd/system"
	"github.com/rezaAmiri123/mallbots/internal/config"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
	"github.com/rs/zerolog"
	"github.com/stackus/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
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
	r.esBasketRepo()

	r.rpcClients()
	r.storeRepo()
	r.productRepo()

	r.application()

	if err = grpcserver.RegisterServerTx(r.container, r.svc.RPC()); err != nil {
		return err
	}

	r.jsStream()
	r.domainEventHandler()
	events.RegisterDomainEventHandlersTx(r.container)

	r.inboxStore()
	r.integrationEventHandler()
	if err = events.RegisterIntegrationEventHandlersTx(r.container); err != nil {
		return err
	}

	return nil
}

func (r *root) registry() {
	r.container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.NewRegistry()
		JsonSerde := serdes.NewJsonSerde(reg)
		if err := domain.RegistrationsWithSerde(JsonSerde); err != nil {
			return nil, err
		}
		if err := basketspb.RegistrationsWithSerde(JsonSerde); err != nil {
			return nil, err
		}
		if err := storespb.RegistrationsWithSerde(JsonSerde); err != nil {
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

func (r *root) esBasketRepo() {
	r.container.AddScoped(constants.BasketsRepoTxKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Basket](
			domain.BasketAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreTxKey).(es.AggregateStore),
		), nil
	})
	r.container.AddScoped(constants.BasketsRepoKey, func(c di.Container) (any, error) {
		return es.NewAggregateRepository[*domain.Basket](
			domain.BasketAggregate,
			c.Get(constants.RegistryKey).(registry.Registry),
			c.Get(constants.AggregateStoreKey).(es.AggregateStore),
		), nil
	})
}

func (r *root) application() {
	// Prometheus counters
	basketsStarted := promauto.NewCounter(prometheus.CounterOpts{
		Name: constants.BasketsStartedCount,
	})
	basketsCheckedOut := promauto.NewCounter(prometheus.CounterOpts{
		Name: constants.BasketsCheckedOutCount,
	})

	r.container.AddScoped(constants.ApplicationTxKey, func(c di.Container) (any, error) {
		return application.NewInstrumentedApp(application.New(
			c.Get(constants.BasketsRepoTxKey).(domain.BasketRepository),
			c.Get(constants.StoresRepoTxKey).(domain.StoreCacheRepository),
			c.Get(constants.ProductsRepoTxKey).(domain.ProductCacheRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), basketsStarted, basketsCheckedOut), nil
	})
	// setup application
	r.container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.NewInstrumentedApp(application.New(
			c.Get(constants.BasketsRepoKey).(domain.BasketRepository),
			c.Get(constants.StoresRepoKey).(domain.StoreCacheRepository),
			c.Get(constants.ProductsRepoKey).(domain.ProductCacheRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.Event]),
		), basketsStarted, basketsCheckedOut), nil
	})
}

func (r *root) domainEventHandler() {
	r.container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return events.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
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

func (r *root) rpcClients() error {
	const (
		backoffLinear  = 100 * time.Millisecond
		backoffRetries = 3
	)

	ctx := context.Background()

	clientErrorUnaryInterceptor := func() grpc.UnaryClientInterceptor {
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return errors.ReceiveGRPCError(invoker(ctx, method, req, reply, cc, opts...))
		}
	}

	getGRPCClient := func(ctx context.Context, addr string, clientTLSConfig *tls.Config) (conn *grpc.ClientConn, err error) {
		retryOpts := []grpc_retry.CallOption{
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(backoffLinear)),
			grpc_retry.WithCodes(codes.NotFound, codes.Aborted),
			grpc_retry.WithMax(backoffRetries),
		}

		var opts []grpc.DialOption
		opts = append(opts, grpc.WithChainUnaryInterceptor(
			grpc_retry.UnaryClientInterceptor(retryOpts...),
			// otelgrpc.UnaryClientInterceptor(),
			clientErrorUnaryInterceptor(),
			edatgrpc.RequestContextUnaryClientInterceptor,
			edatgrpc.WithUnrayClientLogging(r.svc.Logger()),
		))
		if clientTLSConfig != nil {
			clientCreds := credentials.NewTLS(clientTLSConfig)
			opts = append(opts, grpc.WithTransportCredentials(clientCreds))
		} else {
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}

		return grpc.DialContext(ctx, addr, opts...)
	}

	config := r.svc.Config()

	storeClient, err := getGRPCClient(ctx, config.Stores.Address, nil)
	if err != nil {
		return err
	}

	r.container.AddSingleton(constants.StoresServiceName, func(c di.Container) (any, error) {
		return storeClient, nil
	})

	return nil
}

func (r *root) storeRepo() {
	conn := r.container.Get(constants.StoresServiceName).(*grpc.ClientConn)
	client := storespb.NewStoresServiceClient(conn)
	r.container.AddScoped(constants.StoresRepoTxKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresStoreCacheRepository(
			constants.StoresCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx)),
			adapters.NewGrpcStoreRepository(client),
		), nil
	})
	r.container.AddScoped(constants.StoresRepoKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresStoreCacheRepository(
			constants.StoresCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB)),
			adapters.NewGrpcStoreRepository(client),
		), nil
	})
}

func (r *root) productRepo() {
	conn := r.container.Get(constants.StoresServiceName).(*grpc.ClientConn)
	client := storespb.NewStoresServiceClient(conn)
	r.container.AddScoped(constants.ProductsRepoTxKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresProductCacheRepository(
			constants.ProductsCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx)),
			adapters.NewGrpcProductRepository(client),
		), nil
	})
	r.container.AddScoped(constants.ProductsRepoKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresProductCacheRepository(
			constants.ProductsCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB)),
			adapters.NewGrpcProductRepository(client),
		), nil
	})
}

func (r *root) integrationEventHandler() {
	r.container.AddScoped(constants.IntegrationEventHandlersTxKey, func(c di.Container) (any, error) {
		return events.NewIntegrationHandlers(
			c.Get(constants.RegistryKey).(registry.Registry),
			r.amSerializer,
			c.Get(constants.StoresRepoTxKey).(domain.StoreCacheRepository),
			c.Get(constants.ProductsRepoTxKey).(domain.ProductCacheRepository),
			tm.InboxHandler(c.Get(constants.InboxStoreTxKey).(tm.InboxStore)),
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
