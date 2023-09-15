package main

import (
	"GoBackend/config"
	"GoBackend/server"
	"log"
	"os"
)

func main() {
	log.Println("Starting Runners App")

	log.Println("Initializing configuration")
	config := config.InitConfig(getConfigFileName())

	log.Println("Initializing database")
	dbHandler := server.InitDatabase(config)

	log.Println("Initializing HTTP sever")
	httpServer := server.InitHttpServer(config, dbHandler)
	httpServer.Start()
}

func getConfigFileName() string {
	env := os.Getenv("ENV")

	if env != "" {
		return "config-" + env
	}

	return "config"
}
