package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Martins-Iroka/MyGallery-Backend/config"
	"github.com/Martins-Iroka/MyGallery-Backend/docs"
	"github.com/Martins-Iroka/MyGallery-Backend/internal"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/auth"
	"github.com/Martins-Iroka/MyGallery-Backend/internal/env"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	config          config.Configuration
	logger          *zap.SugaredLogger
	otpVerification auth.OTPVerification
	store           internal.Storage
	auth            auth.Authenticator
}

func (app *application) Mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "http://127.0.0.1:4040")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/register", app.registerUserHandler)
			r.Post("/verify", app.verifyUserHandler)
			r.Post("/login", app.loginUserHandler)
			r.Post("/refresh", app.refreshTokenHandler)
			r.Post("/logout", app.logoutHandler)
		})

		r.Route("/photos", func(r chi.Router) {
			r.Use(app.authTokenMiddleware)
			r.Get("/", app.getPhotosHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.photoPostExistMiddleware)
				r.Post("/createcomment", app.createCommentForPostHandler)
				r.Get("/comments", app.getCommentsByPostID)
			})
		})

		r.Route("/videos", func(r chi.Router) {
			r.Use(app.authTokenMiddleware)
			r.Get("/", app.getVideosHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.videoPostExistMiddleware)
				r.Post("/createcomment", app.createVideoCommentHandler)
				r.Get("/comments", app.getVideoCommentByPostID)
			})
		})
	})
	return r
}

func (app *application) Run(mux http.Handler) error {

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.ApiUrl
	docs.SwaggerInfo.BasePath = "/v1"

	srv := &http.Server{
		Addr:         app.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at %s", app.config.Addr)

	return srv.ListenAndServe()
}
