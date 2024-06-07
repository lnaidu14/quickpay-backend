package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Id       string `json:"id" validate:"required,uuid4"`
	Username string `json:"username" validate:"required,min=3,max=25"`
	Ph       string `json:"ph" validate:"required,e164"`
}

var ExampleUser = User{"c85b4c8e-07ab-4c02-849d-71d495d6f905", "Foobar", "+11234567890"}

var validate *validator.Validate

func ValidateUser(user User) string {
	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(user)

	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return "User is valid"
}
