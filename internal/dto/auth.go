package dto

type LoginUserInputDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserOutputDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenOutputDTO struct {
	AccessToken string `json:"access_token"`
}
