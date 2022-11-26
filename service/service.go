package service

import (
	"femalegeek/repository/model"
	"net/http"

	"github.com/labstack/echo"
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
	userUsecase  model.UserUsecase
	eventUsecase model.EventUsecase
}

// NewHTTPService creates intance of service
func NewHTTPService(userUsecase model.UserUsecase, eventUsecase model.EventUsecase) *HTTPService {
	return &HTTPService{
		userUsecase:  userUsecase,
		eventUsecase: eventUsecase,
	}
}

// Routes handle routing for http service
func (s *HTTPService) Routes(route *echo.Echo) {
	route.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "FemaleGeek Web")
	})
	route.GET("user/:id/", s.FindUserByID)
	route.GET("event/:id/", s.FindEventByID)
}
