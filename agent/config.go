package main

import (
	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	DB database `toml:"database"`
}

type database struct {
	Host     string
	Port     string
	UserName string
	PassWord string
}

func GetDBConfig() (host string, port string, userName string, password string) {
	var config tomlConfig
	if _, err := toml.DecodeFile("config/configure.toml", &config); err != nil {
		host, port, userName, password = config.DB.Host, config.DB.Port, config.DB.UserName, config.DB.PassWord
		return
	}
	return
}
