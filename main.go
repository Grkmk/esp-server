package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "server",
		Color: hclog.AutoColor,
		Level: hclog.DefaultLevel,
	})

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading env variables", "error", err)
		return
	}

	portAddr := os.Getenv("PORT_ADDR")

	serveMux := mux.NewRouter()
	corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"localhost:3000"}))
	httpServer := &http.Server{
		Addr:         portAddr,
		Handler:      corsHandler(serveMux),
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		IdleTimeout:  120 * time.Second,
		ErrorLog:     logger.StandardLogger(&hclog.StandardLoggerOptions{}),
	}

	go func() {
		logger.Info("Starting server on", "port", portAddr)

		err := httpServer.ListenAndServe()
		if err == http.ErrServerClosed {
			os.Exit(1)
			return
		}

		if err != nil {
			logger.Error("Error starting http server", "error", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	<-sigChan
	logger.Info("Received terminate, gracefully shutting down")

	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 30*time.Second)
	httpServer.Shutdown(timeoutCtx)
	defer cancelTimeout()

	http.ListenAndServe(portAddr, serveMux)
}
