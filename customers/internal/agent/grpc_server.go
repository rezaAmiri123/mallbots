package agent

import (
	"context"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rezaAmiri123/edatV2/di"
	edatgrpc "github.com/rezaAmiri123/edatV2/grpc"
	// edatpgx "github.com/rezaAmiri123/edatV2/pgx"
	edatlog "github.com/rezaAmiri123/edatV2/log"
	"github.com/rs/zerolog"
	"github.com/stackus/errors"
	// "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	_ "cloud.google.com/go/compute/metadata"
)

// Try go get cloud.google.com/go/compute/metadata and then go mod tidy

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (a *Agent) setupGrpcServer() error {
	logger := a.container.Get(constants.LoggerKey).(zerolog.Logger)
	poolConn := a.container.Get(constants.DatabaseKey).(*pgxpool.Pool)
	var opts []grpc.ServerOption
	opts = append(opts,
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			edatgrpc.RequestContextUnaryServerInterceptor,
			WithServerUnaryEnsureStatus(),
			edatgrpc.WithUnrayServerLogging(logger),
			edatpgx.RpcSessionUnrayInterceptor(poolConn, edatlog.DefaultLogger),
			grpc_ctxtags.UnaryServerInterceptor(),
			////////grpc_opentracing.UnaryServerInterceptor(),
			// otelgrpc.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		)),
		// grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	server := grpc.NewServer(opts...)
	reflection.Register(server)
	grpc_prometheus.Register(server)
	app := a.container.Get(constants.ApplicationKey).(application.ServiceApplication)
	grpcserver.RegisterServer(app, server)

	listener, err := net.Listen(a.config.Rpc.Network, a.config.Rpc.Address)
	if err != nil {
		return err
	}

	a.container.AddSingleton(constants.GRPCServerKey, func(c di.Container) (any, error) {
		return server, nil
	})

	go func() {
		logger.Info().Msgf("grpc server started at %s", listener.Addr())
		err = server.Serve(listener)
		if err != nil {
			_ = a.Shutdown()
		}
	}()
	return err
}

func WithServerUnaryEnsureStatus() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, errors.SendGRPCError(err)
	}
}
