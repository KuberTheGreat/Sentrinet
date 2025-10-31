package main

import (
	"github.com/KuberTheGreat/Sentrinet/internal/api"
	"github.com/KuberTheGreat/Sentrinet/internal/db"
	"github.com/gofiber/fiber/v2"
)

func main(){
	database := db.InitDB()
	app := fiber.New()

	api.SetupRoutes(app, database)

	app.Listen(":8080")
}