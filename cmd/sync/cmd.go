package sync

import (
	"fmt"
	"github.com/johannessarpola/go-lutakko-gigs/cmd/constants"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/api/provider"
	"github.com/johannessarpola/go-lutakko-gigs/pkg/logger"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
	"path"
)

type config struct {
	SourceURL string
	OutputDir string
}

var Config config

var Cmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs data",
	Long:  "Syncs data",
	Run: func(cmd *cobra.Command, args []string) {

		ep := path.Join(v.GetString("output_dir"), constants.EventsFile)
		edp := path.Join(v.GetString("output_dir"), constants.EventsDetailsFile)
		op := v.GetString("input_url")
		placeholderGen := func(_ string, _ string) string { return "" }

		c := provider.Config{
			EventsSourceURL:    op,
			EventSourceFsPath:  ep,
			EventDetailsFsPath: edp,
			AsciiGen:           placeholderGen,
		}

		downloader, err := provider.NewDownloader(&c)
		if err != nil {
			panic(err)
		}
		err = downloader.Download()
		if err != nil {
			logger.Log.Debugf("could not update %s", err.Error())
		}
	},
}

func init() {

	Cmd.Flags().StringVarP(&Config.SourceURL, "input_url", "i", "https://www.jelmu.net", "EventURL to source data")
	Cmd.Flags().StringVarP(&Config.OutputDir, "output_dir", "o", ".data", "Output directory to write to")

	err := v.BindPFlag("input_url", Cmd.Flags().Lookup("input_url"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}
	err = v.BindPFlag("output_dir", Cmd.Flags().Lookup("output_dir"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

}
