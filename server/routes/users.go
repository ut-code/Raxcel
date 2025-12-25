package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetCurrentUserId(c echo.Context) error {
	userId, ok := c.Get("userId").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Failed to get userId from context",
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"userId": userId,
	})
}
