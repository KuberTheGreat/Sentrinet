package main

import (
	// "time"

	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KuberTheGreat/Sentrinet/internal/api"
	"github.com/KuberTheGreat/Sentrinet/internal/db"
	"github.com/KuberTheGreat/Sentrinet/internal/metrics"
	"github.com/KuberTheGreat/Sentrinet/internal/realtime"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)



func main(){
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	database := db.InitDB()
	db.StartCleanupScheduler(database)

	app := fiber.New()

	registry := prometheus.NewRegistry()

	//register custom metrics
	metrics.Register(registry)

	prom := fiberprometheus.NewWithRegistry(registry, "sentrinet", "", "", nil)
	prom.RegisterAt(app, "/metrics")
	app.Use(prom.Middleware)
	
	

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	wsManager := realtime.NewManager()
	
	api.SetupRoutes(app, database, wsManager)

	go func() {
		if err := app.Listen(":8080"); err != nil{
			fmt.Println("Fiber listen error: ", err)
			cancel()
		}
	}()
	
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	select {
	case <-sig:
		fmt.Println("shutdown signal received")
		cancel()
	case <-ctx.Done():
	}

	_, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = app.Shutdown()
	fmt.Println("Server stopped")
}