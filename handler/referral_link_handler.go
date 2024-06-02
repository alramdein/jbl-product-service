package handler

import (
	"net/http"

	"referral-system/usecase"

	"github.com/labstack/echo/v4"
)

type ReferralHandler struct {
	ReferralLinkUsecase usecase.ReferralLinkUsecase
}

func NewReferralHandler(referralLinkUsecase usecase.ReferralLinkUsecase) *ReferralHandler {
	return &ReferralHandler{
		ReferralLinkUsecase: referralLinkUsecase,
	}
}

func (h *ReferralHandler) GenerateReferralLink(c echo.Context) error {
	userID := c.Get("user_id").(string)

	referralLink, err := h.ReferralLinkUsecase.GenerateReferralLink(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to generate referral link")
	}

	return c.JSON(http.StatusOK, referralLink)
}
