package rest

import (
	"net/http"

	"github.com/gabapcia/gameblitz/internal/auth"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrorResponseMissingAuthCredentials = ErrorResponse{Code: "7.0", Message: "missing credentials"}
	ErrorResponseInvalidAuthCredentials = ErrorResponse{Code: "7.1", Message: "invalid credentials"}
)

func buildAuthMiddleware(authenticateFunc auth.AuthenticateFunc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := string(c.Request().Header.Peek("Authorization"))
		if authorization == "" {
			return c.Status(http.StatusUnauthorized).JSON(ErrorResponseMissingAuthCredentials)
		}

		claims, err := authenticateFunc(c.Context(), authorization)
		if err != nil {
			return err
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}
