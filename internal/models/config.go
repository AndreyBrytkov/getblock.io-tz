package models

type Config struct {
	AppConfig      `yaml:"app"`
	ServerConfig   `yaml:"server"`
	GetBlockConfig `yaml:"getblock"`
	StorageConfig  `yaml:"storage"`
}

type AppConfig struct {
	Debug       bool `yaml:"debug"`
	BlockAmount int  `yaml:"blockAmount"`
	Cycle       int  `yaml:"cycle"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

type GetBlockConfig struct {
	Key string `yaml:"api-key"`
}

type StorageConfig struct {
	Type     string `yaml:"type"`
	Address  string `yaml:"address"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
}
