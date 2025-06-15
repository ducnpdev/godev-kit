type Config struct {
	HTTP      HTTP      `mapstructure:"http"`
	DB        DB        `mapstructure:"db"`
	Redis     Redis     `mapstructure:"redis"`
	Logger    Logger    `mapstructure:"logger"`
	Jaeger    Jaeger    `mapstructure:"jaeger"`
	AppConfig AppConfig `mapstructure:"app"`
}

type HTTP struct {
	Port            string        `mapstructure:"port"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

func (c *Config) Validate() error {
	if c.HTTP.Port == "" {
		return errors.New("http port is required")
	}
	if c.HTTP.ShutdownTimeout == 0 {
		return errors.New("http shutdown timeout is required")
	}
	return nil
} 