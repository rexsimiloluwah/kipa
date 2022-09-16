package dto

type CreateUserInputDTO struct {
	Firstname string `json:"firstname" validate:"required,min=2" swaggertype:"string" example:"Similoluwa"`
	Lastname  string `json:"lastname" validate:"required,min=2" swaggertype:"string" example:"Okunowo"`
	Email     string `json:"email" validate:"required,email" swaggertype:"string" example:"me@gmail.com"`
	Username  string `json:"username" validate:"required,min=2" swaggertype:"string" example:""`
	Password  string `json:"password" validate:"required,password" swaggertype:"string" example:"********"`
}

type UpdateUserInputDTO struct {
	Firstname string `json:"firstname" validate:"required,min=2" swaggertype:"string" example:"Similoluwa"`
	Lastname  string `json:"lastname" validate:"required,min=2" swaggertype:"string" example:"Okunowo"`
	Email     string `json:"email" validate:"required,email" swaggertype:"string" example:"me@gmail.com"`
	Username  string `json:"username" validate:"required,min=2" swaggertype:"string" example:""`
}

type UpdateUserPasswordInputDTO struct {
	Password string `json:"password" validate:"required,password" swaggertype:"string" example:"********"`
}
