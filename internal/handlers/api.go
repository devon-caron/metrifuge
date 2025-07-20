package handlers

import (
	"github.com/devon-caron/goapi/internal/middleware"
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func Handler(r *chi.Mux) {
	// Global middleware
	r.Use(chimiddle.StripSlashes)

	r.Route("/account", func(router chi.Router) {

		// Middleware for /account route
		router.Use(middleware.Authorization)

		router.Get("/coins", GetCoinBalance)
	})

	r.Route("/api", func(router chi.Router) {
		router.Use(middleware.OAuth)

		router.Post("/token", TokenHandler)
	})
}
