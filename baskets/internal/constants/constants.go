package constants

// ServiceName The name of this module/service
const ServiceName = "baskets"

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
	DatabaseKey                   = "database"
	AggregateStoreTxKey           = "aggregateStoreTx"
	AggregateStoreKey             = "aggregateStore"
	MessagePublisherTxKey         = "messagePublisherTx"
	MessagePublisherKey           = "messagePublisher"
	MessageSubscriberKey          = "messageSubscriber"
	EventPublisherKey             = "eventPublisher"
	CommandPublisherKey           = "commandPublisher"
	ReplyPublisherKey             = "replyPublisher"
	SagaStoreKey                  = "sagaStore"
	InboxStoreKey                 = "inboxStore"
	ApplicationKey                = "app"
	ApplicationTxKey              = "appTx"
	DomainEventHandlersKey        = "domainEventHandlers"
	IntegrationEventHandlersTxKey = "integrationEventHandlersTx"
	IntegrationEventHandlersKey   = "integrationEventHandlers"
	CommandHandlersKey            = "commandHandlers"
	ReplyHandlersKey              = "replyHandlers"

	BasketsRepoTxKey = "basketsTxRepo"
	BasketsRepoKey   = "basketsRepo"
	StoresRepoTxKey  = "storesRepoTx"
	StoresRepoKey    = "storesRepo"
	ProductsRepoKey  = "productsRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	StoresCacheTableName   = ServiceName + ".stores_cache"
	ProductsCacheTableName = ServiceName + ".products_cache"
)

// Metric Names
const (
	BasketsStartedCount    = "baskets_started_count"
	BasketsCheckedOutCount = "baskets_checked_out_count"
	BaksetsCanceledCount   = "baskets_canceled_count"
)
