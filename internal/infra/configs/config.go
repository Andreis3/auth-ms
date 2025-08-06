package configs

import (
	"time"

	"github.com/spf13/viper"
)

// Conf holds the application configuration loaded from environment variables.
type Configs struct {
	ServerPort              string        `mapstructure:"SERVER_PORT"`                 // HTTP server port
	PostgresHost            string        `mapstructure:"POSTGRES_HOST"`               // PostgreSQL database host
	PostgresPort            string        `mapstructure:"POSTGRES_PORT"`               // PostgreSQL database port
	PostgresUser            string        `mapstructure:"POSTGRES_USER"`               // PostgreSQL database user
	PostgresPassword        string        `mapstructure:"POSTGRES_PASSWORD"`           // PostgreSQL database password
	PostgresDBName          string        `mapstructure:"POSTGRES_DB_NAME"`            // PostgreSQL database name
	PostgresMaxConnections  int32         `mapstructure:"POSTGRES_MAX_CONNECTIONS"`    // Maximum number of database connections
	PostgresMinConnections  int32         `mapstructure:"POSTGRES_MIN_CONNECTIONS"`    // Minimum number of database connections
	PostgresMaxConnLifetime time.Duration `mapstructure:"POSTGRES_MAX_CONN_LIFETIME"`  // Maximum lifetime of a database connection
	PostgresMaxConnIdleTime time.Duration `mapstructure:"POSTGRES_MAX_CONN_IDLE_TIME"` // Maximum idle time for a database connection
	RedisHost               string        `mapstructure:"REDIS_HOST"`                  // Redis host
	RedisPort               string        `mapstructure:"REDIS_PORT"`                  // Redis port
	RedisPassword           string        `mapstructure:"REDIS_PASSWORD"`              // Redis password
	RedisDB                 int           `mapstructure:"REDIS_DB"`                    // Redis database number
	ApplicationName         string        `mapstructure:"APPLICATION_NAME"`            // name of application
	JWTSecret               string        `mapstructure:"JWT_SECRET"`                  // JWT secret
	JWTExpiry               time.Duration `mapstructure:"JWT_EXPIRY"`                  // JWT expiry
	Env                     string        `mapstructure:"ENV"`                         // Environment
}

// LoadConfig loads the application configuration from either a .env file or environment variables.
// It prioritizes the .env file if it exists, otherwise falls back to environment variables.
func LoadConfig() *Configs {
	viper.SetConfigFile(".env") // name of config file (without extension)
	viper.AutomaticEnv()        // Fallback to environment variables if .env is not found

	// Set default values for optional configuration
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("POSTGRES_MAX_CONNECTIONS", 10)
	viper.SetDefault("POSTGRES_MIN_CONNECTIONS", 1)
	viper.SetDefault("POSTGRES_MAX_CONN_LIFETIME", "5m")
	viper.SetDefault("POSTGRES_MAX_CONN_IDLE_TIME", "1m")
	viper.SetDefault("REDIS_DB", 0)
	viper.SetDefault("ENV", "production")

	if err := viper.ReadInConfig(); err != nil {
		// If the .env file is not found, ignore the error and rely on environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil
		}
	}

	var cfg Configs

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil
	}

	return &cfg
}
