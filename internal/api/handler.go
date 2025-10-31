package api

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/KuberTheGreat/Sentrinet/internal/models"
	"github.com/KuberTheGreat/Sentrinet/internal/scan"
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
				fmt.Println("Insert error: ", err)
			} else{
				id, _ := res.LastInsertId()
				fmt.Println("Inserted row ID:", id)
			}
		}

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
}