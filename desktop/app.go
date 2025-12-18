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

type MessageItem struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Content   string `json:"content"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

type GetMessagesResult struct {
	Ok       bool          `json:"ok"`
	Messages []MessageItem `json:"messages"`
	Error    string        `json:"error,omitempty"`
}

func (a *App) GetMessages() GetMessagesResult {
	fmt.Println("GetMessages called")

	apiUrl := getAPIURL()
	fmt.Println("API URL:", apiUrl)

	jwt, err := keyring.Get("Raxcel", "raxcel-user")
	if err != nil {
		fmt.Println("Keyring error:", err)
		return GetMessagesResult{
			Ok:    false,
			Error: fmt.Sprint(err),
		}
	}
	fmt.Println("JWT retrieved successfully")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/messages", apiUrl), nil)
	if err != nil {
		fmt.Println("Request creation error:", err)
		return GetMessagesResult{
			Ok:    false,
			Error: fmt.Sprint(err),
		}
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	fmt.Println("Sending request to:", req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return GetMessagesResult{
			Ok:    false,
			Error: fmt.Sprint(err),
		}
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Read body error:", err)
		return GetMessagesResult{
			Ok:    false,
			Error: fmt.Sprint(err),
		}
	}

	fmt.Println("Response body:", string(body))
	var serverResponse map[string]interface{}
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		return GetMessagesResult{
			Ok:    false,
			Error: fmt.Sprint(err),
		}
	}

	// Parse messages array
	var messages []MessageItem
	if messagesData, ok := serverResponse["messages"].([]interface{}); ok {
		fmt.Printf("Found %d messages\n", len(messagesData))
		for _, msgData := range messagesData {
			if msgMap, ok := msgData.(map[string]interface{}); ok {
				message := MessageItem{
					Id:        msgMap["id"].(string),
					UserId:    msgMap["userId"].(string),
					Content:   msgMap["content"].(string),
					Role:      msgMap["role"].(string),
					CreatedAt: msgMap["createdAt"].(string),
				}
				messages = append(messages, message)
			}
		}
	} else {
		fmt.Println("No messages found in response")
	}

	fmt.Printf("Returning %d messages\n", len(messages))
	return GetMessagesResult{
		Ok:       true,
		Messages: messages,
	}
}

func (a *App) ChatWithAI(message string) ChatResult {
	fmt.Println("ChatWithAI called with message:", message)

	postData := ChatRequest{
		Message: message,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		fmt.Println("Marshal error:", err)
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	apiUrl := getAPIURL()
	fmt.Println("API URL:", apiUrl)

	jwt, err := keyring.Get("Raxcel", "raxcel-user")
	if err != nil {
		fmt.Println("Keyring error:", err)
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	fmt.Println("JWT retrieved successfully")

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/messages", apiUrl), bytes.NewReader(jsonData))
	if err != nil {
		fmt.Println("Request creation error:", err)
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	fmt.Println("Sending request to:", req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.StatusCode)
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
		fmt.Println("Unmarshal error:", err)
		fmt.Println("Response body:", string(body))
		return ChatResult{
			Ok:      false,
			Message: fmt.Sprint(err),
		}
	}
	fmt.Println("AI Message:", serverResponse["aiMessage"])
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
