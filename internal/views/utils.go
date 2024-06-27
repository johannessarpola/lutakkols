package views

import (
	"github.com/johannessarpola/go-lutakko-gigs/internal/views/spinner"
	"time"
)

func formatDatetime(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func newSpinner() spinner.Model {
	n := spinner.New()
	n.Spinner = spinner.LutakkoSpinner
	return n
}
