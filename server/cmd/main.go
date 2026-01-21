package main

import (
	"log"
	"os"

	"github.com/NetlutZ/subscout/internal/database"
	"github.com/NetlutZ/subscout/internal/handler"
	"github.com/NetlutZ/subscout/internal/repository"
	"github.com/NetlutZ/subscout/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File : ", err)
	}

	// Connect to Database
	db, err := database.DatabaseConnect()
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	log.Println("PostgreSQL connected")
	defer db.Close()

	// Create Table
	query := database.Migrate()
	_, err = db.Exec(query)
	if err != nil {
		log.Fatal("Error while creating/migrating database: ", err)
	}

	subscriptionRepositoryDB := repository.NewSubscriptionRepositoryDB(db)
	subscriptionService := service.NewSubscriptionService(subscriptionRepositoryDB)

	app := fiber.New()
	handler.RegisterSubscriptionRoutes(app, subscriptionService)

	userRepo := repository.NewUserRepositoryDB(db)
	authService := service.NewAuthService(userRepo)
	handler.RegisterAuthRoutes(app, authService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
