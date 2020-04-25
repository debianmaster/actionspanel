package sleep

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	// Cmd is the exported cobra command which sleeps
	Cmd = &cobra.Command{
		Use:   "sleep",
		Short: "Sleeps for 30s. Required because scratch containers don't have a sleep binary",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
)

func main() {
	time.Sleep(30 * time.Second)
}
