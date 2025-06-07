package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	config "github.com/vnnyx/employee-management/config/api"
	"github.com/vnnyx/employee-management/internal/auth/entity"
	"github.com/vnnyx/employee-management/internal/constants"
)

func Auth(cfg config.Config) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			return jwt.ErrTokenUnverifiable
		}

		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if authHeader == "" {
			return jwt.ErrTokenUnverifiable
		}

		token, err := jwt.ParseWithClaims(tokenString, &entity.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.App.Key), nil
		})
		if err != nil {
			return jwt.ErrTokenMalformed
		}

		if claims, ok := token.Claims.(*entity.AccessTokenClaims); ok && token.Valid {
			if claims.ExpiresAt.Time.Before(c.Context().Time()) {
				return jwt.ErrTokenExpired
			}

			credential := entity.Credential{
				UserID:    claims.UserID,
				Username:  claims.Username,
				IsAdmin:   &claims.IsAdmin,
				IPAddress: c.IP(),
				RequestID: uuid.NewString(),
			}
			c.Locals(constants.KeyAuthCredential, credential)

			ctx := context.WithValue(c.UserContext(), constants.KeyFiberCtxInformation, entity.FiberCtxInformation{Method: c.Method(), OriginalURL: c.OriginalURL(), Enable: *cfg.Logger.Enable})
			c.SetUserContext(ctx)

			return c.Next()
		} else {
			return jwt.ErrTokenInvalidClaims
		}
	}
}
