package main

import (
	"context"
	"film-management/cmd/server/commands/migrate"
	"film-management/config"
	httpCommonHandler "film-management/internal/common/transport/http"
	domainFilm "film-management/internal/film/domain"
	filmEndpoint "film-management/internal/film/endpoints"
	httpFilmHandler "film-management/internal/film/transport/http"
	domainUser "film-management/internal/user/domain"
	userEndpoint "film-management/internal/user/endpoints"
	httpUserHandler "film-management/internal/user/transport/http"
	"film-management/pkg/database/postgresql"
	"film-management/pkg/logger"
	"film-management/pkg/transport/http/response"
	"film-management/repositories/services"
	filmRepo "film-management/repositories/storage/postgres"
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
		optsForUser []domainUser.OptFunc
		// Init opts for film service
		optsForFilm []domainFilm.OptFunc
	)

	// Init Repositories
	var (
		// User repository
		userRepository = filmRepo.NewUserRepository(postgresClientDB, log)
		// Film repository
		filmRepository = filmRepo.NewFilmRepository(postgresClientDB, log)
		// Password service
		passwordService = services.NewPasswordService(log)
		// Auth service
		authService = services.NewAuthService(cfg.Services.Auth, log)
	)

	// Init services
	//
	// User service
	var userService domainUser.Service
	{
		userService = domainUser.NewService(userRepository, authService, passwordService, optsForUser...)
		userService = domainUser.NewLoggingMiddleware(log)(userService)
	}

	// Film service
	var filmService domainFilm.Service
	{
		filmService = domainFilm.NewService(filmRepository, optsForFilm...)
		filmService = domainFilm.NewLoggingMiddleware(log)(filmService)
	}

	// Init endpoints
	var (
		// User endpoints
		userEndpoints = userEndpoint.NewEndpoints(userService, log)
		// Film endpoints
		filmEndpoints = filmEndpoint.NewEndpoints(filmService, log)
	)

	// Init http handlers
	var httpHandlers *http.ServeMux
	{
		httpHandlers = http.NewServeMux()
		// Common handlers
		httpHandlers.Handle(httpCommonHandler.APIPath, httpCommonHandler.NewHTTPHandlers(cfg, log))
		// User handlers
		httpHandlers.Handle(httpUserHandler.APIPath, httpUserHandler.NewHTTPHandlers(userEndpoints, cfg, log))
		// Film handlers
		httpHandlers.Handle(httpFilmHandler.APIPath, httpFilmHandler.NewHTTPHandlers(filmEndpoints, authService, cfg, log))
		// Base 404 handler
		httpHandlers.HandleFunc("/", response.NotFoundFunc)
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
