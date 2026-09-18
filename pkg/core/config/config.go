package config

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	MongoURI           string `mapstructure:"MONGO_URI"`
	DBName             string `mapstructure:"DB_NAME"`
	Port               string `mapstructure:"PORT"`
	AccessTokenSecret  string `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string `mapstructure:"REFRESH_TOKEN_SECRET"`
	BrevoAPIKey        string `mapstructure:"BREVO_API_KEY"`
	LogLevel           string `mapstructure:"LOG_LEVEL"`
	LogFormat          string `mapstructure:"LOG_FORMAT"`
	GinMode            string `mapstructure:"GIN_MODE"`
}

// globalConfig holds the latest configured values
var globalConfig *Config

// LoadConfig initializes viper, reads the .env file, and sets up watch for changes.
func LoadConfig() *Config {
	viper.SetConfigFile(".env")
	viper.SetDefault("LOG_LEVEL", "info")
	viper.SetDefault("LOG_FORMAT", "json")
	viper.SetDefault("GIN_MODE", "release")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Error reading config file, relying on environment variables: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	globalConfig = &config

	// Watch for changes in the config file
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
		var updatedConfig Config
		if err := viper.Unmarshal(&updatedConfig); err != nil {
			log.Printf("Failed to unmarshal updated config: %v", err)
		} else {
			// Atomic swap or simply updating the pointer (note: in highly concurrent apps, consider atomic.Value or mutex)
			globalConfig = &updatedConfig
			log.Println("Configuration successfully hot-reloaded.")
		}
	})

	return globalConfig
}

// GetConfig returns the latest loaded configuration properties wrapper
func GetConfig() *Config {
	if globalConfig == nil {
		return LoadConfig()
	}
	return globalConfig
}
