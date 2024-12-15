package input

type UserRegisterInput struct {
	FirstName       string `json:"first_name" validate:"required"`
	MiddleName      string `json:"middle_name" `
	LastName        string `json:"last_name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}
