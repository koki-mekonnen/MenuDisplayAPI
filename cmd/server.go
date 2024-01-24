package main

import (
	"foodorderapi/internals/config"
	"foodorderapi/routes"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()

	config.Databaseinit()

	routes.Foodorderroutes(e)

	e.Logger.Fatal(e.Start(":14500"))

}
