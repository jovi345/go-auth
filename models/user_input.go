package models

type UserRegisterInput struct {
	FirstName       string `json:"first_name" validate:"required"`
	MiddleName      string `json:"middle_name" `
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type UserLoginInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}
