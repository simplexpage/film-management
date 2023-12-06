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
	"film-management/pkg/auth"
	"film-management/pkg/database/postgresql"
	"film-management/pkg/logger"
	"film-management/pkg/password"
	authMiddleware "film-management/pkg/transport/http/middlewares/auth"
	"film-management/pkg/transport/http/response"
	filmRepo "film-management/repositories/storage/postgres/film"
	userRepo "film-management/repositories/storage/postgres/user"
	"flag"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/oklog/oklog/pkg/group"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		seedTestData            = flag.Bool("seed-postgres-database", false, "Seed test data")
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

	// Seed test data database
	if *seedTestData {
		err := migrate.SeedTestData(&cfg.Storage.Postgres, log)
		if err != nil {
			log.Error("Failed to seed test data", zap.Error(err))
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
		userRepository = userRepo.NewUserRepository(postgresClientDB, log)
		// Film repository
		filmRepository = filmRepo.NewFilmRepository(postgresClientDB, log)
		// Password service
		passwordService = password.NewPasswordService(log)
		// Auth service
		authService = auth.NewAuthService(cfg.Services.Auth, log)
	)

	// Init services
	//
	// User service
	var userService domainUser.Service
	{
		userService = domainUser.NewService(userRepository, authService, passwordService, optsForUser...)
		userService = domainUser.NewLoggingMiddleware(log)(userService)
		// Init metrics middleware
		fieldKeys := []string{"method", "error"}
		userService = domainUser.NewInstrumentingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "domain",
				Subsystem: fmt.Sprintf("%s_%s", cfg.Name, "user"),
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
				Namespace: "domain",
				Subsystem: fmt.Sprintf("%s_%s", cfg.Name, "user"),
				Name:      "request_duration_seconds",
				Help:      "Total duration of requests in seconds.",
				Buckets: []float64{
					0.1,  // 100 ms
					0.2,  // 200 ms
					0.25, // 250 ms
					0.5,  // 500 ms
					1,    // 1 s
				},
			}, fieldKeys),
		)(userService)
	}

	// Film service
	var filmService domainFilm.Service
	{
		filmService = domainFilm.NewService(filmRepository, optsForFilm...)
		filmService = domainFilm.NewLoggingMiddleware(log)(filmService)
		// Init metrics middleware
		fieldKeys := []string{"method", "error"}
		filmService = domainFilm.NewInstrumentingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "domain",
				Subsystem: fmt.Sprintf("%s_%s", cfg.Name, "film"),
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
				Namespace: "domain",
				Subsystem: fmt.Sprintf("%s_%s", cfg.Name, "film"),
				Name:      "request_duration_seconds",
				Help:      "Total duration of requests in seconds.",
				Buckets: []float64{
					0.1,  // 100 ms
					0.2,  // 200 ms
					0.25, // 250 ms
					0.5,  // 500 ms
					1,    // 1 s
				},
			}, fieldKeys),
		)(filmService)
	}

	// Init endpoints
	var (
		// User endpoints
		userEndpoints = userEndpoint.NewEndpoints(userService, log)
		// Film endpoints
		filmEndpoints = filmEndpoint.NewEndpoints(filmService, log)
	)

	// Init HTTP routers
	var router *gin.Engine
	{
		// Instantiating router
		router = gin.Default()

		// Set error handlers
		response.SetDefaultErrorHandlers(router)

		// Init CORS middleware
		configCORS := cors.DefaultConfig()
		configCORS.AllowOrigins = cfg.HTTP.CorsAllowedOrigins
		router.Use(cors.New(configCORS))

		// Init Auth middleware
		router.Use(authMiddleware.Middleware(cfg.HTTP.NotAuthUrls, authService))

		// Init HTTP routes
		//
		// Common routes
		httpCommonHandler.SetHTTPRoutes(router)
		// User routes
		httpUserHandler.SetHTTPRoutes(router, userEndpoints, log)
		// Film routes
		httpFilmHandler.SetHTTPRoutes(router, filmEndpoints)
	}

	// Init metrics handler
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	// Init group
	var g group.Group
	{
		// Init debug listener server
		debugListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.DebugHTTP.Port))
		if err != nil {
			log.Error("debugListener", zap.Error(err))
		}

		server := &http.Server{
			ReadTimeout:       cfg.DebugHTTP.ReadTimeout * time.Second,
			ReadHeaderTimeout: cfg.DebugHTTP.ReadHeaderTimeout * time.Second,
			WriteTimeout:      cfg.DebugHTTP.WriteTimeout * time.Second,
			Handler:           http.DefaultServeMux,
		}
		g.Add(func() error {
			log.Info("transport debug/HTTP", zap.String("port", fmt.Sprintf(":%d", cfg.DebugHTTP.Port)))

			return server.Serve(debugListener)
		}, func(error) {
			if err := server.Shutdown(context.Background()); err != nil {
				log.Error("transport debug/HTTP during Shutdown", zap.Error(err))
			}
		})
	}
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
			Handler:           router,
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

func httpHandlerForRegister() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Register successful"})
	}
}
