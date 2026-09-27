package configs

type Config struct {
	Environment string
	Port        int
}

func Load() (*Config, error) {
	return &Config{Environment: "development", Port: 8080}, nil
}
