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

// GenerateReferralLink godoc
// @Summary      Generate Referral Link
// @Description  To generate new referral link and expire the old one
// @Tags 		 Referral Link
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  CustomError
// @Failure      401  {object}  CustomError
// @Failure      404  {object}  CustomError
// @Failure      500  {object}  CustomError
// @Router       /referral-link [post]
func (h *ReferralHandler) GenerateReferralLink(c echo.Context) error {
	userID := c.Get("user_id").(string)

	referralLink, err := h.ReferralLinkUsecase.GenerateReferralLink(c.Request().Context(), userID)
	if err != nil {
		return MapErrorToHTTPResponse(err)
	}

	return c.JSON(http.StatusOK, referralLink)
}
