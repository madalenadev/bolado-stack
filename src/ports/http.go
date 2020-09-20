package ports

import (
	"bolado-stack/src/handlers"
	"bolado-stack/src/services"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HTTPPortConfig struct {
	Services   services.Container
	Debug      bool
	HideBanner bool
}

// SetupHTTPServer creates a http service
func SetupHTTPServer(config HTTPPortConfig) {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodOptions, http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{"*"},
	}))

	e.Use(middleware.BodyLimit("2GB"))
	e.HideBanner = config.HideBanner
	e.Debug = config.Debug

	setupRoutes(e, config)

	e.Start(":1234")
}

func setupRoutes(e *echo.Echo, config HTTPPortConfig) {
	// setup handlers
	usersHandler := handlers.NewUserHandler(config.Services)

	usersGroup := e.Group("users")
	usersGroup.GET("/:id", usersHandler.ReadOne)
}
