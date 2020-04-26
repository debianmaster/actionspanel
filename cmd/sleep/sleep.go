// Package sleep is a small application that sleeps
//
// The reason we have this package is because we ultimately deploy into a
// scratch container, which doesn't have /bin/sleep which can be useful to use
// as a preStop hook to allow for zero down time deployments with slow ingress
// controllers.
package sleep

import (
	"time"

	"github.com/spf13/cobra"
)

// Cmd is the exported cobra command which sleeps
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sleep",
		Short: "Sleeps for 30s",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
}

const sleepDuration = 30 * time.Second

func main() {
	time.Sleep(sleepDuration)
}
