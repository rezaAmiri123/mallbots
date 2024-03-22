package di

import "context"

func Get(ctx context.Context, key string) any {
	cnt, ok := ctx.Value(containerKey).(*container)
	if !ok {
		panic("container does not exist on context")
	}

	return cnt.Get(key)
}
