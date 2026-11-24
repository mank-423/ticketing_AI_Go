package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUserBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginUserBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func generateHashedPassword(text string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
}

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var req RegisterUserBody

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Password) < 8 {
		http.Error(w, "Password must be at least 8 bytes", http.StatusBadRequest)
		return
	}

	if len(req.Password) > 72 {
		http.Error(w, "Password must not exceed 72 bytes", http.StatusBadRequest)
		return
	}

	hashedPassword, err := generateHashedPassword(req.Password)

	if err != nil {
		http.Error(w, "Hashing of password failed", http.StatusBadRequest)
		return
	}

	name := req.Name
	email := req.Email
	role := req.Role
	password := hashedPassword

	query := `
		INSERT INTO USER(name, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		returning user_id;
	`

	var userID string
	err = h.DB.QueryRow(r.Context(), query, name, email, password, role).Scan(&userID)

	if err != nil {
		http.Error(w, "Registration failed", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"user_id": userID,
		"name":    name,
		"email":   email,
		"role":    role,
	}

	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req LoginUserBody

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email := req.Email
	password := req.Password

	query := `SELECT password_hash from user where email=$1`

	var hashValueUser string
	err = h.DB.QueryRow(r.Context(), query, email).Scan(&hashValueUser)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashValueUser), []byte(password))

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"email":   email,
		"message": "Login successful!",
	})

}
