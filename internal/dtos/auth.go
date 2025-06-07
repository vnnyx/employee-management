package dtos

import "github.com/invopop/validation"

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *LoginRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Username, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}
