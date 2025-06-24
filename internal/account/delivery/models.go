package delivery

import "github.com/go-playground/validator/v10"

type RegisterRequest struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
}

func (r RegisterRequest) Validate() error {
	validator := validator.New()
	return validator.Struct(r)
}

type TokenResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (r LoginRequest) Validate() error {
	validator := validator.New()
	return validator.Struct(r)
}
