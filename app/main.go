package main

import (
	_handler "avataro/avataro/delivery/http"
	_usecase "avataro/avataro/usecase"
	"avataro/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	e *echo.Echo
)

func init() {
	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	config.InitializeConfig()

}

func main() {

	hushhubUsecase := _usecase.AvataroUsecase()
	_handler.NewAvataroHandler(e, hushhubUsecase)

	e.Logger.Fatal(e.Start(":1111"))
}
