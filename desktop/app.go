package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/zalando/go-keyring"

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

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResult struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	UserId  string `json:"userId,omitempty"`
}

func (a *App) Register(email, password string) RegisterResult {
	postData := RegisterRequest{
		Email:    email,
		Password: password,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return RegisterResult{
			Ok:      false,
			Message: fmt.Sprintf("failed to marshal request: %v", err),
		}
	}
	godotenv.Load()
	apiUrl := os.Getenv("PUBLIC_API_URL")

	resp, err := http.Post(fmt.Sprintf("%s/register", apiUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return RegisterResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to send request: %v", err),
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return RegisterResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to read response: %v", err),
		}
	}
	var serverResponse map[string]string
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return RegisterResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to parse response: %v", err),
		}
	}
	if resp.StatusCode != http.StatusCreated {
		return RegisterResult{
			Ok:      false,
			Message: serverResponse["error"],
		}
	}
	return RegisterResult{
		Ok:      true,
		Message: serverResponse["message"],
		UserId:  serverResponse["userId"],
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResult struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

func (a *App) Login(email, password string) LoginResult {
	postData := LoginRequest{
		Email:    email,
		Password: password,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to marchal request: %v", err),
		}
	}
	godotenv.Load()
	apiUrl := os.Getenv("PUBLIC_API_URL")

	resp, err := http.Post(fmt.Sprintf("%s/login", apiUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to send request: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to read response: %v", err),
		}
	}

	var serverResponse map[string]string
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to parse response %v", err),
		}
	}
	if resp.StatusCode != http.StatusOK {
		return LoginResult{
			Ok:      false,
			Message: serverResponse["error"],
		}
	}
	token := serverResponse["token"]
	err = keyring.Set("Raxcel", email, token)
	if err != nil {
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to store token: %v", err),
		}
	}
	return LoginResult{
		Ok:      true,
		Message: serverResponse["message"],
		Token:   token,
	}
}
