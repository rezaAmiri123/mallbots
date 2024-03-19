package es

import (
	"context"

	"github.com/rezaAmiri123/edatV2/ddd"
)

type (
	EventSourcedAggregate interface {
		ddd.IDer
		ddd.AggregateNamer
		ddd.Eventer
		Versioner
		EventApplier
		EventCommiter
	}
	AggregateStore interface{
		Load(ctx context.Context, aggregate EventSourcedAggregate)error
		Save(ctx context.Context, aggregate EventSourcedAggregate)error
	}

	AggregateStoreMiddleware func(store AggregateStore)AggregateStore
)

func AggregateStoreWithMiddleware(store AggregateStore, mws ...AggregateStoreMiddleware)AggregateStore{
	s := store
	// middleware are applied in reverse; this makes the first middleware
	// in the slice the outermost i.e. first to enter, last to exit
	// given: store, A, B, C
	// result: A(B(C(store)))
	for i:=len(mws)-1;i>=0;i--{
		s=mws[i](s)
	}
	return s
}
