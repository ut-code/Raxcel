package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
