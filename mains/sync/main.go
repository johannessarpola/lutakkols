package main

import (
	"github.com/johannessarpola/lutakkols/cmd/sync"
	"time"
)

func main() {
	c := sync.RunConfig{
		SourceURL:      "https://www.jelmu.net",
		EventsFn:       ".data/events_test.json",
		EventDetailsFn: ".data/event_details_test.json",
		Timeout:        30 * time.Second,
		RateLimit:      1 * time.Second,
		EventLimit:     10,
		Verbose:        true,
	}
	sync.Run(c)
}
