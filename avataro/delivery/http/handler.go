package http

import (
	"avataro/config"
	"avataro/domain"
	"github.com/labstack/echo/v4"
)

type avataroHandler struct {
	avataroUsecase domain.AvataroUsecase
}

func NewAvataroHandler(e *echo.Echo, avataroUsecase domain.AvataroUsecase) {
	handler := &avataroHandler{
		avataroUsecase: avataroUsecase,
	}

	e.GET("/avataro", handler.getAvataro)
}

func (h *avataroHandler) getAvataro(c echo.Context) error {

	text := c.QueryParam("text")
	if text == "" {
		return c.JSON(400, "text is required")
	}

	_set := c.QueryParam("set")
	var set *string
	if _set != "" {
		if _, ok := config.Sets[_set]; !ok {
			return c.JSON(400, "set is invalid")
		}
		set = &_set
	}

	_backgroundSet := c.QueryParam("backgroundSet")
	var backgroundSet *string
	if _backgroundSet != "" {
		if _, ok := config.Backgrounds[_backgroundSet]; !ok {
			return c.JSON(400, "backgroundSet is invalid")
		}
		backgroundSet = &_backgroundSet
	}

	image := h.avataroUsecase.GetAvataro(text, set, backgroundSet)
	if image == nil {
		return c.JSON(500, "error")
	}

	return c.Blob(200, "image/png", image)
}
