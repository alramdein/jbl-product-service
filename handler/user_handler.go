package handler

import (
	"net/http"
	"referral-system/model"
	"referral-system/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUseCase usecase.IUserUsecase
}

func NewUserHandler(userUseCase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{UserUseCase: userUseCase}
}

func (h *UserHandler) RegisterUserGenerator(c echo.Context) error {
	var req model.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	userResp, err := h.UserUseCase.RegisterUser(c.Request().Context(), model.RegisterUserInput{
		Email:    req.Email,
		Password: req.Password,
		Role:     model.GeneratorRole,
	})
	if err != nil {
		return MapErrorToHTTPResponse(err)
	}

	return c.JSON(http.StatusOK, userResp)
}

func (h *UserHandler) RegisterUserContributor(c echo.Context) error {
	code := c.Param("code")
	var req model.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	userResp, err := h.UserUseCase.RegisterUser(c.Request().Context(), model.RegisterUserInput{
		Email:        req.Email,
		Password:     req.Password,
		Role:         model.ContributorRole,
		ReferralCode: code,
	})
	if err != nil {
		return MapErrorToHTTPResponse(err)
	}

	return c.JSON(http.StatusOK, userResp)
}
