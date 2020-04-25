package main

import (
	"github.com/phunki/actionspanel/cmd/api"
	"github.com/phunki/actionspanel/cmd/live"
	"github.com/phunki/actionspanel/cmd/sleep"
	"github.com/spf13/cobra"
)

func main() {
	var (
		rootCmd = &cobra.Command{
			Use:   "actionspanel",
			Short: "Trigger GitHub Actions with a convenient UI",
		}
	)

	rootCmd.AddCommand(api.Cmd)
	rootCmd.AddCommand(live.Cmd)
	rootCmd.AddCommand(sleep.Cmd)
	rootCmd.Execute()
}
