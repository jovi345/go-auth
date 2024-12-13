package input

type UserRegisterInput struct {
	FirstName       string `json:"first_name" binding:"required"`
	MiddleName      string `json:"middle_name"`
	LastName        string `json:"last_name" binding:"required,email"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
