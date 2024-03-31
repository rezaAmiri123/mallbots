package grpcserver

import (
	"context"
	"database/sql"

	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/customers/customerspb"
	"github.com/rezaAmiri123/mallbots/customers/internal/application"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
	"google.golang.org/grpc"
)

type serverTx struct {
	customerspb.UnimplementedCustomersServiceServer
	c di.Container
}

var _ customerspb.CustomersServiceServer = (*serverTx)(nil)

func RegisterServerTx(container di.Container, register grpc.ServiceRegistrar) {
	customerspb.RegisterCustomersServiceServer(register, serverTx{
		c: container,
	})
}

func (s serverTx) RegisterCustomer(ctx context.Context, request *customerspb.RegisterCustomerRequest) (resp *customerspb.RegisterCustomerResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationTxKey).(application.App)}
	
	// next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.RegisterCustomer(ctx, request)
}

func (s serverTx) GetCustomer(ctx context.Context, request *customerspb.GetCustomerRequest) (resp *customerspb.GetCustomerResponse, err error) {
	ctx = s.c.Scoped(ctx)
	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}

	return next.GetCustomer(ctx, request)
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
