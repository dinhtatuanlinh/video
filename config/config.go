package config

import (
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string        `mapstructure:"ENV"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	MaxRetries           int           `mapstructure:"MAX_RETRIES"`
	MigrationUrl         string        `mapstructure:"MIGRATION_URL"`
	RedisAddress         string        `mapstructure:"REDIS_ADDRESS_"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	EmailSenderName      string        `mapstructure:"EMAIL_SENDER_NAME"`
	EmailSenderAddress   string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailSenderPassword  string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	VideoPath            string        `mapstructure:"VIDEO_PATH"`
}

var envKeys = []string{
	"ENV",
	"DB_SOURCE",
}

func LoadConfig(path string) (config Config, err error) {
	env := os.Getenv("ENV")
	if env == "dev" {
		viper.AddConfigPath(path)
		viper.SetConfigName("dev")
		viper.SetConfigType("env")
	} else {
		viper.AddConfigPath(path)
		viper.SetConfigName("local")
		viper.SetConfigType("env")
	}

	viper.AutomaticEnv()
	for _, key := range envKeys {
		_ = viper.BindEnv(key)
	}

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	return
}
