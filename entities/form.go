package entities

type JsonForm struct {
	Name   string `json:"name" validate:"required"`
	Email  string `json:"email" validate:"required,email"`
	Phone  string `json:"phone" validate:"required,phone"`
	Source string `json:"source"`
}
