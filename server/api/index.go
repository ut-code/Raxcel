package api

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/ut-code/Raxcel/server/controllers"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// load environmental variables from .env in developmental environments.
	if os.Getenv("VERCEL_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	e := echo.New()

	e.GET("/", controllers.Greet)
	e.POST("/messages", controllers.ChatWithAI)
	e.ServeHTTP(w, r)
}
