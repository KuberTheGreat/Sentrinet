package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(c *fiber.Ctx) error{
	authHeader := c.Get("Authorization")
	if authHeader == ""{
		return fiber.NewError(fiber.StatusUnauthorized, "Missing token")
	}

	token, err := jwt.Parse(authHeader, func(t *jwt.Token) (interface{}, error){
		return jwtSecret, nil
	})

	if err != nil || !token.Valid{
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user_id", int64(claims["user_id"].(float64)))

	return c.Next()
}