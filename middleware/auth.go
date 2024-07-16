package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// UserData struct represents the user data extracted from the JWT token
type UserData struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    // Add other fields as needed
}

// JWTMiddleware validates JWT token and extracts user information
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        excludedPaths := []string{"/api/login", "/api/register"}

        // Exclude specific paths from token validation
        for _, path := range excludedPaths {
            if c.Path() == path {
                return next(c)
            }
        }

        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Authorization header missing or invalid"})
        }

        // Send token validation request to authentication service
        authURL := os.Getenv("AUTH_SERVICE_URL")
		if authURL == "" {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Auth service URL not set"})
		}
        client := &http.Client{}
        req, err := http.NewRequest("POST", authURL, nil)
        if err != nil {
			log.Error(err)
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error creating token validation request"})
        }
        req.Header.Set("Authorization", authHeader)

        resp, err := client.Do(req)
        if err != nil {
			log.Error(err)
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Error validating token"})
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
            return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Invalid token or authentication service error"})
        }

        // Decode user data from response
        var userData UserData
        if err := json.NewDecoder(resp.Body).Decode(&userData); err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to decode user data"})
        }

		fmt.Println(userData)

        // Add user information to context or wherever needed
        c.Set("user", userData)

        return next(c)
    }
}
