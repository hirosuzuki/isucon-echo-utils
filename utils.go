package isuconechoutils

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Setup(e *echo.Echo) {
	e.Use(middleware.RequestID())
	e.Use(UniqueID)
	e.Use(Logger())
	e.GET("/_/routes", ShowRoutes)
	e.GET("/_/kataribe", ShowKataribe)
}
