package config

type Config struct {
	Server struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		StoragePath string `yaml:"storage_path"`
	} `yaml:"server"`
}

func NewDefaultConfig() *Config {
	cfg := &Config{}
	cfg.Server.Host = "localhost"
	cfg.Server.Port = 1488
	cfg.Server.StoragePath = "./storage"
	return cfg
}
