package dto

type UserInputDTO struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserOutputDTO struct {
	Name string `json:"name"`
}
