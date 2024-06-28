package sync2

import (
	"context"
	"fmt"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
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
			fmt.Println("event error ", err)
		}, ctx)
		e1, e2 := pipes.FanOut(events, ctx)

		var wg sync.WaitGroup
		wg.Add(2) // add two for each output (file)

		go func(events <-chan models.Event) {
			// write events
			subCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			all, err := pipes.Collect(events, subCtx)
			if err != nil {
				fmt.Println(err)
			}

			err = writer.WriteJson(all, ".data/evts.json", writer.PrettyPrint)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(pipes.Materialize(e2, ctx))

		detailResults := as.Details(e1, ctx)

		details := pipes.FilterError(detailResults, func(err error) {
			fmt.Println("details error ", err)
		}, ctx)

		//d1, d2 := fetch.FanOut(details, ctx)

		go func(details <-chan models.EventDetails) {
			// write details
			subCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			all, err := pipes.Collect(details, subCtx)
			if err != nil {
				fmt.Println(err)
			}

			err = writer.WriteJson(all, ".data/dtls.json", writer.PrettyPrint)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(pipes.Materialize(details, ctx))

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
