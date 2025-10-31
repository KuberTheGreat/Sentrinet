package api

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/KuberTheGreat/Sentrinet/internal/handlers"
	"github.com/KuberTheGreat/Sentrinet/internal/models"
	"github.com/KuberTheGreat/Sentrinet/internal/scan"
	"github.com/KuberTheGreat/Sentrinet/internal/scheduler"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type ScanRequest struct {
	Target string `json:"target"`
	StartPort int `json:"start_port"`
	EndPort int `json:"end_port"`
}

func SetupRoutes(app *fiber.App, db *sqlx.DB){
	app.Post("/scan", func(c *fiber.Ctx) error {
		var req ScanRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		results := scan.ScanRange(req.Target, req.StartPort, req.EndPort)
		for _, r := range results{
			res, err := db.NamedExec(
				"INSERT INTO scans (target, port, is_open, duration_ms) VALUES (:target, :port, :is_open, :duration_ms)",
				map[string]interface{}{
					"target": req.Target,
					"port": r.Port,
					"is_open": r.IsOpen,
					"duration_ms": r.Duration,
				},
			)

			if err != nil{
				handlers.CreateNotification(db, 1, 1, "scan_failed", fmt.Sprintf("Scan for %s failed to complete.", req.Target))
				fmt.Println("Insert error: ", err)
			} else{
				id, _ := res.LastInsertId()
				fmt.Println("Inserted row ID:", id)
			}
		}

		handlers.CreateNotification(db, 1, 1, "scan_completed", fmt.Sprintf("Scan for %s is complete!", req.Target))
		return c.JSON(results)
	})

	app.Get("/scans", func(c *fiber.Ctx) error {
		target := c.Query("target", "")
		openOnly := c.Query("open_only", "")
		limitStr := c.Query("limit", "50")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0{
			limit = 50
		}
		
		query := `SELECT * FROM scans WHERE 1=1`
		args := []interface{}{}

		if target != ""{
			query += " AND target LIKE ?"
			args = append(args, "%"+target+"%")
		}

		if strings.ToLower(openOnly) == "true"{
			query += " AND is_open = true"
		}

		query += " ORDER BY created_at DESC LIMIT ?"
		args = append(args, limit)
		
		scans := []models.ScanResult{}
		err = db.Select(&scans, query, args...)
		if err != nil{
			log.Println("DB select error: ", err)
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(scans)
	})

	app.Get("/stats", func(c *fiber.Ctx) error {
		var totalScans int
		var openPorts int
		var avgDuration float64

		if err := db.Get(&totalScans, "SELECT COUNT(*) FROM scans"); err != nil {
			totalScans = 0
		}

		if err := db.Get(&openPorts, "SELECT COUNT(*) FROM scans WHERE is_open = true"); err != nil{
			openPorts = 0
		}

		if err := db.Get(&avgDuration, "SELECT AVG(duration_ms) FROM scans"); err != nil{
			avgDuration = 0
		}

		stats := fiber.Map{
			"total_scans": totalScans,
			"open_ports": openPorts,
			"avg_scan_time_ms": avgDuration,
		}

		return c.JSON(stats)
	})

	app.Delete("/scans/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := db.Exec("DELETE FROM scans WHERE id = ?", id)
		if err != nil{
			log.Println("Delete error: ", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete record"})
		}

		return c.JSON(fiber.Map{"message": fmt.Sprintf("Deleted scan with id %s", id)})
	})

	app.Delete("/scans", func(c *fiber.Ctx) error{
		target := c.Query("target", "")
		if target == ""{
			return c.Status(400).JSON(fiber.Map{"error": "target query parameter required"})
		}

		res, err := db.Exec("DELETE FROM scans WHERE target = ?", target)
		if err != nil{
			log.Println("Delete error: ", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete records"})
		}

		count, _ := res.RowsAffected()
		return c.JSON(fiber.Map{"message": fmt.Sprintf("Delete %d scans for target %s", count, target)})
	})

	schManager := scheduler.NewManager(context.Background(), db)
	if err := schManager.LoadAndStartAll(); err != nil {
		fmt.Println("[Scheduler] failed to load jobs: ", err)
	}

	app.Post("/schedules", func(c *fiber.Ctx) error {
		var req struct{
			Target    string `json:"target"`
			StartPort int    `json:"start_port"`
			EndPort   int    `json:"end_port"`
			IntervalSeconds int `json:"interval_seconds"`
			Active    bool   `json:"active"`
		}

		if err := c.BodyParser(&req); err != nil{
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		if req.Target == "" || req.StartPort <= 0 || req.EndPort <= 0 || req.IntervalSeconds <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "target, ports and interval seconds are required"})
		}
		interval := time.Duration(req.IntervalSeconds) * time.Second
		id, err := schManager.CreateJob(req.Target, req.StartPort, req.EndPort, interval, req.Active)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"id": id})
	})

	app.Get("/schedules", func(c *fiber.Ctx) error {
		rows, err := schManager.ListJobs()
		if err != nil{
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(rows)
	})

	app.Post("/schedules/:id/stop", func(c *fiber.Ctx) error{
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil{
			return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
		}
		if err := schManager.StopJob(id); err != nil{
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"message": "stopped"})
	})

	app.Post("/schedules/:id/start", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error":"invalid id"})
		}
		if err := schManager.StartJobByID(id); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message":"started"})
	})

	app.Delete("/schedules/:id", func(c *fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error":"invalid id"})
		}
		if err := schManager.DeleteJob(id); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"message":"deleted"})
	})

	app.Get("/api/scans", handlers.GetScansHandler(db))
	app.Get("/api/jobs", handlers.GetJobsHandler(db))
	app.Get("/notifications/:userId", handlers.GetUserNotifications(db))
	app.Put("/notifications/:id/read", handlers.MarkNotificationRead(db))
}