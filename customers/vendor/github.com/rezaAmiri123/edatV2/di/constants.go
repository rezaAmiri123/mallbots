package di

// modeled after this excellent DI lib: https://github.com/sarulabs/di
type Scope int

const(
	Singleton Scope = iota + 1
	Scoped
)

type contextKey int

const containerKey contextKey = 1

