package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zalando/go-keyring"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResult struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	UserId  string `json:"userId,omitempty"`
}

func (a *App) Signup(email, password string) RegisterResult {
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

	resp, err := http.Post(fmt.Sprintf("%s/signup", apiUrl), "application/json", bytes.NewReader(jsonData))
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

func (a *App) Signin(email, password string) LoginResult {
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

	resp, err := http.Post(fmt.Sprintf("%s/signin", apiUrl), "application/json", bytes.NewReader(jsonData))
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

func (a *App) GetCurrentUser() CheckResult {
	apiUrl := getAPIURL()
	token, _ := keyring.Get("Raxcel", "raxcel-user")
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/me", apiUrl), nil)
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

type SignOutResult struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func (a *App) SignOut() SignOutResult {
	err := keyring.Delete("Raxcel", "raxcel-user")
	if err != nil {
		return SignOutResult{
			Ok:    false,
			Error: fmt.Sprintf("Failed to sign out: %v", err),
		}
	}
	return SignOutResult{
		Ok: true,
	}
}
