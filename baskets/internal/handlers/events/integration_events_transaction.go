package events

import (
	"context"
	"database/sql"

	"github.com/rezaAmiri123/edatV2/am"
	"github.com/rezaAmiri123/edatV2/di"
	"github.com/rezaAmiri123/mallbots/baskets/internal/constants"
)

func RegisterIntegrationEventHandlersTx(container di.Container) error {
	rawMsgHandler := am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) (err error) {
		ctx = container.Scoped(ctx)
		defer func(tx *sql.Tx) {
			if p := recover(); p != nil {
				_ = tx.Rollback()
				panic(p)
			} else if err != nil {
				_ = tx.Rollback()
			} else {
				err = tx.Commit()
			}
		}(di.Get(ctx, constants.DatabaseTxKey).(*sql.Tx))

		handler := di.Get(ctx, constants.IntegrationEventHandlersTxKey).(am.MessageHandler)
		return handler.HandleMessage(ctx, msg)
	})

	sunscriber := container.Get(constants.MessageSubscriberKey).(am.MessageSubscriber)

	return RegisterIntegrationEventHandlers(sunscriber, rawMsgHandler)
}
