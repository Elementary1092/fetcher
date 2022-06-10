package config

import (
	"github.com/elem1092/fetcher/pkg/logging"
	"github.com/ilyakaznacheev/cleanenv"
	"sync"
)

type Configuration struct {
	ServerConfig ServerConfig `yaml:"server"`
	DBConfig     DBConfig     `yaml:"db_config"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    string `yaml:"port"`
	BaseURL string `yaml:"fetch_from"`
	Pages   string `yaml:"pages"`
}

type DBConfig struct {
	Address     string `yaml:"address"`
	Port        string `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	DBName      string `yaml:"db_name"`
	MaxAttempts int    `yaml:"max_attempts"`
}

var cfg = &Configuration{ServerConfig{}, DBConfig{}}
var once sync.Once

func GetConfiguration() *Configuration {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("Parsing configuration file")
		if err := cleanenv.ReadConfig("config.yml", cfg); err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			logger.Fatal(help)
			panic(err)
		}
		logger.Info("finished reading configuration file")
	})

	return cfg
}
