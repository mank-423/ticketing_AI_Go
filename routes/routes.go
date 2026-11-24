package routes

import (
	"net/http"
	"platform_ticket/backend/handlers"
)

func Router(h *handlers.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheckup)
	mux.HandleFunc("POST /register", h.RegisterUser)
	mux.HandleFunc("POST /login", h.LoginUser)

	return mux
}
