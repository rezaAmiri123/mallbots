package internal

import "github.com/rezaAmiri123/edatV2/sec"

const CreateOrderSagaName = "cosec.CreateOrder"
const CreateOrderReplyChannel = "mallbots.cosec.replies.CreateOrder"

type createOrderSaga struct{
	sec.Saga
}