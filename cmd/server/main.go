package main

import (
	// "time"

	"github.com/KuberTheGreat/Sentrinet/internal/api"
	"github.com/KuberTheGreat/Sentrinet/internal/db"
	// "github.com/KuberTheGreat/Sentrinet/internal/scheduler"
	"github.com/gofiber/fiber/v2"
)

func main(){
	database := db.InitDB()
	app := fiber.New()

	api.SetupRoutes(app, database)

	// job := scheduler.Job{
	// 	Target: "scanme.nmap.org",
	// 	StartPort: 20,
	// 	EndPort: 100,
	// 	Interval: 1 * time.Minute,
	// }
	// scheduler.StartJob(database, job)

	app.Listen(":8080")
}