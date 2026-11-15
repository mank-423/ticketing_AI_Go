package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func connectDB() *pgxpool.Pool {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading DB Url")
	}

	// Getting the DB URL
	dbUrl := os.Getenv("DB_URL")

	// Connection pool
	ctx := context.Background()

	db, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
	}

	defer db.Close()

	// Verify the connection is working or not
	err = db.Ping(ctx)
	if err != nil {
		db.Close()
		log.Fatal("Database is not reachable:", err)
	}

	fmt.Println("Connected to the database!")
	return db

}
