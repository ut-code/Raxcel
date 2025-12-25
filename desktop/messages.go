package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

	var serverResponse map[string]interface{}
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return LoadChatHistoryResult{
			Messages: []Mesaage{},
			Error:    fmt.Sprint(err),
		}
	}

	// Parse messages array
	var messages []Mesaage
	if messagesData, ok := serverResponse["messages"].([]interface{}); ok {
		for _, msgData := range messagesData {
			if msgMap, ok := msgData.(map[string]interface{}); ok {
				message := Mesaage{
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

	return LoadChatHistoryResult{
		Messages: messages,
		Error:    "",
	}
}

type ChatRequest struct {
	Message            string `json:"message"`
	SpreadsheetContext string `json:"spreadsheetContext,omitempty"`
}

type ChatWithAIResult struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (a *App) ChatWithAI(message string, spreadsheetContext string) ChatWithAIResult {
	postData := ChatRequest{
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
	var serverResponse map[string]string
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		fmt.Println("Unmarshal error:", err)
		fmt.Println("Response body:", string(body))
		return ChatWithAIResult{
			Message: "",
			Error:   fmt.Sprint(err),
		}
	}
	return ChatWithAIResult{
		Message: serverResponse["aiMessage"],
		Error:   "",
	}
}
