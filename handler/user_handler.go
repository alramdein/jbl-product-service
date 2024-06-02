package handler

import (
	"net/http"
	"referral-system/model"
	"referral-system/usecase"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUsecase usecase.IUserUsecase
}

func NewUserHandler(userUseCase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: userUseCase}
}

func (h *UserHandler) RegisterUserGenerator(c echo.Context) error {
	var req model.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	userResp, err := h.UserUsecase.RegisterUserGenerator(c.Request().Context(), model.RegisterUserGeneratorInput{
		Email:    req.Email,
		Password: req.Password,
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

	_, err := h.UserUsecase.RegisterUserContributor(c.Request().Context(), model.RegisterUserContributorInput{
		Email:        req.Email,
		ReferralCode: code,
	})
	if err != nil {
		return MapErrorToHTTPResponse(err)
	}

	return c.JSON(http.StatusOK, CustomResponse{
		Message: "Successfully submited the referral",
	})
}

func (h *UserHandler) Login(c echo.Context) error {
	var req model.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	token, err := h.UserUsecase.Login(c.Request().Context(), req.Email, req.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}
