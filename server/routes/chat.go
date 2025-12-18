package routes

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/ut-code/Raxcel/server/db"
	"google.golang.org/genai"
)

type userMessage struct {
	Message string `json:"message"`
}

func Greet(c echo.Context) error {
	return c.String(http.StatusOK, "Hello from Echo!")
}
func ChatWithAI(c echo.Context) error {
	// Get userId from context (set by AuthMiddleware)
	userId, ok := c.Get("userId").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	// Parse user message
	message := new(userMessage)
	if err := c.Bind(message); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}

	// Connect to database
	database, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}

	// Save user message
	userMsg := db.Message{
		Id:      uuid.New().String(),
		UserId:  userId,
		Content: message.Message,
		Role:    "user",
	}
	if err := database.Create(&userMsg).Error; err != nil {
		log.Printf("Failed to save user message: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save message"})
	}

	// Setup gemini client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatal("Error creating Gemini client", err)
	}

	// Generate AI response
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

	// Save AI message
	assistantMsg := db.Message{
		Id:      uuid.New().String(),
		UserId:  userId,
		Content: aiMessage,
		Role:    "assistant",
	}
	if err := database.Create(&assistantMsg).Error; err != nil {
		log.Printf("Failed to save AI message: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save AI message"})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"aiMessage": aiMessage,
	})
}

func GetMessages(c echo.Context) error {
	// Get userId from context (set by AuthMiddleware)
	userId, ok := c.Get("userId").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Connect to database
	database, err := db.ConnectDB()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}

	// Get all messages for the user, ordered by creation time
	var messages []db.Message
	if err := database.Where("user_id = ?", userId).Order("created_at ASC").Find(&messages).Error; err != nil {
		log.Printf("Failed to fetch messages: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch messages"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages": messages,
	})
}
