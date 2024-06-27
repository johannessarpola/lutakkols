package options

// ProviderOption is used to control how the providers operate
type ProviderOption int

// TypeOption to handle provider kind
type TypeOption int

const (
	_ ProviderOption = iota
	SkipCache
)

const (
	_ TypeOption = iota
	UseOffline
	UseOnline
)

// Has check if option is in the option list
func Has[T comparable](option T, opts []T) bool {
	for _, opt := range opts {
		if opt == option {
			return true
		}
	}
	return false
}
