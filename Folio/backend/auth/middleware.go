package auth

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTMiddleware validates JWT tokens and extracts user info
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "missing authorization header",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid authorization header format",
			})
		}

		tokenString := parts[1]
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"error": "invalid or expired token",
			})
		}

		// Store user info in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("is_guest", claims.IsGuest)

		return next(c)
	}
}

// GetUserID extracts user ID from context
func GetUserID(c echo.Context) string {
	userID, _ := c.Get("user_id").(string)
	return userID
}

// GetUserEmail extracts user email from context
func GetUserEmail(c echo.Context) string {
	email, _ := c.Get("user_email").(string)
	return email
}

// IsGuestUser checks if the current user is a guest
func IsGuestUser(c echo.Context) bool {
	isGuest, _ := c.Get("is_guest").(bool)
	return isGuest
}

// OptionalJWTMiddleware validates JWT tokens if present, but doesn't require them
func OptionalJWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			// Extract token from "Bearer <token>"
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				claims, err := ValidateJWT(tokenString)
				if err == nil {
					// Store user info in context if token is valid
					c.Set("user_id", claims.UserID)
					c.Set("user_email", claims.Email)
					c.Set("is_guest", claims.IsGuest)
				}
			}
		}
		
		return next(c)
	}
}

