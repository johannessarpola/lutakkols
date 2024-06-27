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
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*300)
		defer cancel()
		eventResults, errs := as.Events(op, ctx)
		events := make(chan models.EventPartial)
		detailResults := as.Details(events, ctx)

		details := make(chan models.EventDetailsPartial)
		asciiResults := as.Images(details, ctx)

		for {
			select {
			case <-ctx.Done():
				fmt.Println("context done")
				return
			case err := <-errs:
				fmt.Println("main loop - error:", err)
				return
			case e := <-eventResults:
				if e.Err == nil {
					fmt.Printf("forwarding event %s\n", e.Val.ID())
					events <- e.Val
				} else {
					fmt.Println("main loop - event err", e.Err)
				}
			case d := <-detailResults:
				if d.Err == nil {
					fmt.Printf("forwarding details %s\n", d.Val.EventID)
					details <- d.Val
				} else {
					fmt.Println("main loop - details err", d.Err)
				}
			case a := <-asciiResults:
				if a.Err == nil {
					fmt.Printf("got ascii for %s\n", a.Val.EventID)
					fmt.Println(a.Val)
				} else {
					fmt.Println("main loop - ascii err", a.Err)
				}
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
