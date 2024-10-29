package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var config = &Config{}

func init() {
	environment := os.Getenv("MYENV")
	if environment == "" {
		environment = "production"
	}

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("runtime.Caller - it was not possible to recover the information")
	}
	// adjust it if you change the project layout
	baseDir := filepath.Dir(filepath.Dir(filepath.Dir(file)))

	kc, err := setupConfig(filepath.Join(baseDir, fmt.Sprintf("%s.yaml", environment)))
	if err != nil {
		panic(err)
	}

	if err := kc.Unmarshal(config); err != nil {
		panic(err)
	}
}

func setupConfig(in string) (*viper.Viper, error) {
	kc := viper.New()

	kc.SetConfigFile(in)
	if err := kc.ReadInConfig(); err != nil {
		return nil, err
	}

	kc.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
	kc.WatchConfig()

	return kc, nil
}

func GetConfig() *Config {
	return config
}

type Config struct {
	MySQL `mapstructure:"mysql"`
	HTTP
	Logger
	JWT
}

type MySQL struct {
	User     string
	Password string
	Host     string
	Port     int
	Database string
	Debug    bool
}

type HTTP struct {
	Host string
	Port int
}

type Logger struct {
	AddSource bool
	Level     int
}

type JWT struct {
	Expiration time.Duration
}
