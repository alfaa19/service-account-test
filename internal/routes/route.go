package routes

import (
	"github.com/alfaa19/service-account-test/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(h handler.AccountHandler, e *echo.Echo) {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Account routes
	e.POST("/daftar", h.CreateAccount)
	e.GET("/saldo/:noRekening", h.GetSaldo)
	e.POST("/tarik", h.Withdraw)
	e.POST("/tabung", h.Deposit)

}
