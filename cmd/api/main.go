package main

import (
	"github.com/Martins-Iroka/MyGallery-Backend/config"
	"github.com/Martins-Iroka/MyGallery-Backend/internal"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/auth"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/db"
	"go.uber.org/zap"
)

const version = "1.0.0"

//	@title			MiGalaria API
//	@description	API for MiGalaria, an application for pictures and short videos.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host						localhost:8080
// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
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

	twilio := auth.NewTwilioVerification(
		config.Config.TwilioConfig.AccountSID,
		config.Config.TwilioConfig.AuthToken,
		config.Config.TwilioConfig.ServiceID,
	)

	store := internal.NewStorage(db)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		config.Config.AuthConfig.Secret,
		config.Config.AuthConfig.Iss,
		config.Config.AuthConfig.Iss,
	)

	app := &application{
		config:          config.Config,
		logger:          logger,
		otpVerification: twilio,
		store:           store,
		auth:            jwtAuthenticator,
	}

	mux := app.Mount()

	logger.Fatal(app.Run(mux))
}
