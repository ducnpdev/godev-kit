package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type (
	// Config -.
	Config struct {
		App     App     `mapstructure:"APP"`
		HTTP    HTTP    `mapstructure:"HTTP"`
		Log     Log     `mapstructure:"LOG"`
		PG      PG      `mapstructure:"PG"`
		Redis   Redis   `mapstructure:"REDIS"`
		GRPC    GRPC    `mapstructure:"GRPC"`
		RMQ     RMQ     `mapstructure:"RMQ"`
		Kafka   Kafka   `mapstructure:"KAFKA"`
		NATS    NATS    `mapstructure:"NATS"`
		Metrics Metrics `mapstructure:"METRICS"`
		Swagger Swagger `mapstructure:"SWAGGER"`
		JWT     JWT     `mapstructure:"JWT"`
	}

	// App -.
	App struct {
		Name    string `mapstructure:"NAME"`
		MODE    string `mapstructure:"MODE"`
		Version string `mapstructure:"VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port            string        `mapstructure:"PORT"`
		ReadTimeout     time.Duration `mapstructure:"READ_TIMEOUT"`
		WriteTimeout    time.Duration `mapstructure:"WRITE_TIMEOUT"`
		IdleTimeout     time.Duration `mapstructure:"IDLE_TIMEOUT"`
		ShutdownTimeout time.Duration `mapstructure:"SHUTDOWN_TIMEOUT"`
		UsePreforkMode  bool          `mapstructure:"USE_PREFORK_MODE"`
		ApiTimeout      time.Duration `mapstructure:"API_TIMEOUT"`
	}

	// Log -.
	Log struct {
		Level             string       `mapstructure:"LEVEL"`
		InCommingRequest  InOutRequest `mapstructure:"IN_COMMING_REQUEST"`
		OutCommingRequest InOutRequest `mapstructure:"OUT_COMMING_REQUEST"`
	}
	// IN/Out Request
	InOutRequest struct {
		PrintRequest  bool `mapstructure:"PRINT_REQUEST"`
		PrintResponse bool `mapstructure:"PRINT_RESPONSE"`
	}

	// PG -.
	PG struct {
		PoolMax           int           `mapstructure:"POOL_MAX"`
		PoolMin           int           `mapstructure:"POOL_MIN"`
		MaxConnLifetime   time.Duration `mapstructure:"MAX_CONN_LIFETIME"`
		MaxConnIdleTime   time.Duration `mapstructure:"MAX_CONN_IDLE_TIME"`
		HealthCheckPeriod time.Duration `mapstructure:"HEALTH_CHECK_PERIOD"`
		URL               string        `mapstructure:"URL"`
	}

	// Redis -.
	Redis struct {
		URL string `mapstructure:"URL"`
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

	// Kafka -.
	Kafka struct {
		Brokers []string `mapstructure:"BROKERS"`
		GroupID string   `mapstructure:"GROUP_ID"`
		Topics  Topics   `mapstructure:"TOPICS"`
		Control Control  `mapstructure:"CONTROL"`
	}

	// Control -.
	Control struct {
		ProducerEnabled bool `mapstructure:"PRODUCER_ENABLED"`
		ConsumerEnabled bool `mapstructure:"CONSUMER_ENABLED"`
	}

	// Topics -.
	Topics struct {
		UserEvents        string `mapstructure:"USER_EVENTS"`
		TranslationEvents string `mapstructure:"TRANSLATION_EVENTS"`
	}

	// Metrics -.
	Metrics struct {
		Enabled bool `mapstructure:"ENABLED"`
		// SetSkipPaths type array string,
		// Declare string, handle split sep ";"
		SetSkipPaths string `mapstructure:"SKIP_PATHS"`
		// Path metrics, default /metrics
		Path string `mapstructure:"PATH"`
	}

	// Swagger -.
	Swagger struct {
		Enabled bool `mapstructure:"ENABLED"`
	}
	// JWTConfig -.
	JWT struct {
		Secret string `mapstructure:"SECRET"`
	}

	// NATS -.
	NATS struct {
		URL     string        `mapstructure:"URL"`
		Timeout time.Duration `mapstructure:"TIMEOUT"`
		Enable  bool          `mapstructure:"ENABLE"`
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

func (c *Config) Validate() error {
	if c.HTTP.Port == "" {
		return errors.New("http port is required")
	}
	if c.HTTP.ShutdownTimeout == 0 {
		return errors.New("http shutdown timeout is required")
	}
	
	// Validate Kafka configuration
	if err := c.validateKafkaConfig(); err != nil {
		return fmt.Errorf("kafka config validation failed: %w", err)
	}
	
	return nil
}

// validateKafkaConfig validates Kafka-specific configuration
func (c *Config) validateKafkaConfig() error {
	if len(c.Kafka.Brokers) == 0 {
		return errors.New("kafka brokers are required")
	}
	
	if c.Kafka.GroupID == "" {
		return errors.New("kafka group ID is required")
	}
	
	// Log Kafka control settings
	log.Printf("Kafka Control Settings:")
	log.Printf("  Producer Enabled: %v", c.Kafka.Control.ProducerEnabled)
	log.Printf("  Consumer Enabled: %v", c.Kafka.Control.ConsumerEnabled)
	
	return nil
}
