package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
)

type Database struct {
	Host             string `mapstructure:"addr"`
	Port             int    `mapstructure:"port"`
	User             string `mapstructure:"user"`
	Password         string `mapstructure:"password"`
	DbName           string `mapstructure:"dbname"`
	PoolSize         int    `mapstructure:"pool_size"`
	TimeOut          int    `mapstructure:"timeout"`
	PgMigrationsPath string `mapstructure:"PG_MIGRATIONS_PATH"`
}

// Config is a config :).
type Config struct {
	Database *Database `mapstructure:"database"`
}

var (
	once   sync.Once
	config *Config
)

func Get() *Config {
	once.Do(func() {

		viper.SetConfigFile("./config.yaml")
		viper.SetDefault("database.addr", "localhost:5432")
		viper.SetDefault("database.pool_size", 10)
		fmt.Println(viper.ConfigFileUsed())

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		config = &Config{}

		if err := viper.Unmarshal(config); err != nil {
			log.Fatalf("Unable to unmarshall config: %v", err)
		}

		if config.Database == nil {
			log.Fatalf("Database props are missing")
		}

	})

	return config
}
