package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"

	"net/http"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

type UserRequest struct {
	Message string `json:"message"`
}

type ServerResponse struct {
	AiMessage string `json:"aiMessage"`
}

type ChatResult struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (a *App) ChatWithAI(message string) ChatResult {
	postData := UserRequest{
		Message: message,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	godotenv.Load(".env")
	apiUrl := os.Getenv("PUBLIC_API_URL")
	resp, err := http.Post(fmt.Sprintf("%s/messages", apiUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	var serverResponse ServerResponse
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	return ChatResult{
		Ok:      true,
		Message: serverResponse.AiMessage,
	}
}
