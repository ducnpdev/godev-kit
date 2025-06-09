package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type (
	// Config -.
	Config struct {
		App     App     `mapstructure:"APP"`
		HTTP    HTTP    `mapstructure:"HTTP"`
		Log     Log     `mapstructure:"LOG"`
		PG      PG      `mapstructure:"PG"`
		GRPC    GRPC    `mapstructure:"GRPC"`
		RMQ     RMQ     `mapstructure:"RMQ"`
		Metrics Metrics `mapstructure:"METRICS"`
		Swagger Swagger `mapstructure:"SWAGGER"`
	}

	// App -.
	App struct {
		Name    string `mapstructure:"NAME"`
		Version string `mapstructure:"VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port           string `mapstructure:"PORT"`
		UsePreforkMode bool   `mapstructure:"USE_PREFORK_MODE"`
	}

	// Log -.
	Log struct {
		Level string `mapstructure:"LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax int    `mapstructure:"POOL_MAX"`
		URL     string `mapstructure:"URL"`
	}

	// GRPC -.
	GRPC struct {
		Port string `mapstructure:"PORT"`
	}

	// RMQ -.
	RMQ struct {
		ServerExchange string `mapstructure:"RPC_SERVER"`
		ClientExchange string `mapstructure:"RPC_CLIENT"`
		URL            string `mapstructure:"URL"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `mapstructure:"ENABLED"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `mapstructure:"ENABLED"`
	}
)

func (c *Config) binding(v *viper.Viper) error {
	if err := v.Unmarshal(&c); err != nil {
		log.Println("failed to unmarshal config: ", err)
		return err
	}
	return nil
}

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	return loadConfigYAML()
}

func loadConfigYAML() (*Config, error) {
	conf := &Config{}
	vn := viper.New()

	vn.AddConfigPath("config")
	vn.SetConfigName("config")
	vn.SetConfigType("yaml")
	vn.AutomaticEnv()

	vn.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := vn.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	for _, key := range vn.AllKeys() {
		str := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))

		vn.BindEnv(key, str)
	}

	conf.binding(vn)
	vn.WatchConfig()

	vn.OnConfigChange(func(e fsnotify.Event) {
		for _, key := range vn.AllKeys() {
			str := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
			log.Default().Println(key, str, vn.Get(key))
			vn.BindEnv(key, str)
		}

		if err := conf.binding(vn); err != nil {
			log.Println("binding error:", err)
		}
		log.Printf("config: %+v", conf)
	})

	return conf, nil
}
