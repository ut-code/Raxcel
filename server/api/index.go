package api

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/ut-code/Raxcel/server/routes"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// load environmental variables from .env in developmental environments.
	if os.Getenv("VERCEL_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	e := echo.New()

	e.GET("/", routes.Greet)
	e.POST("/messages", routes.ChatWithAI)
	e.ServeHTTP(w, r)
}
