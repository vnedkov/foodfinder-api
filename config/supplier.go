package config

// Supplier is a configuration supplier.
type Supplier interface {
	Get(key string) (string, bool)
}
