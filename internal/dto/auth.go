package dto

type LoginUserInputDTO struct {
	Email    string `json:"email" validate:"required,email" swaggertype:"string" example:"me@gmail.com"`
	Password string `json:"password" validate:"required,password" swaggertype:"string" example:"********"`
}

type LoginUserOutputDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenOutputDTO struct {
	AccessToken string `json:"access_token"`
}
