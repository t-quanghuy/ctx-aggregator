package ctx_aggregator

import (
	"context"
	"errors"
	"sync"
)

type ContextKey string

const (
	contextAggregatorContextKey ContextKey = "ctxAggCtxKey"
)

var (
	InvalidType error = errors.New("Invalid type provided")
)

func NewAggregatorContext[T any](ctx context.Context) context.Context {
	aggr := &Aggregator[T]{
		wg: sync.WaitGroup{},
		ch: make(chan T),
	}

	ctx = context.WithValue(ctx, contextAggregatorContextKey, aggr)

	return ctx
}

type Aggregator[T any] struct {
	wg sync.WaitGroup
	ch chan T
}

func Collect[T any](ctx context.Context, data T) error {
	aggr, ok := ctx.Value(contextAggregatorContextKey).(*Aggregator[T])
	if !ok {
		return InvalidType
	}

	aggr.ch <- data
	return nil
}

func CollectChan[T any](ctx context.Context, ch chan T) error {
	aggr, ok := ctx.Value(contextAggregatorContextKey).(*Aggregator[T])
	if !ok {
		return InvalidType
	}

	aggr.wg.Add(1)
	data := <-ch
	aggr.ch <- data
	aggr.wg.Done()

	return nil
}

func Aggregate[T any](ctx context.Context) ([]T, error) {
	aggr, ok := ctx.Value(contextAggregatorContextKey).(*Aggregator[T])
	if !ok {
		return nil, InvalidType
	}

	aggr.wg.Wait()
	close(aggr.ch)

	result := make([]T, len(aggr.ch))
	for data := range aggr.ch {
		result = append(result, data)
	}

	return result, nil
}
