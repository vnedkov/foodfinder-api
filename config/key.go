package config

// Key is a configuration key.
type Key string

// Get returns the value of the key from the Supplier.
func (name Key) Get(s Supplier) string {
	if value, exists := s.Get(string(name)); exists {
		return value
	}
	return ""
}
