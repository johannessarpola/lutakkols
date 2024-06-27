package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/johannessarpola/lutakkols/cmd/constants"
	"github.com/johannessarpola/lutakkols/cmd/sync"
	"github.com/johannessarpola/lutakkols/internal/views"
	"github.com/johannessarpola/lutakkols/pkg/api/options"
	"github.com/johannessarpola/lutakkols/pkg/api/provider"
	"github.com/spf13/cobra"
	v "github.com/spf13/viper"
	"os"
	"path"
)

type config struct {
	Address string
	Offline bool
	LogFile string
	Verbose bool
}

var Config config

var rootCmd = &cobra.Command{
	Use:   "ui",
	Short: "View Lutakko gigs with CLI",
	Long:  "View Lutakko gigs with CLI",
	Run: func(cmd *cobra.Command, args []string) {

		if v.GetBool("offline") {
			offlineCli("https://www.jelmu.net", ".data")
		} else {
			onlineCli("https://www.jelmu.net")
		}
	},
}

func setupTMUI(p provider.Provider) {

	m := views.NewEventsList(p)
	prog := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := prog.Run(); err != nil {
		fmt.Println("err running program:", err)
		os.Exit(1)
	}
}

func onlineCli(path string) {
	c := provider.Config{
		EventsSourceURL: path,
	}

	p, err := provider.New(&c, options.UseOnline)
	if err != nil {
		panic(err)
	}
	setupTMUI(p)
}

func offlineCli(source string, outputDir string) {
	ep := path.Join(outputDir, constants.EventsFile)
	edp := path.Join(outputDir, constants.EventsDetailsFile)

	config := provider.Config{
		EventSourceFsPath:  ep,
		EventDetailsFsPath: edp,
		AsciiGen:           views.GenerateOfflineAscii,
	}

	p, err := provider.New(&config, options.UseOffline)
	if err != nil {
		panic(err)
	}
	setupTMUI(p)

}

func init() {
	rootCmd.AddCommand(sync.Cmd)

	rootCmd.Flags().StringVarP(&Config.Address, "address", "a", "https://www.jelmu.net", "Server address")
	rootCmd.Flags().BoolVarP(&Config.Offline, "offline", "o", false, "Run in offline mode")
	rootCmd.Flags().StringVarP(&Config.LogFile, "logfile", "l", "debug.log", "File to write log into")
	rootCmd.Flags().BoolVarP(&Config.Verbose, "verbose", "v", false, "Verbose")

	err := v.BindPFlag("address", rootCmd.Flags().Lookup("address"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}
	err = v.BindPFlag("offline", rootCmd.Flags().Lookup("offline"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

	err = v.BindPFlag("logfile", rootCmd.Flags().Lookup("logfile"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

	err = v.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
	if err != nil {
		fmt.Printf("could not bind flag: %v\n", err)
	}

}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
