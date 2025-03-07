package config

import (
	"os"
	"regexp"

	"github.com/Comvoca-AI/comvoca-admin-back/internal/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Auth        AuthConfig
	Application Application
}

type ServerConfig struct {
	Port          int  `mapstructure:"port"`
	MemoryProfile bool `mapstructure:"memory_profile"`
	CPUProfile    bool `mapstructure:"cpu_profile"`
}

type AuthConfig struct {
	Cognito CognitoConfig `mapstructure:"cognito"`
}

type CognitoConfig struct {
	UserPoolId      string `mapstructure:"user_pool_id"`
	ClientId        string `mapstructure:"client_id"`
	Region          string `mapstructure:"region"`
	Secret          string `mapstructure:"secret"`
	RedirectURL     string `mapstructure:"redirect_uri"`
	IssuerURL       string `mapstructure:"issuer_url"`
	JwtSecret       string `mapstructure:"jwt_secret"`
	TokenExpireHour int    `mapstructure:"token_expire_hour"`
}
type DatabaseConfig struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
	DBName        string `mapstructure:"dbname"`
	RunMigrations bool   `mapstructure:"run_migrations"`
}

type Application struct {
	BaseURL         string `mapstructure:"base_url"`
	Debug           bool   `mapstructure:"debug"`
	DocumentsBucket string `mapstructure:"documents_bucket"`
	BodyLimitSize   int    `mapstructure:"body_limit_size"`
	Profile         string `mapstructure:"profile"`
}

var AppConfig Config

func InitConfig() {
	env := os.Getenv("COMVOCA_ENV")
	if env == "" {
		env = "local"
	}
	viper.SetConfigName(env)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("Error reading config file: %v", err)
		return
	}

	resolveEnvVariablesWithDefaults()

	if err := viper.Unmarshal(&AppConfig); err != nil {
		logger.Error("Unable to decode into struct: %v", err)
		return
	}
}

func resolveEnvVariablesWithDefaults() {
	placeholderPattern := regexp.MustCompile(`\${([^:}]+):?([^}]*)}`)

	for _, key := range viper.AllKeys() {
		value := viper.GetString(key)

		// Check if the value matches the placeholder pattern
		if matches := placeholderPattern.FindStringSubmatch(value); len(matches) > 0 {
			envVar := matches[1]       // The environment variable name
			defaultValue := matches[2] // The default value

			// Look up the environment variable
			envValue, exists := os.LookupEnv(envVar)
			if exists {
				viper.Set(key, envValue) // Use environment variable value
			} else {
				viper.Set(key, defaultValue) // Use default value
			}
		}
	}
}
