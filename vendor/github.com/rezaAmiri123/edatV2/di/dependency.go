package di

type (
	DepFactoryFunc func(c Container) (any,error)

	depInfo struct{
		key string
		scope Scope
		factory DepFactoryFunc
	}
)