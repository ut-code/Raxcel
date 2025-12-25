package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zalando/go-keyring"
)

type ChatRequest struct {
	Message            string `json:"message"`
	SpreadsheetContext string `json:"spreadsheetContext,omitempty"`
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

func (a *App) ChatWithAI(message string, spreadsheetContext string) ChatResult {
	fmt.Println("ChatWithAI called with message:", message)
	fmt.Printf("Spreadsheet context length: %d characters\n", len(spreadsheetContext))

	postData := ChatRequest{
		Message:            message,
		SpreadsheetContext: spreadsheetContext,
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
