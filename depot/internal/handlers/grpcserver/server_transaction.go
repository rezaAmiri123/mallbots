package grpcserver

import (
	"context"
	"database/sql"

	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/depot/depotpb"
	"github.com/rezaAmiri123/mallbots/depot/internal/application"
	"github.com/rezaAmiri123/mallbots/depot/internal/constants"
	"google.golang.org/grpc"
)

type serverTx struct {
	c di.Container
	depotpb.UnimplementedDepotServiceServer
}

var _ depotpb.DepotServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, register grpc.ServiceRegistrar) error {
	depotpb.RegisterDepotServiceServer(register, serverTx{
		c: container,
	})
	return nil
}

func (s serverTx) CreateShoppingList(ctx context.Context, request *depotpb.CreateShoppingListRequest) (resp *depotpb.CreateShoppingListResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationTxKey).(application.App)}

	return next.CreateShoppingList(ctx, request)
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
