package handlers

import (
	"github.com/KuberTheGreat/Sentrinet/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func CreateNotification(db *sqlx.DB, userID int, scanID int, notifType, msg string) error{
	_, err := db.Exec(
		`INSERT INTO notifications (user_id, scan_id, type, message)
		VALUES (?, ?, ?, ?)`,
		userID, scanID, notifType, msg)

	return err
}

func GetUserNotifications(db *sqlx.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Params("userId")
		notifs := []models.Notification{}

		err := db.Select(&notifs, "SELECT * FROM notifications WHERE user_id = ? ORDER BY created_at DESC", userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(notifs)
	}
}

func MarkNotificationRead(db *sqlx.DB) fiber.Handler{
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		_, err := db.Exec("UPDATE notifications SET read = 1 WHERE id = ?", id)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.SendStatus(fiber.StatusOK)
	}
}