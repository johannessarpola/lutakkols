package sync2

import (
	"context"
	"fmt"
	"github.com/johannessarpola/lutakkols/pkg/fetch"
	"github.com/johannessarpola/lutakkols/pkg/logger"
	"github.com/johannessarpola/lutakkols/pkg/pipes"
	"github.com/johannessarpola/lutakkols/pkg/writer"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
	"sync"
	"time"
)

type config struct {
	SourceURL string
	OutputDir string
}

var Config config

var TestCmd = &cobra.Command{
	Use:   "abcd",
	Short: "Syncs data",
	Long:  "Syncs data",
	Run: func(cmd *cobra.Command, args []string) {

		logger.SetLogger(&logger.StdOutLogger{})

		op := v.GetString("input_url")
		defaultTimeout := time.Second * 120
		as := fetch.AsyncSource{}
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()

		eventResults := as.Events(op, time.Second*1, ctx)
		events := pipes.FilterError(eventResults, func(err error) {
			logger.Log.Error("event error ", err)
		}, ctx)
		e1, e2 := pipes.FanOut(events, ctx)

		var wg sync.WaitGroup

		wg.Add(1)
		wr1 := writer.WriteChannel(pipes.Materialize(e2, ctx), ".data/events.json", defaultTimeout)

		detailResults := as.Details(e1, ctx)
		details := pipes.FilterError(detailResults, func(err error) {
			logger.Log.Warn("details error ", err)
		}, ctx)

		wg.Add(1)
		wr2 := writer.WriteChannel(pipes.Materialize(details, ctx), ".data/event_details.json", defaultTimeout)

	main:
		for {
			select {
			case r1 := <-wr1:
				if r1.Err != nil {
					logger.Log.Error("could not write events", r1.Err)
				}
				wg.Done()
				break main
			case r2 := <-wr2:
				if r2.Err != nil {
					logger.Log.Error("could not write details", r2.Err)
				}
				wg.Done()
				break main
			}
		}
		wg.Wait()

	},
}

func init() {

	TestCmd.Flags().StringVarP(&Config.SourceURL, "input_url", "i", "https://www.jelmu.net", "EventURL to source data")
	TestCmd.Flags().StringVarP(&Config.OutputDir, "output_dir", "o", ".data", "Output directory to write to")

	err := v.BindPFlag("input_url", TestCmd.Flags().Lookup("input_url"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}
	err = v.BindPFlag("output_dir", TestCmd.Flags().Lookup("output_dir"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

}
