package handlers

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
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

	// Get payload items
	email := req.Email
	password := req.Password

	query := `SELECT user_id, password_hash from user where email=$1`

	var hashValueUser string
	var userID string

	// Look for user
	err = h.DB.QueryRow(r.Context(), query, email).Scan(&userID, &hashValueUser)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}

	// Compare the password
	err = bcrypt.CompareHashAndPassword([]byte(hashValueUser), []byte(password))

	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// create token
	buff := make([]byte, 32)
	_, err = rand.Read(buff)

	if err != nil {
		http.Error(w, "Error token generation", http.StatusInternalServerError)
		return
	}

	rawToken := base64.RawURLEncoding.EncodeToString(buff) //client token
	hash := sha256.Sum256([]byte(rawToken))                // hash the string
	hashStr := hex.EncodeToString(hash[:])                 //encode the hash

	query = `
		insert into user_sessions(user_id, token_hash, expires_at)
		values($1, $2, NOW() + INTERVAL '24 hours')
	`

	// insert into table
	_, err = h.DB.Exec(r.Context(), query, userID, hashStr)

	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"email":        email,
		"access-token": rawToken,
		"message":      "Login successful!",
	})

}
