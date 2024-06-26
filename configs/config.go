package configs

import (
	"github.com/spf13/viper"
)

var cfg *config

type config struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	LocalStoragePath   string `mapstructure:"LOCAL_STORAGE_PATH"`
	ConcurrencyUpload  string `mapstructure:"CONCURRENCY_UPLOAD"`
	ConcurrencyWorkers string `mapstructure:"CONCURRENCY_WORKERS"`
	BucketName         string `mapstructure:"BUCKET_NAME"`
}

func LoadConfig(path string) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
}

func Config() *config {
	return cfg
}
