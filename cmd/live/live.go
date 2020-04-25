package live

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Cmd is the exported cobra command which checks that the service is running
	Cmd = &cobra.Command{
		Use:   "live",
		Short: "Checks for the existence of the liveness marker file",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
)

// LivenessMarkerPath is the file that we check to verify whether or not the
// web service has already exited
const LivenessMarkerPath = "/livemarker"

func main() {
	_, err := os.Open(LivenessMarkerPath)
	if err != nil {
		// Couldn't find file. Application is not likely to be alive right now
		os.Exit(-1)
	}
}
