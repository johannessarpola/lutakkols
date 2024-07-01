package fetch

import "fmt"

type FailedFetch struct {
	url string
	err error
}

func (f FailedFetch) Error() string {
	return fmt.Sprintf("failed to fetch url: %s with %s", f.url, f.err)
}
