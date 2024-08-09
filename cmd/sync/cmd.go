package sync

import (
	"context"
	"fmt"
	"github.com/johannessarpola/lutakkols/cmd/constants"
	"github.com/johannessarpola/lutakkols/pkg/fetch"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"github.com/johannessarpola/lutakkols/pkg/pipes"
	"github.com/johannessarpola/lutakkols/pkg/writer"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
	"path"
	"time"
)

type RunConfig struct {
	SourceURL      string
	EventsFn       string
	EventDetailsFn string
	Timeout        time.Duration
	RateLimit      time.Duration
	EventLimit     int
	Verbose        bool
}

func Run(conf RunConfig) {
	start := time.Now()
	if conf.Verbose {
		logger.SetLogger(&logger.StdOutLogger{})
	}

	timeout := conf.Timeout
	logger.Log.Infof("Starting sync with timeout %v against URL %s writing events to %s and details to %s", timeout, conf.SourceURL, conf.EventsFn, conf.EventDetailsFn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	events := fetch.Async.Events(conf.SourceURL, conf.EventLimit, ctx)
	e1, e2 := pipes.FanOut(events, ctx)

	logger.Log.Infof("Writing events into %s", conf.EventsFn)

	eventWriteChan := make(chan pipes.Result[bool])
	detailsWriteChan := make(chan pipes.Result[bool])

	eventWriteChan = writer.WriteChannel(e1, conf.EventsFn, timeout)

	rateLimitedEvents := pipes.ThrottleChannel(e2, time.Second, ctx)
	detailResults := fetch.Async.Details(rateLimitedEvents, ctx)
	details := pipes.FilterError(detailResults, func(err error) {
		logger.Log.Warn("details error ", err)
	}, ctx)
	detailsWriteChan = writer.WriteChannel(details, conf.EventDetailsFn, timeout)

	var dwr1, dwr2 bool

	for {
		select {
		case eventWrite := <-eventWriteChan:
			if eventWrite.Err != nil {
				logger.Log.Error("Could not write events", eventWrite.Err)
			}
			if !dwr1 {
				logger.Log.Infof("Events written successfully to %s", conf.EventsFn)
			}
			dwr1 = true
		case detailsWrite := <-detailsWriteChan:
			if detailsWrite.Err != nil {
				logger.Log.Error("Could not write event details", detailsWrite.Err)
			}
			if !dwr2 {
				logger.Log.Infof("Event details written successfully to %s", conf.EventDetailsFn)
			}
			dwr2 = true
		}
		if dwr1 && dwr2 {
			break
		}
	}
	logger.Log.Infof("Doneso in %d ms!", time.Since(start).Milliseconds())
}

var Cmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs data",
	Long:  "Syncs data",
	Run: func(cmd *cobra.Command, args []string) {

		ep := path.Join(v.GetString("output_dir"), constants.EventsFile)
		edp := path.Join(v.GetString("output_dir"), constants.EventsDetailsFile)
		op := v.GetString("input_url")
		to := v.GetDuration("timeout")
		rl := v.GetDuration("rate_limit")
		el := v.GetInt("event_limit")
		verbose := v.GetBool("verbose")

		c := RunConfig{
			SourceURL:      op,
			EventsFn:       ep,
			EventDetailsFn: edp,
			Timeout:        to,
			RateLimit:      rl,
			Verbose:        verbose,
			EventLimit:     el,
		}

		Run(c)
	},
}

func init() {

	Cmd.Flags().StringP("input_url", "i", "https://www.jelmu.net", "EventURL to source data")
	Cmd.Flags().StringP("output_dir", "o", ".data", "Output directory to write to")
	Cmd.Flags().DurationP("timeout", "t", time.Second*120, "timeout for synchronization task")
	Cmd.Flags().DurationP("rate_limit", "r", time.Second*1, "ratelimiter for requests")
	Cmd.Flags().IntP("event_limit", "l", 0, "limit on how mnay events to fetch")

	err := v.BindPFlag("input_url", Cmd.Flags().Lookup("input_url"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

	err = v.BindPFlag("output_dir", Cmd.Flags().Lookup("output_dir"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

	err = v.BindPFlag("timeout", Cmd.Flags().Lookup("timeout"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

	err = v.BindPFlag("rate_limit", Cmd.Flags().Lookup("rate_limit"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

	err = v.BindPFlag("event_limit", Cmd.Flags().Lookup("event_limit"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

}
