package service

import (
	"femalegeek/repository/model"
	"femalegeek/usecase"
	"net/http"

	"github.com/kumparan/go-utils"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// FindEventByID find event detail
func (s *HTTPService) FindEventByID(e echo.Context) error {
	var event *model.Event
	ctx := e.Request().Context()
	id := e.Param("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, AnotherErrorResponse("Mohon masukkan id"))
	}

	logger := log.WithFields(log.Fields{
		"id": id,
	})

	event, err := s.eventUsecase.FindByID(ctx, utils.StringToInt64(id))
	switch {
	case err == usecase.ErrNotFound:
		return e.JSON(http.StatusNotFound, NotFoundResponse())
	case err != nil:
		logger.Error(err)
		return e.JSON(http.StatusInternalServerError, AnotherErrorResponse(err.Error()))
	case event == nil:
		return e.JSON(http.StatusNotFound, NotFoundResponse())
	}

	return e.JSON(http.StatusOK, event)
}
