package constants

// ServiceName The name of this module/service
const ServiceName = "ordering"

// GRPC Service Names
const (
	StoresServiceName    = "STORES"
	CustomersServiceName = "CUSTOMERS"
)

// Dependency Injection Keys
const (
	RegistryKey                 = "registry"
	DomainDispatcherKey         = "domainDispatcher"
	DatabaseTxKey               = "tx"
	DatabaseKey                 = "db"
	MessagePublisherTxKey      = "messagePublisherTx"
	MessagePublisherKey         = "messagePublisher"
	MessageSubscriberKey        = "messageSubscriber"
	EventPublisherTxKey           = "eventPublisherTx"
	EventPublisherKey           = "eventPublisher"
	CommandPublisherKey         = "commandPublisher"
	ReplyPublisherKey           = "replyPublisher"
	SagaStoreKey                = "sagaStore"
	InboxStoreKey               = "inboxStore"
	ApplicationTxKey            = "appTx"
	ApplicationKey              = "app"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"

	AggregateStoreTxKey = "aggregateStoreTx"
	AggregateStoreKey   = "aggregateStore"
	OrdersRepoTxKey     = "ordersTxRepo"
	OrdersRepoKey       = "ordersRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"
)
