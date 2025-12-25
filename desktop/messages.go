package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ut-code/Raxcel/server/types"
	"github.com/zalando/go-keyring"
)

type Mesaage struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Content   string `json:"content"`
	Role      string `json:"role"`
	CreatedAt string `json:"createdAt"`
}

type LoadChatHistoryResult struct {
	Messages []Mesaage `json:"messages"`
	Error    string    `json:"error"`
}

func (a *App) LoadChatHistory() LoadChatHistoryResult {
	apiUrl := getAPIURL()

	jwt, err := keyring.Get("Raxcel", "raxcel-user")
	if err != nil {
		return LoadChatHistoryResult{
			Messages: []Mesaage{},
			Error:    fmt.Sprint(err),
		}
	}
	fmt.Println("JWT retrieved successfully")

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/messages", apiUrl), nil)
	if err != nil {
		return LoadChatHistoryResult{
			Messages: []Mesaage{},
			Error:    fmt.Sprint(err),
		}
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	fmt.Println("Sending request to:", req.URL.String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return LoadChatHistoryResult{
			Messages: []Mesaage{},
			Error:    fmt.Sprint(err),
		}
	}
	defer resp.Body.Close()

	fmt.Println("Response status:", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LoadChatHistoryResult{
			Messages: []Mesaage{},
			Error:    fmt.Sprint(err),
		}
	}

	var serverResponse types.LoadChatHistoryResponse
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return LoadChatHistoryResult{
			Messages: []Mesaage{},
			Error:    fmt.Sprintf("Failed to parse response: %v", err),
		}
	}

	// Check middleware error
	if serverResponse.AuthMiddlewareReturn != nil && serverResponse.MiddlewareError != "" {
		return LoadChatHistoryResult{
			Messages: []Mesaage{},
			Error:    serverResponse.MiddlewareError,
		}
	}

	// Check handler error
	if serverResponse.Error != "" {
		return LoadChatHistoryResult{
			Messages: []Mesaage{},
			Error:    serverResponse.Error,
		}
	}

	// Convert db.Message to Mesaage
	messages := make([]Mesaage, len(serverResponse.Messages))
	for i, msg := range serverResponse.Messages {
		messages[i] = Mesaage{
			Id:        msg.Id,
			UserId:    msg.UserId,
			Content:   msg.Content,
			Role:      msg.Role,
			CreatedAt: msg.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return LoadChatHistoryResult{
		Messages: messages,
		Error:    "",
	}
}

type ChatWithAIResult struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (a *App) ChatWithAI(message string, spreadsheetContext string) ChatWithAIResult {
	postData := types.ChatWithAIRequest{
		Message:            message,
		SpreadsheetContext: spreadsheetContext,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return ChatWithAIResult{
			Message: "",
			Error:   fmt.Sprint(err),
		}
	}
	apiUrl := getAPIURL()

	jwt, err := keyring.Get("Raxcel", "raxcel-user")
	if err != nil {
		return ChatWithAIResult{
			Message: "",
			Error:   fmt.Sprint(err),
		}
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/messages", apiUrl), bytes.NewReader(jsonData))
	if err != nil {
		return ChatWithAIResult{
			Message: "",
			Error:   fmt.Sprint(err),
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ChatWithAIResult{
			Message: "",
			Error:   fmt.Sprint(err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChatWithAIResult{
			Message: "",
			Error:   fmt.Sprint(err),
		}
	}

	var serverResponse types.ChatWithAIResponse
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return ChatWithAIResult{
			Message: "",
			Error:   fmt.Sprintf("Failed to parse response: %v", err),
		}
	}

	// Check middleware error
	if serverResponse.AuthMiddlewareReturn != nil && serverResponse.MiddlewareError != "" {
		return ChatWithAIResult{
			Message: "",
			Error:   serverResponse.MiddlewareError,
		}
	}

	// Check handler error
	if serverResponse.Error != "" {
		return ChatWithAIResult{
			Message: "",
			Error:   serverResponse.Error,
		}
	}

	return ChatWithAIResult{
		Message: serverResponse.AiMessage,
		Error:   "",
	}
}
