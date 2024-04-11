package grpcserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/stores/internal/application"
	"github.com/rezaAmiri123/mallbots/stores/internal/constants"
	"github.com/rezaAmiri123/mallbots/stores/storespb"
	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	storespb.UnimplementedStoresServiceServer
}

var _ storespb.StoresServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	storespb.RegisterStoresServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateStore(ctx context.Context, request *storespb.CreateStoreRequest) (resp *storespb.CreateStoreResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationTxKey).(application.App)}

	return next.CreateStore(ctx, request)
}

func (s serverTx) AddProduct(ctx context.Context, request *storespb.AddProductRequest) (resp *storespb.AddProductResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationTxKey).(application.App)}

	return next.AddProduct(ctx, request)
}

func (s serverTx) IncreaseProductPrice(ctx context.Context, request *storespb.IncreaseProductPriceRequest) (resp *storespb.IncreaseProductPriceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationTxKey).(application.App)}

	return next.IncreaseProductPrice(ctx, request)
}

func (s serverTx) GetStore(ctx context.Context, request *storespb.GetStoreRequest) (resp *storespb.GetStoreResponse, err error) {
	ctx = s.c.Scoped(ctx)
	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetStore(ctx, request)
}

func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		fmt.Println("rollback")
		_ = tx.Rollback()
		return err
	} else {
		fmt.Println("commit")
		return tx.Commit()
	}
}
