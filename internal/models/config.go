package models

type Config struct {
	AppConfig      `yaml:"app"`
	GetBlockConfig `yaml:"getblock"`
	StorageConfig  `yaml:"storage"`
}

type AppConfig struct {
	Debug bool `yaml:"debug"`
	Cycle int  `yaml:"cycle"`
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
