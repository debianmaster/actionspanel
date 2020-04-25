package api

import (
	"fmt"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
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

type environmentConfiguration struct {
	APIPort int `default:"8080" envconfig:"api_port"`
}

func main() {
	var cfg environmentConfiguration

	// Load environment variables
	err := envconfig.Process("ap", &cfg)
	if err != nil {
		log.Err(err, "couldn't load config")
		os.Exit(1)
	}

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.APIPort),
		Handler: router,
	}

	log.Infof("Server is listening on port %d...", cfg.APIPort)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Err(err, "couldn't shutdown cleanly")
	}
}
