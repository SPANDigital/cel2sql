package dialect

import (
	"fmt"
	"sync"
)

var (
	registryMu sync.RWMutex
	registry   = make(map[Name]func() Dialect)
)

// Register registers a dialect factory function by name.
// This is typically called in an init() function by each dialect package.
func Register(name Name, factory func() Dialect) {
	registryMu.Lock()
	defer registryMu.Unlock()
	registry[name] = factory
}

// Get returns a new dialect instance by name.
// Returns an error if the dialect is not registered.
func Get(name Name) (Dialect, error) {
	registryMu.RLock()
	defer registryMu.RUnlock()
	factory, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("%w: dialect %q is not registered", ErrUnsupportedFeature, name)
	}
	return factory(), nil
}

// Registered returns the names of all registered dialects.
func Registered() []Name {
	registryMu.RLock()
	defer registryMu.RUnlock()
	names := make([]Name, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
