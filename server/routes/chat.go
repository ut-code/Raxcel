package routes

import (
	"context"
	"fmt"
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
	log.Println("ChatWithAI called")

	// Get userId from context (set by AuthMiddleware)
	userId, ok := c.Get("userId").(string)
	log.Println("UserId from context:", userId, "ok:", ok)
	if !ok {
		log.Println("Unauthorized: userId not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	// Parse user message
	message := new(userMessage)
	if err := c.Bind(message); err != nil {
		log.Println("Failed to bind message:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON"})
	}
	log.Println("Received message:", message.Message)

	// Connect to database
	database, err := db.ConnectDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}
	log.Println("Database connected")

	// Save user message
	userMsg := db.Message{
		Id:      uuid.New().String(),
		UserId:  userId,
		Content: message.Message,
		Role:    "user",
	}
	log.Println("Saving user message to database...")
	if err := database.Create(&userMsg).Error; err != nil {
		log.Printf("Failed to save user message: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save message"})
	}
	log.Println("User message saved successfully")

	// Get recent messages (最新7件取得して、最後の1件=現在のメッセージを除外)
	log.Println("Fetching recent messages for context...")
	var recentMessages []db.Message
	err = database.Where("user_id = ?", userId).Order("created_at DESC").Limit(7).Find(&recentMessages).Error
	if err != nil {
		log.Printf("Failed to get recent messages: %v", err)
	}
	log.Printf("Retrieved %d messages", len(recentMessages))

	// 現在のメッセージを除外して、時系列順に並び替え
	var contextMessages []db.Message
	for i := len(recentMessages) - 1; i >= 1; i-- { // i >= 1 で最新(現在のメッセージ)を除外
		contextMessages = append(contextMessages, recentMessages[i])
	}
	log.Printf("Using %d messages as context", len(contextMessages))

	// Build prompt with conversation history
	prompt := ""
	if len(contextMessages) > 0 {
		prompt = "Previous conversation:\n"
		for _, m := range contextMessages {
			if m.Role == "user" {
				prompt += fmt.Sprintf("User: %s\n", m.Content)
			} else {
				prompt += fmt.Sprintf("Assistant: %s\n", m.Content)
			}
		}
		prompt += "\nCurrent message:\n"
	}
	prompt += fmt.Sprintf("User: %s", message.Message)
	log.Printf("Prompt length: %d characters", len(prompt))

	// Setup gemini client
	log.Println("Creating Gemini client...")
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Fatal("Error creating Gemini client", err)
	}

	// Generate AI response
	log.Println("Generating AI response...")
	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(prompt), nil)
	if err != nil {
		log.Printf("Gemini API error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate content"})
	}
	log.Println("AI response generated")

	var aiMessage string
	if result != nil && len(result.Candidates) > 0 {
		candidate := result.Candidates[0]
		if len(candidate.Content.Parts) > 0 {
			aiMessage = candidate.Content.Parts[0].Text
		}
	}
	log.Println("AI Message content:", aiMessage)

	// Save AI message
	log.Println("Saving AI message to database...")
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
	log.Println("AI message saved successfully")

	log.Println("Returning response to client")
	return c.JSON(http.StatusCreated, map[string]string{
		"aiMessage": aiMessage,
	})
}

func GetMessages(c echo.Context) error {
	log.Println("GetMessages called")

	// Get userId from context (set by AuthMiddleware)
	userId, ok := c.Get("userId").(string)
	log.Println("UserId from context:", userId, "ok:", ok)
	if !ok {
		log.Println("Unauthorized: userId not found in context")
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	// Connect to database
	log.Println("Connecting to database...")
	database, err := db.ConnectDB()
	if err != nil {
		log.Println("Database connection error:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database connection failed"})
	}
	log.Println("Database connected")

	// Get all messages for the user, ordered by creation time
	log.Println("Fetching messages for user:", userId)
	var messages []db.Message
	if err := database.Where("user_id = ?", userId).Order("created_at ASC").Find(&messages).Error; err != nil {
		log.Printf("Failed to fetch messages: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch messages"})
	}
	log.Printf("Found %d messages", len(messages))

	log.Println("Returning messages to client")
	return c.JSON(http.StatusOK, map[string]interface{}{
		"messages": messages,
	})
}
