package main

import (
	"GoBackend/config"
	"os"
)

func main() {
	config := config.InitConfig(getConfigFileName())
}

func getConfigFileName() string {
	env := os.Getenv("ENV")

	if env != "" {
		return "config-" + env
	}

	return "config"
}
