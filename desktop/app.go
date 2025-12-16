package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/99designs/keyring"
	"github.com/joho/godotenv"

	"net/http"
)

// ANSIカラーコード
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
)

var apiURL string

func getAPIURL() string {
	if apiURL == "" {
		godotenv.Load()
		apiUrl := os.Getenv("PUBLIC_API_URL")
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

func getKeyring() keyring.Keyring {
	fmt.Println(ColorYellow + "[Keyring] Getting user home directory..." + ColorReset)
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(ColorRed + "[Keyring] Failed to get home directory: " + err.Error() + ColorReset)
		return nil
	}
	fmt.Println(ColorGreen + "[Keyring] Home directory: " + ColorReset + homedir)

	configPath := homedir + "/.config/raxcel"
	fmt.Println(ColorYellow + "[Keyring] Opening keyring with config path: " + ColorReset + configPath)

	ring, err := keyring.Open(keyring.Config{
		ServiceName: "Raxcel",
		AllowedBackends: []keyring.BackendType{
			keyring.SecretServiceBackend,
			keyring.KeychainBackend,
			keyring.WinCredBackend,
			keyring.FileBackend,
		},
		FileDir: configPath,
	})

	if err != nil {
		fmt.Println(ColorRed + "[Keyring] Failed to open keyring: " + err.Error() + ColorReset)
		return nil
	}

	fmt.Println(ColorGreen + "[Keyring] Successfully opened keyring" + ColorReset)
	return ring
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
	apiUrl := getAPIURL()
	ring := getKeyring()
	item, err := ring.Get("raxcel-user")
	jwt := item.Data
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
	fmt.Println(ColorYellow + "[Login] Starting login process for: " + ColorReset + email)

	postData := LoginRequest{
		Email:    email,
		Password: password,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		fmt.Println(ColorRed + "[Login] Failed to marshal request" + ColorReset)
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to marshal request: %v", err),
		}
	}

	apiUrl := getAPIURL()
	fmt.Println(ColorYellow + "[Login] Sending request to API: " + ColorReset + apiUrl + "/login")

	resp, err := http.Post(fmt.Sprintf("%s/login", apiUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println(ColorRed + "[Login] Failed to send request" + ColorReset)
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to send request: %v", err),
		}
	}
	defer resp.Body.Close()

	fmt.Println(ColorYellow + "[Login] Response status: " + ColorReset + resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(ColorRed + "[Login] Failed to read response body" + ColorReset)
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to read response: %v", err),
		}
	}

	var serverResponse map[string]string
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		fmt.Println(ColorRed + "[Login] Failed to parse response JSON" + ColorReset)
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to parse response %v", err),
		}
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(ColorRed + "[Login] Login failed: " + serverResponse["error"] + ColorReset)
		return LoginResult{
			Ok:      false,
			Message: serverResponse["error"],
		}
	}

	token := serverResponse["token"]
	fmt.Println(ColorGreen + "[Login] Successfully received token: " + ColorReset + ColorCyan + token + ColorReset)

	fmt.Println(ColorYellow + "[Login] Storing token in keyring..." + ColorReset)
	ring := getKeyring()
	if ring == nil {
		fmt.Println(ColorRed + "[Login] Failed to get keyring" + ColorReset)
		return LoginResult{
			Ok:      false,
			Message: "Failed to access keyring",
		}
	}

	err = ring.Set(keyring.Item{
		Key:  "raxcel-user",
		Data: []byte(token),
	})
	if err != nil {
		fmt.Println(ColorRed + "[Login] Failed to store token in keyring: " + err.Error() + ColorReset)
		return LoginResult{
			Ok:      false,
			Message: fmt.Sprintf("Failed to store token: %v", err),
		}
	}

	fmt.Println(ColorGreen + "[Login] Successfully stored token in keyring" + ColorReset)
	fmt.Println(ColorGreen + "✓ Login completed successfully!" + ColorReset)
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
	apiUrl := getAPIURL()
	ring := getKeyring()
	item, _ := ring.Get("raxcel-user")
	token := item.Data
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
