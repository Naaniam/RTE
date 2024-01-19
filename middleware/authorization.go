package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AdminAuthorize(jwtSecret []byte, next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := ExtractTokenFromHeader(c)
		if tokenString == "" {
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		isAdmin, ok := claims["admin"].(string)
		if !ok || strings.ToLower(isAdmin) != "admin" {
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Locals("token", token)
		return next(c)
	}
}

func MemberAuthorize(jwtSecret []byte, next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := ExtractTokenFromHeader(c)
		if tokenString == "" {
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return fiber.NewError(http.StatusInternalServerError, err.Error())
		}

		isAdmin, ok := claims["member"].(string)
		if !ok || strings.ToLower(isAdmin) != "member" {
			return fiber.NewError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Locals("token", token)
		return next(c)
	}
}

// ExtractTokenFromHeader function is used to extract the token part after the Bearer part
func ExtractTokenFromHeader(c *fiber.Ctx) string {
	header := c.Get("Authorization")
	if header == "" {
		return ""
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}
