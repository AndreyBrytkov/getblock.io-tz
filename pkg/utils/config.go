package utils

import (
	"io/ioutil"
	"path/filepath"

	"github.com/AndreyBrytkov/getblock.io-tz/internal/models"
	yaml "gopkg.in/yaml.v2"
)

func LoadConfig() (*models.Config, error) {
	buf, err := ioutil.ReadFile(filepath.Join(".", "config", "config.yaml"))
	if err != nil {
		return nil, WrapErr("app", "Read config file error", err)
	}
	config := &models.Config{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		return nil, WrapErr("app", "Unmarshal config error", err)
	}
	return config, nil
}
