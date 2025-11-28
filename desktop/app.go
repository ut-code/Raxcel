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

var apiURL string

func getAPIURL() string {
	if apiURL == "" {
		apiUrl := os.Getenv("PUBLIC_API_KEY")
		return apiUrl
	}
	return apiURL
}

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

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResult struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (a *App) ChatWithAI(message string) ChatResult {
	postData := ChatRequest{
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
	apiUrl := getAPIURL()
	jwt, err := keyring.Get("Raxcel", "raxcel-user")
	if err != nil {
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/messages", apiUrl), bytes.NewReader(jsonData))
	if err != nil {
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	client := &http.Client{}
	resp, err := client.Do(req)
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
	var serverResponse map[string]string
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	return ChatResult{
		Ok:      true,
		Message: serverResponse["aiMessage"],
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
	apiUrl := getAPIURL()

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
	Token   string `json:"token,omitempty"`
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
			Message: fmt.Sprintf("Failed to marshal request: %v", err),
		}
	}
	godotenv.Load()
	apiUrl := getAPIURL()

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
	err = keyring.Set("Raxcel", "raxcel-user", token)
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

type CheckResult struct {
	Ok     bool   `json:"ok"`
	UserId string `json:"userId,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (a *App) CheckUser() CheckResult {
	godotenv.Load()
	apiUrl := getAPIURL()
	token, _ := keyring.Get("Raxcel", "raxcel-user")
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/user", apiUrl), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return CheckResult{
			Ok:    false,
			Error: fmt.Sprint(err),
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CheckResult{
			Ok:    false,
			Error: fmt.Sprint(err),
		}
	}
	var serverResponse map[string]string
	if err := json.Unmarshal(body, &serverResponse); err != nil {
		return CheckResult{
			Ok:    false,
			Error: fmt.Sprint(err),
		}
	}
	if serverResponse["error"] != "" {
		return CheckResult{
			Ok:    false,
			Error: serverResponse["error"],
		}
	}
	return CheckResult{
		Ok:     true,
		UserId: serverResponse["userId"],
	}
}
