package handler

import (
	"context"
	"log"
	"net/http"
	"os"

	"google.golang.org/genai"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

type UserMessage struct {
	Message string `json:"message"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// load api key
	// can't load environmental variables from .env file in production
	if os.Getenv("VERCEL_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	// setup echo
	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Echo!")
	})

	e.POST("/messages", func(c echo.Context) error {
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			log.Fatal("GEMINI_API_KEY is not set")
		}
		// setup gemini client
		ctx := context.Background()
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey: apiKey,
		})
		if err != nil {
			log.Fatal("Error creating Gemini client", err)
		}
		message := new(UserMessage)
		if err := c.Bind(message); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
		}
		result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(message.Message), nil)
		if err != nil {
			log.Printf("Gemini API error: %v", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate content"})
		}
		var aiMessage string
		if result != nil && len(result.Candidates) > 0 {
			candidate := result.Candidates[0]
			if len(candidate.Content.Parts) > 0 {
				aiMessage = candidate.Content.Parts[0].Text
			}
		}
		return c.JSON(http.StatusCreated, map[string]string{
			"aiMessage": aiMessage,
		})
	})
	e.ServeHTTP(w, r)
}
