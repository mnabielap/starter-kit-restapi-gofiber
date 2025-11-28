package main

import (
	"log"

	"starter-kit-restapi-gofiber/internal/config"
	"starter-kit-restapi-gofiber/internal/database"
	"starter-kit-restapi-gofiber/internal/routes"

	_ "starter-kit-restapi-gofiber/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// @title           GoFiber Starter Kit API
// @version         1.0
// @description     Full Gofiber Boilerplate
// @host            localhost:3000
// @BasePath        /v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg)

	app := fiber.New(fiber.Config{AppName: "GoFiber App"})

	app.Use(logger.New())
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(recover.New())

	routes.SetupRoutes(app, db, cfg)

	log.Printf("Server running on port %s", cfg.Port)
	log.Printf("Swagger: http://localhost:%s/swagger/", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}