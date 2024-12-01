package ctx_aggregator

import (
	"errors"
	"strings"
)

type contextKey string

const (
	contextAggregatorContextKey contextKey = "ctxAggCtxKey"
)

// Define errors
var (
	ErrNotFoundOrInvalid error = errors.New("missing aggregator or invalid type of data provided")
)

// buildContextKey builds context key from default context key and input keys
func buildContextKey(keys ...string) contextKey {
	if len(keys) == 0 {
		return contextAggregatorContextKey
	}

	keys = append([]string{string(contextAggregatorContextKey)}, keys...)
	return contextKey(strings.Join(keys, "_"))
}
