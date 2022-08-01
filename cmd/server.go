package main

import (
	"fmt"
	"keeper/internal/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.New()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", cfg.Port)))
}
