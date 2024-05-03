package grpcserver

import (
	"context"
	"database/sql"

	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/baskets/basketspb"
	"github.com/rezaAmiri123/mallbots/baskets/internal/application"
	"github.com/rezaAmiri123/mallbots/baskets/internal/constants"
	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	basketspb.UnimplementedBasketServiceServer
}

var _ basketspb.BasketServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	basketspb.RegisterBasketServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) StartBasket(ctx context.Context, request *basketspb.StartBasketRequest) (resp *basketspb.StartBasketResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationTxKey).(application.App)}

	return next.StartBasket(ctx, request)
}

func (s serverTx) CheckoutBasket(ctx context.Context, request *basketspb.CheckoutBasketRequest) (resp *basketspb.CheckoutBasketResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationTxKey).(application.App)}

	return next.CheckoutBasket(ctx, request)
}

func (s serverTx) AddItem(ctx context.Context, request *basketspb.AddItemRequest) (resp *basketspb.AddItemResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationTxKey).(application.App)}

	return next.AddItem(ctx, request)
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
