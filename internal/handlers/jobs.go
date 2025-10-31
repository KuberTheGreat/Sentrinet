package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type JobRow struct {
	ID              int64  `db:"id" json:"id"`
	Target          string `db:"target" json:"target"`
	StartPort       int    `db:"start_port" json:"start_port"`
	EndPort         int    `db:"end_port" json:"end_port"`
	IntervalSeconds int    `db:"interval_seconds" json:"interval_seconds"`
	Active          int    `db:"active" json:"active"`
	CreatedAt       string `db:"created_at" json:"created_at"`
}

func GetJobsHandler(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limitStr := c.Query("limit", "10")
		offsetStr := c.Query("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			offset = 0
		}

		var jobs []JobRow
		query := `SELECT * FROM jobs ORDER BY created_at DESC LIMIT ? OFFSET ?`
		if err := db.Select(&jobs, query, limit, offset); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"limit":  limit,
			"offset": offset,
			"data":   jobs,
		})
	}
}
