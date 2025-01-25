package config

import (
	"os"
)

type envSupplier struct{}

// Get retrieves the value of the environment variable with the given key.
func (e *envSupplier) Get(key string) (string, bool) {
	return os.LookupEnv(key)
}

// NewEnvSupplier creates a new Supplier that reads from the environment.
func NewEnvSupplier() Supplier {
	return &envSupplier{}
}
