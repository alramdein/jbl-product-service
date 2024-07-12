package handler

import (
	"net/http"
	"referral-system/model"
	"referral-system/usecase"

	_ "referral-system/docs"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUsecase usecase.IUserUsecase
}

func NewUserHandler(userUseCase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{UserUsecase: userUseCase}
}

// @Summary Register User Generator
// @Description Register a user who can generate referral links.
// @Tags 		User
// @Accept json
// @Produce json
// @Param email formData string true "User's email"
// @Param password formData string true "User's password"
// @Success 200 {object} model.RegisterUserGeneratorResponse
// @Failure 400 {object} CustomError
// @Failure 409 {object} CustomError
// @Failure 404  {object}  CustomError
// @Failure 412 {object} CustomError
// @Failure 404 {object} CustomError
// @Failure 500 {object} CustomError
// @Router /register [post]
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

// RegisterUserContributor godoc
// @Summary      Register User Contributor
// @Description  Register a user who can contribute with a referral code.
// @Tags 		 User
// @Accept       json
// @Produce      json
// @Param        code   path      string  true  "Referral code"
// @Param        email   formData   string  true  "User's email"
// @Success      200  {object}  CustomResponse
// @Failure      400  {object} 	CustomError
// @Failure      404  {object}  CustomError
// @Failure      409  {object} 	CustomError
// @Failure      412  {object} 	CustomError
// @Failure      404  {object} 	CustomError
// @Failure      500  {object}  CustomError
// @Router       /register/{code} [post]
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

// Login godoc
// @Summary      User Login
// @Description  Authenticate a user and generate a JWT token.
// @Tags Auth
// @Accept       json
// @Produce      json
// @Param        email   formData   string  true  "User's email"
// @Param        password   formData   string  true  "User's password"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  CustomError
// @Failure      404  {object}  CustomError
// @Failure      401  {object}  CustomError
// @Failure      500  {object}  CustomError
// @Router       /login [post]
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
