package ctx_aggregator

import (
	"context"
)

/*
RegisterBaseContextAggregator register new baseAggregator[T] into context
by package key combined with input keys.
In order to use multiple aggregator, use different input keys.
Example:

	ctx = RegisterBaseContextAggregator[int](ctx, key)
	func () {err := Collect(ctx, value1)}()
	func () {err := Collect(ctx, value2)}()
	result, err = Aggregate[int](ctx, key) // result = []int{value1, value2}
*/
func RegisterBaseContextAggregator[T any](ctx context.Context, keys ...string) context.Context {
	agg := &baseAggregator[T]{
		datas: make([]T, 0),
	}
	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

type baseAggregator[T any] struct {
	datas []T
}

// Collect[T] collect data into aggregator by key
func Collect[T any](ctx context.Context, data T, keys ...string) error {
	ctxKey := buildContextKey(keys...)
	agg, ok := ctx.Value(ctxKey).(*baseAggregator[T])
	if !ok {
		return ErrNotFoundOrInvalid
	}
	agg.datas = append(agg.datas, data)
	return nil
}

// Aggregate[T] get data from aggregator by key
func Aggregate[T any](ctx context.Context, keys ...string) ([]T, error) {
	ctxKey := buildContextKey(keys...)
	agg, ok := ctx.Value(ctxKey).(*baseAggregator[T])
	if !ok {
		return nil, ErrNotFoundOrInvalid
	}
	return agg.datas, nil
}
