package options

// ProviderOption is used to control how the providers operate
type ProviderOption int

// WriteOpton to handle output controls
type WriteOpton int

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

const (
	_ WriteOpton = iota
	PrettyPrint
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
