package config

import (
	"film-management/pkg/auth"
	"film-management/pkg/database/postgresql"
	"film-management/pkg/logger"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

var configInstance *Config
var configOnce sync.Once

type Config struct {
	Name string
	HTTP struct {
		Port               int
		CorsAllowedOrigins []string
		NotAuthUrls        []string
		ReadTimeout        time.Duration
		ReadHeaderTimeout  time.Duration
		WriteTimeout       time.Duration
	}
	DebugHTTP struct {
		Port              int
		ReadTimeout       time.Duration
		ReadHeaderTimeout time.Duration
		WriteTimeout      time.Duration
	}
	Log     logger.Config
	Storage struct {
		Postgres postgresql.Config
	}
	Services struct {
		Auth auth.Config
	}
}

func GetConfig(configPath string) *Config {
	configOnce.Do(func() {
		v := viper.New()

		// Set config file
		setConfigFile(v, configPath)

		// Set default values
		setDefaults(v)

		if err := v.ReadInConfig(); err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}

		err := v.Unmarshal(&configInstance)
		if err != nil {
			log.Fatalf("unable to decode into struct, %v", err)
		}
	})

	return configInstance
}

func setConfigFile(v *viper.Viper, configPath string) {
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
}

func setDefaults(v *viper.Viper) {
	// Main
	v.SetDefault("name", "film_management")
	// Http
	v.SetDefault("http.port", 8080)
	v.SetDefault("http.readTimeout", 5)
	v.SetDefault("http.readHeaderTimeout", 3)
	v.SetDefault("http.writeTimeout", 10)
	v.SetDefault("http.corsAllowedOrigins", []string{"*"})
	v.SetDefault("http.notAuthUrls", []string{
		"/api/v1/health",
		"/api/v1/swagger",
		"/api/v1/user/register",
		"/api/v1/user/login",
	})
	// Debug Http
	v.SetDefault("debugHttp.port", 8081)
	v.SetDefault("debugHttp.readTimeout", 5)
	v.SetDefault("debugHttp.readHeaderTimeout", 3)
	v.SetDefault("debugHttp.writeTimeout", 10)
	// Services
	// User
	v.SetDefault("services.auth.authDurationMin", 60)
	v.SetDefault("services.auth.pathPublicKeyFile", "config/ssl/jwtRS256.key.pub")
	v.SetDefault("services.auth.pathPrivateKeyFile", "config/ssl/jwtRS256.key")
	// Storage
	v.SetDefault("storage.postgres.host", "db_film_management")
	v.SetDefault("storage.postgres.port", 5432)
	v.SetDefault("storage.postgres.user", "film")
	v.SetDefault("storage.postgres.password", "film")
	v.SetDefault("storage.postgres.database", "db")
	// Log
	v.SetDefault("log.json", false)
	v.SetDefault("log.level", "debug")
	v.SetDefault("log.colored", true)
	v.SetDefault("log.development", true)
}
