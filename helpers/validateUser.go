package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID       string `json:"id" validate:"required,uuid4"`
	Username string `json:"username" validate:"required,min=3,max=25"`
	PH       string `json:"ph" validate:"required,e164"`
}

var validate *validator.Validate

func ValidateUser(user User) {
	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(user)

	if err != nil {
		fmt.Println("User is invalid")
		fmt.Println("Error: ", err)
	}

}
