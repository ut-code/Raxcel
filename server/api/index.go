package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	middleware "github.com/ut-code/Raxcel/server/middlewares"
	"github.com/ut-code/Raxcel/server/routes"
)

func SetupRouter() *echo.Echo {
	router := echo.New()

	router.GET("/", routes.Greet)

	messageGroup := router.Group("/messages")
	{
		messageGroup.Use(middleware.AuthMiddleware)
		messageGroup.POST("", routes.ChatWithAI)
		messageGroup.GET("", routes.LoadChatHistory)
	}

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", routes.Signup)
		authGroup.POST("/signin", routes.Signin)
		authGroup.GET("/verify-email", routes.VerifyEmail)
	}

	userGroup := router.Group("/users")
	{
		userGroup.Use(middleware.AuthMiddleware)
		userGroup.GET("/me", routes.GetCurrentUserId)
	}

	return router
}

func VercelHandler(w http.ResponseWriter, r *http.Request) {
	router := SetupRouter()
	router.ServeHTTP(w, r)
}
