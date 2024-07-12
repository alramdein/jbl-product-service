package middleware

import (
	"net/http"
	"strings"

	"referral-system/handler"
	"referral-system/util"

	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, handler.CustomError{
					StatusCode: http.StatusUnauthorized,
					Message:    "missing or invalid token",
				})
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := util.ValidateJWT(tokenString, secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, handler.CustomError{
					StatusCode: http.StatusUnauthorized,
					Message:    "invalid token",
				})
			}

			c.Set("user_id", claims["user_id"])
			c.Set("role_id", claims["role_id"])

			return next(c)
		}
	}
}
