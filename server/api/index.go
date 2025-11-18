package api

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/ut-code/Raxcel/server/db"
	middleware "github.com/ut-code/Raxcel/server/middlewares"
	"github.com/ut-code/Raxcel/server/routes"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// load environmental variables from .env in developmental environments.
	if os.Getenv("VERCEL_ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
		db.Migrate()
	}
	e := echo.New()

	e.GET("/", routes.Greet)
	messages := e.Group("/messages")
	messages.Use(middleware.AuthMiddleware)
	messages.POST("/", routes.ChatWithAI)
	e.GET("/user", routes.CheckUser)
	e.POST("/register", routes.Register)
	e.GET("/verify-email", routes.VerifyEmail)
	e.POST("/login", routes.Login)
	e.ServeHTTP(w, r)
}
