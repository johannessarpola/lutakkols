package sync2

import (
	"context"
	"fmt"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/fetch"
	"github.com/johannessarpola/lutakkols/pkg/logger"
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
		defaultTimeout := time.Second * 8
		as := fetch.AsyncSource{}
		ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
		eventResults, errs := as.Events(op, ctx)
		events := fetch.FilterError(eventResults, func(err error) {
			fmt.Println("event error ", err)
		}, ctx)
		e1, e2 := fetch.FanOut(events, ctx)

		var wg sync.WaitGroup

		go func(events <-chan models.Event) {
			wg.Add(1)
			// write events
			subCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			fmt.Println("collecting ...")
			all, err := fetch.Collect(events, subCtx)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Printf("writing into file size %d\n", len(all))
			err = writer.WriteJson(all, ".data/evts.json", writer.PrettyPrint)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(e2)

		detailResults := as.Details(e1, ctx)

		details := fetch.FilterError(detailResults, func(err error) {
			fmt.Println("details error ", err)
		}, ctx)

		d1, d2 := fetch.FanOut(details, ctx)

		go func(details <-chan models.EventDetails) {
			wg.Add(1)
			// write details
			subCtx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			fmt.Println("collecting ...")
			all, err := fetch.Collect(details, subCtx)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("writing into file size %d\n", len(all))
			err = writer.WriteJson(all, ".data/dtls.json", writer.PrettyPrint)
			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}(d2)

		asciiResults := as.Images(d1, ctx)
		ascii := fetch.FilterError(asciiResults, func(err error) {
			fmt.Println("ascii error ", err)
		}, ctx)

	consume:
		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done")
				break consume
			case err := <-errs:
				fmt.Println("main loop - error:", err)
				break consume
			case a := <-ascii:
				fmt.Printf("main loop - got ascii %s\n", a.EventID)
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
