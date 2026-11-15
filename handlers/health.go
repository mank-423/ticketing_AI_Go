package handlers

import (
	"encoding/json"
	"net/http"
)

func HealthCheckup(w http.ResponseWriter, r *http.Request) {

	response := map[string]string{
		"status":  "ok",
		"message": "Server is healthy!",
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)

}
