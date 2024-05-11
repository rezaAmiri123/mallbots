package depot

import (
	"context"
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/go-chi/chi/v5"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/nats-io/nats.go"
	"github.com/rezaAmiri123/edatV2/am"
	amserializer "github.com/rezaAmiri123/edatV2/am/serializer"
	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/di"
	edatgrpc "github.com/rezaAmiri123/edatV2/grpc"
	"github.com/rezaAmiri123/edatV2/postgresotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/registry/serdes"
	"github.com/rezaAmiri123/mallbots/cmd/system"
	"github.com/rezaAmiri123/mallbots/depot/internal/adapters"
	"github.com/rezaAmiri123/mallbots/depot/internal/application"
	"github.com/rezaAmiri123/mallbots/depot/internal/constants"
	"github.com/rezaAmiri123/mallbots/depot/internal/domain"
	"github.com/rezaAmiri123/mallbots/depot/internal/handlers/grpcserver"
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
	r.shoppingListRepo()
	r.rpcClients()
	r.storeRepo()
	r.productRepo()

	r.application()

	if err = grpcserver.RegisterServerTx(r.container, r.svc.RPC()); err != nil {
		return err
	}

	return nil
}

func (r *root) registry() {
	r.container.AddSingleton(constants.RegistryKey, func(c di.Container) (any, error) {
		reg := registry.NewRegistry()
		jsonSerde := serdes.NewJsonSerde(reg)
		_ = jsonSerde
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

func (r *root) shoppingListRepo() {
	r.container.AddScoped(constants.ShoppingListsRepoTxKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresShoppingListRepository(
			constants.ShoppingListsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx)),
		), nil
	})
	r.container.AddSingleton(constants.ShoppingListsRepoKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresShoppingListRepository(
			constants.ShoppingListsTableName,
			postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB)),
		), nil
	})

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
	r.container.AddScoped(constants.StoresCacheRepoTxKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresStoreCacheRepository(
			constants.StoresCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx)),
			adapters.NewGrpcStoreRepository(client),
		), nil
	})
	r.container.AddSingleton(constants.StoresCacheRepoKey, func(c di.Container) (any, error) {
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
	r.container.AddScoped(constants.ProductsCacheRepoTxKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresProductCacheRepository(
			constants.ProductsCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseTxKey).(*sql.Tx)),
			adapters.NewGrpcProductRepository(client),
		), nil
	})
	r.container.AddSingleton(constants.ProductsCacheRepoKey, func(c di.Container) (any, error) {
		return adapters.NewPostgresProductCacheRepository(
			constants.ProductsCacheTableName,
			postgresotel.Trace(c.Get(constants.DatabaseKey).(*sql.DB)),
			adapters.NewGrpcProductRepository(client),
		), nil
	})
}

func (r *root) application() {
	r.container.AddScoped(constants.ApplicationTxKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.ShoppingListsRepoTxKey).(domain.ShoppingListRepository),
			c.Get(constants.StoresCacheRepoTxKey).(domain.StoreCacheRepository),
			c.Get(constants.ProductsCacheRepoTxKey).(domain.ProductCacheRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.AggregateEvent]),
		), nil
	})
	// setup application
	r.container.AddScoped(constants.ApplicationKey, func(c di.Container) (any, error) {
		return application.New(
			c.Get(constants.ShoppingListsRepoKey).(domain.ShoppingListRepository),
			c.Get(constants.StoresCacheRepoKey).(domain.StoreCacheRepository),
			c.Get(constants.ProductsCacheRepoKey).(domain.ProductCacheRepository),
			c.Get(constants.DomainDispatcherKey).(ddd.EventPublisher[ddd.AggregateEvent]),
		), nil
	})
}
