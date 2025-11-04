package config

import "github.com/Martins-Iroka/MyGallery-Backend/internal/env"

type Configuration struct {
	Addr string
	DB   dbConfig
}

type dbConfig struct {
	Addr                       string
	MaxOpenConns, MaxIdleConns int
	MaxIdleTime                string
}

var Config = initConfig()

func initConfig() Configuration {
	return Configuration{
		Addr: env.GetString("ADDR", ":8080"),
		DB: dbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/mygallery-backend?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}
