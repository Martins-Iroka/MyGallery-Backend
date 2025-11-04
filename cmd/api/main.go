package main

import (
	"github.com/Martins-Iroka/MyGallery-Backend/internal/env"
	"go.uber.org/zap"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	app := &application{
		config: cfg,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
