package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

var APP_PORT="3000"

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
      "ip_addr" : c.RealIP(),
    })
	})
	e.Logger.Fatal(e.Start(":" + APP_PORT))
}
