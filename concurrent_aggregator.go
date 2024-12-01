package ctx_aggregator

import (
	"context"
	"runtime"
	"sync"
)

/*
RegisterCContextAggregator register new concurrentAggregator[T] into context
by package key combined with input keys.
In order to use multiple aggregator, use different input keys.
Example:

	ctx = RegisterBaseContextAggregator[int](ctx, key)
	go Collect(ctx, value1)
	go Collect(ctx, value2)
	result, err = Aggregate[int](ctx, key) // result = []int{value1, value2}
*/
func RegisterCContextAggregator[T any](ctx context.Context, keys ...string) context.Context {
	agg := &concurrentAggregator[T]{
		wg:    sync.WaitGroup{},
		datas: make([]T, 0),
	}
	ctxKey := buildContextKey(keys...)
	return context.WithValue(ctx, ctxKey, agg)
}

type concurrentAggregator[T any] struct {
	wg    sync.WaitGroup
	datas []T
}

func CWait[T any](ctx context.Context, keys ...string) context.Context {
	ctxKey := buildContextKey(keys...)
	agg, ok := ctx.Value(ctxKey).(*concurrentAggregator[T])
	if !ok {
		return ctx
	}

	agg.wg.Add(1)
	runtime.SetFinalizer(ctx, func() {
		agg.wg.Done()
	})
	return ctx
}

// CCollect[T] collect data into aggregator by key
func CCollect[T any](ctx context.Context, data T, keys ...string) error {
	ctxKey := buildContextKey(keys...)
	agg, ok := ctx.Value(ctxKey).(*concurrentAggregator[T])
	if !ok {
		return ErrNotFoundOrInvalid
	}
	agg.datas = append(agg.datas, data)
	return nil
}

// CAggregate[T] get data from aggregator by key
func CAggregate[T any](ctx context.Context, keys ...string) ([]T, error) {
	ctxKey := buildContextKey(keys...)
	agg, ok := ctx.Value(ctxKey).(*concurrentAggregator[T])
	if !ok {
		return nil, ErrNotFoundOrInvalid
	}
	agg.wg.Wait()
	return agg.datas, nil
}
