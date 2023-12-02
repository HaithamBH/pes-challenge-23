package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitializeServer() {
	app := echo.New()

	app.Use(middleware.Logger())
	app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	app.POST("/register", Register)
	app.POST("/login", Login)
	// app.GET("/checkToken", CheckToken)

	// protectedRoutes := app.Group("/")
	// protectedRoutes.Use(middlewareFc.Authentificate)

	app.Logger.Fatal(app.Start(":5000"))
}
