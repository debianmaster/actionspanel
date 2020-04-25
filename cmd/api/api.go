package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	health "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/julienschmidt/httprouter"
	"github.com/kelseyhightower/envconfig"
	"github.com/phunki/actionspanel/pkg/config"
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
	var cfg config.Config

	// Load environment variables
	err := envconfig.Process("ap", &cfg)
	if err != nil {
		log.Err(err, "couldn't load config")
		os.Exit(1)
	}

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.Write([]byte("Hello World!"))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}

	h := health.New()
	healthRouter := httprouter.New()
	healthRouter.HandlerFunc("GET", "/health", healthhttp.HandleHealthJSON(h))

	healthSrv := &http.Server{
		Handler: healthRouter,
		Addr:    fmt.Sprintf(":%d", cfg.HealthServerPort),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint
		healthSrv.Shutdown(context.Background())
		srv.Shutdown(context.Background())
		close(idleConnsClosed)
	}()

	go func() {
		log.Infof("Starting health server on port %d...", cfg.HealthServerPort)
		if err := healthSrv.ListenAndServe(); err != http.ErrServerClosed {
			log.Infof("Failed to shutdown health server gracefully: %v", err)
		}
		log.Infof("Shutting down health server...")
	}()

	livenessFile, err := os.Create("/livemarker")
	if err != nil {
		log.Infof("failed to create liveness marker: %v", err)
		os.Exit(1)
	}
	defer os.Remove(livenessFile.Name())

	log.Infof("Server is listening on port %d...", cfg.ServerPort)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Err(err, "couldn't shutdown cleanly")
	}

	<-idleConnsClosed
	log.Infof("Shutting down server...")
}
