package routes

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/resend/resend-go/v3"
	"github.com/ut-code/Raxcel/server/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusConflict, map[string]string{
				"error": "the email is already used",
			})
		}
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
		fmt.Println(err)
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
	//TODO: send html instead of json
	reqToken := c.QueryParam("token")
	database, err := db.ConnectDB()
	if err != nil {
		return c.String(http.StatusInternalServerError, "failed to connect to database")
	}
	var token db.Token
	if err := database.Where("token = ?", reqToken).First(&token).Error; err != nil {
		return c.String(http.StatusNotFound, "invalid verification token")
	}

	if time.Now().After(token.ExpiresAt) {
		return c.String(http.StatusBadRequest, "verification token has expired")
	}
	if err := database.Model(&db.User{}).Where("id = ?", token.UserId).Update("is_verified", true).Error; err != nil {
		return c.String(http.StatusInternalServerError, "failed to verify user")
	}
	database.Delete(&token)
	return c.String(http.StatusOK, "email verified!")
}

func Login(c echo.Context) error {
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
			"error": "email not verified",
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid email or password",
		})
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Id,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	secretKey := os.Getenv("SECRET_KEY")
	signedToken, _ := claims.SignedString([]byte(secretKey))
	return c.JSON(http.StatusOK, map[string]string{
		"message": "logged in",
		"token":   signedToken,
	})
}

func generateSecureToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func sendVerificationEmail(email, token string) error {
	apiKey := os.Getenv("RESEND_API_KEY")
	apiUrl := os.Getenv("API_URL")
	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		From:    "Raxcel <noreply@raxcel.utcode.net>",
		To:      []string{email},
		Subject: "Verify your account",
		Html:    fmt.Sprintf("<p>Click the link below to verify your email</p><a href=%s/verify-email?token=%s>Click here!</a>", apiUrl, token),
	}
	_, err := client.Emails.Send(params)
	return err
}

func CheckUser(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "missing authorization header",
		})
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid authorization format",
		})
	}

	secretKey := os.Getenv("SECRET_KEY")

	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprint(err),
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)
	userId := claims.Issuer

	database, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError,
			map[string]string{
				"error": fmt.Sprint(err),
			})
	}
	var user db.User
	err = database.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": fmt.Sprint(err),
		})
	}
	if !user.IsVerified {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "unauthorized",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"userId": user.Id,
	})
}
