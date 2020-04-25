package api

import (
	"net/http"

	"github.com/phunki/actionspanel/pkg/log"
	"github.com/spf13/cobra"
)

var (
	// Cmd is the exported cobra command which starts the webhook handler service
	Cmd = &cobra.Command{
		Use:   "api",
		Short: "Runs the web service that runs our web api for handling GitHub events",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
)

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Infof("Server is listening...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Err(err, "couldn't shutdown cleanly")
	}
}
