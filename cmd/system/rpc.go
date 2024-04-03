package system

import (
	"context"
	"fmt"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	edatgrpc "github.com/rezaAmiri123/edatV2/grpc"
	"github.com/stackus/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type RPC interface {
	RPC() *grpc.Server
}

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s *System) initRpc() {
	// logger := a.container.Get(constants.LoggerKey).(zerolog.Logger)
	// poolConn := a.container.Get(constants.DatabaseKey).(*pgxpool.Pool)
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
			serverErrorUnaryInterceptor(),
			edatgrpc.WithUnrayServerLogging(s.Logger()),
			// edatpgx.RpcSessionUnrayInterceptor(poolConn, edatlog.DefaultLogger),
			grpc_ctxtags.UnaryServerInterceptor(),
			////////grpc_opentracing.UnaryServerInterceptor(),
			// otelgrpc.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpc_recovery.UnaryServerInterceptor(),
		)),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)

	grpc.ChainUnaryInterceptor(
		otelgrpc.UnaryServerInterceptor(),
		serverErrorUnaryInterceptor(),
	)
	// If there are streaming endpoints also add
	// grpc.StreamInterceptor(
	// 	otelgrpc.StreamServerInterceptor(),
	// ),
	s.rpc = grpc.NewServer(opts...)
	reflection.Register(s.rpc)
	grpc_prometheus.Register(s.rpc)
}

func (s *System) RPC() *grpc.Server {
	return s.rpc
}

func (s *System) WaitForRPC(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.cfg.Rpc.Address)
	if err != nil {
		return err
	}

	group, gCtx := errgroup.WithContext(ctx)
	group.Go(func() error {
		fmt.Println("rpc server started")
		defer fmt.Println("rpc server shutdown")
		if err := s.RPC().Serve(listener); err != nil && err != grpc.ErrServerStopped {
			return err
		}
		return nil
	})
	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("rpc server to be shutdown")
		stopped := make(chan struct{})
		go func() {
			s.RPC().GracefulStop()
			close(stopped)
		}()
		timeout := time.NewTimer(s.cfg.ShutdownTimeout)
		select {
		case <-timeout.C:
			// Force it to stop
			s.RPC().Stop()
			return fmt.Errorf("rpc server failed to stop gracefully")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}

func serverErrorUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		return resp, errors.SendGRPCError(err)
	}
}
