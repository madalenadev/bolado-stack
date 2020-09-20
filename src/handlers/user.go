package handlers

import (
	"bolado-stack/src/services"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	services services.Container
}

func NewUserHandler(services services.Container) UserHandler {
	return UserHandler{services}
}

func (uh UserHandler) ReadOne(c echo.Context) error {
	return nil
}
