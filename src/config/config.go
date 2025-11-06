package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	defaultLocalConfigPath = "config/local.yaml"
)

type Config struct {
	Env         string         `yaml:"env" env-default:"local"`
	StoragePath string         `yaml:"storage_path"` //пока убрал env-required:"true"`
	HttpConfig  HTTPConfig `yaml:"http_server"`
}

type HTTPConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		// В идеале передавать флаг, но пока пусть будет так
		// panic("config path is empty")
		configPath = defaultLocalConfigPath
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config path does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func MustLoadByPath(configPath string) *Config {
	if configPath == "" {
		// В идеале передавать флаг, но пока пусть будет так
		// panic("config path is empty")
		configPath = defaultLocalConfigPath
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config path does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
