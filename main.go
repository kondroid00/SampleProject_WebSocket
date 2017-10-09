package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func serve(c echo.Context) error {
	err := RoomManagerInstance().Serve(c)
	return err
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/:id", serve)
	e.Logger.Fatal(e.Start("0.0.0.0:1323"))
}
