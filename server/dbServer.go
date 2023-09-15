package server

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func InitDatabase(config *viper.Viper) *gorm.DB {

	dsn := config.GetString("database.connection_string")
	maxIdleConnections := config.GetInt("database.max_idle_connections")
	maxOpenConnections := config.GetInt("database.max_open_connections")
	connectionMaxLifetime := config.GetDuration("database.connection_max_lifetime")

	if dsn == "" {
		log.Fatalf("Database source name is missing")
	}

	// automatically ping database after initialized to check database availability
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Error while initializing database: %v", err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxIdleConns(maxIdleConnections)

	sqlDB.SetMaxOpenConns(maxOpenConnections)
	sqlDB.SetConnMaxLifetime(connectionMaxLifetime)
	return db

}
