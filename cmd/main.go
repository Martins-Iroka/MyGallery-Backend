package main

import (
	"github.com/Martins-Iroka/MyGallery-Backend/cmd/api"
	"github.com/Martins-Iroka/MyGallery-Backend/config"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/db"
	"go.uber.org/zap"
)

func main() {
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.NewPostgreInstance(
		config.Config.DB.Addr,
		config.Config.DB.MaxOpenConns,
		config.Config.DB.MaxIdleConns,
		config.Config.DB.MaxIdleTime,
	)

	if err != nil {
		logger.Fatalf("db error - %s", err)
	}
	defer db.Close()
	logger.Info("data connection pool established")

	app := &api.Application{
		Config: config.Config,
		Logger: logger,
	}

	mux := app.Mount()

	logger.Fatal(app.Run(mux))
}
