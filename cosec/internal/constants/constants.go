package constants

// ServiceName The name of this module/service
const ServiceName = "cosec"

// GRPC Service Names
const (
	StoresServiceName    = "STORES"
	CustomersServiceName = "CUSTOMERS"
)

// Dependency Injection Keys
const (
	RegistryKey                   = "registry"
	DomainDispatcherKey           = "domainDispatcher"
	DatabaseTxKey                 = "tx"
	DatabaseKey                   = "db"
	MessagePublisherTxKey           = "messagePublisherTx"
	MessagePublisherKey           = "messagePublisher"
	MessageSubscriberKey          = "messageSubscriber"
	EventPublisherKey             = "eventPublisher"
	CommandPublisherKey           = "commandPublisher"
	ReplyPublisherTxKey             = "replyPublisherTx"
	ReplyPublisherKey             = "replyPublisher"
	SagaStoreKey                  = "sagaStore"
	InboxStoreTxKey               = "inboxStoreTx"
	InboxStoreKey                 = "inboxStore"
	ApplicationKey                = "app"
	DomainEventHandlersKey        = "domainEventHandlers"
	IntegrationEventHandlersTxKey = "integrationEventHandlersTx"
	IntegrationEventHandlersKey   = "integrationEventHandlers"
	CommandHandlersKey            = "commandHandlers"
	ReplyHandlersKey              = "replyHandlers"

	SagaKey         = "saga"
	OrchestratorKey = "orchestrator"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"
)
