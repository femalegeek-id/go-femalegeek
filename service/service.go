package service

import (
	"femalegeek/repository/model"
	"femalegeek/usecase"
	"net/http"

	"github.com/kumparan/go-utils"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// NotFoundResponse returns when bad response return
func NotFoundResponse() *model.ErrorResponse {
	errBadResponse := &model.ErrorResponse{
		Code:    0,
		Message: "Tidak ditemukan.",
	}

	return errBadResponse
}

// AnotherErrorResponse returns when bad response return
func AnotherErrorResponse(e string) *model.ErrorResponse {
	errResponse := &model.ErrorResponse{
		Code:    0,
		Message: e,
	}

	return errResponse
}

// Service defines usecases http service
type Service interface {
	Routes(route *echo.Echo)
}

// HTTPService implements Service
type HTTPService struct {
	userUsecase model.UserUsecase
}

// FindUserByID find user detail
func (s *HTTPService) FindUserByID(e echo.Context) error {
	var user *model.User
	ctx := e.Request().Context()
	id := e.Param("id")
	if id == "" {
		return e.JSON(http.StatusBadRequest, AnotherErrorResponse("Mohon masukkan id"))
	}

	logger := log.WithFields(log.Fields{
		"id": id,
	})

	user, err := s.userUsecase.FindByID(ctx, utils.StringToInt64(id))
	switch {
	case err == usecase.ErrNotFound:
		return e.JSON(http.StatusNotFound, NotFoundResponse())
	case err != nil:
		logger.Error(err)
		return e.JSON(http.StatusInternalServerError, AnotherErrorResponse(err.Error()))
	case user == nil:
		return e.JSON(http.StatusNotFound, NotFoundResponse())
	}

	return e.JSON(http.StatusOK, user)
}

// NewHTTPService creates intance of service
func NewHTTPService(userUsecase model.UserUsecase) *HTTPService {
	return &HTTPService{
		userUsecase: userUsecase,
	}
}

// Routes handle routing for http service
func (s *HTTPService) Routes(route *echo.Echo) {
	route.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "FemaleGeek Web")
	})
	route.GET("user/:id/", s.FindUserByID)
}
