package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/netosts/goledger-challenge-besu/internal/database"
	"github.com/netosts/goledger-challenge-besu/internal/handlers"
	"github.com/netosts/goledger-challenge-besu/internal/repositories"
	"github.com/netosts/goledger-challenge-besu/internal/routes"
	"github.com/netosts/goledger-challenge-besu/internal/usecases"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	dbConfig := database.NewConfig()
	db, err := dbConfig.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.InitializeSchema(db); err != nil {
		log.Fatalf("Failed to initialize database schema: %v", err)
	}

	dbRepo := repositories.NewPostgresRepository(db)

	contractUseCase := usecases.NewContractUseCase(dbRepo)

	handler := handlers.NewHandler(contractUseCase)

	router := routes.SetupRoutes(handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
