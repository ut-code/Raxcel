package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	middleware "github.com/ut-code/Raxcel/server/middlewares"
	"github.com/ut-code/Raxcel/server/routes"
)

func SetupRouter() *echo.Echo {
	e := echo.New()

	e.GET("/", routes.Greet)
	messages := e.Group("/messages")
	messages.Use(middleware.AuthMiddleware)
	messages.POST("", routes.ChatWithAI)
	e.GET("/user", routes.CheckUser)
	e.POST("/register", routes.Register)
	e.GET("/verify-email", routes.VerifyEmail)
	e.POST("/login", routes.Login)
	return e
}

func VercelHandler(w http.ResponseWriter, r *http.Request) {
	e := SetupRouter()
	e.ServeHTTP(w, r)
}
