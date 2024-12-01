package ctx_aggregator

import (
	"context"
)

// RegisterBaseContextAggregator init new baseAggregator[T] pointer
// and register value to context by package key. To use multiple context
// aggregators, use different keys in input parameters.
func RegisterBaseContextAggregator[T any](ctx context.Context, keys ...string) context.Context {
	aggr := &baseAggregator[T]{
		datas: make([]T, 0),
	}
	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, aggr)
}

type baseAggregator[T any] struct {
	datas []T
}

// Collect[T] collect data and append to datas slice
// of baseAggregator
func Collect[T any](ctx context.Context, data T, keys ...string) error {
	ctxKey := buildContextKey(keys...)
	agg, ok := ctx.Value(ctxKey).(*baseAggregator[T])
	if !ok {
		return ErrNotFoundOrInvalid
	}
	agg.datas = append(agg.datas, data)
	return nil
}

// Aggregate[T] return collected datas
func Aggregate[T any](ctx context.Context, keys ...string) ([]T, error) {
	ctxKey := buildContextKey(keys...)
	agg, ok := ctx.Value(ctxKey).(*baseAggregator[T])
	if !ok {
		return nil, ErrNotFoundOrInvalid
	}
	return agg.datas, nil
}
