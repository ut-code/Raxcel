package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var e *echo.Echo

func init() {
	e = echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Echo!")
	})
}

func Handler(w http.ResponseWriter, r *http.Request) {
	e.ServeHTTP(w, r)
}

func Local() {
	e.Logger.Fatal(e.Start(":1323"))
}
