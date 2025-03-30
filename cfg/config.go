package config

type Config struct {
	Server struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		StoragePath string `yaml:"storage_path"`
	} `yaml:"server"`

	Limits struct {
		UploadDownload int `yaml:"upload_download"`
		ListFiles      int `yaml:"list_files"`
	} `yaml:"limits"`
}

func NewDefaultConfig() *Config {
	cfg := &Config{}
	cfg.Server.Host = "localhost"
	cfg.Server.Port = 1337
	cfg.Server.StoragePath = "./storage"
	cfg.Limits.UploadDownload = 10
	cfg.Limits.ListFiles = 100
	return cfg
}
