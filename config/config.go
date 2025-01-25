package config

import (
	"context"
	"sync"
)

var (
	globalOnce       sync.Once
	singletonDefault Supplier
)

// Global returns the default Supplier.
func Global() Supplier {
	globalOnce.Do(func() {
		singletonDefault = NewEnvSupplier()
	})
	return singletonDefault
}

// configCtxKey is the context key for the Supplier
type configCtxKey struct{}

// WithContext returns a new context with the given Supplier attached.
func WithContext(ctx context.Context, s Supplier) context.Context {
	return context.WithValue(ctx, configCtxKey{}, s)
}

// FromContext returns the Supplier attached to the context, or the default Supplier if none is attached.
func FromContext(ctx context.Context) Supplier {
	if s, ok := ctx.Value(configCtxKey{}).(Supplier); ok {
		return s
	}
	return Global()
}
