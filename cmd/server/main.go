package main

import (
	"context"
	"film-management/cmd/server/commands/migrate"
	"film-management/config"
	httpCommonHandler "film-management/internal/common/transport/http"
	"film-management/internal/user/domain"
	"film-management/internal/user/endpoints"
	httpUserHandler "film-management/internal/user/transport/http"
	"film-management/pkg/database/postgresql"
	"film-management/pkg/logger"
	userRepo "film-management/repositories/storage/postgres"
	"flag"
	"fmt"
	"github.com/oklog/oklog/pkg/group"
	"go.uber.org/zap"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Init flags
	var (
		configPath              = flag.String("config-path", "./config", "Path to config file")
		migratePostgresDatabase = flag.Bool("migrate-postgres-database", false, "Migrate postgres database")
	)

	// Parse flags
	flag.Parse()

	// Init config
	cfg := config.GetConfig(*configPath)

	// Init logger
	log := logger.GetZapLogger(&cfg.Log)

	// Migrate database
	if *migratePostgresDatabase {
		err := migrate.PostgresDatabase(&cfg.Storage.Postgres, log)
		if err != nil {
			log.Error("Failed to migrate Postgres database", zap.Error(err))
		}

		return
	}

	// Init postgres client
	postgresClientDB, errClientDB := postgresql.Connect(&cfg.Storage.Postgres, log)
	if errClientDB != nil {
		log.Error("Failed to connect to postgres database", zap.Error(errClientDB))
	}

	// Init opts for services and repositories
	var (
		// Init opts for user service
		optsForUser []domain.OptFunc
	)

	// Add config to params
	optsForUser = append(optsForUser, domain.WithConfig(cfg))

	// Init Repositories
	var (
		// User repository
		userRepository = userRepo.NewUserRepository(postgresClientDB, log)
	)

	// Init services
	//
	// User service
	var userService domain.Service
	{
		userService = domain.NewService(userRepository, optsForUser...)
		userService = domain.NewLoggingMiddleware(log)(userService)
	}

	// Init endpoints
	//
	// User endpoints
	userEndpoints := endpoints.NewEndpoints(userService, log)

	// Init http handlers
	var httpHandlers *http.ServeMux
	{
		httpHandlers = http.NewServeMux()
		// Common handlers
		httpHandlers.Handle(httpCommonHandler.APIPath, httpCommonHandler.NewHTTPHandlers(cfg, log))
		// User handlers
		httpHandlers.Handle(httpUserHandler.APIPath, httpUserHandler.NewHTTPHandlers(userEndpoints, cfg, log))
	}

	// Init group
	var g group.Group
	{
		// Init http server
		httpListener, errHTTPListener := net.Listen("tcp", fmt.Sprintf(":%d", cfg.HTTP.Port))
		if errHTTPListener != nil {
			log.Error("Exiting due to HTTP listener error", zap.Error(errHTTPListener))
		}

		server := &http.Server{
			ReadTimeout:       cfg.HTTP.ReadTimeout * time.Second,
			ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout * time.Second,
			WriteTimeout:      cfg.HTTP.WriteTimeout * time.Second,
			Handler:           httpHandlers,
		}

		g.Add(func() error {
			log.Info("transport HTTP", zap.String("port", fmt.Sprintf(":%d", cfg.HTTP.Port)))

			return server.Serve(httpListener)
		}, func(error) {
			if err := server.Shutdown(context.Background()); err != nil {
				log.Error("transport HTTP during Shutdown", zap.Error(err))
			}
		})
	}
	{
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	log.Error("exit", zap.Error(g.Run()))
}
