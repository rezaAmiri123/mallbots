package grpcserver

import (
	"context"
	"database/sql"

	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/ordering/internal/application"
	"github.com/rezaAmiri123/mallbots/ordering/internal/constants"
	"github.com/rezaAmiri123/mallbots/ordering/orderingpb"
	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	orderingpb.UnimplementedOrderingServiceServer
}

var _ orderingpb.OrderingServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, register grpc.ServiceRegistrar) error {
	orderingpb.RegisterOrderingServiceServer(register, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateOrder(ctx context.Context, request *orderingpb.CreateOrderRequest) (_ *orderingpb.CreateOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.CreateOrder(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
