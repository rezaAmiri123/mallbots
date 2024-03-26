package agent

import (
	"context"
	"database/sql"

	"github.com/rezaAmiri123/edatV2/am"
	amserializer "github.com/rezaAmiri123/edatV2/am/serializer"
	"github.com/rezaAmiri123/edatV2/amotel"
	"github.com/rezaAmiri123/edatV2/amprom"
	"github.com/rezaAmiri123/edatV2/di"
	edatlog "github.com/rezaAmiri123/edatV2/log"
	"github.com/rezaAmiri123/edatV2/postgres"
	"github.com/rezaAmiri123/edatV2/postgresotel"
	"github.com/rezaAmiri123/edatV2/registry"
	"github.com/rezaAmiri123/edatV2/tm"
	"github.com/rezaAmiri123/mallbots/customers/internal/constants"
	"github.com/rezaAmiri123/mallbots/customers/internal/handlers/events"
)

func (a *Agent) setupEventHandler() (err error) {
	sentCounter := amprom.SentMessagesCounter(constants.ServiceName)
	a.container.AddScoped(constants.MessagePublisherKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		outboxStore := postgres.NewOutboxStore(constants.OutboxTableName, tx)
		return am.NewMessagePublisher(
			c.Get(constants.StreamKey).(am.MessageStream),
			amotel.OtelMessageContextInjector(),
			sentCounter,
			tm.OutboxPublisher(outboxStore),
		), nil
	})

	
	a.container.AddScoped(constants.EventPublisherKey, func(c di.Container) (any, error) {
		return am.NewEventPublisher(
			c.Get(constants.RegistryKey).(registry.Registry),
			a.getAMSerializer(),
			c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
		), nil
	})

	a.container.AddScoped(constants.InboxStoreKey, func(c di.Container) (any, error) {
		tx := postgresotel.Trace(c.Get(constants.DatabaseTransactionKey).(*sql.Tx))
		return postgres.NewInboxStore(constants.InboxTableName, tx), nil
	})

	a.container.AddScoped(constants.DomainEventHandlersKey, func(c di.Container) (any, error) {
		return events.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
	})
	events.RegisterDomainEventHandlersTx(a.container)

	go func() {
		outboxProcessor := tm.NewOutboxProcessor(
			a.container.Get(constants.StreamKey).(am.MessageStream),
			postgres.NewOutboxStore(
				constants.OutboxTableName,
				a.container.Get(constants.DatabaseTransactionKey).(*sql.Tx),
			),
		)
		// TODO make a gracefull shutdown
		err := outboxProcessor.Start(context.Background())
		logger := edatlog.DefaultLogger
		if err != nil {
			logger.Error("customers outbox processor encountered an error", edatlog.Error(err))
		}
	}()

	return nil
}

func(a *Agent)getAMSerializer()am.MessageSerializer{
	switch a.config.SerdeType{
	default:
		return amserializer.NewJsonSerializer()
	}
}