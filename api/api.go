package api

import (
	"github.com/devon-caron/metrifuge/api/internal/handlers"
	"github.com/devon-caron/metrifuge/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/go-chi/cors"
	"github.com/sirupsen/logrus"
	"net/http"
)

var log *logrus.Logger

func StartApi() {
	log = logger.Get()
	var router = chi.NewRouter()

	// CORS middleware with debug settings
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for debugging
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers for debugging
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	handlers.RouterHandler(router)

	log.Info("API initialized, starting server...")
	go func() {
		err := http.ListenAndServe("0.0.0.0:8080", router)
		if err != nil {
			log.Error(err)
		}
	}()
}
