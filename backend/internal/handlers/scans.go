package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func GetScansHandler(db *sqlx.DB) fiber.Handler{
	return func(c *fiber.Ctx) error {
		limitStr := c.Query("limit", "10")
		offsetStr := c.Query("offset", "0")

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0{
			limit = 10
		}

		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0{
			offset = 0
		}

		var scans []struct {
			ID         int64  `db:"id" json:"id"`
			Target     string `db:"target" json:"target"`
			Port       int    `db:"port" json:"port"`
			IsOpen     bool   `db:"is_open" json:"is_open"`
			DurationMs int64  `db:"duration_ms" json:"duration_ms"`
			CreatedAt  string `db:"created_at" json:"created_at"`
			UserId int64 `db:"user_id" json:"user_id"`
		}

		query := `SELECT * FROM scans ORDER BY created_at DESC LIMIT ? OFFSET ?;`

		if err := db.Select(&scans, query, limit, offset); err != nil{
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"limit": limit,
			"offset": offset,
			"data": scans,
		})
	}
}