package sync2

import (
	"context"
	"fmt"
	"github.com/johannessarpola/lutakkols/pkg/api/models"
	"github.com/johannessarpola/lutakkols/pkg/fetch"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
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
		op := v.GetString("input_url")

		as := fetch.AsyncSource{}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		eventResults, errs := as.Events(op, ctx)
		events := fetch.FilterError(eventResults, func(err error) {
			fmt.Println("event error ", err)
		}, ctx)
		detailResults := as.Details(events, ctx)

		details := fetch.FilterError(detailResults, func(err error) {
			fmt.Println("details error ", err)
		}, ctx)

		d1, d2 := fetch.FanOut(details, ctx)

		go func(details <-chan *models.EventDetails) {
			// write details
			all, err := fetch.Collect(d2)
			if err != nil {
				fmt.Println(err)
			}

		}(d2)

		asciiResults := as.Images(d1, ctx)
		ascii := fetch.FilterError(asciiResults, func(err error) {
			fmt.Println("ascii error ", err)
		}, ctx)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done")
				return
			case err := <-errs:
				fmt.Println("main loop - error:", err)
				return
			case a := <-ascii:
				fmt.Printf("main loop - got ascii %s\n", a.Ascii)
			}

		}

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
