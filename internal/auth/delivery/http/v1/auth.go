package v1

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/auth"
	"github.com/vnnyx/employee-management/internal/dtos"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
)

type AuthHandler struct {
	uc auth.UseCase
}

func NewAuthHandler(uc auth.UseCase) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}

// @Summary      Login
// @Description  User login to get access token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dtos.LoginRequest true "Login Request"
// @Success      200 {object} dtos.Response{data=map[string]string} "Access Token Response"
// @Failure      400 {object} apperror.Error "Bad Request"
// @Router       /v1/auth/login [POST]
// @Security     NoAuth
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	ctx, span := instrumentation.NewTraceSpan(
		c.UserContext(),
		"AuthHandler.Login()",
	)
	defer span.End()

	var request dtos.LoginRequest
	if err := c.BodyParser(&request); err != nil {
		return errors.Wrap(err, "AuthHandler().Login().BodyParser()")
	}

	err := request.Validate()
	if err != nil {
		return errors.Wrap(err, "AuthHandler().Login().Validate()")
	}

	accessToken, err := h.uc.Login(ctx, request.Username, request.Password)
	if err != nil {
		return errors.Wrap(err, "AuthHandler().Login().uc.Login()")
	}

	return c.Status(http.StatusOK).JSON(
		dtos.Response{
			RequestID: uuid.NewString(),
			Data: map[string]any{
				"access_token": accessToken,
			},
		},
	)
}
