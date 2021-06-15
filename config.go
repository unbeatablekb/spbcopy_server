package main

import (
	"github.com/go-ini/ini"
)

type Config struct {
	Port int
	CPFile string
}

func initConfig(configFile string) (*Config, string) {
	cfg, err := ini.Load(configFile)
	if err != nil {
		return nil, "Failed to parse config file: " + err.Error()
	}

	config := new(Config)
	config.Port, err = cfg.Section("base").Key("port").Int()
	if err != nil {
		return nil, "Failed to read port from config file: " + err.Error()
	}

	config.CPFile = cfg.Section("base").Key("copyFile").String()

	return config, ""
}