package config

import "os"

type AppConfig struct {
	Port string
}

var App AppConfig

func loadAppConfig() {
	App = AppConfig{
		Port: os.Getenv("APP_PORT"),
	}
}
