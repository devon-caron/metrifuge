package api

import (
	"github.com/go-chi/cors"
	"net/http"

	"github.com/devon-caron/metrifuge/api/internal/handlers"
	"github.com/go-chi/chi"
	_ "github.com/go-chi/cors"
	log "github.com/sirupsen/logrus"
)

func StartApi() {
	log.SetReportCaller(true)
	var r *chi.Mux = chi.NewRouter()

	// CORS middleware with debug settings
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins for debugging
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // Allow all headers for debugging
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	handlers.RouterHandler(r)

	log.Info("API initialized, starting server...")
	err := http.ListenAndServe("0.0.0.0:8000", r)
	if err != nil {
		log.Error(err)
	}
}
