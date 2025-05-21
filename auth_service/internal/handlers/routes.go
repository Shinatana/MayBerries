package handlers

import (
	"net/http"

	"auth_service/config"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, cfg *config.Config) {
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterHandler(w, cfg)
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		LoginHandler(w, r, cfg)
	})

	// Пример публичного хелсчека
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
}
