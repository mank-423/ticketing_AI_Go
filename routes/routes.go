package routes

import (
	"net/http"
	"platform_ticket/backend/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.HealthCheckup)

	return mux
}
