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

type SignupResult struct {
	UserId string `json:"userId"`
	Error  string `json:"error"`
}

func (a *App) Signup(email, password string) SignupResult {
	postData := RegisterRequest{
		Email:    email,
		Password: password,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return SignupResult{
			UserId: "",
			Error:  fmt.Sprintf("failed to marshal request: %v", err),
		}
	}
	apiUrl := getAPIURL()

	resp, err := http.Post(fmt.Sprintf("%s/auth/signup", apiUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return SignupResult{
			UserId: "",
			Error:  fmt.Sprintf("Failed to send request: %v", err),
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SignupResult{
			UserId: "",
			Error:  fmt.Sprintf("Failed to read response: %v", err),
		}
	}
	var serverResponse map[string]string
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return SignupResult{
			UserId: "",
			Error:  fmt.Sprintf("Failed to parse response: %v", err),
		}
	}
	if resp.StatusCode != http.StatusCreated {
		return SignupResult{
			UserId: "",
			Error:  serverResponse["error"],
		}
	}
	return SignupResult{
		UserId: serverResponse["userId"],
		Error:  "",
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SigninResult struct {
	Token string `json:"token"`
	Error string `json:"error"`
}

func (a *App) Signin(email, password string) SigninResult {
	postData := LoginRequest{
		Email:    email,
		Password: password,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		return SigninResult{
			Token: "",
			Error: fmt.Sprintf("Failed to marshal request: %v", err),
		}
	}
	apiUrl := getAPIURL()

	resp, err := http.Post(fmt.Sprintf("%s/auth/signin", apiUrl), "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return SigninResult{
			Token: "",
			Error: fmt.Sprintf("Failed to send request: %v", err),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return SigninResult{
			Token: "",
			Error: fmt.Sprintf("Failed to read response: %v", err),
		}
	}

	var serverResponse map[string]string
	err = json.Unmarshal(body, &serverResponse)
	if err != nil {
		return SigninResult{
			Token: "",
			Error: fmt.Sprintf("Failed to parse response %v", err),
		}
	}
	if resp.StatusCode != http.StatusOK {
		return SigninResult{
			Token: "",
			Error: serverResponse["error"],
		}
	}
	token := serverResponse["token"]
	err = keyring.Set("Raxcel", "raxcel-user", token)
	if err != nil {
		return SigninResult{
			Token: "",
			Error: fmt.Sprintf("Failed to store token: %v", err),
		}
	}
	return SigninResult{
		Token: token,
		Error: "",
	}
}

type GetCurrentUserResult struct {
	UserId string `json:"userId"`
	Error  string `json:"error"`
}

func (a *App) GetCurrentUser() GetCurrentUserResult {
	apiUrl := getAPIURL()
	token, _ := keyring.Get("Raxcel", "raxcel-user")
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/me", apiUrl), nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return GetCurrentUserResult{
			UserId: "",
			Error:  fmt.Sprint(err),
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetCurrentUserResult{
			UserId: "",
			Error:  fmt.Sprint(err),
		}
	}
	var serverResponse map[string]string
	if err := json.Unmarshal(body, &serverResponse); err != nil {
		return GetCurrentUserResult{
			UserId: "",
			Error:  fmt.Sprint(err),
		}
	}
	if serverResponse["error"] != "" {
		return GetCurrentUserResult{
			UserId: "",
			Error:  serverResponse["error"],
		}
	}
	return GetCurrentUserResult{
		UserId: serverResponse["userId"],
		Error:  "",
	}
}

type SignOutResult struct {
	Error string `json:"error"`
}

func (a *App) SignOut() SignOutResult {
	err := keyring.Delete("Raxcel", "raxcel-user")
	if err != nil {
		return SignOutResult{
			Error: fmt.Sprintf("Failed to sign out: %v", err),
		}
	}
	return SignOutResult{
		Error: "",
	}
}
