package internal

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
	"github.com/rezaAmiri123/edatV2/sec"
	"github.com/rezaAmiri123/mallbots/cosec/internal/models"
	"github.com/rezaAmiri123/mallbots/customers/customerspb"
)

const CreateOrderSagaName = "cosec.CreateOrder"
const CreateOrderReplyChannel = "mallbots.cosec.replies.CreateOrder"

type createOrderSaga struct {
	sec.Saga[*models.CreateOrderData]
}

func NewCreateOrderSaga() sec.Saga[*models.CreateOrderData] {
	saga := createOrderSaga{
		Saga: sec.NewSaga[*models.CreateOrderData](CreateOrderSagaName, CreateOrderReplyChannel),
	}

	// 0. -RejectOrder
	saga.AddStep().
		Compensation(saga.rejectOrder)

	// 1. AuthorizeCustomer
	saga.AddStep().
		Action(saga.authorizeCustomer)

	return saga
}

func (s createOrderSaga) rejectOrder(ctx context.Context, data *models.CreateOrderData) (string, ddd.Command, error) {
	return "", nil, nil
}

func (s createOrderSaga) authorizeCustomer(ctx context.Context, data *models.CreateOrderData) (string, ddd.Command, error) {
	cmd := ddd.NewCommand(customerspb.AuthorizeCustomerCommand, &customerspb.AuthorizeCustomer{Id: data.CustomerID})

	return customerspb.CommandChannel, cmd, nil
}
