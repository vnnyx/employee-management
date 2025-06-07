package entity

import (
	"github.com/invopop/validation"
)

type Credential struct {
	Username  string
	UserID    string
	IsAdmin   *bool
	IPAddress string
	RequestID string
}

func (c Credential) Validate() error {
	return validation.ValidateStruct(
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.UserID, validation.Required),
		validation.Field(&c.IsAdmin, validation.NotNil),
		validation.Field(&c.IPAddress, validation.Required, validation.Length(7, 15)),
		validation.Field(&c.RequestID, validation.Required),
	)
}

type FiberCtxInformation struct {
	Method, OriginalURL string
	Enable              bool
}
