// Package api is the package that houses our api application
//
// This is where we configure and instantiate our running web application service.
package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	health "github.com/AppsFlyer/go-sundheit"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/julienschmidt/httprouter"
	"github.com/palantir/go-githubapp/githubapp"
	"github.com/phunki/actionspanel/pkg/config"
	"github.com/phunki/actionspanel/pkg/gh"
	"github.com/phunki/actionspanel/pkg/log"
	"github.com/spf13/cobra"
)

const shutddownSignalBufferSize = 1

// Cmd is the exported cobra command which starts the webhook handler service
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "api",
		Short: "Runs the web service that runs our web api for handling GitHub events",
		Run: func(cmd *cobra.Command, args []string) {
			main()
		},
	}
}

func newHealthServer(cfg config.Config) *http.Server {
	h := health.New()
	healthRouter := httprouter.New()
	healthRouter.HandlerFunc("GET", "/health", healthhttp.HandleHealthJSON(h))

	return &http.Server{
		Handler: healthRouter,
		Addr:    fmt.Sprintf(":%d", cfg.HealthServerPort),
	}
}

func newServer(cfg config.Config) *http.Server {
	githubConfig := gh.NewGitHubConfig(cfg)
	githubClientCreator := gh.NewGitHubClientCreator(cfg)

	sessionManager := config.NewSessionManagerFactory(cfg).CreateSessionManager()
	loginHandler := gh.NewLoginHandler(sessionManager, true, githubConfig, githubClientCreator)

	installationHandler := gh.NewInstallationHandler(githubClientCreator)
	webhookHandler := githubapp.NewDefaultEventDispatcher(githubConfig, installationHandler)
	router := httprouter.New()
	router.Handler("POST", "/webhook", webhookHandler)
	router.Handler("GET", "/api/auth/github", loginHandler)

	router.HandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World!"))
		if err != nil {
			log.Err(err, "failed to write response")
		}
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
		Handler: router,
	}
}

func waitForShutdown(healthSrv *http.Server, srv *http.Server, idleConnsClosed chan struct{}) {
	sigint := make(chan os.Signal, shutddownSignalBufferSize)

	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)

	<-sigint

	err := healthSrv.Shutdown(context.Background())
	if err != nil {
		log.Err(err, "couldn't shutdown health server")
	}

	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Err(err, "couldn't shutdown server")
	}

	close(idleConnsClosed)
}

func main() {
	cfg := config.NewConfig()
	healthSrv := newHealthServer(cfg)
	srv := newServer(cfg)
	idleConnsClosed := make(chan struct{})

	go waitForShutdown(healthSrv, srv, idleConnsClosed)

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
		panic("Couldn't create liveness marker")
	}

	defer func() {
		err := os.Remove(livenessFile.Name())
		if err != nil {
			log.Err(err, "Failed to remove liveness marker")
		}
	}()

	log.Infof("Server is listening on port %d...", cfg.ServerPort)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Err(err, "couldn't shutdown cleanly")
	}

	<-idleConnsClosed
	log.Infof("Shutting down server...")
}
