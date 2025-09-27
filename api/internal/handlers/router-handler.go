package handlers

import (
	"github.com/go-chi/chi"
)

func RouterHandler(router *chi.Mux) {
	// Global Log Receiver
	// lr := receiver.GetLogReceiver()

	// // Global middleware
	// router.Use(chimiddle.StripSlashes)

	// router.Route("/api", func(router chi.Router) {
	// 	router.Get("/health", lr.HealthHandler)
	// 	router.Post("/ingest", lr.ReceiveLogsHandler)
	// })
}
