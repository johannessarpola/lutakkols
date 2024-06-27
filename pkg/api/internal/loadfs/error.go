package loadfs

import "fmt"

type notFoundException struct {
	ID   string
	Path string
}

func (e notFoundException) Error() string {
	return fmt.Sprintf("Could not find element with id %s in path %s", e.ID, e.Path)
}
