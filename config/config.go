package config

import (
	"encoding/json"
	"github.com/miiy/save-page/file"
)

type Config struct {
	Debug bool
	Proxy string
	StoragePath string `json:"storage-path"`
}

func NewConfig(name string) (*Config, error) {
	fileByte, err := file.ReadAll(name)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := json.Unmarshal(fileByte, config); err != nil {
		return nil, err
	}
	return config, nil
}