package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetCurrentUserResponse struct {
	Error  string `json:"error,omitempty"`
	UserId string `json:"userId,omitempty"`
}

func GetCurrentUser(c echo.Context) error {
	userId, ok := c.Get("userId").(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, GetCurrentUserResponse{
			Error: "Failed to get userId from context",
		})
	}
	return c.JSON(http.StatusOK, GetCurrentUserResponse{
		UserId: userId,
	})
}
