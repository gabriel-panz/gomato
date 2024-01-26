package main

import (
	"github.com/gabriel-panz/gomato/handler"
	"github.com/labstack/echo/v4"
)

func main() {
	app := echo.New()

	app.GET("/", handler.HandleShowHome)
	app.GET("/focus", handler.HandleShowFocus)
	app.POST("/pause", handler.HandleShowPause)

	app.Logger.Fatal(app.Start(":9000"))
}
