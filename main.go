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
	db "github.com/grkmk/esp-server/database"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Color: hclog.AutoColor,
		Level: hclog.DefaultLevel,
	})

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading env variables", "error", err)
		return
	}

	portAddr := os.Getenv("PORT_ADDR")
	httpServer := initHttpServer(logger, portAddr)

	var hostAddr string
	if os.Getenv("ENV") == "dev" {
		hostAddr = "0.0.0.0" // docker container ip
	} else {
		hostAddr = os.Getenv("POSTGRES_CONTAINER")
	}

	dbInstance := db.InitDB(&db.DBSettings{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     hostAddr,
		DBName:   os.Getenv("POSTGRES_DB"),
		Port:     os.Getenv("POSTGRES_PORT"),
	}, logger)

	trapShutdownSig(logger)
	shutdownHttpServer(httpServer)
	dbInstance.Close()
}

func initHttpServer(logger hclog.Logger, portAddr string) *http.Server {
	serveMux := mux.NewRouter()

	getRouter := serveMux.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) { rw.Write([]byte("hello\n")) })

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
			return
		}

		http.ListenAndServe(portAddr, serveMux)
	}()

	return httpServer
}

func trapShutdownSig(logger hclog.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	<-sigChan
	logger.Info("Received terminate, gracefully shutting down")
}

func shutdownHttpServer(httpServer *http.Server) {
	timeoutCtx, cancelTimeout := context.WithTimeout(context.Background(), 30*time.Second)
	httpServer.Shutdown(timeoutCtx)
	defer cancelTimeout()
}
