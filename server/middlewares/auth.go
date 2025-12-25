package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/ut-code/Raxcel/server/utils"
)

type AuthMiddlewareReturnType struct {
	Error string `json:"error"`
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, AuthMiddlewareReturnType{
				Error: "missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.JSON(http.StatusUnauthorized, AuthMiddlewareReturnType{
				Error: "invalid authorization format",
			})
		}

		userId, err := utils.ValidateJWT(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, AuthMiddlewareReturnType{
				Error: "invalid token",
			})
		}

		c.Set("userId", userId)
		return next(c)
	}
}
