package dto

type LoginUserInputDTO struct {
	Email    string `json:"email" validate:"required,email" swaggertype:"string" example:"me@gmail.com"`
	Password string `json:"password" validate:"required" swaggertype:"string" example:"********"`
}

type LoginUserOutputDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenOutputDTO struct {
	AccessToken string `json:"access_token"`
}

type ForgotPasswordInputDTO struct {
	Email string `json:"email" validate:"required,email" swaggertype:"string" example:"user@gmail.com"`
}

type ResetPasswordInputDTO struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,password" swaggertype:"string" example:"********"`
}
