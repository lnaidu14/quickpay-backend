package helpers

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"quickpay/main/types"
)

var validate *validator.Validate

func ValidateUser(user types.User) string {
	validate = validator.New(validator.WithRequiredStructEnabled())

	err := validate.Struct(user)

	if err != nil {
		return fmt.Sprintf("%s", err)
	}
	return "User is valid"
}
