package config

import (
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/internal/env"
)

type Configuration struct {
	Addr         string
	Env          string
	DB           dbConfig
	TwilioConfig twilioConfig
	AuthConfig   authConfig
}

type dbConfig struct {
	Addr                       string
	MaxOpenConns, MaxIdleConns int
	MaxIdleTime                string
}

type twilioConfig struct {
	AccountSID, AuthToken, ServiceID string
}

type authConfig struct {
	Secret string
	Exp    time.Duration
	Iss    string
}

var Config = initConfig()

func initConfig() Configuration {
	return Configuration{
		Addr: env.GetString("ADDR", ":8080"),
		Env:  env.GetString("ENV", "development"),
		DB: dbConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/mygallery?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		TwilioConfig: twilioConfig{
			AccountSID: env.GetString("TWILIO_ACCOUNT_SID", ""),
			AuthToken:  env.GetString("TWILIO_AUTH_TOKEN", ""),
			ServiceID:  env.GetString("TWILIO_SID", ""),
		},
		AuthConfig: authConfig{
			Secret: env.GetString("AUTH_TOKEN_SECRET", "test"),
			Exp:    time.Hour * 24 * 3,
			Iss:    "MiGalaria",
		},
	}
}
