package views

import (
	"github.com/johannessarpola/lutakkols/internal/views/spinner"
)

func newSpinner() spinner.Model {
	n := spinner.New()
	n.Spinner = spinner.LutakkoSpinner
	return n
}
