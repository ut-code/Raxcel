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

var e *echo.Echo
var client *genai.Client

type UserMessage struct {
	Message string `json:"message"`
}

func init() {
	// load api key
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	// setup gemini client
	ctx := context.Background()
	client, err = genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatal("Error creating Gemini client", err)
	}
	// setup echo
	e = echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Echo!")
	})

	e.POST("/messages", func(c echo.Context) error {
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
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"aiMessage": aiMessage,
		})
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	e.ServeHTTP(w, r)
}

func Local() {
	e.Logger.Fatal(e.Start(":1323"))
}
