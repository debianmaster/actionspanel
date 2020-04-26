package main

import (
	"github.com/phunki/actionspanel/cmd/api"
	"github.com/phunki/actionspanel/cmd/live"
	"github.com/phunki/actionspanel/cmd/sleep"
	"github.com/phunki/actionspanel/pkg/log"
	"github.com/spf13/cobra"
)

func main() {
	var (
		rootCmd = &cobra.Command{
			Use:   "actionspanel",
			Short: "Trigger GitHub Actions with a convenient UI",
		}
	)

	rootCmd.AddCommand(api.Cmd())
	rootCmd.AddCommand(live.Cmd())
	rootCmd.AddCommand(sleep.Cmd())

	err := rootCmd.Execute()
	if err != nil {
		log.Err(err, "failed to execute rootCmd")
	}
}
