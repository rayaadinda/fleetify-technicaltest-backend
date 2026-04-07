package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/rayaadinda/fleetify-technicaltest-backend/internal/middleware"
)

type AuthHandler struct {
	jwtSecret string
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type staticUser struct {
	ID       uint
	Role     string
	Password string
}

var staticUsers = map[string]staticUser{
	"admin":  {ID: 1, Role: "admin", Password: "admin123"},
	"kerani": {ID: 2, Role: "kerani", Password: "kerani123"},
}

func NewAuthHandler(jwtSecret string) *AuthHandler {
	return &AuthHandler{jwtSecret: jwtSecret}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid request body"})
	}

	user, ok := staticUsers[req.Username]
	if !ok || user.Password != req.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "invalid username or password"})
	}

	token, err := middleware.GenerateToken(user.ID, user.Role, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "failed to generate token"})
	}

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":       user.ID,
			"username": req.Username,
			"role":     user.Role,
		},
	})
}
