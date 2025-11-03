package auth

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret-key")

type AuthHandler struct{
	DB *sqlx.DB
}

func NewAuthHandler(db *sqlx.DB) *AuthHandler{
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error{
	var data struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&data); err!=nil{
		return fiber.ErrBadRequest
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil{
		return fiber.ErrInternalServerError
	}

	_, err = h.DB.Exec(`INSERT INTO users (username, password_hash) 
						VALUES (?, ?)`, 
						data.Username, string(hash))
	if err != nil{
		log.Fatal(err)
		return fiber.NewError(fiber.StatusConflict, "Username already taken")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully!"})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error{
	var data struct{
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&data); err != nil{
		return fiber.ErrBadRequest
	}

	var user User
	err := h.DB.Get(&user, `SELECT * FROM users WHERE username = ?`, data.Username)
	if err != nil{
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(data.Password)); err != nil{
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil{
		return fiber.ErrInternalServerError
	}

	return c.JSON(fiber.Map{"token": tokenString})
}