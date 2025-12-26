package main

import (
	"log"
	"os"

	"github.com/NetlutZ/subscout/internal/auth"
	"github.com/NetlutZ/subscout/internal/db"
	"github.com/NetlutZ/subscout/internal/subscriptions"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File : ", err)
	}

	if err := db.Connect(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}
	log.Println("PostgreSQL connected")
	defer db.DB.Close()

	err = db.Migrate()
	if err != nil {
		log.Fatal("error while migrating to database: ", err)
	}

	app := fiber.New()
	authGroup := app.Group("/auth")
	auth.NewService().RegisterRoute(authGroup)

	api := app.Group("/api", auth.Protected())
	subscriptions.NewService().RegisterRoute(api)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen("0.0.0.0:" + port))
}
