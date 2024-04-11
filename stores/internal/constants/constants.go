package constants

// ServiceName The name of this module/service
const ServiceName = "stores"

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
	DatabaseKey                 = "database"
	MessagePublisherTxKey       = "messagePublisherTx"
	MessageSubscriberKey        = "messageSubscriber"
	EventPublisherKey           = "eventPublisher"
	CommandPublisherKey         = "commandPublisher"
	ReplyPublisherKey           = "replyPublisher"
	AggregateStoreTxKey         = "aggregateStoreTx"
	AggregateStoreKey           = "aggregateStore"
	SagaStoreKey                = "sagaStore"
	InboxStoreTxKey             = "inboxStoreTx"
	ApplicationKey              = "app"
	ApplicationTxKey            = "appTx"
	DomainEventHandlersKey      = "domainEventHandlers"
	IntegrationEventHandlersKey = "integrationEventHandlers"
	CommandHandlersKey          = "commandHandlers"
	ReplyHandlersKey            = "replyHandlers"

	CatalogHandlersTxKey = "catalogHandlersTx"
	CatalogHandlersKey   = "catalogHandlers"
	MallHandlersTxKey    = "mallHandlersTx"
	MallHandlersKey      = "mallHandlers"

	StoresRepoTxKey   = "storesRepoTx"
	StoresRepoKey     = "storesRepo"
	ProductsRepoTxKey = "productsRepoTx"
	ProductsRepoKey   = "productsRepo"
	CatalogRepoTxKey  = "catalogRepoTx"
	CatalogRepoKey    = "catalogRepo"
	MallRepoTxKey     = "mallRepoTx"
	MallRepoKey       = "mallRepo"
)

// Repository Table Names
const (
	OutboxTableName    = ServiceName + ".outbox"
	InboxTableName     = ServiceName + ".inbox"
	EventsTableName    = ServiceName + ".events"
	SnapshotsTableName = ServiceName + ".snapshots"
	SagasTableName     = ServiceName + ".sagas"

	CatalogTableName = ServiceName + ".products"
	MallTableName    = ServiceName + ".stores"
)
