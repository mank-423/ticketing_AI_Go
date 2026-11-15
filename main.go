package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"platform_ticket/backend/database"
	"platform_ticket/backend/routes"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error identifying env variables")
	}

	// Getting port
	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	if port == "" {
		port = "8000"
	}

	db := database.ConnectDB(dbUrl)
	defer db.Close()

	// importing routes
	router := routes.Router()

	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}
