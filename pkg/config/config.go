package config

import (
	"encoding/json"
	"github.com/miiy/save-page/pkg/file"
)

type DialContext struct {
	Timeout int
	KeepAlive int
}

type Config struct {
	Debug       bool
	Timeout     int
	Proxy       string
	StoragePath string      `json:"storage-path"`
	DialContext DialContext `json:"dial_context"`
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