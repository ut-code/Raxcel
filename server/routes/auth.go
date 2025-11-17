package routes

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ut-code/Raxcel/server/db"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyEmailRequest struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	req := new(RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid format",
		})
	}
	// validation
	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "email and password are required",
		})
	}
	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "password must be at least 8 characters",
		})
	}
	database, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to connect to database",
		})
	}
	var existingUser db.User
	result := database.Where("email = ?", req.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "email already registerd",
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to hash password",
		})
	}
	user := db.User{
		Id:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		IsVerified:   false,
	}
	if err := database.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create user",
		})
	}
	tokenString := generateSecureToken()
	token := db.Token{
		Id:        uuid.New().String(),
		UserId:    user.Id,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := database.Create(&token).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create verification token",
		})
	}

	if err := sendVerificationEmail(user.Email, tokenString); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to send verification email",
		})
	}
	return c.JSON(http.StatusCreated, map[string]string{
		"message": "user created",
		"userId":  user.Id,
	})
}

func VerifyEmail(c echo.Context) error {
	req := new(VerifyEmailRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}
	database, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to connect to database",
		})
	}
	var token db.Token
	if err := database.Where("token = ?", req.Token).First(&token).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "invalid verification token",
		})
	}

	if time.Now().After(token.ExpiresAt) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "verification token has expired",
		})
	}
	if err := database.Model(&db.User{}).Where("id = ?", token.UserId).Update("is_verified", true).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to verify user",
		})
	}
	database.Delete(&token)
	return c.JSON(http.StatusOK, map[string]string{
		"message": "email verified",
	})
}

func Login(c echo.Context) error {
	// Validate Email and Password
	// Store JWT in cookie
	req := new(LoginRequest)
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}
	database, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to connect to database",
		})
	}
	var user db.User
	if err := database.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}
	if !user.IsVerified {
		return c.JSON(http.StatusForbidden, map[string]string{
			"error": "email not verifies",
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(req.Password), []byte(user.PasswordHash)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid email or password",
		})
	}
	return nil
}

func generateSecureToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sendVerificationEmail(email, token string) error {
	// Send email with resend
	return nil
}
